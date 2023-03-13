package sh.logfire.flink.filter

import io.confluent.kafka.serializers.{ AbstractKafkaSchemaSerDeConfig, KafkaAvroDeserializer, KafkaAvroSerializer }
import io.github.embeddedkafka.schemaregistry.{ EmbeddedKafka, EmbeddedKafkaConfig }
import org.apache.avro.generic.GenericRecord
import org.apache.flink.runtime.testutils.MiniClusterResourceConfiguration
import org.apache.flink.streaming.api.scala.StreamExecutionEnvironment
import org.apache.flink.table.api.bridge.scala.StreamTableEnvironment
import org.apache.flink.test.util.MiniClusterWithClientResource
import org.apache.kafka.clients.producer.{ Callback, ProducerRecord, RecordMetadata }
import org.apache.kafka.common.serialization.{ Deserializer, Serializer }
import org.json.JSONObject
import org.scalatest.BeforeAndAfterAll
import org.scalatest.flatspec.AnyFlatSpec
import org.scalatest.matchers.should.Matchers
import sh.logfire.request.{ FilterRequest, SeverityLevel, Source }

import java.time.Duration
import java.util.Collections
import scala.jdk.CollectionConverters.mapAsJavaMapConverter

class FlinkAvroTestMain extends AnyFlatSpec with Matchers with BeforeAndAfterAll {
  private val flinkCluster = new MiniClusterWithClientResource(
    new MiniClusterResourceConfiguration.Builder()
      .setNumberSlotsPerTaskManager(1)
      .setNumberTaskManagers(1)
      .build)

  private val sourceTopic = "filter-test-raw"

  implicit val config: EmbeddedKafkaConfig = EmbeddedKafkaConfig(
    customBrokerProperties = Map(
      "transaction.max.timeout.ms" -> "3600000",
    ),
    kafkaPort = 9092,
    schemaRegistryPort = 8081
  )

  private val serializer = {
    val serializerProps = Map(AbstractKafkaSchemaSerDeConfig.SCHEMA_REGISTRY_URL_CONFIG -> "http://localhost:8081")
    val serializer      = new KafkaAvroSerializer()
    serializer.configure(serializerProps.asJava, false)
    serializer.asInstanceOf[Serializer[GenericRecord]]
  }

  implicit val genericSerializer: Serializer[GenericRecord] = serializer
  implicit val anySerializer: Serializer[AnyRef]            = serializer.asInstanceOf[Serializer[AnyRef]]

  private val deserializer = {
    val serializerProps = Map(AbstractKafkaSchemaSerDeConfig.SCHEMA_REGISTRY_URL_CONFIG -> "http://localhost:8081")
    val deserializer    = new KafkaAvroDeserializer()
    deserializer.configure(serializerProps.asJava, false)
    deserializer.asInstanceOf[Deserializer[GenericRecord]]
  }

  implicit val genericDeserializer: Deserializer[GenericRecord] = deserializer
  implicit val anyDeserializer: Deserializer[AnyRef]            = deserializer.asInstanceOf[Deserializer[AnyRef]]

  override def beforeAll(): Unit = {
    super.beforeAll()
    EmbeddedKafka.start()
    flinkCluster.before()

    producerRecords(sourceTopic, RawData.getGenericRecords())
  }

  override def afterAll(): Unit = {
    flinkCluster.after()
    EmbeddedKafka.stop()
    super.afterAll()
  }

  private def producerRecords(topic: String, records: Seq[GenericRecord]): Unit =
    EmbeddedKafka.withProducer[AnyRef, GenericRecord, Unit](producer => {
      records.foreach(record => {
        producer.send(
          new ProducerRecord(topic, 0, 1600000000000L, null, record),
          new Callback {
            override def onCompletion(r1: RecordMetadata, e: Exception): Unit =
              if (e != null) {
                e.printStackTrace()
              }
          }
        )
      })
      producer.flush()
    })

  private def consumeRecords(topic: String): Unit =
    EmbeddedKafka.withConsumer[AnyRef, GenericRecord, Unit](consumer => {
      consumer.subscribe(Collections.singletonList(topic))
      var numberOfRecordsConsumed = 0
      var time                    = 0L
      val constTimeInSec          = 2L
      while (time < 10) {
        val records = consumer.poll(Duration.ofSeconds(constTimeInSec))
        numberOfRecordsConsumed += records.count()
        records.forEach(consumerRecord => {
          println("consumed message:" + consumerRecord.value())
          println("schema: " + consumerRecord.value().getSchema)
        })
        time += constTimeInSec
      }
    })

  val env: StreamExecutionEnvironment = StreamExecutionEnvironment.getExecutionEnvironment
  private val tableEnv                = StreamTableEnvironment.create(env)

  "Filter request without any filter condition " should
  "not filter" in {
    val brokerAddress         = "localhost:9092"
    val schemaRegistryAddress = "http://localhost:8081"
    val rowsCount = Application
      .filter(FilterRequest(sources = Seq(Source(sourceID = sourceTopic))),
              env,
              brokerAddress,
              schemaRegistryAddress,
              tableEnv,
      "filter")
      .executeAndCollect(93)
      .count(r => r != null)
    assert(rowsCount == 93)
  }

