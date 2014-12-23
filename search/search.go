package search

import (
	"log"
	"os"

	"github.com/blevesearch/bleve"

	"github.com/TheDistributedBay/TheDistributedBay/database"
)

func NewSearcher(db database.Database) (*Searcher, error) {
	mapping := bleve.NewIndexMapping()
	mapping.TypeMapping["Hash"] = bleve.NewDocumentDisabledMapping()
	mapping.TypeMapping["PublicKey"] = bleve.NewDocumentDisabledMapping()
	mapping.TypeMapping["Signature"] = bleve.NewDocumentDisabledMapping()
	os.RemoveAll("search.bleve")
	index, err := bleve.New("search.bleve", mapping)
	if err != nil {
		return nil, err
	}
	s := &Searcher{index, db}
	go db.AddClient(s)
	return s, nil
}

type Searcher struct {
	i  bleve.Index
	db database.Database
}

func (s *Searcher) NewTorrent(t *database.Torrent) {
	s.i.Index(t.Hash, t)
}

func (s *Searcher) Search(term string, from int, size int) ([]*database.Torrent, uint64, error) {
	q := bleve.NewQueryStringQuery(term)
	r := bleve.NewSearchRequestOptions(q, size, from, false)
	result, err := s.i.Search(r)
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
