import argparse
import json
import logging
import time
import sys
# import requests

from confluent_kafka import Producer


logging.basicConfig(
  format='%(asctime)s %(name)-12s %(levelname)-8s %(message)s',
  datefmt='%Y-%m-%d %H:%M:%S',
  level=logging.INFO,
  handlers=[
      logging.FileHandler("logs_producer.log"),
      logging.StreamHandler(sys.stdout)
  ]
)

logger = logging.getLogger()



class ProducerCallback:
    def __init__(self, record, log_success=False):
        self.record = record
        self.log_success = log_success

    def __call__(self, err, msg):
        if err:
            logger.error('Error producing record {}'.format(self.record))
        elif self.log_success:
            logger.info('Produced {} to topic {} partition {} offset {}'.format(
                self.record,
                msg.topic(),
                msg.partition(),
                msg.offset()
            ))


def main(args):
    logger.info('Starting logs producer')
    conf = {
        'bootstrap.servers': args.bootstrap_server,
        'linger.ms': 200,
        'client.id': 'logs-1',
        'partitioner': 'murmur2_random'
    }

    producer = Producer(conf)

    while True:
        with open('logs.json', 'r') as f:
            data = json.load(f)

        

        for logs in data["data"]:
            producer.produce(topic=args.topic,
            value=json.dumps(logs),
            on_delivery=ProducerCallback(logs, log_success=True))
        #     x = requests.post('http://localhost/api/logfire.sh', json=logs)
        #     print(x.status_code)
            time.sleep(0.01)

        # time.sleep(5)


if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('--bootstrap-server', default='localhost:9092')
    parser.add_argument('--topic', default='logs-topic')
    args = parser.parse_args()
    main(args)