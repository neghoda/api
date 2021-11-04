# API Challenge

Postman collection located at the repo root [api.postman_collection.json]

## Requirements

- docker
- docker-compose
- PostgreSQL
- golang:1.17

## How to setup API for development

```make dep``` - load dependencies\
```make dc-up``` - start empty PostgreSQL database  
```make dc-show``` - check PostgresSQL status\
```make migrate-macos-install``` OR ```make migrate-linux-install``` - install migrate\
```make migrate-up``` - roll-up migrations  
```make run``` - start service

## How to setup API for development in Docker

```make dc-up``` - start empty PostgreSQL database  
```make dc-show``` - check PostgresSQL status\
```make migrate-macos-install``` OR ```make migrate-linux-install``` - install migrate\
```make migrate-up``` - roll-up migrations  
```docker build -t app .``` - build image
```docker container run -p 8080:8080 --network="host" app``` - start service

## TODO

- Token blacklisting
