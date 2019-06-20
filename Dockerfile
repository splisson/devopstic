# build stage
FROM golang:1.12.6-alpine3.9 AS build-env
RUN apk add bash make curl git gcc libc-dev && curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
ADD . /go/src/github.com/splisson/devopstic
RUN cd /go/src/github.com/splisson/devopstic && export GOPATH=/go && make

# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /go/src/github.com/splisson/devopstic/bin/devopstic /app/
RUN apk update && apk add ca-certificates curl && rm -rf /var/cache/apk/*
RUN update-ca-certificates
