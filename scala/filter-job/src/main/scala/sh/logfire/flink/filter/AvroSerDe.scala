package sh.logfire.flink.filter

import org.apache.avro.Schema.Type._
import org.apache.avro.generic.{ GenericData, GenericFixed, GenericRecord }
import org.apache.avro.util.Utf8
import org.apache.avro.{ LogicalTypes, Schema, SchemaBuilder }
import org.apache.flink.api.common.typeinfo._
import org.apache.flink.api.java.typeutils.{ MapTypeInfo, ObjectArrayTypeInfo, RowTypeInfo }
import org.apache.flink.table.api.DataTypes
import org.apache.flink.table.types.logical._
import org.apache.flink.table.types.{ logical, DataType }
import org.apache.flink.types.{ Row, RowKind }

import java.math.{ BigDecimal, BigInteger }
import java.nio.ByteBuffer
import java.sql.{ Date, Time, Timestamp }
import java.time.temporal.ChronoField
import java.time.{ Instant, LocalDate, LocalDateTime, LocalTime }
import java.util
import java.util.TimeZone
import java.util.concurrent.TimeUnit
import javax.annotation.Nullable
import scala.jdk.CollectionConverters.collectionAsScalaIterableConverter
import scala.util.{ Failure, Success, Try }

object AvroSerDe {

  private val LOCAL_TZ          = TimeZone.getDefault
  private val MICROS_PER_SECOND = 1000000L

  @Nullable @transient private val jodaConverter = JodaConverter.getConverter

  def getRowSchema(topic: String, schemaRegistryAddress: String): TypeInformation[Row] = {
    val avroUtils = new AvroUtils(schemaRegistryAddress)

    val valueSchemaString = avroUtils.getValueSchema(topic)

    convertToTypeInfo(valueSchemaString).asInstanceOf[TypeInformation[Row]]

  }

  def getRowDataType(topic: String, schemaRegistryAddress: String): DataType = {
    val avroUtils = new AvroUtils(schemaRegistryAddress)

    val valueSchemaString = avroUtils.getValueSchema(topic)

    convertToDataType(valueSchemaString)
  }

  def convertToTypeInfo(fieldSchema: Schema): TypeInformation[_] =
    fieldSchema.getType match {
      case STRING  => Types.STRING
      case NULL    => Types.VOID
      case BOOLEAN => Types.BOOLEAN
      case BYTES => {
        fieldSchema.getLogicalType match {
          case _: LogicalTypes.Decimal => Types.BIG_DEC
          case _                       => Types.PRIMITIVE_ARRAY(Types.BYTE)
        }
      }
      case ENUM => Types.STRING
      case UNION => {
        if (fieldSchema.getTypes.size == 2 && (fieldSchema.getTypes
              .get(0)
              .getType == Schema.Type.NULL)) convertToTypeInfo(fieldSchema.getTypes.get(1))
        else if (fieldSchema.getTypes.size == 2 && (fieldSchema.getTypes
                   .get(1)
                   .getType == Schema.Type.NULL)) convertToTypeInfo(fieldSchema.getTypes.get(0))
        else if (fieldSchema.getTypes.size == 1) convertToTypeInfo(fieldSchema.getTypes.get(0))
        else {
          val fieldsTypeInfos = (for {
            y <- fieldSchema.getTypes.asScala
            if !y.getType.equals(NULL)
          } yield convertToTypeInfo(y)).toList

          val memberNames = fieldsTypeInfos.indices.map(i => "member" + i).toArray

          Types.ROW_NAMED(memberNames, fieldsTypeInfos: _*)
        }
      }
      case FIXED => {
        fieldSchema.getLogicalType match {
          case t: LogicalTypes.Decimal => Types.BIG_DEC
          case _                       => Types.PRIMITIVE_ARRAY(Types.BYTE)
        }
      }
      case INT => {
        val logicalType = fieldSchema.getLogicalType
        if (logicalType == LogicalTypes.date) Types.SQL_DATE
        else if (logicalType == LogicalTypes.timeMillis) Types.SQL_TIME
        else Types.INT
      }
      case LONG => {
        if ((fieldSchema.getLogicalType == LogicalTypes.timestampMillis)
            || (fieldSchema.getLogicalType == LogicalTypes.timestampMicros)
            || (fieldSchema.getLogicalType == LogicalTypes.localTimestampMillis)
            || (fieldSchema.getLogicalType == LogicalTypes.localTimestampMicros))
          Types.SQL_TIMESTAMP
        else if (fieldSchema.getLogicalType == LogicalTypes.timeMicros) Types.SQL_TIME
        else Types.LONG
      }
      case FLOAT  => Types.FLOAT
      case DOUBLE => Types.DOUBLE
      case ARRAY  => Types.OBJECT_ARRAY(convertToTypeInfo(fieldSchema.getElementType))
      case MAP    => Types.MAP(Types.STRING, convertToTypeInfo(fieldSchema.getValueType))
      //      TODO all logical types

      case RECORD =>
        val fields = fieldSchema.getFields

        val types = new Array[TypeInformation[_]](fields.size)
        val names = new Array[String](fields.size)
        for (i <- 0 until fields.size) {
          val field = fields.get(i)
          types(i) = convertToTypeInfo(field.schema)
          names(i) = field.name
        }
        Types.ROW_NAMED(names, types: _*)
    }

