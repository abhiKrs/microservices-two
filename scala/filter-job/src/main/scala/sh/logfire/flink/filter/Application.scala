package sh.logfire.flink.filter

import org.apache.avro.generic.GenericData
import org.apache.flink.api.common.eventtime.{SerializableTimestampAssigner, WatermarkStrategy}
import org.apache.flink.api.common.typeinfo.TypeInformation
import org.apache.flink.connector.kafka.source.enumerator.initializer.OffsetsInitializer
import org.apache.flink.connector.kafka.source.reader.deserializer.KafkaRecordDeserializationSchema
import org.apache.flink.connector.kafka.source.{KafkaSource, KafkaSourceBuilder}
import org.apache.flink.formats.avro.typeutils.GenericRecordAvroTypeInfo
import org.apache.flink.streaming.api.scala.{DataStream, StreamExecutionEnvironment}
import org.apache.flink.table.api.bridge.scala.StreamTableEnvironment
import org.apache.flink.table.api.{DataTypes, Schema}
import org.apache.flink.table.types.DataType
import org.apache.flink.table.types.logical.{LogicalTypeRoot, RowType}
import org.apache.flink.types.Row
import org.slf4j.{Logger, LoggerFactory}
import sh.logfire.request.{DateTimeFilter, FilterRequest, Source}

import java.time.{Duration, Instant}
import java.util.Properties
import scala.util.{Failure, Success, Try}

object Application {
  val logger: Logger = LoggerFactory.getLogger(getClass.getSimpleName)
  def main(args: Array[String]): Unit = {
    val brokerAddress         = Configs.KAFKA_BROKERS
    val schemaRegistryAddress = Configs.SCHEMA_REGISTRY_ADDRESS

    val env: StreamExecutionEnvironment = StreamExecutionEnvironment.getExecutionEnvironment
    env.setParallelism(1)
    val tEnv = StreamTableEnvironment.create(env)

    "test-go-logs"
    "vector-docker-logs"
    "vector-demo-logs"
    val ts = com.google.protobuf.timestamp.Timestamp.fromJavaProto(com.google.protobuf.util.Timestamps.fromMillis(1678383025968L))

    filter(
      FilterRequest(
        sources = Seq(
                      Source(sourceID = "vector-demo-logs", sourceName = "demo"))
        , dateTimeFilter = Some(DateTimeFilter(startTimeStamp = Some(ts))
      )),
      env,
      brokerAddress,
      schemaRegistryAddress,
      tEnv,
      "test"
    ).print()

    env.execute()
  }

