# build stage
FROM 906629767317.dkr.ecr.us-east-1.amazonaws.com/weyv/golang-base:v1.12.6-alpine3.9 AS build-env
ADD . /go/src/github.com/splisson/opstic
RUN cd /go/src/github.com/splisson/opstic && export GOPATH=/go && make

# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /go/src/github.com/splisson/opstic/bin/opstic /app/
RUN apk update && apk add ca-certificates curl && rm -rf /var/cache/apk/*
RUN update-ca-certificates
# ENTRYPOINT ./feedGrab
