# API Challenge

Postman collection located at the repo root [api.postman_collection.json]

## Requirements

- docker
- docker-compose
- PostgreSQL
- golang:1.14.4

## How to setup DB for development

```make dep``` - load dependencies
```make dc-up``` - start empty PostgreSQL database  
```make dc-show``` - check PostgresSQL status\
```migrate-macos-install``` OR ```migrate-linux-install``` - install migrate
```make migrate-up``` - roll-up migrations  
```make run``` - start service

## TODO

- Tokens blacklisting
