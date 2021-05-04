FROM golang:1.16.3-alpine3.13

COPY ./src /src

WORKDIR /src/

RUN apk update \
  && apk add --no-cache git \
  && go get github.com/gin-gonic/gin \
  && go get github.com/jinzhu/gorm \
  && go get github.com/go-sql-driver/mysql
