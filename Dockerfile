FROM golang:1.17 AS builder
WORKDIR /bin
COPY . .
RUN CGO_ENABLED=0 go build -mod=mod -o api -a ./src

FROM alpine
COPY --from=builder ./bin/api /
COPY ./src/server/http/static ./src/server/http/static
CMD ["/api"]