  def convertAvroType(schema: Schema, info: TypeInformation[_], `object`: Any): Any = {
    if (`object` == null) return null
    schema.getType match {
      case RECORD =>
        `object` match {
          case genericRecord: GenericRecord =>
            return convertAvroRecordToRow(schema, info.asInstanceOf[RowTypeInfo], genericRecord)
          case _ => null
        }
        throw new IllegalStateException("IndexedRecord expected but was: " + `object`.getClass)
      case ENUM => {
        if (`object`.isInstanceOf[org.apache.avro.generic.GenericData.EnumSymbol]) `object`.toString
        else null
      }
      case NULL => null
      case STRING => {
        `object` match {
          case _: Utf8   => `object`.toString
          case _: String => `object`
          case _         => null
        }
      }
      case ARRAY =>
        info match {
          case value: BasicArrayTypeInfo[_, _] =>
            val elementInfo = value.getComponentInfo
            convertToObjectArray(schema.getElementType, elementInfo, `object`)
          case x: ObjectArrayTypeInfo[_, _] =>
            val elementInfo = info.asInstanceOf[ObjectArrayTypeInfo[_, _]].getComponentInfo
            convertToObjectArray(schema.getElementType, elementInfo, `object`)
        }
      case MAP =>
        val mapTypeInfo  = info.asInstanceOf[MapTypeInfo[_, _]]
        val convertedMap = new util.HashMap[String, Any]
        val map          = `object`.asInstanceOf[util.HashMap[_, _]]
        import scala.collection.JavaConversions._
        for (entry <- map.entrySet) {
          convertedMap.put(
            entry.getKey.toString,
            convertAvroType(schema.getValueType, mapTypeInfo.getValueTypeInfo, entry.getValue)
          )
        }
        convertedMap
      case UNION =>
        val types = schema.getTypes
        val size  = types.size
        if (size == 2 && (types.get(0).getType == Schema.Type.NULL))
          convertAvroType(types.get(1), info, `object`)
        else if (size == 2 && (types.get(1).getType == Schema.Type.NULL))
          convertAvroType(types.get(0), info, `object`)
        else if (size == 1) convertAvroType(types.get(0), info, `object`)
        else {
          val row       = Row.withNames(RowKind.INSERT)
          val fieldInfo = info.asInstanceOf[RowTypeInfo].getFieldTypes

          schema.getTypes.asScala
            .filter(!_.getType.equals(NULL))
            .zipWithIndex
            .foreach(x => row.setField("member" + x._2, convertAvroType(x._1, fieldInfo(x._2), `object`)))

          row
        }
      case FIXED =>
        val fixedBytes = `object`.asInstanceOf[GenericFixed].bytes
        if (info == Types.BIG_DEC) return convertToDecimal(schema, fixedBytes)
        fixedBytes
      case BYTES =>
        val byteBuffer = `object`.asInstanceOf[ByteBuffer]
        val bytes      = new Array[Byte](byteBuffer.remaining)
        byteBuffer.get(bytes)
        if (info == Types.BIG_DEC) convertToDecimal(schema, bytes)
        else bytes
      case INT =>
        if (info == Types.SQL_DATE) convertToDate(`object`)
        else if (info == Types.SQL_TIME) convertToTime(`object`)
        else if (`object`.isInstanceOf[java.lang.Integer]) `object`
        else null
      case LONG =>
        if (info == Types.SQL_TIMESTAMP)
          convertToTimestamp(
            `object`,
            schema.getLogicalType == LogicalTypes.timestampMicros || schema.getLogicalType == LogicalTypes.localTimestampMicros,
            schema.getLogicalType == LogicalTypes.localTimestampMillis || schema.getLogicalType == LogicalTypes.localTimestampMicros
          )
        else if (info == Types.SQL_TIME) convertToTime(`object`)
        else if (`object`.isInstanceOf[java.lang.Long]) `object`
        else null
      case FLOAT => {
        if (`object`.isInstanceOf[java.lang.Float]) `object`
        else null
      }
      case DOUBLE => {
        if (`object`.isInstanceOf[java.lang.Double]) `object`
        else null
      }
      case BOOLEAN => {
        if (`object`.isInstanceOf[java.lang.Boolean]) `object`
        else null
      }
    }
    //    throw new RuntimeException("Unsupported Avro type:" + schema)
  }

