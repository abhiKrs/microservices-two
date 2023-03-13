package sh.logfire.flink.filter

import io.confluent.kafka.schemaregistry.client.CachedSchemaRegistryClient
import org.apache.avro.{LogicalTypes, Schema, SchemaBuilder}
import org.apache.avro.Schema.Field

import java.util

class AvroUtils(schemaRegistryAddress: String) {
  private val schemaRegistry = new CachedSchemaRegistryClient(schemaRegistryAddress, 1000)

  def getValueSchema(topic: String): Schema = {
    val valueSchemaString = schemaRegistry
      .getLatestSchemaMetadata(topic + "-value")
      .getSchema

    val newTsField = new Field(Configs.LOGFIRE_TIMESTAMP_FIELD, SchemaBuilder.builder.stringType(), null, null)

    val schema = new Schema.Parser().parse(valueSchemaString)

    val newFields: util.ArrayList[Field] = new util.ArrayList[Field]()

    schema.getFields.forEach(f => {
      if (f.name() != Configs.LOGFIRE_TIMESTAMP_FIELD) newFields.add(new Schema.Field(f.name(), f.schema(), f.doc(), f.defaultVal()))
    })

    newFields.add(newTsField)


    Schema.createRecord("record", null, null, false, newFields)
  }
}
