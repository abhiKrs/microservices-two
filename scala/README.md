
## Setup Flink Using Operator

```ps
# Deploy cert manager and wait for it to become active
kubectl create -f https://github.com/jetstack/cert-manager/releases/download/v1.8.2/cert-manager.yaml

helm repo add flink-operator-repo https://downloads.apache.org/flink/flink-kubernetes-operator-1.4.0/
helm repo update

helm install flink-kubernetes-operator flink-operator-repo/flink-kubernetes-operator --namespace logfire-local
```

## Setup Filter Cluster

```ps
# Install sbt by following the link
# https://www.scala-sbt.org/download.html

(cd scala && sbt assemblyPackageDependency)
(cd scala && sbt assembly)

docker build -t logfire/flink -f scala/flink.Dockerfile scala/

# For local building
docker build -t logfire/flink -f scala/flink.dev.Dockerfile scala/

kubectl apply -f scala/flink-session-cluster.yaml -n logfire-local
```

## Setup schema registry

```
helm upgrade --install cp-schema-registry kafka/cp-schema-registry --namespace logfire-local --set kafka.external=true --set kafka.address=my-cluster-kafka-brokers:9092
```

## Setup flink service
```
docker build -t logfire/filter-service -f scala/Dockerfile scala/

# For local building
docker build -t logfire/filter-service -f scala/filterService.dev.Dockerfile scala/


helm upgrade --install filter-service scala/helm/filter-service --namespace logfire-local --set redis.host=redis-service.redis.svc.cluster.local

# SETUP FOR ROHAN ONLY
helm upgrade --install filter-service scala/helm/filter-service --values scala/helm/filter-service/values-rohan.yaml
```

## Setup redis
```
helm upgrade --install redis bitnami/redis --set auth.enabled=false
#redis-master
export REDIS_PASSWORD=$(kubectl get secret --namespace izac redis -o jsonpath="{.data.redis-password}" | base64 -d)
REDISCLI_AUTH="$REDIS_PASSWORD" redis-cli -h 127.0.0.1 -p 6379
```

## Run test cases
```
sbt "project filterJob" test
```