  def convertAvroRecordToRow(schema: Schema, typeInfo: RowTypeInfo, record: GenericRecord): Row = {
    val fields    = schema.getFields
    val fieldInfo = typeInfo.getFieldTypes
    val length    = fields.size
    val row       = Row.withNames(RowKind.INSERT)
    for (i <- 0 until length) {
      val field = fields.get(i)
      val fieldValue = Try(record.get(field.name())) match {
        case Success(value) => value
        case Failure(_)     => null
      }
      val c = convertAvroType(field.schema, fieldInfo(i), fieldValue)
      row.setField(field.name(), c)
    }
    row
  }

  private def convertToObjectArray(
      elementSchema: Schema,
      elementInfo: TypeInformation[_],
      `object`: Any
  ) = {
    val list = `object`.asInstanceOf[util.List[_]]
    val convertedArray: Array[Any] = java.lang.reflect.Array
      .newInstance(elementInfo.getTypeClass, list.size)
      .asInstanceOf[Array[Any]]
    for (i <- 0 until list.size) {
      convertedArray(i) = convertAvroType(elementSchema, elementInfo, list.get(i))
    }
    convertedArray
  }

  private def convertToDecimal(schema: Schema, bytes: Array[Byte]) = {
    val decimalType = schema.getLogicalType.asInstanceOf[LogicalTypes.Decimal]
    val bigDecimal  = new BigDecimal(new BigInteger(bytes), decimalType.getScale)
    bigDecimal
  }

  private def convertToDate(`object`: Any) = {
    var millis = 0L
    `object` match {
      case value: Integer =>
        // adopted from Apache Calcite
        val t = value.toLong * 86400000L
        millis = t - LOCAL_TZ.getOffset(t).toLong
      case date: LocalDate =>
        val t = date.toEpochDay * 86400000L
        millis = t - LOCAL_TZ.getOffset(t).toLong
      case _ =>
        if (jodaConverter != null) millis = jodaConverter.convertDate(`object`)
        //        TODO
        else
          throw new IllegalArgumentException(
            "Unexpected object type for DATE logical type. Received: " + `object`
          )
    }
    new Date(millis)
  }

  private def convertToTime(`object`: Any): java.sql.Time = {
    var millis = 0L
    `object` match {
      case integer: Integer => millis = integer.toLong
      case l: Long          => millis = l / 1000L
      case time: LocalTime  => millis = time.get(ChronoField.MILLI_OF_DAY)
      case _ =>
        if (jodaConverter != null) millis = jodaConverter.convertTime(`object`)
        else
          throw new IllegalArgumentException(
            "Unexpected object type for DATE logical type. Received: " + `object`
          )
    }
    new Time(millis - LOCAL_TZ.getOffset(millis))
  }

