FROM golang:1.16.3-alpine3.13 as builder

WORKDIR /opt/go-find-duplicates

ADD . ./

RUN go build

FROM alpine:3.13

RUN apk --no-cache add bash

COPY --from=builder /opt/go-find-duplicates/go-find-duplicates /bin
