package search

import (
	"testing"

	"github.com/TheDistributedBay/TheDistributedBay/database"
)

func TestSearcherCreation(t *testing.T) {
	db, err := database.NewTorrentDB("test.db")
	if err != nil {
		t.Fatal(err)
	}
	s, err := NewSearcher(db)
	if err != nil {
		t.Fatal(err)
	}
	if s.b == nil {
		t.Fatal("SearchProvider is nil")
	}
}
