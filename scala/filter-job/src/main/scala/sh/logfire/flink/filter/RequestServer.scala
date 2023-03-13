package sh.logfire.flink.filter

import io.confluent.kafka.schemaregistry.client.CachedSchemaRegistryClient
import io.grpc.protobuf.services.ProtoReflectionService
import io.grpc.stub.StreamObserver
import io.grpc.{ Server, ServerBuilder }
import org.apache.flink.client.program.rest.RestClusterClient
import org.apache.flink.configuration.Configuration
import org.apache.flink.runtime.highavailability.nonha.standalone.StandaloneClientHAServices
import org.apache.flink.runtime.rpc.FatalErrorHandler
import org.apache.flink.streaming.api.scala.StreamExecutionEnvironment
import org.apache.flink.table.api.bridge.scala.StreamTableEnvironment
import org.slf4j.{ Logger, LoggerFactory }
import redis.clients.jedis.Jedis
import scalapb.json4s.JsonFormat
import sh.logfire.flink.filter.RequestServer.jedis
import sh.logfire.request._

import java.util.UUID
import java.util.concurrent.{ ScheduledThreadPoolExecutor, TimeUnit }
import scala.annotation.tailrec
import scala.concurrent.{ ExecutionContext, Future }
import scala.jdk.CollectionConverters.collectionAsScalaIterableConverter
import scala.util.{ Failure, Success, Try }

object RequestServer {
  private val logger: Logger = LoggerFactory.getLogger(getClass.getSimpleName)
  private lazy val env: StreamExecutionEnvironment =
    StreamExecutionEnvironment.createRemoteEnvironment(Configs.FLINK_CLUSTER_HOST, Configs.FLINK_CLUSTER_PORT.toInt)
  private lazy val tEnv = StreamTableEnvironment.create(env)
  private lazy val jedis = {
    val jedis = new Jedis(Configs.REDIS_HOST, Configs.REDIS_PORT.toInt)
    Try(jedis.auth(Configs.REDIS_PASSWORD)) match {
      case Failure(exception) =>
      case Success(value)     =>
    }
    jedis
  }

  private val schemaRegistryAddress     = Configs.SCHEMA_REGISTRY_ADDRESS
  private val brokerAddress             = Configs.KAFKA_BROKERS
  private lazy val schemaRegistryClient = new CachedSchemaRegistryClient(schemaRegistryAddress, 1000)

  def main(args: Array[String]): Unit = {
    periodicallyCheckRedis()

    val server = new RequestServer(ExecutionContext.global)
    server.start()
    server.blockUntilShutdown()
  }

  private lazy val clusterClient: RestClusterClient[String] = {
    val clusterId         = ""
    val restServerAddress = s"http://${Configs.FLINK_CLUSTER_HOST}:${Configs.FLINK_CLUSTER_PORT}"
    new RestClusterClient[String](
      env.getConfiguration.asInstanceOf[Configuration],
      clusterId,
      (c: Configuration, e: FatalErrorHandler) => new StandaloneClientHAServices(restServerAddress))
  }

  private val port = 50051

  private def periodicallyCheckRedis() = {

    val ex = new ScheduledThreadPoolExecutor(1)
    val task = new Runnable {
      def run() =
        Try {
          jedis
            .hgetAll(Configs.FLINK_JOBS_REDIS_KEY)
            .entrySet()
            .asScala
            .foreach(e => {
              if (e.getValue.toLong < System.currentTimeMillis()) {
                val jobs = clusterClient.listJobs().get().asScala.filter(_.getJobName == e.getKey)
                if (jobs.nonEmpty) {
                  val jobID = jobs.head.getJobId
                  clusterClient.cancel(jobID)
                }
                jedis.hdel(Configs.FLINK_JOBS_REDIS_KEY, e.getValue)
              }
            })
        } match {
          case Failure(exception) => {
            logger.error("unable to list and cancel filter jobs")
            exception.printStackTrace()
          }
          case Success(value) =>
        }
    }
    val f = ex.scheduleAtFixedRate(task, 1, 120, TimeUnit.SECONDS)
  }

}

class RequestServer(executionContext: ExecutionContext) { self =>
  private[this] var server: Server = null

