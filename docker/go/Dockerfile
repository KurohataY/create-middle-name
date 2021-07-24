FROM golang:1.16.3-alpine3.13

ARG GCC_VERSION="10.2.0"
ENV GCC_VERSION=${GCC_VERSION}
ENV LC_ALL="C"

COPY ./src /src

WORKDIR /src/

RUN apk update \
  && apk add build-base \
  && apk --no-cache add make git gcc libtool musl-dev ca-certificates dumb-init \
  && go get github.com/gin-gonic/gin \
  && go get github.com/jinzhu/gorm \
  && go get github.com/go-sql-driver/mysql

RUN apk add build-base