  "Filter request with info security warning level filters " should
  " filter only warning records" in {
    val brokerAddress         = "localhost:9092"
    val schemaRegistryAddress = "http://localhost:8081"
    val rows = Application
      .filter(FilterRequest(sources = Seq(Source(sourceID = sourceTopic)), severityLevels = Seq(SeverityLevel.WARNING)),
              env,
              brokerAddress,
              schemaRegistryAddress,
              tableEnv,
        "filter")
      .executeAndCollect(15)

    val rowsCount = rows.count(r => r != null)
    assert(rowsCount == 15)

    rows.foreach(r => {
      val json = new JSONObject(r)
      assert(json.get(Configs.SEVERITY_LEVEL_FIELD) == "warning")
    })

  }

  "Filter request with info security info level filters " should
  " filter only info records" in {
    val brokerAddress         = "localhost:9092"
    val schemaRegistryAddress = "http://localhost:8081"
    val rows = Application
      .filter(FilterRequest(sources = Seq(Source(sourceID = sourceTopic)), severityLevels = Seq(SeverityLevel.INFO)),
              env,
              brokerAddress,
              schemaRegistryAddress,
              tableEnv,
        "filter")
      .executeAndCollect(11)

    val rowsCount = rows.count(r => r != null)
    assert(rowsCount == 11)

    rows.foreach(r => {
      val json = new JSONObject(r)
      assert(json.get(Configs.SEVERITY_LEVEL_FIELD) == "info")
    })
  }

  "Filter request with info security debug level filters " should
  " return only debug level records" in {
    val brokerAddress         = "localhost:9092"
    val schemaRegistryAddress = "http://localhost:8081"
    val rows = Application
      .filter(FilterRequest(sources = Seq(Source(sourceID = sourceTopic)), severityLevels = Seq(SeverityLevel.DEBUG)),
              env,
              brokerAddress,
              schemaRegistryAddress,
              tableEnv,
        "filter")
      .executeAndCollect(12)

    val rowsCount = rows.count(r => r != null)
    assert(rowsCount == 12)

    rows.foreach(r => {
      val json = new JSONObject(r)
      assert(json.get(Configs.SEVERITY_LEVEL_FIELD) == "debug")
    })
  }

  "Filter request with info security debug and info level filters " should
  " return debug and info records" in {
    val brokerAddress         = "localhost:9092"
    val schemaRegistryAddress = "http://localhost:8081"
    val rows = Application
      .filter(
        FilterRequest(sources = Seq(Source(sourceID = sourceTopic)),
                      severityLevels = Seq(SeverityLevel.DEBUG, SeverityLevel.INFO)),
        env,
        brokerAddress,
        schemaRegistryAddress,
        tableEnv,
        "filter"
      )
      .executeAndCollect(23)

    val rowsCount = rows.count(r => r != null)
    assert(rowsCount == 23)

    rows.foreach(r => {
      val json = new JSONObject(r)
      assert(Array("debug", "info").contains(json.get(Configs.SEVERITY_LEVEL_FIELD)))
    })
  }

  "Filter request without any search query filter condition " should
  "not filter" in {
    val brokerAddress         = "localhost:9092"
    val schemaRegistryAddress = "http://localhost:8081"
    val rowsCount = Application
      .filter(FilterRequest(sources = Seq(Source(sourceID = sourceTopic))),
              env,
              brokerAddress,
              schemaRegistryAddress,
              tableEnv,
        "filter")
      .executeAndCollect(93)
      .count(r => r != null)
    assert(rowsCount == 93)
  }

  "Filter request with 'never' search query filter condition " should
  "should filter" in {
    val brokerAddress         = "localhost:9092"
    val schemaRegistryAddress = "http://localhost:8081"
    val rows = Application
      .filter(FilterRequest(sources = Seq(Source(sourceID = sourceTopic)), searchQueries = Seq("never")),
              env,
              brokerAddress,
              schemaRegistryAddress,
              tableEnv,
        "filter")
      .executeAndCollect(12)

    val rowsCount = rows.count(r => r != null)
    assert(rowsCount == 12)

    rows.foreach(r => {
      val json = new JSONObject(r)
      assert(json.get(Configs.MESSAGE_FIELD).asInstanceOf[String].contains("never"))
    })
  }

  "Filter request with 'gonna' search query filter condition " should
  "should filter" in {
    val brokerAddress         = "localhost:9092"
    val schemaRegistryAddress = "http://localhost:8081"
    val rows = Application
      .filter(FilterRequest(sources = Seq(Source(sourceID = sourceTopic)), searchQueries = Seq("gonna")),
              env,
              brokerAddress,
              schemaRegistryAddress,
              tableEnv,
        "filter")
      .executeAndCollect(29)

    val rowsCount = rows.count(r => r != null)
    assert(rowsCount == 29)

    rows.foreach(r => {
      val json = new JSONObject(r)
      assert(json.get(Configs.MESSAGE_FIELD).asInstanceOf[String].contains("gonna"))
    })
  }

  "Filter request with 'gonna' and 'never' search query filter condition " should
  "should filter" in {
    val brokerAddress         = "localhost:9092"
    val schemaRegistryAddress = "http://localhost:8081"
    val rows = Application
      .filter(FilterRequest(sources = Seq(Source(sourceID = sourceTopic)), searchQueries = Seq("gonna", "never")),
              env,
              brokerAddress,
              schemaRegistryAddress,
              tableEnv,
        "filter")
      .executeAndCollect(29)

    val rowsCount = rows.count(r => r != null)
    assert(rowsCount == 29)

    rows.foreach(r => {
      val json = new JSONObject(r)
      assert(
        json.get(Configs.MESSAGE_FIELD).asInstanceOf[String].contains("gonna") || json
          .get(Configs.MESSAGE_FIELD)
          .asInstanceOf[String]
          .contains("never"))
    })
  }

}
