package database

import (
	"os"
	"testing"

	"github.com/TheDistributedBay/TheDistributedBay/core"
)

func NumTorrents(d core.Database) int {
	c := make(chan string)
	go d.GetTorrents(c)
	count := 0
	for range c {
		count += 1
	}
	return count
}

func TestDB(t *testing.T) {
	db, err := NewTorrentDB("test.db")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll("test.db")

	if n := NumTorrents(db); n != 0 {
		t.Fatal("Expected 0 torrents, got %d", n)
	}

	r := core.Torrent{}
	r.Hash = "foo"
	if err := db.Add(&r); err != nil {
		t.Fatal(err)
	}

	n, err := db.Get("foo")
	if err != nil {
		t.Fatal(err)
	}

	if n.Hash != r.Hash {
		t.Fatal("%v != %v", n, r)
	}
}