  private def convertToTimestamp(`object`: Any, isMicros: Boolean, isLocal: Boolean): Timestamp = {
    var millis = 0L
    `object` match {
      case micros: Long =>
        if (isMicros) {
          if (!isLocal) {
            val offsetMillis = LOCAL_TZ.getOffset(micros / 1000L)
            val seconds      = micros / MICROS_PER_SECOND - offsetMillis / 1000
            val nanos        = (micros % MICROS_PER_SECOND).toInt * 1000 - offsetMillis % 1000 * 1000
            val timestamp    = new Timestamp(seconds * 1000L)
            timestamp.setNanos(nanos)
            return timestamp
          } else {
//          TODO not sure
            val seconds   = micros / MICROS_PER_SECOND
            val nanos     = (micros % MICROS_PER_SECOND).toInt * 1000
            val timestamp = new Timestamp(seconds * 1000L)
            timestamp.setNanos(nanos)
            return timestamp
          }
        } else millis = micros
      case instant: Instant =>
        if (!isLocal) {
          val offsetMillis = LOCAL_TZ.getOffset(instant.toEpochMilli)
          val seconds      = instant.getEpochSecond - offsetMillis / 1000
          val nanos        = instant.getNano - offsetMillis % 1000 * 1000
          val timestamp    = new Timestamp(seconds * 1000L)
          timestamp.setNanos(nanos)
          return timestamp
        } else {
//          TODO not sure
          val seconds   = instant.getEpochSecond
          val nanos     = instant.getNano
          val timestamp = new Timestamp(seconds * 1000L)
          timestamp.setNanos(nanos)
          return timestamp
        }

      case _ =>
        `object` match {
          case _ =>
            if (jodaConverter != null) millis = jodaConverter.convertTimestamp(`object`)
            else
              throw new IllegalArgumentException(
                "Unexpected object type for DATE logical type. Received: " + `object`
              )
        }
    }
    if (isLocal) new Timestamp(millis) else new Timestamp(millis - LOCAL_TZ.getOffset(millis))
  }

  def convertToDataType(schema: Schema, insideUnion: Boolean = false): DataType = {
    val fieldDataType: DataType = schema.getType match {
      case RECORD =>
        val schemaFields = schema.getFields
        val fields       = new Array[DataTypes.Field](schemaFields.size)
        for (i <- 0 until schemaFields.size) {
          val field = schemaFields.get(i)
          fields(i) = DataTypes.FIELD(
            field.name,
            convertToDataType(
              field.schema,
              if (field.schema.getType.equals(UNION)) true else insideUnion
            )
          )
        }
        DataTypes.ROW(fields: _*).notNull
      case ENUM =>
        DataTypes.STRING.notNull
      case ARRAY =>
        DataTypes.ARRAY(convertToDataType(schema.getElementType)).notNull
      case MAP =>
        DataTypes.MAP(DataTypes.STRING.notNull, convertToDataType(schema.getValueType)).notNull
      case UNION =>
        if (schema.getTypes.size == 2 && (schema.getTypes.get(0).getType == Schema.Type.NULL)) {
          convertToDataType(schema.getTypes.get(1), true)
        } else if (schema.getTypes.size == 2 && (schema.getTypes.get(1).getType == Schema.Type.NULL)) {
          convertToDataType(schema.getTypes.get(0), true)
        } else if (schema.getTypes.size == 1) {
          convertToDataType(schema.getTypes.get(0), insideUnion)
        } else {
          val fieldsTypeInfos = (for {
            y <- schema.getTypes.asScala
            if !y.getType.equals(NULL)
          } yield convertToDataType(y, true)).toList

          val fields = fieldsTypeInfos.zipWithIndex.map(f => DataTypes.FIELD("member" + f._2, f._1))

          DataTypes.ROW(fields: _*).notNull
        }
      case FIXED =>
        // logical decimal type
        schema.getLogicalType match {
          case decimalType: LogicalTypes.Decimal =>
            DataTypes
              .DECIMAL(decimalType.getPrecision, decimalType.getScale)
              .notNull
              .bridgedTo(classOf[java.math.BigDecimal])
          case _ => DataTypes.VARBINARY(schema.getFixedSize).notNull
        }
      case STRING =>
        // convert Avro's Utf8/CharSequence to String
        DataTypes.STRING
      case BYTES =>
        schema.getLogicalType match {
          case decimalType: LogicalTypes.Decimal =>
            DataTypes
              .DECIMAL(decimalType.getPrecision, decimalType.getScale)
              .bridgedTo(classOf[java.math.BigDecimal])
              .notNull()
          case _ => DataTypes.BYTES()
        }
      case INT =>
        // logical date and time type
        val logicalType = schema.getLogicalType
        if (logicalType == LogicalTypes.date)
          DataTypes.DATE.notNull.bridgedTo(classOf[java.sql.Date])
        else if (logicalType == LogicalTypes.timeMillis)
          DataTypes.TIME(3).notNull.bridgedTo(classOf[java.sql.Time])
        else DataTypes.INT.notNull
      case LONG =>
        // logical timestamp type
        if (schema.getLogicalType == LogicalTypes.timestampMillis)
          DataTypes.TIMESTAMP(3).bridgedTo(classOf[java.sql.Timestamp])
        else if (schema.getLogicalType == LogicalTypes.timestampMicros)
          DataTypes.TIMESTAMP(6).bridgedTo(classOf[java.sql.Timestamp])
        else if (schema.getLogicalType == LogicalTypes.timeMillis)
          DataTypes.TIME(3).bridgedTo(classOf[java.sql.Time])
        else if (schema.getLogicalType == LogicalTypes.timeMicros)
          DataTypes.TIME(6).bridgedTo(classOf[java.sql.Time])
        else if (schema.getLogicalType == LogicalTypes.localTimestampMicros)
          DataTypes.TIMESTAMP_LTZ(6).bridgedTo(classOf[java.sql.Timestamp])
        else if (schema.getLogicalType == LogicalTypes.localTimestampMillis)
          DataTypes.TIMESTAMP_LTZ(3).bridgedTo(classOf[java.sql.Timestamp])
        else DataTypes.BIGINT
      case FLOAT =>
        DataTypes.FLOAT.notNull
      case DOUBLE =>
        DataTypes.DOUBLE.notNull
      case BOOLEAN =>
        DataTypes.BOOLEAN.notNull
      case NULL =>
        DataTypes.NULL
    }
    if (insideUnion) fieldDataType.nullable else fieldDataType.notNull
  }

