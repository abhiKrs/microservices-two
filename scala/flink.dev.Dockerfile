FROM flink:1.16.1-scala_2.12-java11
COPY flink-cluster/target/scala-2.12/flink-cluster-assembly-0.1-SNAPSHOT.jar /opt/flink/lib
COPY filter-job/target/scala-2.12/filter-job-assembly-0.1-SNAPSHOT.jar /opt/flink/lib