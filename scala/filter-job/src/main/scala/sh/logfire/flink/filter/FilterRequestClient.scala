package sh.logfire.flink.filter

import io.grpc.{ ManagedChannel, ManagedChannelBuilder, StatusRuntimeException }
import scalapb.json4s.JsonFormat
import sh.logfire.request.FlinkServiceGrpc.FlinkServiceBlockingStub
import sh.logfire.request.{ DateTimeFilter, FilterRequest, FlinkServiceGrpc, Source }

import java.util.concurrent.TimeUnit
import java.util.logging.{ Level, Logger }

object FilterRequestClient {
  def apply(host: String, port: Int): FilterRequestClient = {
    val channel      = ManagedChannelBuilder.forAddress(host, port).usePlaintext().build
    val blockingStub = FlinkServiceGrpc.blockingStub(channel)
    new FilterRequestClient(channel, blockingStub)
  }

  def main(args: Array[String]): Unit = {

//    val client = FilterRequestClient("localhost", 50051)
    val client = FilterRequestClient("g14", 31821)

    try {
      client.greet()
    } finally {
      client.shutdown()
    }
  }
}

class FilterRequestClient private (
    private val channel: ManagedChannel,
    private val blockingStub: FlinkServiceBlockingStub
) {
  private[this] val logger = Logger.getLogger(classOf[FilterRequestClient].getName)

  def shutdown(): Unit =
    channel.shutdown.awaitTermination(5, TimeUnit.SECONDS)

  def greet(): Unit = {

    val request2: FilterRequest = JsonFormat.fromJsonString[FilterRequest](s"""
         |{
         |  "sources": [
         |    {
         |      "sourceID": "test-go-docker-logs"
         |    }
         |  ]
         |}
         |""".stripMargin)




    val ts = com.google.protobuf.timestamp.Timestamp.fromJavaProto(com.google.protobuf.util.Timestamps.fromMillis(1678383025968L))

    val request: FilterRequest = FilterRequest(
      sources = Seq(Source(sourceID = "vector-docker-logs", sourceName = "docker"),
                    Source(sourceID = "vector-demo-logs", sourceName = "demo")),
      dateTimeFilter = Some(DateTimeFilter(startTimeStamp = Some(ts)))
    )

    try {

      blockingStub.submitFilterRequest(request).foreach(e => println(e.record))
    } catch {
      case e: StatusRuntimeException =>
        logger.log(Level.WARNING, "RPC failed: {0}", e.getStatus)
    }
  }
}
