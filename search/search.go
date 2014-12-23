package search

import (
	"log"

	"github.com/TheDistributedBay/TheDistributedBay/database"
)

func NewSearcher(db database.Database) (*Searcher, error) {
	b, err := NewBleve("search.bleve")
	b.BatchSize = 100
	if err != nil {
		return nil, err
	}
	s := &Searcher{b, db}
	go db.AddClient(s)
	return s, nil
}

type Searcher struct {
	b  *Bleve
	db database.Database
}

func (s *Searcher) NewTorrent(t *database.Torrent) {
	s.b.NewBatchedTorrent(t)
}

func (s *Searcher) Search(term string, from int, size int) ([]*database.Torrent, uint64, error) {
	result, err := s.b.Search(term, from, size)
	if err != nil {
		return nil, 0, err
	}
	torrents := make([]*database.Torrent, 0, len(result.Hits))
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
