lazy val commonSettings = Seq(
  resolvers ++= Seq(
    Resolver.mavenLocal,
    "Local Maven Repository" at "file://" + Path.userHome.absolutePath + "/.m2/repository",
    Resolver.sonatypeRepo("releases"),
    Resolver.sonatypeRepo("snapshots"),
    //    "Apache Development Snapshot Repository" at "https://repository.apache.org/content/repositories/snapshots/",
    "confluent" at "https://packages.confluent.io/maven/",
    "jitpack" at "https://jitpack.io",
    "mapr" at "https://repository.mapr.com/maven/",
    "jboss" at "https://repository.jboss.org/nexus/content/repositories/thirdparty-releases/",
    "github" at "https://maven.pkg.github.com/argoproj/argo-workflows/",
    "gitlab" at "https://gitlab.com/api/v4/projects/42331414/packages/maven/",
    //    "Apache Development Candidate Repository" at "https://repository.apache.org/content/repositories/orgapacheflink-1488/",
    //    "Cascading repo" at "https://conjars.org/repo"
  ),
  credentials +=
    Credentials("Gitlab Package Registry",
                "https://gitlab.com/api/v4/projects/42331414/packages/maven",
                "gitlab-maven",
                "glpat-4aSB6sHzqNuULannKaxc"),
  Compile / run := Defaults
    .runTask(Compile / fullClasspath, Compile / run / mainClass, Compile / run / runner)
    .evaluated,
  Compile / run / fork := true,
  Global / cancelable := true,
  assembly / assemblyMergeStrategy := {
    case path if path.contains("META-INF/services") => MergeStrategy.concat
    case "META-INF/services/org.apache.flink.table.factories.Factory" =>
      MergeStrategy.concat
    case "META-INF/services/org.apache.flink.table.factories.TableFactory" =>
      MergeStrategy.concat
    case "META-INF/services/org.apache.flink.table.factories.DynamicTableFactory" =>
      MergeStrategy.concat
    case "META-INF/services/org.apache.iceberg.flink.FlinkDynamicTableFactory" => MergeStrategy.concat
    case "META-INF/services/org.apache.iceberg.flink.FlinkCatalogFactory"      => MergeStrategy.concat
    case "META-INF/services/org/apache/flink/table/gateway/api/endpoint/SqlGatewayEndpointFactory" =>
      MergeStrategy.concat
    case "META-INF/services/org.apache.flink.connector.jdbc.table.JdbcDynamicTableFactory" =>
      MergeStrategy.concat
    case PathList("META-INF", _*) => MergeStrategy.discard
    case "reference.conf"         => MergeStrategy.concat
    case _                        => MergeStrategy.first
  },
  ThisBuild / scalaVersion := "2.12.12",
  version := "0.1-SNAPSHOT",
  avroStringType := "String"
)

val scalaPbSettings = Seq(
  Compile / PB.targets := Seq(
    scalapb.validate
      .preprocessor()      -> (Compile / sourceManaged).value / "scalapb",
    scalapb.gen()          -> (Compile / sourceManaged).value / "scalapb",
    scalapb.validate.gen() -> (Compile / sourceManaged).value / "scalapb"
  ),
  Compile / PB.protoSources += file("protobuf")
)

lazy val settings =
commonSettings ++ scalaPbSettings

name := "scala"

version := "0.1-SNAPSHOT"

organization := "sh.logfire"

lazy val global = project
  .in(file("."))
  .settings {
    settings
  }
  .disablePlugins(AssemblyPlugin)
  .aggregate(
    filterJob,
    flinkCluster,
  )

val depsSeparateSettings = assembly / assemblyOption := (assembly / assemblyOption).value
  .copy(includeScala = false, includeDependency = false)

lazy val filterJob = (project in file("filter-job"))
  .settings(
    name := "filter-job",
    mainClass := Some("sh.logfire.flink.filter.Application"),
    settings,
    test in assembly := {},
    libraryDependencies ++= flinkProvidedDependencies ++ extraDependenciesFromFlink ++ scalaPbDeps ++ filterJobDeps ++ flowTestDeps ++ flinkProvidedDependenciesCompile ++ filterJobsRequiredDeps,
    depsSeparateSettings
  )

lazy val flinkCluster = (project in file("flink-cluster"))
  .settings(
    name := "flink-cluster",
    commonSettings,
    test in assembly := {},
    libraryDependencies ++= flinkProvidedDependencies ++ extraDependenciesFromFlink ++ filterJobDeps ++ flowTestDeps
  )

