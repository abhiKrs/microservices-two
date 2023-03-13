FROM centos:centos7 as builder

RUN yum update -y && \
    yum install -y sudo && \
    yum install -y java-11-openjdk-devel && \
    sudo rm -f /etc/yum.repos.d/bintray-rpm.repo && \
    curl -L https://www.scala-sbt.org/sbt-rpm.repo > sbt-rpm.repo && \
    sudo mv sbt-rpm.repo /etc/yum.repos.d/ && \
    sudo yum install -y sbt && \
    yum clean all

COPY build.sbt /app/
COPY project /app/project
COPY protobuf /app/protobuf

RUN (cd /app/; sbt -J-Xmx2048m assemblyPackageDependency)

COPY filter-job /app/filter-job

RUN (cd /app/; sbt -J-Xmx4096m -J-Xss4M -J-Xms1024m -J-Xss100M -mem 2096  assembly)

FROM registry.access.redhat.com/ubi8/openjdk-11

COPY --from=builder /app/filter-job/target/scala-2.12/filter-job-assembly-0.1-SNAPSHOT-deps.jar /opt/filter-job-assembly-0.1-SNAPSHOT-deps.jar
COPY --from=builder /app/filter-job/target/scala-2.12/filter-job-assembly-0.1-SNAPSHOT.jar /opt/filter-job-assembly-0.1-SNAPSHOT.jar

ENTRYPOINT ["java", "-cp", "/opt/filter-job-assembly-0.1-SNAPSHOT-deps.jar:/opt/filter-job-assembly-0.1-SNAPSHOT.jar", "sh.logfire.flink.filter.RequestServer"]