package frontend

import (
	"github.com/TheDistributedBay/TheDistributedBay/database"
	"github.com/blevesearch/bleve"
	"log"
	"os"
)

func NewTorrentClient(db database.Database) TorrentClient {
	mapping := bleve.NewIndexMapping()
	os.RemoveAll("search.bleve")
	index, err := bleve.New("search.bleve", mapping)
	if err != nil {
		log.Fatal(err)
	}
	return TorrentClient{index, db}
}

type IndexableTorrent struct {
	Name, Description string
	Tags              []string
}

type TorrentClient struct {
	bleve.Index
	database.Database
}

func (tc TorrentClient) NewTorrent(t *database.Torrent) {
	indexable := IndexableTorrent{
		Name:        t.Name,
		Description: t.Description,
		Tags:        t.Tags,
	}
	tc.Index.Index(t.Hash, indexable)
}

func (tc TorrentClient) Search(queryStr string) (*bleve.SearchResult, error) {
	query := bleve.NewQueryStringQuery(queryStr)
	searchRequest := bleve.NewSearchRequest(query)
	result, err := tc.Index.Search(searchRequest)
	if err != nil {
		return nil, err
	}
	return result, nil
}
