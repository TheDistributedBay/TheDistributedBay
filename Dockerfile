FROM golang:onbuild
MAINTAINER Tristan Rice, rice@fn.lc

RUN apt-get install -y libleveldb-dev

EXPOSE 8080 7654
