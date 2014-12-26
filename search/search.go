package search

import (
	"log"

	"github.com/mattbaird/elastigo/lib"

	"github.com/TheDistributedBay/TheDistributedBay/core"
)

type SearchProvider interface {
	Exists(hash string) error
	NewBatchedTorrent(t *core.Torrent)
	Search(term string, from, size int, categories []uint8, sort string) (*elastigo.Hits, error)
}

func NewSearcher(db core.Database) (*Searcher, error) {
	b, err := NewElastic("localhost")
	if err != nil {
		return nil, err
	}
	s := &Searcher{b, db}
	go s.shovel()
	return s, nil
}

type Searcher struct {
	b  SearchProvider
	db core.Database
}

func (s *Searcher) shovel() {
	for {
		log.Print("Starting index run")
		c := make(chan string)
		go s.db.GetTorrents(c)
		count := 0
		for h := range c {
			if s.b.Exists(h) == nil {
				continue
			}
			t, err := s.db.Get(h)
			if err != nil {
				log.Print(err)
				continue
			}
			count += 1
			s.b.NewBatchedTorrent(t)
		}
		log.Printf("%d new torrents indexed", count)
	}
}

func (s *Searcher) Search(term string, from, size int, categories []uint8, sort string) ([]*core.Torrent, int, error) {
	result, err := s.b.Search(term, from, size, categories, sort)
	if err != nil {
		return nil, 0, err
	}
	torrents := make([]*core.Torrent, 0, size)
	for _, e := range result.Hits {
		t, err := s.db.Get(e.Id)
		if err != nil {
			log.Print("Stale torrent in index ", e.Id)
			continue
		}
		torrents = append(torrents, t)
	}
	return torrents, result.Total, nil
}
