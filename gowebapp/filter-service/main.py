# import argparse
import json
import logging
import sys
from typing import Dict

from createJob import NewJob
from models import Filters
import os
# from typing import List, Optional
from pyflink.datastream import StreamExecutionEnvironment
from pyflink.table import StreamTableEnvironment, EnvironmentSettings

logging.basicConfig(
  format='%(asctime)s %(name)-12s %(levelname)-8s %(message)s',
  datefmt='%Y-%m-%d %H:%M:%S',
  level=logging.INFO,
  handlers=[
      logging.StreamHandler(sys.stdout)
  ]
)

logger = logging.getLogger()


def main():
    print("in main")
    
    with open('/tmp/tmpfile-filtered-job-aaplication.txt', "r") as f:
        f.seek(0)
        fileObj = f.read()
        print("file data : ", fileObj)

        obj: Dict[str, any] = json.loads(fileObj)

    print(obj)
    
    filterData = Filters().from_json(obj)
    logger.debug(filterData)
    # NewJob(sourceTopic=args.sourcetopic, sinkTopic=args.sinktopic, filter=filterData)
    NewJob(appData=filterData)



if __name__ == '__main__':
    main()


# -----------------------------------

# def write_to_kafka(env):

# def main():
#     # Create streaming environment
#     env = StreamExecutionEnvironment.get_execution_environment()

#     settings = EnvironmentSettings.new_instance(
#             ).in_streaming_mode(
#             ).build()

#     # create table environment
#     tbl_env = StreamTableEnvironment.create(stream_execution_environment=env,
#                                             environment_settings=settings)

#     # add kafka connector dependency
#     kafka_jar = os.path.join(
#         os.path.abspath(os.path.dirname(__file__)),
#         'bin/flink-sql-connector-kafka-1.16.1.jar'
#     )

#     tbl_env.get_config()\
#             .get_configuration()\
#             .set_string("pipeline.jars", "file://{}".format(kafka_jar))

#     #######################################################################
#     # Create Kafka Source Table with DDL
#     #######################################################################
#     logs_topic = "my_topic"
#     src_ddl = f"""
#         CREATE TABLE logs_rawdata (
#             platform VARCHAR,
#             message VARCHAR,
#             msgid VARCHAR,
#             level VARCHAR,
#             dt VARCHAR,
#             pid INT,
#             version INT
#         ) WITH (
#             'connector' = 'kafka',
#             'topic' = '{logs_topic}',
#             'properties.bootstrap.servers' = 'my-cluster-kafka-bootstrap.kafka:9092',
#             'properties.group.id' = 'sales-usd',
#             'properties.auto.offset.reset' = 'latest',
#             'format' = 'json'
#         )
#     """
#     # src_ddl = """
#     #     CREATE TABLE logs_rawdata (
#     #         seller_id VARCHAR,
#     #         amount_usd DOUBLE,
#     #         sale_ts BIGINT,
#     #         proctime AS PROCTIME()
#     #     ) WITH (
#     #         'connector' = 'kafka',
#     #         'topic' = 'logs-topic',
#     #         'properties.bootstrap.servers' = 'localhost:9092',
#     #         'properties.group.id' = 'sales-usd',
#     #         'properties.auto.offset.reset' = 'latest',
#     #         'format' = 'json'
#     #     )
#     # """

#     tbl_env.execute_sql(src_ddl)

#     # create and initiate loading of source Table
#     tbl = tbl_env.from_path('logs_rawdata')

#     print('\nSource Schema')
#     tbl.print_schema()

#     #####################################################################
#     # Define Tumbling Window Aggregate Calculation (Seller Sales Per Minute)
#     #####################################################################
#     sql = """
#         SELECT
#           *
#         FROM logs_rawdata
#         WHERE message='we.org'
#     """
#     # sql = """
#     #     SELECT
#     #       seller_id,
#     #       TUMBLE_END(proctime, INTERVAL '60' SECONDS) AS window_end,
#     #       SUM(amount_usd) * 0.85 AS window_sales
#     #     FROM logs_rawdata
#     #     WHERE 
#     #     GROUP BY
#     #       TUMBLE(proctime, INTERVAL '60' SECONDS),
#     #       seller_id
#     # """
#     revenue_tbl = tbl_env.sql_query(sql)

#     print('\nProcess Sink Schema')
#     revenue_tbl.print_schema()

#     ###############################################################
#     # Create Kafka Sink Table
#     ###############################################################
#     sink_ddl = """
#         CREATE TABLE filter_data (
#             platform VARCHAR,
#             message VARCHAR,
#             msgid VARCHAR,
#             level VARCHAR,
#             dt VARCHAR,
#             pid INT,
#             version INT
#         ) WITH (
#             'connector' = 'kafka',
#             'topic' = 'flink-topic',
#             'properties.bootstrap.servers' = 'my-cluster-kafka-bootstrap.kafka:9092',
#             'format' = 'json'
#         )
#     """
#     # sink_ddl = """
#     #     CREATE TABLE sales_euros (
#     #         seller_id VARCHAR,
#     #         window_end TIMESTAMP(3),
#     #         window_sales DOUBLE
#     #     ) WITH (
#     #         'connector' = 'kafka',
#     #         'topic' = 'flink-topic',
#     #         'properties.bootstrap.servers' = 'localhost:9092',
#     #         'format' = 'json'
#     #     )
#     # """
#     tbl_env.execute_sql(sink_ddl)

#     # write time windowed aggregations to sink table
#     revenue_tbl.execute_insert('filter_data').wait()

#     tbl_env.execute('windowed-sales-euros')

# if __name__ == '__main__':
#     print("here")
#     main()