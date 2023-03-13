FROM registry.access.redhat.com/ubi8/openjdk-11

COPY filter-job/target/scala-2.12/filter-job-assembly-0.1-SNAPSHOT-deps.jar /opt/filter-job-assembly-0.1-SNAPSHOT-deps.jar
COPY filter-job/target/scala-2.12/filter-job-assembly-0.1-SNAPSHOT.jar /opt/filter-job-assembly-0.1-SNAPSHOT.jar

ENTRYPOINT ["java", "-cp", "/opt/filter-job-assembly-0.1-SNAPSHOT-deps.jar:/opt/filter-job-assembly-0.1-SNAPSHOT.jar", "sh.logfire.flink.filter.RequestServer"]