package main

import (
	"flag"
	"log"
	"runtime"

	_ "expvar"
	_ "net/http/pprof"

	"github.com/TheDistributedBay/TheDistributedBay/database"
	"github.com/TheDistributedBay/TheDistributedBay/frontend"
	"github.com/TheDistributedBay/TheDistributedBay/importer"
	"github.com/TheDistributedBay/TheDistributedBay/network"
	"github.com/TheDistributedBay/TheDistributedBay/search"
	"github.com/TheDistributedBay/TheDistributedBay/stats"
	"github.com/TheDistributedBay/TheDistributedBay/tls"
	"github.com/TheDistributedBay/TheDistributedBay/torrent"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	listen := flag.String("listen", ":7654", "Address to listen on")
	connect := flag.String("connect", "", "Address to connect to")
	httpAddress := flag.String("http", ":8080", "Address to listen on for HTTP")
	dumpPath := flag.String("databasedump", "", "The path to a database dump that can be loaded. This should end in .csv.gz")
	devAssets := flag.Bool("devassets", false, "Tells the front end to serve development assests instead of the precompiled production ones.")
	flag.Parse()

	db, err := database.NewTorrentDB("torrent.db")
	if err != nil {
		log.Fatal("Error opening torrent database", err)
	}
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

	if *dumpPath != "" {
		go importer.Import(*dumpPath, db)
	}

	go stats.ReportStats(db, cm)
	go torrent.DiscoverPeers(cm, *listen)

	s, err := search.NewSearcher(db)
	if err != nil {
		log.Fatal("Error opening index", err)
	}

	log.Println("Running...")
	frontend.Serve(httpAddress, db, s, *devAssets)
	cm.Close()
}
