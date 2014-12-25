package search

import (
	"log"
	"time"

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
		time.Sleep(30 * time.Second)
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

func (s *Searcher) Search(term string, from, size int, categories []uint8) ([]*core.Torrent, int, error) {
	result, err := s.b.Search(term, 0, 100000000)
	if err != nil {
		return nil, 0, err
	}
	torrents := make([]*core.Torrent, 0, size)
	matchCount := 0
	for _, e := range result.Hits {
		t, err := s.db.Get(e.ID)
		if err != nil {
			log.Print("Stale torrent in index %s", e.ID)
			continue
		}
    included := false
    for _, category := range categories {
      if category == t.CategoryID {
        included = true
        break
      }
    }
		if len(categories) == 0 || categories[0] == 0 || included {
			if matchCount >= from && matchCount < (from+size) {
				torrents = append(torrents, t)
			}
			matchCount += 1
		}
	}
	return torrents, matchCount, nil
}
