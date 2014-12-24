FROM golang
MAINTAINER Tristan Rice, rice@fn.lc

RUN apt-get install -y libleveldb-dev

RUN go get github.com/TheDistributedBay/TheDistributedBay
RUN go install github.com/TheDistributedBay/TheDistributedBay

ENTRYPOINT /go/bin/TheDistributedBay

EXPOSE 8080 7654
