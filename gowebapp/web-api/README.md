# web-api




# Scripts to Run



## To dockerise and run the web-api in development 

<br>

- ###  With logs
```
docker compose up
```

<br>

- ### OR, in detach mod - no access to logs from container

```
docker compose up -d
```

access at 
> [https:localhost](http://localhost/)
>
<br>

- ### To run the whole backend-stack

```
docker compose -f backend.yaml up -d
```

### Done..


<br>


## To run web-api locally on your machine outside docker
# 




first run postgresdb inside docker
```
docker compose -f postgres-admin.yaml up -d
```
then to run web-api
```
go run main.go
```
- [https:localhost:8080](http://localhost:8080/)

<br>


## To dockerise and run the web-api in Production(not now) 
#

```
docker compose -f production.yaml up -d
```


<!-- for kafka-producer: 

```
docker exec --interactive --tty broker kafka-console-consumer --bootstrap-server broker:9092 --topic quickstart --from-beginning
```

for kafka-consumer: 

```
docker exec --interactive --tty broker kafka-console-producer --bootstrap-server broker:9092 --topic quickstart
```
<br>

## PostgresAdmin @

<br>

  > [https:localhost:5050](http://localhost:5050/)

[to configure:](https://hevodata.com/learn/docker-postgresql/)

user user, password values from docker compose -->
