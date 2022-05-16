FROM golang:1.17-alpine3.15 as builder

RUN apk --no-cache add build-base

WORKDIR /opt/go-find-duplicates

COPY ./go.mod ./go.sum ./

RUN go mod download -x

COPY . .

RUN go build

RUN go test ./...

FROM alpine:3.15

RUN apk --no-cache add bash

COPY --from=builder /opt/go-find-duplicates/go-find-duplicates /bin