  private def start(): Unit = {
    server = ServerBuilder
      .forPort(RequestServer.port)
      .addService(FlinkServiceGrpc.bindService(new FlinkServiceImpl, executionContext))
      .addService(ProtoReflectionService.newInstance())
      .build
      .start
    RequestServer.logger.info("Server started, listening on " + RequestServer.port)
    sys.addShutdownHook {
      System.err.println("*** shutting down gRPC server since JVM is shutting down")
      self.stop()
      System.err.println("*** server shut down")
    }
  }

  private def stop(): Unit =
    if (server != null) {
      server.shutdown()
    }

  private def blockUntilShutdown(): Unit =
    if (server != null) {
      server.awaitTermination()
    }

  private class FlinkServiceImpl extends FlinkServiceGrpc.FlinkService {
    override def submitFilterRequest(filterRequest: FilterRequest,
                                     responseObserver: StreamObserver[FilterResponse]): Unit = {
      if (filterRequest.sources.size < 1)
        responseObserver.onError(new Exception("There should be at least one source"))

      RequestServer.logger.info("received filter request: %s".formatted(JsonFormat.toJsonString(filterRequest)))

      if (filterRequest.sqlQuery.nonEmpty) {
        filterRequest.sources.foreach(
          s => {
            if (s.sourceName.isEmpty) {
              responseObserver.onError(new Exception("sourceName should not be empty for source: " + s.sourceID))
            }
            if (s.sourceID.isEmpty) responseObserver.onError(new Exception("sourceID should not be empty"))
          }
        )
      }

      val jobName = getAvailableUUID(jedis)

      Try {
        filterRequest.sources.foreach(s =>
          RequestServer.schemaRegistryClient.getLatestSchemaMetadata(s.sourceID + "-value"))
      } match {
        case Failure(exception) => responseObserver.onError(exception)
        case Success(value)     =>
      }

      Try {

        val env: StreamExecutionEnvironment =
          StreamExecutionEnvironment.createRemoteEnvironment(Configs.FLINK_CLUSTER_HOST,
                                                             Configs.FLINK_CLUSTER_PORT.toInt)
        val tEnv = StreamTableEnvironment.create(env)

        val dataStream = Application.filter(filterRequest,
                                            env,
                                            RequestServer.brokerAddress,
                                            RequestServer.schemaRegistryAddress,
                                            tEnv,
                                            jobName)

        val iterator = dataStream.executeAndCollect(jobName)

        while (iterator.hasNext) {
          val record = iterator.next

          jedis.hset(Configs.FLINK_JOBS_REDIS_KEY, jobName, (System.currentTimeMillis() + 600000).toString)

          responseObserver.onNext(FilterResponse(record, jobName))
        }
      } match {
        case Failure(exception) => {
          RequestServer.logger.error("Error while streaming:")
          exception.printStackTrace()
          responseObserver.onError(exception)
        }
        case Success(value) => {
          responseObserver.onCompleted()
        }
      }

      Try {
        val jobID = RequestServer.clusterClient
          .listJobs()
          .get()
          .asScala
          .filter(_.getJobName == jobName)
          .head
          .getJobId

        RequestServer.clusterClient.cancel(jobID)
      } match {
        case Failure(exception) =>
          RequestServer.logger.error("error while canceling job :\n" + exception.printStackTrace())
        case Success(value) =>
      }

    }

    override def cancelFilterRequest(request: FilterCancellationRequest): Future[FilterCancellationResponse] =
      Try {
        val jobID =
          RequestServer.clusterClient.listJobs().get().asScala.filter(_.getJobName == request.jobName).head.getJobId

        RequestServer.clusterClient.cancel(jobID)
      } match {
        case Failure(exception) =>
          RequestServer.logger.error("error while canceling job :\n" + exception.printStackTrace())
          Future.fromTry(Success(FilterCancellationResponse()))
        case Success(value) => Future.fromTry(Success(FilterCancellationResponse(success = true)))
      }

    @tailrec
    private def getAvailableUUID(jedis: Jedis): String = {
      val uuid = UUID.randomUUID().toString
      if (jedis.hmget("flink-jobs", uuid).get(0) == null) {
        jedis.hset("flink-jobs", uuid, (System.currentTimeMillis() + 600000).toString)
        uuid
      } else getAvailableUUID(jedis)
    }
  }
}
