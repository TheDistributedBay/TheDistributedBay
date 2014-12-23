package search

import (
	"github.com/TheDistributedBay/TheDistributedBay/database"
	"github.com/blevesearch/bleve"
	"os"
)

func NewSearcher(db database.Database) (*Searcher, error) {
	mapping := bleve.NewIndexMapping()
	os.RemoveAll("search.bleve")
	index, err := bleve.New("search.bleve", mapping)
	if err != nil {
		return nil, err
	}
	s := &Searcher{index}
	go db.AddClient(s)
	return s, nil
}

type IndexableTorrent struct {
	Name, Description string
	Tags              []string
}

type Searcher struct {
	bleve.Index
}

func (s *Searcher) NewTorrent(t *database.Torrent) {
	indexable := IndexableTorrent{
		Name:        t.Name,
		Description: t.Description,
		Tags:        t.Tags,
	}
	s.Index.Index(t.Hash, t)
	s.Index.Index(t.Hash, indexable)
}

func (s *Searcher) Search(queryStr string) (*bleve.SearchResult, error) {
	query := bleve.NewQueryStringQuery(queryStr)
	searchRequest := bleve.NewSearchRequest(query)
	result, err := s.Index.Search(searchRequest)
	if err != nil {
		return nil, err
	}
	return result, nil
}