lazy val dependencies =
  new {
    val flinkVersion               = "1.16.1"
    val confluentVersion           = "7.3.2"
    val jacksonVersion             = "2.12.5"
    val kafkaVersion               = "2.8.0"
    val avroVersion                = "1.10.2"
    val awsJavaVersion             = "1.11.951"
    val testcontainersScalaVersion = "0.40.2"
    val akkaHttpVersion            = "10.2.7"
    val akkaVersion                = "2.7.0"
    val hadoopVersion              = "3.3.4"
    val deltaConnectorsVersion     = "0.6.0"
    val sparkVersion               = "3.2.2"
    val icebergVersion             = "1.1.0"
    val circeVersion               = "0.14.4"

    val flinkClients               = "org.apache.flink" % "flink-clients"                 % flinkVersion
    val flinkScala                 = "org.apache.flink" %% "flink-scala"                  % flinkVersion
    val flinkStreamingScala        = "org.apache.flink" %% "flink-streaming-scala"        % flinkVersion
    val flinkTableApiScalaBridge   = "org.apache.flink" %% "flink-table-api-scala-bridge" % flinkVersion
    val flinkStateBackendRocksdb   = "org.apache.flink" % "flink-statebackend-rocksdb"    % flinkVersion
    val flinkConnectorJdbc         = "org.apache.flink" % "flink-connector-jdbc"          % flinkVersion
    val flinkTablePlanner          = "org.apache.flink" %% "flink-table-planner"          % flinkVersion
    val flinkTableCommon           = "org.apache.flink" % "flink-table-common"            % flinkVersion
    val flinkPython                = "org.apache.flink" % "flink-python"                  % flinkVersion
    val flinkCore                  = "org.apache.flink" % "flink-core"                    % flinkVersion
    val flinkSqlParquet            = "org.apache.flink" % "flink-sql-parquet"             % flinkVersion
    val flinkConnectorKafka        = "org.apache.flink" % "flink-connector-kafka"         % flinkVersion
    val flinkAvro                  = "org.apache.flink" % "flink-avro"                    % flinkVersion
    val flinkAvroConfluentRegistry = "org.apache.flink" % "flink-avro-confluent-registry" % flinkVersion
    val flinkCsv                   = "org.apache.flink" % "flink-csv"                     % flinkVersion
    val flinkOrc                   = "org.apache.flink" % "flink-orc"                     % flinkVersion
    val flinkMaprFs                = "org.apache.flink" % "flink-mapr-fs"                 % "1.14.4"
    val flinkCEP                   = "org.apache.flink" % "flink-cep"                     % flinkVersion
    val flinkConnectorHive         = "org.apache.flink" %% "flink-connector-hive"         % flinkVersion exclude ("org.apache.hive", "hive-exec") exclude ("org.apache.hive", "hive-metastore") exclude ("org.apache.avro", "avro") exclude ("org.antlr", "antlr-runtime") exclude ("org.apache.hadoop", "hadoop-common") exclude ("org.apache.hadoop", "hadoop-hdfs") exclude ("org.apache.hadoop", "hadoop-mapreduce-client-core") exclude ("org.apache.hadoop", "hadoop-yarn-common") exclude ("org.apache.hadoop", "hadoop-yarn-client")
    val flinkMetricsPrometheus     = "org.apache.flink" % "flink-metrics-prometheus"      % flinkVersion
    val flinkMetricsInfluxDB       = "org.apache.flink" % "flink-metrics-influxdb"        % flinkVersion
    val flinkStateProcessorApi     = "org.apache.flink" % "flink-state-processor-api"     % flinkVersion
    val flinkS3Presto              = "org.apache.flink" % "flink-s3-fs-presto"            % flinkVersion
    val flinkS3Hadoop              = "org.apache.flink" % "flink-s3-fs-hadoop"            % flinkVersion
    val flinkConnectorFiles        = "org.apache.flink" % "flink-connector-files"         % flinkVersion
    val flinkSqlGateway            = "org.apache.flink" % "flink-sql-gateway-api"         % flinkVersion
    val flinkHbase                 = "org.apache.flink" % "flink-connector-hbase-1.4"     % flinkVersion exclude ("org.apache.hbase", "hbase-client") exclude ("org.mortbay.jetty", "jetty-util") exclude ("org.mortbay.jetty", "jetty") exclude ("org.mortbay.jetty", "jetty-sslengine") exclude ("org.mortbay.jetty", "jsp-2.1") exclude ("org.mortbay.jetty", "jsp-api-2.1") exclude ("org.mortbay.jetty", "servlet-api-2.5") exclude ("org.apache.hbase", "hbase-annotations") exclude ("com.sun.jersey", "jersey-core") exclude ("org.apache.hadoop", "hadoop-common") exclude ("org.apache.hadoop", "hadoop-auth") exclude ("org.apache.hadoop", "hadoop-annotations") exclude ("org.apache.hadoop", "hadoop-mapreduce-client-core") exclude ("org.apache.hadoop", "hadoop-client") exclude ("org.apache.hadoop", "hadoop-hdfs") exclude ("log4j", "log4j") exclude ("org.slf4j", "slf4j-log4j12")
    val flinkK8sOperator           = "org.apache.flink" % "flink-kubernetes-operator"     % "1.3.1"

    val grpcNetty           = "io.grpc"                    % "grpc-netty"             % scalapb.compiler.Version.grpcJavaVersion
    val grpcServices        = "io.grpc"                    % "grpc-services"          % scalapb.compiler.Version.grpcJavaVersion
    val scalaPbRuntimeGrpc  = "com.thesamet.scalapb"       %% "scalapb-runtime-grpc"  % scalapb.compiler.Version.scalapbVersion
    val scalaPbValidateCore = "com.thesamet.scalapb"       %% "scalapb-validate-core" % scalapb.validate.compiler.BuildInfo.version % "protobuf"
    val scalaPbValidateCats = "com.thesamet.scalapb"       %% "scalapb-validate-cats" % scalapb.validate.compiler.BuildInfo.version
    val catsCore            = "org.typelevel"              %% "cats-core"             % "2.6.1"
    val scalaPbRuntime      = "com.thesamet.scalapb"       %% "scalapb-runtime"       % scalapb.compiler.Version.scalapbVersion % "protobuf"
    val scalaPbJson4s       = "com.thesamet.scalapb"       %% "scalapb-json4s"        % "0.11.1"
    val jacksonCore         = "com.fasterxml.jackson.core" % "jackson-core"           % jacksonVersion
    val jacksonDatabind     = "com.fasterxml.jackson.core" % "jackson-databind"       % jacksonVersion
    val jodaTime            = "joda-time"                  % "joda-time"              % "2.10.13"
    val kafkaAvroSerializer = "io.confluent"               % "kafka-avro-serializer"  % confluentVersion
    val slf4j               = "org.slf4j"                  % "slf4j-simple"           % "2.0.0-alpha5"
    val circeCore           = "io.circe"                   %% "circe-core"            % circeVersion
    val circeGeneric        = "io.circe"                   %% "circe-generic"         % circeVersion
    val circeParser         = "io.circe"                   %% "circe-parser"          % circeVersion
    val jedis               = "redis.clients"              % "jedis"                  % "4.4.0-m2"
    val json                = "org.json"                   % "json"                   % "20230227"

    //    TEST
    val scalaTest               = "org.scalatest"                 %% "scalatest"                      % "3.2.9"
    val flinkTestUtils          = "org.apache.flink"              % "flink-test-utils"                % flinkVersion
    val embeddedSchemaRegistryy = "io.github.embeddedkafka"       %% "embedded-kafka-schema-registry" % confluentVersion
    val gson                    = "com.google.code.gson"          % "gson"                            % "2.8.8"
    val jsonToAvro              = "tech.allegro.schema.json2avro" % "converter"                       % "0.2.10"
    val testContainersScalaTest = "com.dimafeng"                  %% "testcontainers-scala-scalatest" % testcontainersScalaVersion

  }

