# syntax=docker/dockerfile:1
FROM ubuntu:20.04
# ubuntu
MAINTAINER  eddie@hessegg.com

WORKDIR     klaatoo-explorer/

COPY        build/ ./

RUN ./klaatoo-explorer client start