package main

import (
	"flag"
	"log"
	"runtime"

	"github.com/TheDistributedBay/TheDistributedBay/database"
	"github.com/TheDistributedBay/TheDistributedBay/frontend"
	"github.com/TheDistributedBay/TheDistributedBay/network"
	"github.com/TheDistributedBay/TheDistributedBay/tls"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	listen := flag.String("listen", ":7654", "Address to listen on")
	connect := flag.String("connect", "", "Address to connect to")
	http_address := flag.String("http", ":8080", "Address to listen on for HTTP")
	flag.Parse()

	db := database.NewTorrentDB()
	cm := network.NewConnectionManager(db)

	if *connect != "" {
		c, err := tls.Dial(*connect)
		if err != nil {
			log.Fatalf("Error trying to connect to %v : %v", *connect, err)
		}
		cm.Handle(tls.Wrap(c))
	} else {
		l, err := tls.Listen(*listen)
		if err != nil {
			log.Fatalf("Error trying to listen to %v : %v", *listen, err)
		}
		go cm.Listen(l)
	}

	log.Println("Running...")
	frontend.Serve(http_address, db)
	cm.Close()
}
