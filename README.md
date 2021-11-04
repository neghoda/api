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
```make migrate-up``` - roll-up migrations\
```docker build -t app .``` - build image\
```docker container run -p 8080:8080 --network="host" app``` - start service

## ENV variables used for configuration

```REFRESH_TOKEN_LEN```    - default 32 bytes\
```ACCESS_TOKEN_TTL_SEC``` - default 90 sec\
```ACCESS_TOKEN_SECRET```  - default ""\
```USER_SESSION_TTL_SEC``` - default 86400 sec\

```CRON_CONFIG``` - when to sync fund data. default "46 16 * * *"(every day in 4:46PM)\

```HTTP_SERVER_PORT```       - default :8080\
```HTTP_SERVER_URL_PREFIX``` - default "/api"\
```HTTP_SWAGGER_ENABLE```    - enable serving of swagger json and yml files \
```HTTP_SWAGGER_SERVE_DIR``` - location of swagger files. default "src/server/http/static"\
```HTTP_CORS_ALLOWED_HOST``` - CORS host settings. default "*"\

```POSTGRES_MASTER_HOST```     - default localhost\
```POSTGRES_MASTER_NAME```     - default postgres\
```POSTGRES_MASTER_USER```     - default postgres\
```POSTGRES_MASTER_PASSWORD``` - default 12345\
```POSTGRES_POOL_SIZE```       - default 10\
```POSTGRES_MAX_RETRIES```     - default 5\
```POSTGRES_READ_TIMEOUT```    - default 10 sec\
```POSTGRES_WRITE_TIMEOUT```   - default 10 sec\

## TODO

- Token blacklisting
