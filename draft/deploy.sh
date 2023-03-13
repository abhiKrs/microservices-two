#! /bin/bash
until read -r -p "Enter the Environment (prod or dev) " environment && test "$environment" != ""; do
  continue
done
#until read -r -p "Enter the path of the helm-chart folder in microservices (ex : /op/microservices/helm-chart) " path && test "$path" != ""; do
#  continue
#done
#echo $path
#"$configmap" = "$path/configmap/"
#echo $configmap
if [ $environment = "dev" ]
then
echo "$environment"
echo "starting development operator installation"
kubectl create ns logfire-local
pwd
kubectl apply -f configmap/ -n logfire-local
#sudo chmod -R o+x /opt/microservices/
#cd helm-chart
helm repo add strimzi https://strimzi.io/charts/
helm install my-strimzi-operator strimzi/strimzi-kafka-operator --namespace logfire-local --wait
helm install crds operators/crds/ -n logfire-local
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update
helm upgrade --install api-gateway ingress-nginx/ingress-nginx --namespace logfire-local --wait
helm install es-operator operators/es-operator/ -n logfire-local --wait

##################Installing application########3#
echo "starting development application installation"
helm install postgres psql-chart/ -n logfire-local
helm install flink flink-chart/ -n logfire-local
helm install kafka kafka-chart/ -n logfire-local
helm install webapp webapp-chart/ -n logfire-local
echo "development application installation completed"
exit 1
else
#############Installing production operator #############
echo "starting production operator installation"
kubectl create ns logfire-local
kubectl apply -f configmap/ -n logfire-local
cd /helm-chart
helm repo add strimzi https://strimzi.io/charts/
helm install my-strimzi-operator strimzi/strimzi-kafka-operator --namespace logfire-local --wait
helm install crds operators/crds/ -n logfire-local
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update
helm upgrade --install api-gateway ingress-nginx/ingress-nginx --namespace logfire-local --wait
helm install es-operator operators/es-operator/ -n logfire-local --wait

##########Installing Production application ###############
echo "starting production application installation"
helm install postgres postgres-prodchart/ -n logfire-local
helm install flink flink-chart/ -n logfire-local
helm install kafka kafka-chart/ -n logfire-local
helm install webapp webapp-chart/ -n logfire-local
helm install prod-ingress ingress-prodchart/ -n logfire-local
echo "production application installation completed"
exit 1
fi
