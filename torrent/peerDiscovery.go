package torrent

import (
	"log"
	"net"
	"strconv"
	"time"

	"github.com/nictuku/dht"

	"github.com/TheDistributedBay/TheDistributedBay/network"
	"github.com/TheDistributedBay/TheDistributedBay/tls"
)

func DiscoverPeers(cm *network.ConnectionManager, address string) {
	// Hex encoded string: "The Distributed Bay!"
	ih, err := dht.DecodeInfoHash("5468652044697374726962757465642042617921")
	if err != nil {
		log.Fatalf("DHT DecodeInfoHash error: %v\n", err)
	}

	_, portStr, err := net.SplitHostPort(address)
	if err != nil {
		log.Fatal("Bind address error!", err)
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatal(err)
	}

	config := dht.NewConfig()
	config.Port = port
	config.NumTargetPeers = 10
	node, err := dht.New(config)
	if err != nil {
		log.Fatal("Error creating DHT node!")
	}

	go node.Run()
	go checkPeers(node, cm)

	log.Println("Requesting peers from the DHT!")

	for {
		node.PeersRequest(string(ih), true)
		time.Sleep(15 * time.Second)
	}

}

func checkPeers(n *dht.DHT, cm *network.ConnectionManager) {
	for r := range n.PeersRequestResults {
		for _, peers := range r {
			for _, x := range peers {
				address := dht.DecodePeerAddress(x)
				log.Printf("DHT Found Peer: %v\n", address)
				go attemptConnectionToPeer(address, cm)
			}
		}
	}
}

func attemptConnectionToPeer(address string, cm *network.ConnectionManager) {
	c, err := tls.Dial(address)
	if err != nil {
		log.Printf("Error trying to connect to %v : %v", address, err)
	} else {
		cm.Handle(c)
	}
}
