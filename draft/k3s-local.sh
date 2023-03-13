sudo /usr/local/bin/k3s-uninstall.sh
curl -sfL https://get.k3s.io | sh -s - --docker
sudo chmod 777 /etc/rancher/k3s/k3s.yaml
cp -r /etc/rancher/k3s/k3s.yaml ~/.kube/config
kubectl get nodes
kubectl create namespace logfire-local
kubectl get ns
helm upgrade --install api-gateway ingress-nginx/ingress-nginx --namespace logfire-local  
kubectl apply -f secrets -n logfire-local
kubectl apply -f postgres -n logfire-local  
docker build -t logfire/gowebapi:latest gowebapp/web-api
docker build -t logfire/notification:latest gowebapp/notification
docker build -t logfire/livetail:latest gowebapp/livetail
docker build -t logfire/filter-service:latest gowebapp/filter-service
kubectl apply -f deploy-local/operators -n logfire-local --wait
kubectl apply -f deploy-local/deployment -n logfire-local --wait
sleep 5m
# Execute a shell in the container using kubectl
kubectl exec -it postgres-deployment-0 -n logfire-local sh <<EOF

# Run psql in the container with the appropriate parameters and execute the SQL command
psql -d logfire-dev -U postgresuser -c "CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";"



