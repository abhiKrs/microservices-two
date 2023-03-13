package sh.logfire.flink.filter

import io.confluent.kafka.schemaregistry.client.CachedSchemaRegistryClient
import io.confluent.kafka.serializers.{AbstractKafkaSchemaSerDeConfig, KafkaAvroDeserializer, KafkaAvroDeserializerConfig}
import org.apache.avro.Schema
import org.apache.avro.generic.{GenericData, GenericRecord}
import org.apache.flink.api.common.typeinfo.TypeInformation
import org.apache.flink.api.java.typeutils.RowTypeInfo
import org.apache.flink.streaming.connectors.kafka.KafkaDeserializationSchema
import org.apache.flink.types.Row
import org.apache.kafka.clients.consumer.ConsumerRecord

import java.sql.Timestamp
import java.time.Instant
import java.util

@SerialVersionUID(1000L)
class KafkaGenericAvroDeserializationSchema(topic: String, schemaRegistryAddress: String) extends KafkaDeserializationSchema[Row] {

  @transient private var valueDeSerializer: KafkaAvroDeserializer = null

  override val getProducedType: TypeInformation[Row] =
    AvroSerDe.getRowSchema(topic, schemaRegistryAddress)
  private def rowTypeInfo: RowTypeInfo = getProducedType.asInstanceOf[RowTypeInfo]

  private val valueSchema: Schema = new AvroUtils(schemaRegistryAddress).getValueSchema(topic)

  private def checkInitialized(): Unit =
    if (valueDeSerializer == null) {
      val props = new util.HashMap[String, String]
      props.put(
        AbstractKafkaSchemaSerDeConfig.SCHEMA_REGISTRY_URL_CONFIG,
        schemaRegistryAddress
      )
      props.put(KafkaAvroDeserializerConfig.SPECIFIC_AVRO_READER_CONFIG, "false")
      val client = new CachedSchemaRegistryClient(
        schemaRegistryAddress,
        AbstractKafkaSchemaSerDeConfig.MAX_SCHEMAS_PER_SUBJECT_DEFAULT
      )
      valueDeSerializer = new KafkaAvroDeserializer(client, props)
      valueDeSerializer.configure(props, false)
    }

  override def isEndOfStream(nextElement: Row): Boolean = false

  override def deserialize(record: ConsumerRecord[Array[Byte], Array[Byte]]): Row = {
    checkInitialized()

    if (record.value() == null) {
      return null
    }

    val kafkaTs = Instant.ofEpochMilli(record.timestamp).toString

    val valueRecord =
      valueDeSerializer.deserialize(record.topic, record.value).asInstanceOf[GenericRecord]

    val newGenericRecord = new GenericData.Record(valueSchema)
    newGenericRecord.put(Configs.LOGFIRE_TIMESTAMP_FIELD, kafkaTs)
    valueSchema.getFields.forEach(f => if (f.name() != Configs.LOGFIRE_TIMESTAMP_FIELD) newGenericRecord.put(f.name(), valueRecord.get(f.name())))


    val row = AvroSerDe.convertAvroRecordToRow(valueSchema, rowTypeInfo, newGenericRecord)
    row.setField(Configs.LOGFIRE_TIMESTAMP_FIELD, kafkaTs)
    row
  }
}
