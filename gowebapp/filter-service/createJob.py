import os
from typing import List, Optional
from pyflink.datastream import StreamExecutionEnvironment
from pyflink.table import StreamTableEnvironment, EnvironmentSettings
# from pyflink.datastream.connectors.kafka import *

from models import Filters

# BootstrapServer = 'my-cluster-kafka-bootstrap.kafka.svc.cluster.local:9092'
# BootstrapServer = 'my-cluster-kafka-bootstrap.logfire-local.svc.cluster.local:9092'
BootstrapServer = 'my-cluster-kafka-bootstrap:9092'
# groupId = 'source_consumer_group'
groupId = 'go-kafka-consumer'


def CreateSourceTable(kafkaTopic: str, bootstrapServer: str, groupId: str) -> str:
    src_ddl = f"""
        CREATE TABLE logs_rawdata (
            platform VARCHAR,
            message VARCHAR,
            msgid VARCHAR,
            level VARCHAR,
            dt TIMESTAMP_LTZ(3),
            pid INT,
            version INT
        ) WITH (
            'connector' = 'kafka',
            'topic' = '{kafkaTopic}',
            'properties.bootstrap.servers' = '{bootstrapServer}',
            'properties.group.id' = '{groupId}',
            'properties.auto.offset.reset' = 'latest',
            'format' = 'json',
            'json.timestamp-format.standard' = 'ISO-8601'
        )
    """
    print("source sql table: ", src_ddl)
    return src_ddl


# env = StreamExecutionEnvironment.get_execution_environment()

# def AddSources(env: StreamExecutionEnvironment, topicLists: List[str], bootstrapserver: str, groupId: str):
#     for topic in topicLists:
#         env.add_source()


def CreateSinkTable(kafkat: str, bootstrapServer: str) -> str:
    sink_ddl = f"""
        CREATE TABLE filter_data (
            platform VARCHAR,
            message VARCHAR,
            msgid VARCHAR,
            level VARCHAR,
            dt VARCHAR,
            pid INT,
            version INT
        ) WITH (
            'connector' = 'kafka',
            'topic' = '{kafkat}',
            'properties.bootstrap.servers' = '{bootstrapServer}',
            'format' = 'json'
        )
    """
    print("sink sql table: ", sink_ddl)
    return sink_ddl


# def FilterSqlQuery(kafkat: str, bootstrapServer: str) -> str:
#     sql = """
#         SELECT
#           *
#         FROM logs_rawdata
#         WHERE hostname='we.org'
#     """
#     return sink_ddl

def parseList(l: List[str]):
    return "(" + ", ".join(repr(e) for e in l) + ")"


def MyFilter(platform: Optional[List[str]], level: Optional[List[str]], start_time: Optional[str],
             end_time: Optional[str], message_contains: Optional[str]) -> str:
    # sql = ""
    #
    # if platform and level:
    #     sql = f"SELECT * FROM logs_rawdata WHERE platform IN {parseList(platform)} AND level IN {parseList(level)};"
    # elif platform and not level:
    #     sql = f"SELECT * FROM logs_rawdata WHERE platform IN {parseList(platform)};"
    # elif not platform and level:
    #     sql = f"SELECT * FROM logs_rawdata WHERE level IN {parseList(level)};"

    list_of_filters = []
    if platform:
        list_of_filters.append(f"platform IN {parseList(platform)}")
    if level:
        list_of_filters.append(f"level IN {parseList(level)}")
    if start_time:
        list_of_filters.append(f"dt >= '{start_time}'")
    if end_time:
        list_of_filters.append(f"dt <= '{end_time}'")
    if message_contains:
        list_of_filters.append(f"message like '%{message_contains}%'")

    conditions = " AND ".join(list_of_filters)
    sql = "SELECT platform, message , msgid, level, cast(dt as string), pid, version FROM logs_rawdata WHERE " + conditions
    print("my filter sql: ", sql)
    return sql


def NewJob(appData: Filters):
    # Create streaming environment
    print("started creating new job")
    env = StreamExecutionEnvironment.get_execution_environment()

    settings = EnvironmentSettings.new_instance(
    ).in_streaming_mode(
    ).build()

    # create table environment
    tbl_env = StreamTableEnvironment.create(stream_execution_environment=env,
                                            environment_settings=settings)

    # add kafka connector dependency
    kafka_jar = os.path.join(
        os.path.abspath(os.path.dirname(__file__)),
        'bin/flink-sql-connector-kafka-1.16.1.jar'
    )

    tbl_env.get_config() \
        .get_configuration() \
        .set_string("pipeline.jars", "file://{}".format(kafka_jar))

    #######################################################################
    # Create Kafka Source Table with DDL
    #######################################################################
    # logs_topic = topic
    src_ddl = CreateSourceTable(kafkaTopic=appData.sourcetopic, bootstrapServer=BootstrapServer, groupId=groupId)

    tbl_env.execute_sql(src_ddl)

    # create and initiate loading of source Table
    tbl = tbl_env.from_path('logs_rawdata')

    print('\nSource Schema')
    tbl.print_schema()

    #####################################################################
    # Define Tumbling Window Aggregate Calculation
    #####################################################################
    sql = MyFilter(appData.platforms, appData.levels, appData.start_time, appData.end_time, appData.message_contains)

    filter_tbl = tbl_env.sql_query(sql)

    print('\nProcess Sink Schema')
    filter_tbl.print_schema()

    ###############################################################
    # Create Kafka Sink Table
    ###############################################################
    sink_ddl = CreateSinkTable(appData.sinktopic, bootstrapServer=BootstrapServer)

    tbl_env.execute_sql(sink_ddl)

    # write time windowed aggregations to sink table
    # filter_tbl.execute_insert('filter_data').wait()
    filter_tbl.execute_insert('filter_data')
    # job_client = filter_tbl.execute_insert('filter_data').get_job_client()

    # jobId = job_client.get_job_id()
    # print("my job id: "+ jobId)

