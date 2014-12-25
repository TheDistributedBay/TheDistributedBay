package search

import (
	"log"

	"github.com/TheDistributedBay/TheDistributedBay/core"
)

func NewSearcher(db core.Database, dir string) (*Searcher, error) {
	b, err := NewBleve(dir)
	if err != nil {
		return nil, err
	}
	b.BatchSize = 100
	s := &Searcher{b, db}
	go s.shovel()
	return s, nil
}

type Searcher struct {
	b  *Bleve
	db core.Database
}

func (s *Searcher) shovel() {
	for {
		c := make(chan string)
		go s.db.GetTorrents(c)
		for h := range c {
			if s.b.Exists(h) == nil {
				continue
			}
			t, err := s.db.Get(h)
			if err != nil {
				log.Print(err)
				continue
			}
			s.b.NewBatchedTorrent(t)
		}
	}
}

func (s *Searcher) Search(term string, from int, size int) ([]*core.Torrent, uint64, error) {
	result, err := s.b.Search(term, from, size)
	if err != nil {
		return nil, 0, err
	}
	torrents := make([]*core.Torrent, 0, len(result.Hits))
	for _, e := range result.Hits {
		t, err := s.db.Get(e.ID)
		if err != nil {
			log.Print("Stale torrent in index %s", e.ID)
			continue
		}
		torrents = append(torrents, t)
	}
	return torrents, result.Total, nil
}
