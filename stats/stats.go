package stats

import (
	"log"
	"time"

	"github.com/TheDistributedBay/TheDistributedBay/database"
	"github.com/TheDistributedBay/TheDistributedBay/network"
)

func ReportStats(db *database.TorrentDB, cm *network.ConnectionManager) {
	for range time.Tick(time.Minute) {
		log.Printf("TorrentDB(%d torrents) PeerManager(%d peers)", db.NumTorrents(), cm.NumPeers())
	}
}
