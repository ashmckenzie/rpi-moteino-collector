FROM golang:latest
MAINTAINER ash@the-rebellion.net

ENV GOPATH /gopath
ENV APP_PATH ${GOPATH}/src/github.com/ashmckenzie/moteino-collector

RUN apt-get -qq update && apt-get -qq install -y rsync

RUN mkdir -p ${APP_PATH}
COPY app/ ${APP_PATH}

WORKDIR ${APP_PATH}
RUN make releases
