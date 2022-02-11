FROM golang:1.16.14-alpine3.14 as builder

RUN apk --no-cache add build-base

WORKDIR /opt/go-find-duplicates

COPY ./go.mod ./go.sum ./

RUN go mod download

COPY . .

RUN go build

RUN go test ./...

FROM alpine:3.14

RUN apk --no-cache add bash

COPY --from=builder /opt/go-find-duplicates/go-find-duplicates /bin