  def convertRowToAvroRecord(schemaString: String, row: Row) = {
    val schema = Schema.parse(schemaString)
    val fields = schema.getFields
    val length = fields.size
    val record = new GenericData.Record(schema)
    for (i <- 0 until length) {
      val field = fields.get(i)
      record.put(i, convertFlinkType(field.schema, row.getField(field.name())))
    }
    record
  }

  def convertRowToAvroRecord(schema: Schema, row: Row) = {
    val fields = schema.getFields
    val length = fields.size
    val record = new GenericData.Record(schema)
    for (i <- 0 until length) {
      val field = fields.get(i)
      record.put(i, convertFlinkType(field.schema, row.getField(field.name())))
    }
    record
  }

  private def convertFlinkType(schema: Schema, `object`: Any): Any = {
//    TODO jvm types can be different than the defaults because of the udfs
    if (`object` == null) return null
    schema.getType match {
      case RECORD =>
        `object` match {
          case row: Row => return convertRowToAvroRecord(schema.toString(), row)
          case _        =>
        }
        throw new IllegalStateException("Row expected but was: " + `object`.getClass)
      case ENUM =>
        return new GenericData.EnumSymbol(schema, `object`.toString)
      case ARRAY =>
        val elementSchema  = schema.getElementType
        val array          = `object`.asInstanceOf[Array[Any]]
        val convertedArray = new GenericData.Array[Any](array.length, schema)
        for (element <- array) {
          convertedArray.add(convertFlinkType(elementSchema, element))
        }
        return convertedArray
      case MAP =>
        val map          = `object`.asInstanceOf[util.Map[_, _]]
        val convertedMap = new util.HashMap[Utf8, Any]
        for (entry <- map.entrySet.asScala) {
          convertedMap.put(
            new Utf8(entry.getKey.toString),
            convertFlinkType(schema.getValueType, entry.getValue)
          )
        }
        return convertedMap
      case UNION =>
        val types                = schema.getTypes
        val size                 = types.size
        var actualSchema: Schema = null
        if (size == 2 && (types.get(0).getType == Schema.Type.NULL)) actualSchema = types.get(1)
        else if (size == 2 && (types.get(1).getType == Schema.Type.NULL))
          actualSchema = types.get(0)
        else if (size == 1) actualSchema = types.get(0)
//          TODO
        else { // generic type
          return `object`
        }
        return convertFlinkType(actualSchema, `object`)
      case FIXED =>
        // check for logical type
        `object` match {
          case bigDecimal: BigDecimal =>
            return new GenericData.Fixed(schema, convertFromDecimal(schema, bigDecimal))
          case _ => return new GenericData.Fixed(schema, `object`.asInstanceOf[Array[Byte]])
        }
      case STRING =>
        return new Utf8(`object`.toString)
      case BYTES =>
        `object` match {
          case bigDecimal: BigDecimal =>
            return ByteBuffer.wrap(convertFromDecimal(schema, bigDecimal))
          case _ => return ByteBuffer.wrap(`object`.asInstanceOf[Array[Byte]])
        }
      case INT =>
        // check for logical types
        `object` match {
          case date: Date                => return convertFromDate(schema, date)
          case localDate: LocalDate      => return convertFromDate(schema, Date.valueOf(localDate))
          case time: Time                => return convertFromTimeMillis(schema, time)
          case localTime: LocalTime      => return convertFromTimeMillis(schema, Time.valueOf(localTime))
          case tinyInt: java.lang.Byte   => return tinyInt.toInt
          case smallInt: java.lang.Short => return smallInt.toInt
          case _                         => return `object`
        }
      case LONG =>
        `object` match {
          case timestamp: Timestamp => return convertFromTimestamp(schema, timestamp)
          case time: LocalDateTime  => return convertFromTimestamp(schema, Timestamp.valueOf(time))
          case time: Time           => return convertFromTimeMicros(schema, time)
//          TODO implement for  TIMESTAMP WITH TIME ZONE
//          case offsetDateTime: java.time.OffsetDateTime => {
//            offsetDateTime
//          }

          case time: java.time.Instant =>
            schema.getLogicalType match {
              case _: LogicalTypes.LocalTimestampMillis => return time.toEpochMilli
              case _: LogicalTypes.LocalTimestampMicros =>
                return TimeUnit.SECONDS.toMicros(time.getEpochSecond) + TimeUnit.NANOSECONDS
                  .toMicros(time.getNano)
            }

          case _ => return `object`
        }
      case FLOAT   => return `object`
      case DOUBLE  => return `object`
      case BOOLEAN => return `object`
      //      TODO NULL
      case NULL => return `object`
    }
    throw new RuntimeException("Unsupported Avro type:" + schema.getType)
  }

