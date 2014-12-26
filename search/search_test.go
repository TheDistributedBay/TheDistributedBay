package search

import (
	"errors"
	"testing"
	"time"

	"github.com/mattbaird/elastigo/lib"

	"github.com/TheDistributedBay/TheDistributedBay/core"
)

type TestDB struct {
	t *testing.T
}

func (db *TestDB) Add(t *core.Torrent) error {
	return nil
}
func (db *TestDB) AddSignature(s *core.Signature)        {}
func (db *TestDB) AddTorrentClient(c chan *core.Torrent) {}
func (db *TestDB) Get(hash string) (*core.Torrent, error) {
	if hash == "search get" {
		return &core.Torrent{Name: "search torrent"}, nil
	} else if hash == "search bad" {
		return nil, errors.New("stale torrent")
	} else if hash != "test hash" {
		db.t.Fatal("Received incorrect hash!", hash)
	}
	return &core.Torrent{Hash: "test hash"}, nil
}
func (db *TestDB) GetTorrents(c chan string) {
	c <- "test hash"
	c <- "search bad"
	c <- "test hash found"
	close(c)
}

type TestSearchProvider struct {
	t *testing.T
	c chan bool
}

func (sp *TestSearchProvider) Exists(hash string) error {
	if hash == "test hash" {
		sp.c <- true
		return errors.New("it doesn't exist. :(")
	} else {
		return nil
	}
}
func (sp *TestSearchProvider) NewBatchedTorrent(t *core.Torrent) {
	sp.c <- true
}
func (sp *TestSearchProvider) Search(term string, from, size int, categories []uint8, sort string) (*elastigo.Hits, error) {
	if term == "bad" {
		return nil, errors.New("bad query")
	}
	if term != "test" || from != 5 || size != 10 || categories[0] != 1 || sort != "banana:asc" {
		sp.t.Fatal("Search parameters do not match the ones passed in.")
	}
	return &elastigo.Hits{
		Total: 10000,
		Hits: []elastigo.Hit{elastigo.Hit{
			Id: "search get",
		}, elastigo.Hit{
			Id: "search bad",
		}},
	}, nil
}

func TestSearcherCreation(t *testing.T) {
	db := TestDB{}
	s, err := NewSearcher(&db)
	if err != nil {
		t.Fatal(err)
	}
	if s.b == nil {
		t.Fatal("SearchProvider is nil")
	}
	if s.db != &db {
		t.Fatal("DB was not set!")
	}
}
func TestSearcherShovel(t *testing.T) {
	c := make(chan bool)
	b := &TestSearchProvider{t: t, c: c}
	db := &TestDB{t}
	s := &Searcher{b, db}
	go s.shovel()
	select {
	case <-c:
	case <-time.After(time.Second):
		t.Fatal("exist failed to call")
	}
	select {
	case <-c:
	case <-time.After(time.Second):
		t.Fatal("new batched torrent failed to be called")
	}
}
func TestSearcherSearch(t *testing.T) {
	c := make(chan bool, 2)
	b := &TestSearchProvider{t: t, c: c}
	db := &TestDB{t}
	s := &Searcher{b, db}
	torrents, count, err := s.Search("test", 5, 10, []uint8{1, 2}, "banana:asc")
	if err != nil {
		t.Fatal(err)
	}
	if count != 10000 {
		t.Fatal("Incorrect results count")
	}
	if len(torrents) != 1 {
		t.Fatal("Wrong number of returned torrents.")
	}
	if torrents[0].Name != "search torrent" {
		t.Fatal("Wrong torrent")
	}

	// Bad search
	_, _, err = s.Search("bad", 5, 10, []uint8{1, 2}, "banana:asc")
	if err == nil {
		t.Fatal("Search should have caused an error")
	}
}
