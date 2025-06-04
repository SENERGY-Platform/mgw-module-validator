FROM golang:1.24 AS builder

ARG VERSION=dev

COPY . /go/src/app
WORKDIR /go/src/app

RUN CGO_ENABLED=0 GOOS=linux go build -o bin -ldflags="-X 'main.version=$VERSION'" main.go

FROM alpine:3.22

RUN mkdir -p /opt/module-validator
RUN mkdir -p /mnt/data
WORKDIR /opt/module-validator
COPY --from=builder /go/src/app/bin bin

ENTRYPOINT ["./bin", "-b", "/mnt/data"]
