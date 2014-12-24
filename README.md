# [![TheDistributedBay](https://cdn.rawgit.com/TheDistributedBay/TheDistributedBay/master/frontend/angular/app/images/The_Distributed_Bay_logo_black.svg)](https://github.com/TheDistributedBay/TheDistributedBay)
An implementation of The Distributed Bay.
One monolithic binary that bootstraps into the network.

## Running with Docker
On Ubuntu 14.10:
```
sudo apt-get install docker.io
sudo docker run thedistributedbay/thedistributedbay
```

[Docker](https://www.docker.com/) is a container deploying environment that provides automated container images from git.

The Distributed Bay Registry page: https://registry.hub.docker.com/u/thedistributedbay/thedistributedbay/

## Running Manually
Assuming you have [Go](http://golang.org/) and the [GOPATH](https://golang.org/doc/code.html#GOPATH) correctly configured all you have to do is run:
```sh
go get github.com/TheDistributedBay/TheDistributedBay
go install github.com/TheDistributedBay/TheDistributedBay

$GOPATH/bin/TheDistributedBay
```


## Development

### Backend/Core
The Distributed Bay is primarily written in [Go](http://golang.org/).

[![GoDoc](https://godoc.org/github.com/TheDistributedBay/TheDistributedBay?status.svg)](https://godoc.org/github.com/TheDistributedBay/TheDistributedBay)

To get the source code install Go, configure the [GOPATH](https://golang.org/doc/code.html#GOPATH) and then run:
```sh
go get github.com/TheDistributedBay/TheDistributedBay
```

The source code will be available in `$GOPATH/src/github.com/TheDistributedBay/TheDistributedBay`.

### Frontend
The frontend is written in [AngularJS](https://angularjs.org/) and located in `frontend/angular/app`.

The compiled frontend assets are checked into Git. This is done so the backend developers don't need to worry about the frontend and to make deployments easier.

To edit the frontend you'll need to first install [Node.js](https://nodejs.org/) or [io.js](http://iojs.org/).

You'll also need to install Ruby and the gem `compass`.

Once done, navigate into the `frontend/angular` folder and install the dependencies.

```sh
cd frontend/angular
npm install
npm install -g bower grunt-cli
bower install
gem install compass
```

Grunt is used to handle compilation of the frontend.
To tell Grunt to automatically recompile the SCSS files run:
```sh
grunt serve
```


You'll also need to tell the backend to serve the development assets from `frontend/angular/app` instead of the production ones in `frontend/angular/dist`. You can do this by running:
```sh
go run main -devassets=true
```