  private def convertFromDecimal(schema: Schema, decimal: BigDecimal) = {
    val logicalType = schema.getLogicalType
    logicalType match {
      case decimalType: LogicalTypes.Decimal =>
        // rescale to target type
        val rescaled = decimal.setScale(decimalType.getScale, BigDecimal.ROUND_UNNECESSARY)
        // byte array must contain the two's-complement representation of the
        // unscaled integer value in big-endian byte order
        decimal.unscaledValue.toByteArray
      case _ => throw new RuntimeException("Unsupported decimal type.")
    }
  }

  private def convertFromDate(schema: Schema, date: Date) = {
    val logicalType = schema.getLogicalType
    if (logicalType == LogicalTypes.date) { // adopted from Apache Calcite
      val converted = toEpochMillis(date)
      (converted / 86400000L).toInt
    } else throw new RuntimeException("Unsupported date type.")
  }

  private def toEpochMillis(date: java.util.Date) = {
    val time = date.getTime
    time + LOCAL_TZ.getOffset(time).toLong
  }

  private def convertFromTimeMillis(schema: Schema, date: Time) = {
    val logicalType = schema.getLogicalType
    if (logicalType == LogicalTypes.timeMillis) { // adopted from Apache Calcite
      val converted = toEpochMillis(date)
      (converted % 86400000L).toInt
    } else throw new RuntimeException("Unsupported time type.")
  }

  private def convertFromTimestamp(schema: Schema, date: Timestamp) = {
    val logicalType = schema.getLogicalType
    if (logicalType == LogicalTypes.timestampMillis) { // adopted from Apache Calcite
      val time = date.getTime
      time + LOCAL_TZ.getOffset(time).toLong
    } else if (logicalType == LogicalTypes.timestampMicros) {
      val millis = date.getTime
      val micros = millis * 1000 + (date.getNanos % 1000000 / 1000)
      val offset = LOCAL_TZ.getOffset(millis) * 1000L
      micros + offset
    } else if (logicalType == LogicalTypes.localTimestampMillis) {
      val time = date.getTime
      time
    } else if (logicalType == LogicalTypes.localTimestampMicros) {
      val millis = date.getTime
      val micros = millis * 1000 + (date.getNanos % 1000000 / 1000)
      micros
    } else throw new RuntimeException("Unsupported timestamp type.")
  }

