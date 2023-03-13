# (TODO)
> 1. convert yaml to helm
> 2. remove Dapr
> 3. Test redis 2 
> 4. bash script for build + (dynamic)
> 5. clean up cmd.md
> 6. minikube vs k3s
> 7. postgres automate extension
> 8. separate secret files
> 9. separate ingress rules
> 10. starrocks, postgres, kafka, redis, elasticsearch (separate kubernetes cluster?) - for stateful deployments
for production and staging 
> 11. Setup devspace
> 12. Github actions


<br>

# Local depoyment

- ## setup k3s
  - provide namespace in the current context

- ## build-images.sh
  - 
<!-- eval $(minikube -p minikube docker-env) -->
<!-- docker build -t logfire/gowebapi:latest gowebapp/web-api
docker build -t logfire/notification:latest gowebapp/notification
docker build -t logfire/livetail:latest gowebapp/livetail
docker build -t logfire/filter-service:latest gowebapp/filter-service -->

- ## add nginx-ingress-controller
```ps
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update
helm upgrade --install api-gateway ingress-nginx/ingress-nginx --namespace $namespace --wait
```
<br>

- ## Deploy Application
### TODO Replace with helm
```ps
helm install logfire-apps helm-deploy/local-deploy --namespace $namespace --wait
helm install -f localvalues.yaml logfirecharts ./helm-deploy --namespace $namespace --wait
```
<!-- kubectl apply -f deploy-local/operators -n $namespace --wait
kubectl apply -f deploy-local/deployment -n $namespace --wait -->

# Staging





# Production