  def filter(filterRequest: FilterRequest,
             env: StreamExecutionEnvironment,
             brokerAddress: String,
             schemaRegistryAddress: String,
             tEnv: StreamTableEnvironment,
             jobName: String): DataStream[String] = {

    val dataStreams = filterRequest.sources.map(source => {

      val topic = source.sourceID

      val consumerProperties = {
        val properties = new Properties()
        properties.setProperty("isolation.level", "read_committed")
        properties.setProperty("commit.offsets.on.checkpoint", "true")
        properties
      }

      val startingOffsets = filterRequest.dateTimeFilter match {
        case Some(dateTimeFilter) =>
          dateTimeFilter.startTimeStamp match {
            case Some(startTimeStamp) =>
              OffsetsInitializer.timestamp(
                com.google.protobuf.util.Timestamps
                  .toMillis(com.google.protobuf.timestamp.Timestamp.toJavaProto(startTimeStamp)))
            case None => OffsetsInitializer.earliest()
          }
        case None => OffsetsInitializer.earliest()
      }

      val kafkaSourceBuilder: KafkaSourceBuilder[Row] = KafkaSource
        .builder()
        .setBootstrapServers(brokerAddress)
        .setTopics(topic)
        .setGroupId(jobName + "-filter")
        .setStartingOffsets(startingOffsets)
        .setClientIdPrefix(jobName + "-" + source.sourceID + "-filter-")
        .setProperties(consumerProperties)
        .setDeserializer(KafkaRecordDeserializationSchema.of(
          new KafkaGenericAvroDeserializationSchema(topic, schemaRegistryAddress)))

      val kafkaSource: KafkaSource[Row] = filterRequest.dateTimeFilter match {
        case Some(dateTimeFilter) =>
          dateTimeFilter.endTimeStamp match {
            case Some(startTimeStamp) =>
              val stoppingOffsetsInitializer = OffsetsInitializer.timestamp(
                com.google.protobuf.util.Timestamps
                  .toMillis(com.google.protobuf.timestamp.Timestamp.toJavaProto(startTimeStamp)))

              kafkaSourceBuilder.setBounded(stoppingOffsetsInitializer).build()

            case None => kafkaSourceBuilder.build()
          }
        case None => kafkaSourceBuilder.build()
      }

      val sourceInfo = AvroSerDe.getRowSchema(topic, schemaRegistryAddress)

//      TODO advance watermark when no new records
      val watermarkStrategy = WatermarkStrategy
        .forBoundedOutOfOrderness(Duration.ofSeconds(5))
        .withTimestampAssigner(new SerializableTimestampAssigner[Row] {
          override def extractTimestamp(element: Row, recordTimestamp: Long): Long =
            Instant.parse(element.getFieldAs[String](Configs.LOGFIRE_TIMESTAMP_FIELD)).toEpochMilli
        })
        .withIdleness(Duration.ofSeconds(5))

      val dataStream = env.fromSource(kafkaSource, watermarkStrategy, source.sourceName + "_kafka")(sourceInfo)

      val dataType = AvroSerDe.getRowDataType(topic, schemaRegistryAddress)

      val rowType    = dataType.getLogicalType.asInstanceOf[RowType]
      val fieldNames = rowType.getFieldNames

      val logLevelsToFilter = filterRequest.severityLevels.map(_.toString().toLowerCase)

      val withLevelFilter =
        if (logLevelsToFilter.nonEmpty && fieldNames
              .contains(Configs.SEVERITY_LEVEL_FIELD) && Array(LogicalTypeRoot.CHAR, LogicalTypeRoot.VARCHAR)
              .contains(rowType.getTypeAt(rowType.getFieldIndex(Configs.SEVERITY_LEVEL_FIELD)).getTypeRoot))
          dataStream.filter(r =>
            logLevelsToFilter.contains(r.getFieldAs[String](Configs.SEVERITY_LEVEL_FIELD).toLowerCase))
        else dataStream

      val withSearchFilter =
        if (filterRequest.searchQueries.nonEmpty && fieldNames.contains(Configs.MESSAGE_FIELD) && Array(LogicalTypeRoot.CHAR,
                                                                                            LogicalTypeRoot.VARCHAR)
              .contains(rowType.getTypeAt(rowType.getFieldIndex(Configs.MESSAGE_FIELD)).getTypeRoot)) {
          withLevelFilter.filter(r => {
            val fieldValue = r.getFieldAs[String](Configs.MESSAGE_FIELD)
            if (fieldValue != null) {
              val firstFilter = fieldValue
                .contains(filterRequest.searchQueries.head.trim)

              filterRequest.searchQueries.tail.foldLeft(firstFilter)((a, b) => a || fieldValue.contains(b.trim))
            } else {
              false
            }
          })

        } else {
          withLevelFilter
        }

      val (newDataStream, newDataType): (DataStream[Row], DataType) = Try {
        if (source.sourceName.nonEmpty) {
          val table = tEnv.fromChangelogStream(
            withSearchFilter,
            Schema
              .newBuilder()
              .fromRowDataType(dataType)
              .columnByMetadata(Configs.LOGFIRE_ROWTIME_FIELD, DataTypes.TIMESTAMP_LTZ(3))
              .watermark(Configs.LOGFIRE_ROWTIME_FIELD, "SOURCE_WATERMARK()")
              .build()
          )
          tEnv.createTemporaryView(source.sourceName, table)
          (tEnv.toChangelogStream(table), table.getResolvedSchema.toSinkRowDataType)
        } else {
          (withSearchFilter, tEnv.fromDataStream(withSearchFilter).getResolvedSchema.toSinkRowDataType)
        }
      } match {
        case Failure(exception) =>
          exception.printStackTrace()
          logger.error("Error while creating temporary view")
          (withSearchFilter, tEnv.fromDataStream(withSearchFilter).getResolvedSchema.toSinkRowDataType)
        case Success(t) => t
      }

      val avroSchema = AvroSerDe
        .convertToSchema(newDataType.getLogicalType.asInstanceOf[RowType], "record")

      val genericInfo = new GenericRecordAvroTypeInfo(avroSchema).asInstanceOf[TypeInformation[GenericData.Record]]

      newDataStream
        .map(row => AvroSerDe.convertRowToAvroRecord(avroSchema.toString, row))(genericInfo)
        .map(_.toString)(TypeInformation.of(classOf[String]))

    })

    if (filterRequest.sqlQuery.nonEmpty) {
      Try(tEnv.sqlQuery(filterRequest.sqlQuery)) match {
        case Failure(exception) =>
          logger.error("Unable to run the sql query")
          exception.printStackTrace()
          unionAndSortStreams(dataStreams)
        case Success(table) =>
//          TODO add the datetime col
          val avroSchema = AvroSerDe
            .convertToSchema(table.getResolvedSchema.toSinkRowDataType.getLogicalType.asInstanceOf[RowType], "record")
          val genericInfo = new GenericRecordAvroTypeInfo(avroSchema).asInstanceOf[TypeInformation[GenericData.Record]]
          tEnv
            .toChangelogStream(table)
            .map(row => AvroSerDe.convertRowToAvroRecord(avroSchema.toString, row))(genericInfo)
            .map(_.toString)(TypeInformation.of(classOf[String]))
      }
    } else {
      unionAndSortStreams(dataStreams)
    }

  }

  private def unionAndSortStreams(dataStreams: Seq[DataStream[String]]): DataStream[String] =
    if (dataStreams.size > 1) {
      dataStreams
        .foldLeft(dataStreams.head)((c, n) => c.union(n))
        .keyBy(_ => 1: java.lang.Integer)(TypeInformation.of(classOf[java.lang.Integer]))
        .process(new SortingFunction())(TypeInformation.of(classOf[String]))
    } else {
      dataStreams.head
    }

}
