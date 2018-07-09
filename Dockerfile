# OFFICIAL REPOSITORY: https://hub.docker.com/_/golang/
FROM golang:1.10.2

MAINTAINER Yuichi Watanabe

ENV DEBIAN_FRONTEND noninteractive

ENV SRC_DIR /go/src/github.com/vvatanabe/backlog-bulk-issue-register
RUN mkdir -p $SRC_DIR
WORKDIR $SRC_DIR