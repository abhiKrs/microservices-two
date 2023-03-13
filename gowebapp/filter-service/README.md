

```ps
source .venv/bin/activate
pip install -r requirements.txt


docker exec -it broker kafka-topics --create \
    --bootstrap-server localhost:9092 \
    --topic logs-topic

docker exec -it broker kafka-console-consumer --from-beginning \
    --bootstrap-server localhost:9092 \
    --topic logs-topic

docker exec broker kafka-topics --list --zookeeper zookeeper:2181
docker exec broker kafka-topics --delete --zookeeper zookeeper:2181 --topic logs-topic
```