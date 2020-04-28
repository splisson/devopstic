# build stage
FROM golang:1.13.10-alpine3.11 AS build-env
RUN mkdir -p /opt && apk add bash make curl git gcc libc-dev && curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
ADD . /opt/devopstic
RUN cd /opt/devopstic && make

# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /opt/devopstic/bin/devopstic /app/
RUN apk update && apk add ca-certificates curl && rm -rf /var/cache/apk/*
RUN update-ca-certificates
