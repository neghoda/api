FROM alpine
COPY ./bin/api /
COPY ./src/server/http/static /swaggerui
CMD ["/api"]
