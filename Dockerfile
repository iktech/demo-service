FROM golang:1.19.2-alpine3.16 as build

ARG VERSION=1.0.0

RUN mkdir /app

WORKDIR /app

ADD *.go go.mod go.sum /app/
ADD cmd /app/cmd
ADD core /app/core
ADD data /app/data
ADD internal /app/internal

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w -s -X main.version=${VERSION}" -o demo-service cmd/demo/main.go

FROM scratch
LABEL maintainer="Igor Kolomiyets <igor.kolomiyets@iktech.io>"

COPY --from=build /app/demo-service /app/demo-service

CMD ["/app/demo-service"]