  private def convertFromTimeMicros(schema: Schema, date: Time) = {
    val logicalType = schema.getLogicalType
    if (logicalType == LogicalTypes.timeMicros) { // adopted from Apache Calcite
      val converted = toEpochMillis(date)
      (converted % 86400000L) * 1000L
    } else throw new RuntimeException("Unsupported time type.")
  }

  def convertToSchema(logicalType: logical.LogicalType, rowName: String): Schema = {
//    TODO handle rowName clash in different fields
    var precision: Int    = 0
    val nullable: Boolean = logicalType.isNullable
    logicalType.getTypeRoot match {
      case org.apache.flink.table.types.logical.LogicalTypeRoot.NULL =>
        SchemaBuilder.builder.nullType
      case org.apache.flink.table.types.logical.LogicalTypeRoot.BOOLEAN =>
        val bool: Schema = SchemaBuilder.builder.booleanType
        if (nullable) {
          nullableSchema(bool)
        } else {
          bool
        }
      case logical.LogicalTypeRoot.INTEGER | logical.LogicalTypeRoot.TINYINT | logical.LogicalTypeRoot.SMALLINT =>
        val integer: Schema = SchemaBuilder.builder.intType
        if (nullable) {
          nullableSchema(integer)
        } else {
          integer
        }
      case logical.LogicalTypeRoot.BIGINT =>
        val bigint: Schema = SchemaBuilder.builder.longType
        if (nullable) {
          nullableSchema(bigint)
        } else {
          bigint
        }
      case logical.LogicalTypeRoot.FLOAT =>
        val f: Schema = SchemaBuilder.builder.floatType
        if (nullable) {
          nullableSchema(f)
        } else {
          f
        }
      case logical.LogicalTypeRoot.DOUBLE =>
        val d: Schema = SchemaBuilder.builder.doubleType
        if (nullable) {
          nullableSchema(d)
        } else {
          d
        }
      case logical.LogicalTypeRoot.CHAR | logical.LogicalTypeRoot.VARCHAR =>
        val str: Schema = SchemaBuilder.builder.stringType
        if (nullable) {
          nullableSchema(str)
        } else {
          str
        }
      case logical.LogicalTypeRoot.BINARY | logical.LogicalTypeRoot.VARBINARY =>
        val binary: Schema = SchemaBuilder.builder.bytesType
        if (nullable) {
          nullableSchema(binary)
        } else {
          binary
        }
      case logical.LogicalTypeRoot.TIMESTAMP_WITHOUT_TIME_ZONE =>
        // use long to represents Timestamp
        val timestampType: TimestampType = logicalType.asInstanceOf[TimestampType]
        precision = timestampType.getPrecision
        var avroLogicalType: org.apache.avro.LogicalType = null
        if (precision == 3) {
          avroLogicalType = LogicalTypes.timestampMillis
        } else if (precision == 6) {
          avroLogicalType = LogicalTypes.timestampMicros()
        } else {
//          TODO without logical type
          avroLogicalType = LogicalTypes.timestampMicros()
        }
        val timestamp: Schema = avroLogicalType.addToSchema(SchemaBuilder.builder.longType)
        if (nullable) {
          nullableSchema(timestamp)
        } else {
          timestamp
        }
      case logical.LogicalTypeRoot.TIMESTAMP_WITH_LOCAL_TIME_ZONE =>
        // use long to represents Timestamp
        val timestampType: LocalZonedTimestampType =
          logicalType.asInstanceOf[LocalZonedTimestampType]
        precision = timestampType.getPrecision
        var avroLogicalType: org.apache.avro.LogicalType = null
        if (precision == 3) {
          avroLogicalType = LogicalTypes.localTimestampMillis
        } else if (precision == 6) {
          avroLogicalType = LogicalTypes.localTimestampMicros
        } else {
          avroLogicalType = LogicalTypes.localTimestampMicros()
        }
        val timestamp: Schema = avroLogicalType.addToSchema(SchemaBuilder.builder.longType)
        if (nullable) {
          nullableSchema(timestamp)
        } else {
          timestamp
        }
//        TODO TIMESTAMP_WITH_TIME_ZONE
//      case logical.LogicalTypeRoot.TIMESTAMP_WITH_TIME_ZONE => ???
//      TODO INTERVAL

      case logical.LogicalTypeRoot.DATE =>
        // use int to represents Date
        val date: Schema = LogicalTypes.date.addToSchema(SchemaBuilder.builder.intType)
        if (nullable) {
          nullableSchema(date)
        } else {
          date
        }
      case logical.LogicalTypeRoot.TIME_WITHOUT_TIME_ZONE =>
        precision = (logicalType.asInstanceOf[TimeType]).getPrecision
        var time: Schema = LogicalTypes.timeMillis.addToSchema(SchemaBuilder.builder.intType)
        if (precision == 6) {
          time = LogicalTypes.timeMicros.addToSchema(SchemaBuilder.builder.longType)
        } else if (precision == 3) {
          time = LogicalTypes.timeMillis.addToSchema(SchemaBuilder.builder.intType)
        } else {
//          TODO don't use avro logical type
          throw new IllegalArgumentException(
            "Avro does not support TIME type with precision: " + precision + ", it only supports precisions 3 and 6"
          )
        }
        if (nullable) {
          nullableSchema(time)
        } else {
          time
        }
      case logical.LogicalTypeRoot.DECIMAL =>
        val decimalType: DecimalType = logicalType.asInstanceOf[DecimalType]
        // store BigDecimal as byte[]
        val decimal: Schema = LogicalTypes
          .decimal(decimalType.getPrecision, decimalType.getScale)
          .addToSchema(SchemaBuilder.builder.bytesType)
        if (nullable) {
          nullableSchema(decimal)
        } else {
          decimal
        }
      case logical.LogicalTypeRoot.ROW =>
        val rowType: RowType                   = logicalType.asInstanceOf[RowType]
        val fieldNames: java.util.List[String] = rowType.getFieldNames
        // we have to make sure the record name is different in a Schema
        var builder: SchemaBuilder.FieldAssembler[Schema] =
          SchemaBuilder.builder.record(rowName).fields
        for (i <- 0 until rowType.getFieldCount) {
          val fieldName: String              = fieldNames.get(i)
          val fieldType: logical.LogicalType = rowType.getTypeAt(i)
          val fieldBuilder: SchemaBuilder.GenericDefault[Schema] =
            builder.name(fieldName).`type`(convertToSchema(fieldType, rowName + "_" + fieldName))
          if (fieldType.isNullable) {
            builder = fieldBuilder.withDefault(null)
          } else {
            builder = fieldBuilder.noDefault
          }
        }
        val record: Schema = builder.endRecord
        if (nullable) {
          nullableSchema(record)
        } else {
          record
        }
      case logical.LogicalTypeRoot.MULTISET | logical.LogicalTypeRoot.MAP =>
        val map: Schema = SchemaBuilder.builder.map.values(
          convertToSchema(extractValueTypeToAvroMap(logicalType), rowName)
        )
        if (nullable) {
          nullableSchema(map)
        } else {
          map
        }
      case logical.LogicalTypeRoot.ARRAY =>
        val arrayType: ArrayType = logicalType.asInstanceOf[ArrayType]
        val array: Schema =
          SchemaBuilder.builder.array.items(convertToSchema(arrayType.getElementType, rowName))
        if (nullable) {
          nullableSchema(array)
        } else {
          array
        }
      case logical.LogicalTypeRoot.RAW =>
        throw new UnsupportedOperationException(
          "Unsupported to derive Schema for type: " + logicalType
        )
      case _ => {
        throw new UnsupportedOperationException(
          "Unsupported to derive Schema for type: " + logicalType
        )
      }

    }
  }

  def extractValueTypeToAvroMap(`type`: logical.LogicalType): logical.LogicalType = {
    var keyType: logical.LogicalType   = null
    var valueType: logical.LogicalType = null
    `type` match {
      case mapType: MapType =>
        keyType = mapType.getKeyType
        valueType = mapType.getValueType
      case _ =>
        val multisetType = `type`.asInstanceOf[MultisetType]
        keyType = multisetType.getElementType
        valueType = new IntType
    }
    if (!keyType.getTypeRoot.getFamilies.contains(LogicalTypeFamily.CHARACTER_STRING))
      throw new UnsupportedOperationException(
        "Avro format doesn't support non-string as key type of map. " + "The key type is: " + keyType.asSummaryString
      )
    valueType
  }

  private def nullableSchema(schema: Schema) =
    if (schema.isNullable) schema
    else Schema.createUnion(SchemaBuilder.builder.nullType, schema)

}
