package sh.logfire.flink.filter

object Constants {
  final val FLINK_CLUSTER_HOST      = "FLINK_CLUSTER_HOST"
  final val FLINK_CLUSTER_PORT      = "FLINK_CLUSTER_PORT"
  final val KAFKA_BROKERS           = "KAFKA_BROKERS"
  final val SCHEMA_REGISTRY_ADDRESS = "SCHEMA_REGISTRY_ADDRESS"
  final val SEVERITY_LEVEL_FIELD    = "SEVERITY_LEVEL_FIELD"
  final val MESSAGE_FIELD           = "MESSAGE_FIELD"
  final val LOGFIRE_TIMESTAMP_FIELD = "LOGFIRE_TIMESTAMP_FIELD"
  final val LOGFIRE_ROWTIME_FIELD   = "LOGFIRE_ROWTIME_FIELD"
  final val REDIS_HOST              = "REDIS_HOST"
  final val REDIS_PORT              = "REDIS_PORT"
  final val REDIS_PASSWORD          = "REDIS_PASSWORD"
  final val FLINK_JOBS_REDIS_KEY    = "FLINK_JOBS_REDIS_KEY"
}
