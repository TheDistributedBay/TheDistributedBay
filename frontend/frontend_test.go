package frontend

import (
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/TheDistributedBay/TheDistributedBay/core"
	"github.com/TheDistributedBay/TheDistributedBay/search"
)

func setupHttpTest(t *testing.T) (core.Database, *search.Searcher, string) {
	db := &TestDB{t}
	searcher, err := search.NewSearcher(db, "thedistributedbay_test")
	if err != nil {
		t.Fatal(err)
	}
	router := ApiRouter(db, searcher)
	server := httptest.NewServer(router)

	return db, searcher, server.URL
}

type TestDB struct {
	t *testing.T
}

func (db *TestDB) Add(t *core.Torrent) error {
	return nil
}
func (db *TestDB) AddSignature(s *core.Signature)        {}
func (db *TestDB) AddTorrentClient(c chan *core.Torrent) {}
func (db *TestDB) Get(hash string) (*core.Torrent, error) {
	if hash == "bad hash" {
		return nil, errors.New("invalid hash")
	}
	return &core.Torrent{
		Hash:        "test hash",
		Name:        "test hash",
		Description: "test description",
		CategoryID:  1,
		CreatedAt:   time.Unix(0, 0),
		Files:       100,
		Size:        100,
		Seeders:     core.NewRange(100),
		Leechers:    core.NewRange(100),
		Completed:   core.NewRange(100),
	}, nil
}
func (db *TestDB) GetTorrents(c chan string) {}
