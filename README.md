# microservices
```ps
## Local-Installation

##Installing k3s for develoment using script 
sh scripts/install.sh
##for MAC to start cluster
k3d cluster create cluster_name
##for Mac to delete cluster
k3d cluster delete vluster_name

##Installing k3s for development manually
curl -sfL https://get.k3s.io | sh -s - --docker
sudo chmod 777 /etc/rancher/k3s/k3s.yaml
cp -r /etc/rancher/k3s/k3s.yaml ~/.kube/config (or)
add the below line in ~/.bashrc
export KUBECONFIG=/etc/rancher/k3s/k3s.yaml
#check nodes in k3
kubectl get nodes
kubectl get ns
############################


#Installing application using script : 
sh scripts/deploy.sh

#Installing application manually:
kubectl create ns logfire-local

#Apply cofigmap and secrets:
kubectl apply -f configmap/ -n logfire-local

#Apply operators for applications using helm:
#get into helm-chart directory
helm repo add strimzi https://strimzi.io/charts/
helm install my-strimzi-operator strimzi/strimzi-kafka-operator --namespace logfire-local --wait
helm install crds helm-chart/operators/crds/ -n logfire-local
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update
helm upgrade --install api-gateway ingress-nginx/ingress-nginx --namespace logfire-local --wait
helm install es-operator operators/es-operator/ -n logfire-local --wait

#Install application  in development
#get into helm-chart directory
helm install postgres psql-chart/ -n logfire-local
helm install flink flink-chart/ -n logfire-local
helm install kafka kafka-chart/ -n logfire-local
helm install webapp webapp-devchart/ -n logfire-local
helm install ingress-local ingress-localchart/ -n logfire-local
```

## Production Deployment
```ps
##Installing application using script : 
sh scripts/deploy.sh

##Installing application manually:
kubectl create ns logfire-local

#Apply cofigmap and secrets:
kubectl apply -f configmap/ -n logfire-local

#Apply operators for applications using helm:
#get into helm-chart directory
helm repo add strimzi https://strimzi.io/charts/
helm install my-strimzi-operator strimzi/strimzi-kafka-operator --namespace logfire-local --wait
helm install crds helm-chart/operators/crds/ -n logfire-local
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update
helm upgrade --install api-gateway ingress-nginx/ingress-nginx --namespace logfire-local --wait
helm install es-operator operators/es-operator/ -n logfire-local --wait

#Install application  in development
#get into helm-chart directory
helm install postgres postgres-prodchart/ -n logfire-local
helm install flink flink-chart/ -n logfire-local
helm install kafka kafka-chart/ -n logfire-local
helm install webapp webapp-prodchart/ -n logfire-local
helm install prod-ingress ingress-prodchart/ -n logfire-local
###Ending application install###############
```
