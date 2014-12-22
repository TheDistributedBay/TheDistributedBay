package main

import (
	"flag"
	"fmt"
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
	flag.Parse()

	db := database.NewTorrentDB()
	cm := network.NewConnectionManager(db)

	if *connect != "" {
		c, err := tls.Dial(*connect)
		if err != nil {
			log.Fatalf("Error trying to connect to %v : %v", *connect, err)
		}
		cm.Handle(c)
	} else {
		l, err := tls.Listen(*listen)
		if err != nil {
			log.Fatalf("Error trying to listen to %v : %v", *listen, err)
		}
		go cm.Listen(l)
	}

	fmt.Println("Running...")
	frontend.Serve()
	cm.Close()
}
