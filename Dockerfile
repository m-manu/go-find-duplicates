FROM golang:1.22-alpine3.19 AS builder

RUN apk --no-cache add build-base

WORKDIR /opt/go-find-duplicates

COPY ./go.mod ./go.sum ./

RUN go mod download -x

COPY . .

RUN go build

RUN go test ./...

FROM alpine:3.19

RUN apk --no-cache add bash

COPY --from=builder /opt/go-find-duplicates/go-find-duplicates /bin
