eval $(minikube -p minikube docker-env)
docker build -t logfire/gowebapi:latest gowebapp/web-api
docker build -t logfire/notification:latest gowebapp/notification
docker build -t logfire/livetail:latest gowebapp/livetail
docker build -t logfire/filter-service:latest gowebapp/filter-service
