# go build
FROM golang:alpine AS builder
ENV CGO_ENABLED=0
COPY . /go/src/go-gin-jwt
RUN cd /go/src/go-gin-jwt && go build

# pack up
FROM scratch
COPY --from=builder /go/src/go-gin-jwt/go-gin-jwt /
ENTRYPOINT ["/go-gin-jwt"]

ENV DBUSER root
ENV DBPWD   xiaomingdemima
ENV DBHOST  10.100.0.2:3306
ENV DBNAME  go_jwt