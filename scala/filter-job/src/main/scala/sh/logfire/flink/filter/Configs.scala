package sh.logfire.flink.filter

object Configs {
  val FLINK_CLUSTER_HOST      = System.getenv().getOrDefault(Constants.FLINK_CLUSTER_HOST, "g14")
  val FLINK_CLUSTER_PORT      = System.getenv().getOrDefault(Constants.FLINK_CLUSTER_PORT, "32219")
  val KAFKA_BROKERS           = System.getenv().getOrDefault(Constants.KAFKA_BROKERS, "g14:9092")
  val SCHEMA_REGISTRY_ADDRESS = System.getenv().getOrDefault(Constants.SCHEMA_REGISTRY_ADDRESS, "http://g14:30081")
  val SEVERITY_LEVEL_FIELD    = System.getenv().getOrDefault(Constants.SEVERITY_LEVEL_FIELD, "level")
  val MESSAGE_FIELD           = System.getenv().getOrDefault(Constants.MESSAGE_FIELD, "message")
  val LOGFIRE_TIMESTAMP_FIELD = System.getenv().getOrDefault(Constants.LOGFIRE_TIMESTAMP_FIELD, "dt")
  val LOGFIRE_ROWTIME_FIELD   = System.getenv().getOrDefault(Constants.LOGFIRE_ROWTIME_FIELD, "rowtime")
  val REDIS_HOST              = System.getenv().getOrDefault(Constants.REDIS_HOST, "g14")
  val REDIS_PORT              = System.getenv().getOrDefault(Constants.REDIS_PORT, "32302")
  val REDIS_PASSWORD          = System.getenv().getOrDefault(Constants.REDIS_PASSWORD, "")
  val FLINK_JOBS_REDIS_KEY    = System.getenv().getOrDefault(Constants.FLINK_JOBS_REDIS_KEY, "flink-jobs")

}
