# Setup

<br>

## TODO

1. convert yaml to helm
2. remove Dapr
3. Test redis 2 
4. bash script for build + (dynamic)
5. clean up cmd.md
6. minikube vs k3s
7. postgres automate extension



## MacOS
#
<br>

### Step 1 : Setup Docker

- [Install](https://docs.docker.com/desktop/install/mac-install/)
- Run Docker-desktop (Double-click Docker.app in the Applications folder to start Docker)
- log in to docker hub
- check docker running
```ps
docker version
docker ps
```
<br>

### Step 2 : Install Minikube
- [Installation-Instruction](https://minikube.sigs.k8s.io/docs/start/)

```ps
# using Homebrew
brew install minikube

# using binaries
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-darwin-amd64
sudo install minikube-darwin-amd64 /usr/local/bin/minikube

# start minikube pod in docker
minikube start
sudo minikube tunnel
minikube status
kubectl get all
```
<br>

### Step 3 : Install Helm
[Installation-Instruction](https://helm.sh/docs/intro/install/)

```ps
# using homebrew
brew install helm

# verify installation
helm
helm version
helm repo add [chart-name]
helm repo list
helm list --all-namespaces
```
<br>

### Step 4 : Create Namespace

```bash
namespace=logfire-local
kubectl create namespace $namespace
kubectl get ns
```
<br>

### Step 5 : Install Dapr
```ps
# Add the official Dapr Helm chart.
helm repo add dapr https://dapr.github.io/helm-charts/

helm repo update

# Install the Dapr chart on your cluster in the dapr-system namespace.
helm upgrade --install dapr dapr/dapr \
--version=1.9 \
--namespace dapr-system \
--create-namespace \
--wait

# to verify installation
kubectl get all -n dapr-system

# See which chart versions are available
helm search repo dapr --devel --versions

# to Uninstall
helm uninstall dapr --namespace dapr-system
```

<br>

### Step 6 : Setup Ingress
- [Quick-start Link](https://kubernetes.github.io/ingress-nginx/deploy/#quick-start)

```ps
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update

# with dapr-annotations
helm upgrade --install api-gateway ingress-nginx/ingress-nginx --values ingress/dapr-annotations.yaml --namespace $namespace --wait

minikube service list

kubectl apply -f ingress/local-ingress.yaml
<!-- kubectl get ingress
kubectl describe ing webapp-ingress -->
```

<br>

### Step 7: Deploy Application
- switch to docker-daemon running inside minikube cluster
- create local image of application
- first apply secrets
- next apply application config manifesta
- next app deployment manifests

```ps 
minikube docker-env
eval $(minikube -p minikube docker-env)

# build docker images
docker build -t logfire/gowebapi:latest gowebapp/web-api
docker build -t logfire/notification:latest gowebapp/notification
docker build -t logfire/livetail:latest gowebapp/livetail
docker build -t logfire/filter-service:latest gowebapp/filter-service

# apply secrets
kubectl apply -f secrets -n $namespace

# deploy postgres
kubectl apply -f postgres -n $namespace
# DO some changes to postgres container
# to get postgres cli running inside its container
docker ps
# get your container-id of pod running postgresdb and replace here
kubectl exec -it postgres-deployment-0 sh /
psql -d logfire-dev -U postgresuser
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
# exit from the container


# deploy app
kubectl apply -f gowebapp -n $namespace

# Rebuild a container. 
# make sure you are in minikube docker daemon (eval $(minikube -p minikube docker-env))
kubectl delete -f gowebapp/goweb-api.yaml -n logfire-local
docker build -t logfire/gowebapi:latest gowebapp/web-api 
kubectl apply -f gowebapp/goweb-api.yaml -n logfire-local

kubectl delete -f gowebapp/livetail.yaml -n logfire-local
docker build -t logfire/livetail:latest gowebapp/livetail 
kubectl apply -f gowebapp/livetail.yaml -n logfire-local

kubectl delete -f gowebapp/filter-service.yaml -n logfire-local
docker build -t logfire/filter-service:latest gowebapp/filter-service 
kubectl apply -f gowebapp/filter-service.yaml -n logfire-local

<!-- DONE -->
```
#

## Setup Kafka

### Using Strimzi
<!-- kubectl create -f 'https://strimzi.io/install/latest?namespace=kafka' -n kafka
kubectl get pod -n kafka --watch
# Apply the `Kafka` Cluster CR file
kubectl apply -f https://strimzi.io/examples/latest/kafka/kafka-persistent-single.yaml -n kafka 
# or -->
```ps
kubectl create namespace kafka
# comment out single node or multi node cluster yaml
kubectl apply -f kafka/strimzi-operator.yaml -n kafka --wait
kubectl apply -f kafka/kafka-persistent-single.yaml -n kafka --wait


#With the cluster running, run a simple producer to send messages to a Kafka topic (the topic is automatically created):
kubectl -n kafka run kafka-producer -ti --image=quay.io/strimzi/kafka:0.33.1-kafka-3.3.2 --rm=true --restart=Never -- bin/kafka-console-producer.sh --bootstrap-server my-cluster-kafka-bootstrap:9092 --topic my-topic

# And to receive them in a different terminal, run:
kubectl -n kafka run kafka-consumer -ti --image=quay.io/strimzi/kafka:0.33.1-kafka-3.3.2 --rm=true --restart=Never -- bin/kafka-console-consumer.sh --bootstrap-server my-cluster-kafka-bootstrap:9092 --topic my-topic --from-beginning

kubectl -n kafka run kafka-consumer -ti --image=quay.io/strimzi/kafka:0.33.1-kafka-3.3.2 --rm=true --restart=Never -- bin/kafka-console-consumer.sh --bootstrap-server my-cluster-kafka-bootstrap:9092 --topic filter_topic_21a67b7e-3057-45cc-8da1-cd4521fc7530 --from-beginning

# Deleting your Apache Kafka cluster
kubectl -n kafka delete $(kubectl get strimzi -o name -n kafka)
# This will remove all Strimzi custom resources, including the Apache Kafka cluster and any KafkaTopic custom resources but leave the Strimzi cluster operator running so that it can respond to new Kafka custom resources.

# Deleting the Strimzi cluster operator
kubectl -n kafka delete -f kafka
```
<!-- kubectl -n kafka delete -f 'https://strimzi.io/install/latest?namespace=kafka' -->


#

## Setup ElasticSearch


<!-- helm repo add elastic https://helm.elastic.co
helm repo update
helm install elasticsearch --version 7.17.3 elastic/elasticsearch \
--namespace elasticsearch \
--create-namespace \
--wait -->

<!-- kubectl create -f https://download.elastic.co/downloads/eck/2.6.1/crds.yaml
kubectl apply -f https://download.elastic.co/downloads/eck/2.6.1/operator.yaml
kubectl apply -f elastic-search/elasticsearch.yaml -n logfire-local -->

```ps
kubectl apply -f elastic-search
kubectl get pods --selector='elasticsearch.k8s.elastic.co/cluster-name=my-elasticsearch' -n logfire-local
# to test connection
kubectl port-forward -n logfire-local service/my-elasticsearch-es-http 9200
PASSWORD=$(kubectl get secret -n logfire-local my-elasticsearch-es-elastic-user -o go-template='{{.data.elastic | base64decode}}')
curl -u "elastic:$PASSWORD" -k "https://localhost:9200"

```
#

## Setup Flink

```ps
kubectl apply -f flink -n logfire-local --wait

kubectl port-forward -n logfire-local  --address 0.0.0.0 service/flink-jobmanager 8081:8081
<!-- Forwarding from 0.0.0.0:8081 -8081 -->
```

## Setup Flink Using Operator

```ps
# Deploy cert manager and wait for it to become active
kubectl create -f https://github.com/jetstack/cert-manager/releases/download/v1.8.2/cert-manager.yaml

helm repo add flink-operator-repo https://downloads.apache.org/flink/flink-kubernetes-operator-1.4.0/
helm repo update

helm install flink-kubernetes-operator flink-operator-repo/flink-kubernetes-operator
```

## Setup Filter Service 2

```ps
# Install sbt by following the link
# https://www.scala-sbt.org/download.html

(cd scala && sbt assemblyPackageDependency)

docker build -t logfire/flink -f scala/Dockerfile scala/

kubectl apply -f scala/flink-session-cluster.yaml
```


#
#

 # minikube installation
https://minikube.sigs.k8s.io/docs/start/
```ps
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64

sudo install minikube-linux-amd64 /usr/local/bin/minikube

minikube start

# we get error due to access denied to docker.  so we create a usergroup docker and use it

sudo usermod -aG docker $USER && newgrp docker

minikube start --driver=docker

# done
```
# Kubectl commands to interact with cluster

```ps
# to check control-plane or master node
kubectl get node

# to get pods running
kubectl get pod

none running at the moment. So we need to create them.
#But before that we need our ConfigMap and Secret to exist. As these are used by deployments to create pods
# So we create our manifests and apply them o k8s cluster
# to apply configMap

<!-- UPDATE access_requests SET approved = 'true' WHERE id = '7acd52cb-0928-45cb-a5dd-1a07efd7876c';
7acd52cb-0928-45cb-a5dd-1a07efd7876c -->


kubectl apply -f postgres-config.yaml
kubectl apply -f livetail-config.yaml
kubectl apply -f postgres-pvc-pv.yaml
kubectl apply -f webapp-ingress.yaml

# to apply secrets
kubectl apply -f postgres-secret.yaml
kubectl apply -f aws-secret.yaml
kubectl apply -f google-secret.yaml
kubectl apply -f docker-secret.yaml

#  we first start our database pod as web-api depends on it
kubectl apply -f postgres.yaml

# finally we deploy our web-application
kubectl apply -f goweb-api.yaml

# now to get all the pods and services...
kubectl get all

# we don't see configMap and secret here. to get them we need separate command
kubectl get configmap
kubectl get secret

# to get kube services
kubectl get svc

# to get persistent volume and persistent volume claim
kubectl get pv
kubectl get pvc

NAME               TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)          AGE
gowebapi-service   NodePort    10.106.67.30    <none       8080:30100/TCP   6h4m
kubernetes         ClusterIP   10.96.0.1       <none       443/TCP          9h
postgres-service   ClusterIP   10.98.198.120   <none       5432/TCP         6h20m

    # here we have gowebapi-service of type "NodePort" which is accessible externally at port "30100"
    # but we need minicube ip address, which is running this kubernetes cluster

# to get cluster's ip address
minikube ip
# or
kubectl get node
# or, for more details
kubectl get node -o wide

NAME       STATUS   ROLES           AGE   VERSION   INTERNAL-IP    EXTERNAL-IP   OS-IMAGE             KERNEL-VERSION      CONTAINER-RUNTIME
minikube   Ready    control-plane   9h    v1.25.3   192.168.49.2   <none       Ubuntu 20.04.5 LTS   5.15.0-46-generic   docker://20.10.20


# to open browser dashboard
minikube dashboard
```

- in secrets we need to put base64 encoded values
```ps
# to encode
echo -n postgresuser | base64
# to decode
echo 'bGludXhoaW50LmNvbQo=' | base64 --decode
```
- nodeport: to access externally (between 30000-32767)



<br>

# To setup ingress

- We can use the default nginx-ingress-controller provided by minikube by enabling it(if want to use default settings provided by minikube)

minikube addons enable ingress (don't use this)
- automatically starts the k8s nginx implementation of ingress controller

Since K8s doesn't provide LoadBalancer, for bare-metal setup we will be using MetalLb.

https://kubernetes.github.io/ingress-nginx/deploy/baremetal/

<br>

## Helm Setup

[Install](https://helm.sh/docs/intro/install/)

[Installation](https://devopscube.com/install-configure-helm-kubernetes/)


```ps
# Download your desired version, to home directory and clean after installation
wget -O helm.tar.gz https://get.helm.sh/helm-v3.10.3-darwin-amd64.tar.gz

# Unpack it 
tar -zxvf helm-v3.10.3-darwin-amd64.tar.gz

# Find the helm binary in the unpacked directory, and move it to its desired destination 
mv linux-amd64/helm /usr/local/bin/helm

# verify installation
helm
helm version
helm repo add stable https://charts.helm.sh/stable
helm repo list
helm list --all-namespaces
```

<br>

## MetalLb Setup(for bare-metal servers)

(no need if using minikube, provides as addon)

https://metallb.universe.tf/installation/

https://metallb.universe.tf/configuration/

- we will also need Helm to easily install nginx-controller & to configure it as type-`LoadBalancer` to be used with MetalLb(can't use default one of minikube, as it configures its' as type-`NodePort`)

<br>

## Minikube + Metallb + Nginx-ingress(for bare-metal servers)

```ps
# start minikube
minikube status

# if minikube not running
minikube start

kubectl get nodes
minikube service list
minikube addons list

# if not enabled metallb
minikube addons enable metallb
minikube addons list
minikube ip
# a.b.c.d
# to be used to assign ip range to metallb

# to verify metal-lb is running
kubectl get ns
minikube addons configure metallb
# provide range of ip which metal-lb will own and assign to LoadBalancer services in the k8s cluster automatically
# -- Enter Load Balancer Start IP: <minikube ipor a.b.c.240
# -- Enter Load Balancer End IP: <minikube ipor a.b.c.250
kubectl get all -n metallb-system
kubectl get configmap/config -n metallb-system -o yaml

# to test metallb installation working poperly
kubectl create deploy nginx --image nginx
kubectl expose deploy nginx --port 80 --type LoadBalancer
kubectl get all
curl $INGRESS_EXTERNAL_IP
kubectl delete deployment nginx
kubectl delete svc nginx

```

<br>

## Nginx-ingress Setup
#

[Quick-start Link](https://kubernetes.github.io/ingress-nginx/deploy/#quick-start)

```ps
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update

helm upgrade --install api-gateway ingress-nginx/ingress-nginx --values ingress/dapr-annotations.yaml --namespace logfire-local --wait 

# or
    ```bash
    namespace=logfire-staging
    kubectl create namespace $namespace
    ```
helm upgrade --install api-gateway ingress-nginx/ingress-nginx --values ingress/dapr-annotations.yaml --namespace $namespace


minikube service -p logfire-staging list
helm list -n ingress-nginx


kubectl apply -f ingress/webapp-ingress.yaml
kubectl get ingress
kubectl describe ing webapp-ingress
```

<br>


# DAPR

## On Kubernetes using helm-chart

```ps
# Add the official Dapr Helm chart.
helm repo add dapr https://dapr.github.io/helm-charts/

helm repo update

# See which chart versions are available
helm search repo dapr --devel --versions

# Install the Dapr chart on your cluster in the dapr-system namespace.
helm upgrade --install dapr dapr/dapr \
--version=1.9 \
--namespace dapr-system \
--create-namespace \
--wait

# to verify installation
kubectl get all -n dapr-system

# to Uninstall
helm uninstall dapr --namespace dapr-system
```
<br>


## Self-Hosted(on local machine)
#

```ps
# first install dapr on local machine
wget -q https://raw.githubusercontent.com/dapr/cli/master/install/install.sh -O - | /bin/bash
```

### On Kubernetes:
#
```ps
#Make sure the correct “target” cluster is set. Check kubectl context to verify.
kubectl config get-contexts
 
#You can set a different context using 
kubectl config use-context <CONTEXT>

# to initialise:
dapr init -k

# to verify:
kubectl get all -n dapr-system

# check dashboard on local browser(port-forward)
kubectl port-forward service/dapr-dashboard 8080:8080 -n dapr-system

# to unistall
dapr uninstall -k
```
<br>

### Initialisation on local environment
```ps
dapr init

#Dapr runtime installed to /home/altrova/.dapr/bin, you may run the following to add it to your path if you want to run daprd directly:

export PATH=$PATH:/home/altrova/.dapr/bin

dapr --version

CLI version: 1.9.1 
Runtime version: 1.9.5

#config stored at:
ls $HOME/.dapr

# to unistall(clean)
dapr uninstall --all
```

<br>

# KAFKA

```ps
kubectl apply -f kafka -n kafka

# go into broker container
kubectl exec broker-deployment-66cdffdbb6-d8r2r -it -n kafka sh

# to start a kafka-topic run this from inside broker
kafka-topics --create --bootstrap-server localhost:29092 --replication-factor 1 --partitions 1 --topic minikube-topic

# to enter into consumer mode execute this in the broker container console
kafka-console-consumer --bootstrap-server localhost:29092 --topic minikube-topic


# testing connection
kubectl run python-repl --rm -i --tty --image python:3.9 sh
pip install confluent-kafka
python
from confluent_kafka import Producer
import socket
conf = {"bootstrap.servers": "kafka-service.kafka.svc.cluster.local:9092", "client.id": socket.gethostname()}
conf = {"bootstrap.servers": "kafka-service:9092", "client.id": socket.gethostname()}
producer = Producer(conf)
producer.produce("minikube-topic", key="message", value="message_from_python_producer")



```
<br>

## Using Strimzi
```ps
kubectl create namespace kafka
kubectl create -f 'https://strimzi.io/install/latest?namespace=kafka' -n kafka
kubectl get pod -n kafka --watch
# Apply the `Kafka` Cluster CR file
kubectl apply -f https://strimzi.io/examples/latest/kafka/kafka-persistent-single.yaml -n kafka 
# or
kubectl apply -f kafka/kafka-persistent-single.yaml -n kafka

#With the cluster running, run a simple producer to send messages to a Kafka topic (the topic is automatically created):
kubectl -n kafka run kafka-producer -ti --image=quay.io/strimzi/kafka:0.33.1-kafka-3.3.2 --rm=true --restart=Never -- bin/kafka-console-producer.sh --bootstrap-server my-cluster-kafka-bootstrap:9092 --topic my-topic

# And to receive them in a different terminal, run:
kubectl -n kafka run kafka-consumer -ti --image=quay.io/strimzi/kafka:0.33.1-kafka-3.3.2 --rm=true --restart=Never -- bin/kafka-console-consumer.sh --bootstrap-server my-cluster-kafka-bootstrap:9092 --topic my-topic --from-beginning

# Deleting your Apache Kafka cluster
kubectl -n kafka delete $(kubectl get strimzi -o name -n kafka)
# This will remove all Strimzi custom resources, including the Apache Kafka cluster and any KafkaTopic custom resources but leave the Strimzi cluster operator running so that it can respond to new Kafka custom resources.

# Deleting the Strimzi cluster operator
kubectl -n kafka delete -f 'https://strimzi.io/install/latest?namespace=kafka'

```



### https://docs.bitnami.com/tutorials/deploy-scalable-kafka-zookeeper-cluster-kubernetes/
## Step 1: Deploy Apache Zookeeper

``` ps 

# deploy an Apache Zookeeper cluster with three node.
# will not be exposed publicly, it is deployed with authentication disabled.
helm install zookeeper bitnami/zookeeper \
    --set replicaCount=3 \
    --set auth.enabled=false \
    --set allowAnonymousLogin=true \
    --namespace kafka \
    --create-namespace \
    --wait

ZooKeeper can be accessed via port 2181 on the following DNS name from within your cluster:

    zookeeper.kafka.svc.cluster.local

# To connect to your ZooKeeper server run the following commands:
export POD_NAME=$(kubectl get pods --namespace kafka -l "app.kubernetes.io/name=zookeeper,app.kubernetes.io/instance=zookeeper,app.kubernetes.io/component=zookeeper" -o jsonpath="{.items[0].metadata.name}")
kubectl -n kafka exec -it $POD_NAME -- zkCli.sh

# To connect to your ZooKeeper server from outside the cluster execute the following commands:
  kubectl port-forward --namespace kafka svc/zookeeper 2181:2181 &
    zkCli.sh 127.0.0.1:2181
```
    NOTES:
        CHART NAME: zookeeper
        CHART VERSION: 11.0.3
        APP VERSION: 3.8.0


## Step 2: Deploy Apache Kafka
```ps
# In this case, you will provide the name of the Apache Zookeeper service as a parameter to the Helm chart. Apache Zookeeper service name obtained at the end of Step 1.
helm install kafka bitnami/kafka \
  --set zookeeper.enabled=false \
  --set replicaCount=3 \
  --set externalZookeeper.servers=zookeeper \
  --namespace kafka \
  --wait

# Kafka can be accessed by consumers via port 9092 on the following DNS name from within your cluster:

    kafka.kafka.svc.cluster.local

# Each Kafka broker can be accessed by producers via port 9092 on the following DNS name(s) from within your cluster:

    kafka-0.kafka-headless.kafka.svc.cluster.local:9092
    kafka-1.kafka-headless.kafka.svc.cluster.local:9092
    kafka-2.kafka-headless.kafka.svc.cluster.local:9092

# To create a pod that you can use as a Kafka client run the following commands:

kubectl run kafka-client --restart='Never' --image docker.io/bitnami/kafka:3.3.1-debian-11-r34 --namespace kafka --command -- sleep infinity
kubectl exec --tty -i kafka-client --namespace kafka -- bash

#    PRODUCER:
      kafka-console-producer.sh \
            --broker-list kafka-0.kafka-headless.kafka.svc.cluster.local:9092,kafka-1.kafka-headless.kafka.svc.cluster.local:9092,kafka-2.kafka-headless.kafka.svc.cluster.local:9092 \
            --topic test

#    CONSUMER:
      kafka-console-consumer.sh \
            --bootstrap-server kafka.kafka.svc.cluster.local:9092 \
            --topic test \
            --from-beginning

```
    NOTES:
        CHART NAME: kafka
        CHART VERSION: 20.0.4
        APP VERSION: 3.3.1

<br>

# Step 3: Test Apache Kafka

```ps

kubectl get pods --namespace kafka -l "app.kubernetes.io/name=kafka,app.kubernetes.io/instance=kafka,app.kubernetes.io/component=kafka" -o jsonpath="{.items[0].metadata.name}"

# Create a topic named mytopic using the commands below
kubectl --namespace kafka exec -it kafka-0 -- kafka-topics.sh --create --zookeeper zookeeper:2181 --replication-factor 1 --partitions 1 --topic mytopic

export POD_NAME=$(kubectl get pods --namespace kafka -l "app.kubernetes.io/name=kafka,app.kubernetes.io/instance=kafka,app.kubernetes.io/component=kafka" -o jsonpath="{.items[0].metadata.name}")

kubectl --namespace kafka exec -it $POD_NAME -- kafka-topics.sh --create --zookeeper ZOOKEEPER-SERVICE-NAME:2181 --replication-factor 1 --partitions 1 --topic mytopic
```
# Optional cert-manager

## Install cert-manager
#
[Installation](https://cert-manager.io/docs/installation/)

[Tutorial](https://cert-manager.io/docs/tutorials/acme/nginx-ingress/)

[Certificate](https://cert-manager.io/docs/concepts/certificate/)

```ps
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.10.1/cert-manager.yaml
kubectl get ns
kubectl -n cert-manager get all

# to uninstall
kubectl delete -f https://github.com/cert-manager/cert-manager/releases/download/v1.10.1/cert-manager.yaml

# to enable cert-issuer
kubectl apply -f cert-issuer/staging-issuer.yaml

#  to use tls-enabled ingress
kubectl apply -f ingress/staging-ingress.yaml
kubectl describe ingress webapp-ingress
curl apibeta.logfire.sh

# we get error for redirect. to follow redirect
curl -L apibeta.logfire.sh

# we get private certificate error. to ignore
curl -Lk apibeta.logfire.sh
```

## Get ssl-certificate
#

```ps
sudo apt-get -y install certbot python3-certbot-apache
sudo certbot certonly

# copy to other folder change permissions create secret and then remove it from this folder
#sudo chmod 777 priv.key
#sudo chmod 777 crt.key
#ksudo ubectl create secret tls apibeta-logfire-tls-secret --key priv.key --cert crt.key

# Port-forward 80 and 443 to listen to host network from ingress

kubectl port-forward service/ingress-nginx-controller -n ingress-nginx 80:80
error:
Unable to listen on port 80: Listeners failed to create with the following errors: [unable to create listener: Error listen tcp4 127.0.0.1:80: bind: permission denied unable to create listener: Error listen tcp6 [::1]:80: bind: permission denied]
error: unable to listen on any of the requested ports: [{80 80}]

# to resolve (a dirty trick as it opens a terminal which needs to remain open)

which kubectl
sudo setcap CAP_NET_BIND_SERVICE=+eip /usr/local/bin/kubectl
kubectl port-forward -n ingress-nginx  --address 0.0.0.0 service/ingress-nginx-controller 80:80 443:443
kubectl port-forward -n logfire-local  --address 0.0.0.0 service/api-gateway-ingress-nginx-controller 80:80 443:443

```

<br>


<br>

# DEBUG commands

```ps
kubectl get all -n logfire-staging

# get logs from pods in consol
kubectl logs pod/postgres-deployment-5558648d6-hsgjr -n logfire-staging

# describe pods
kubectl describe pod/testgo-866c77449d-5jkj7 -n logfire-staging

# get dapr logs from a pod
kubectl logs testgo-866c77449d-5jkj7 -c daprd -n logfire-staging


minikube docker-env
eval $(minikube -p minikube docker-env)
docker build -t logfire/gowebapi:latest gowebapp/web-api
docker build -t logfire/notification:latest gowebapp/notification
docker build -t logfire/livetail:latest gowebapp/livetail
docker logs 947ba8a3c25d -f
# to get postgres cli running inside its container
docker exec -it a6936f62e455 psql -d logfire-dev -U postgresuser
select * 
from pg_stat_activity
where (state = 'idle in transaction')
    and xact_start is not null;
kubectl describe nodes
kubectl delete -f gowebapp/livetail.yaml -n logfire-local

#StarRocks Configuration
Reference Document : https://github.com/StarRocks/starrocks-kubernetes-operator
#Initializing starrock 
kubectl apply -f https://raw.githubusercontent.com/StarRocks/starrocks-kubernetes-operator/main/deploy/starrocks.com_starrocksclusters.yaml
kubectl apply -f https://raw.githubusercontent.com/StarRocks/starrocks-kubernetes-operator/main/deploy/operator.yaml
kubectl apply -f starrocks-fe-and-be.yaml
kubectl get svc -n starrocks
#deleting starrocks
kubectl delete -f starrocks-fe-and-be.yaml
kubectl delete -f  https://raw.githubusercontent.com/StarRocks/starrocks-kubernetes-operator/main/deploy/operator.yaml



```