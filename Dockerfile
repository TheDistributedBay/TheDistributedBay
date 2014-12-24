FROM debian:wheezy
MAINTAINER Tristan Rice, rice@fn.lc

RUN apt-get update -y
RUN apt-get install -y golang golang-src git mercurial

ENV GOPATH /srv/go

RUN go get github.com/TheDistributedBay/TheDistributedBay

RUN cd /srv/go/github.com/TheDistributedBay/TheDistributedBay; go build

ENTRYPOINT /srv/go/github.com/TheDistributedBay/TheDistributedBay/TheDistributedBay

EXPOSE 8080 7654
