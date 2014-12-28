package torrent

import (
	"log"

	"github.com/TheDistributedBay/TheDistributedBay/core"
	"github.com/TheDistributedBay/TheDistributedBay/search"
)

type StatsUpdater struct {
	s  *search.Searcher
	db core.Database
	c  chan *core.Torrent
}

func NewStatsUpdater(s *search.Searcher, db core.Database) *StatsUpdater {
	updater := &StatsUpdater{s, db, make(chan *core.Torrent)}
	go updater.run()
	return updater
}

func (updater *StatsUpdater) run() {
	batch := make([]*core.Torrent, 70)
	i := 0
	for torrent := range updater.c {
		batch[i] = torrent
		i += 1
		if i == 70 {
			go updater.processBatch(batch)
			batch = make([]*core.Torrent, 70)
			i = 0
		}
	}
}

func (updater *StatsUpdater) processBatch(torrents []*core.Torrent) {
	hashes := make([]string, len(torrents))
	for i, t := range torrents {
		hashes[i] = t.NiceInfoHash()
	}
	resp, err := ScrapeTrackers(core.Trackers, hashes)
	if err != nil {
		log.Println(err)
		return
	}
	for i, details := range resp {
		t := torrents[i]
		t.Seeders = details.Seeders
		t.Leechers = details.Leechers
		t.Completed = details.Completed
		updater.s.NewBatchedTorrent(t)
	}
	log.Println("Updated tracker information for", len(resp), "torrents.")
}

func (updater *StatsUpdater) QueueTorrent(torrent *core.Torrent) {
	updater.c <- torrent
}