val flinkProvidedDependencies = Seq(
  dependencies.flinkClients             % Provided,
  dependencies.flinkScala               % Provided,
  dependencies.flinkStreamingScala      % Provided,
  dependencies.flinkStateBackendRocksdb % Provided,
  dependencies.flinkTablePlanner        % Provided,
  dependencies.flinkPython              % Provided,
  dependencies.flinkCore                % Provided,
  dependencies.flinkCEP                 % Provided,
  dependencies.flinkConnectorFiles      % Provided,
)

val flinkProvidedDependenciesCompile = Seq(
  dependencies.flinkClients,
  dependencies.flinkScala,
  dependencies.flinkStreamingScala,
  dependencies.flinkStateBackendRocksdb,
  dependencies.flinkTablePlanner,
  dependencies.flinkPython,
  dependencies.flinkCore,
  dependencies.flinkCEP,
  dependencies.flinkConnectorFiles,
)

val extraDependenciesFromFlink = Seq(
  dependencies.flinkSqlParquet,
  dependencies.flinkConnectorKafka,
  dependencies.flinkAvro,
  dependencies.flinkAvroConfluentRegistry,
  dependencies.flinkCsv,
  dependencies.flinkOrc,
  dependencies.flinkConnectorJdbc,
  dependencies.flinkTableApiScalaBridge,
  dependencies.flinkHbase,
  dependencies.flinkConnectorHive,
  dependencies.flinkSqlGateway,
)

val scalaPbDeps = Seq(
  dependencies.grpcNetty,
  dependencies.grpcServices,
  dependencies.scalaPbRuntimeGrpc,
  dependencies.scalaPbValidateCore,
  dependencies.scalaPbValidateCats,
  dependencies.catsCore,
  dependencies.scalaPbRuntime,
  dependencies.scalaPbJson4s
)

val filterJobDeps = Seq(
  dependencies.jodaTime,
  dependencies.kafkaAvroSerializer,
  dependencies.jedis,
  dependencies.json,
  dependencies.slf4j % Provided,
)

val filterJobsRequiredDeps = Seq(
  dependencies.slf4j
)

val flowTestDeps = Seq(
  dependencies.scalaTest,
  dependencies.flinkTestUtils,
  dependencies.embeddedSchemaRegistryy,
  dependencies.gson,
  dependencies.jsonToAvro,
  dependencies.testContainersScalaTest,
  dependencies.circeCore,
  dependencies.circeGeneric,
  dependencies.circeParser,
).map(_ % Test)
