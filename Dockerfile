FROM golang
MAINTAINER Tristan Rice, rice@fn.lc

RUN apt-get update && apt-get install -y \
		libleveldb-dev \
    elasticsearch \
	&& rm -rf /var/lib/apt/lists/*

RUN go get github.com/TheDistributedBay/TheDistributedBay
RUN go install github.com/TheDistributedBay/TheDistributedBay

ENTRYPOINT /usr/bin/bash -c "service elasticsearch start; /go/bin/TheDistributedBay"

EXPOSE 8080 7654
