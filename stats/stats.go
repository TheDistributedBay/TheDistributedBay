package stats

import (
	"log"
	"time"

	"github.com/TheDistributedBay/TheDistributedBay/database"
	"github.com/TheDistributedBay/TheDistributedBay/network"
)

func ReportStats(db *database.TorrentDB, cm *network.ConnectionManager) {
	for range time.Tick(time.Minute) {
		c := make(chan string)
		go db.GetTorrents(c)
		count := 0
		for range c {
			count += 1
		}
		log.Printf("TorrentDB(%d torrents) PeerManager(%d peers)", count, cm.NumPeers())
	}
}
