FROM golang:1.17 AS builder
WORKDIR /bin/api
COPY . .
RUN CGO_ENABLED=1 go build -mod=mod -o ./bin/api -a ./src

FROM alpine
COPY --from=builder ./bin/api /
COPY ./src/server/http/static /swaggerui
CMD ["/api"]
