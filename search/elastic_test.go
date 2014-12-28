package search

import (
	"math/rand"
	"testing"

	"github.com/TheDistributedBay/TheDistributedBay/core"
)

func TestNewElastic(t *testing.T) {
	host := "localhost"
	e, err := NewElastic(host, "thedistributedbay_test")
	if err != nil {
		t.Fatal(err)
	}
	if e.c.Domain != host {
		t.Fatal("Host failed to set")
	}
}

func TestElasticNewTorrent(t *testing.T) {
	e, err := NewElastic("localhost", "thedistributedbay_test")
	if err != nil {
		t.Fatal(err)
	}
	t1 := &core.Torrent{
		Hash: randSeq(40),
		Name: randSeq(40),
	}
	e.NewTorrent(t1)

	exists := e.Exists(t1.Hash)
	if exists != nil {
		t.Fatal("Hash doesn't exist!", exists)
	}

	t2 := &core.Torrent{
		Hash: randSeq(40),
		Name: randSeq(40),
	}
	e.NewBatchedTorrent(t2)
	e.Flush()
	e.b.Stop()

	exists = e.Exists(t2.Hash)
	if exists != nil {
		t.Fatal("Hash doesn't exist!", exists)
	}

	// Test no hash
	exists = e.Exists("this hash shouldn't exist ever")
	if exists == nil {
		t.Fatal("Hash shouldn't exist!")
	}
}

func TestElasticSearch(t *testing.T) {

	e, err := NewElastic("localhost", "thedistributedbay_test")
	if err != nil {
		t.Fatal(err)
	}

	desc := randSeq(40)
	t1 := &core.Torrent{
		Hash:        randSeq(40),
		Name:        randSeq(40),
		Description: desc,
		CategoryID:  1,
		Seeders:     core.NewRange(100),
		Leechers:    core.NewRange(10),
	}
	t2 := &core.Torrent{
		Hash:        randSeq(40),
		Name:        randSeq(40),
		Description: desc,
		CategoryID:  2,
		Seeders:     core.NewRange(10),
		Leechers:    core.NewRange(100),
	}
	e.NewTorrent(t1)
	e.NewTorrent(t2)

	// Search with seeders desc
	results, err := e.Search(desc, 0, 10, []uint8{}, "Seeders.Min:desc")
	if err != nil {
		t.Fatal(err)
	}
	if len(results.Hits) != 2 {
		t.Fatal("Wrong number of results")
	}
	if results.Hits[0].Id != t1.Hash {
		t.Fatal("Sorting by seeders failed.")
	}

	// Search with seeders asc
	results, err = e.Search(desc, 0, 10, []uint8{}, "Seeders.Min:asc")
	if err != nil {
		t.Fatal(err)
	}
	if len(results.Hits) != 2 {
		t.Fatal("Wrong number of results")
	}
	if results.Hits[0].Id != t2.Hash {
		t.Fatal("Sorting by seeders failed.")
	}

	// Search with category
	results, err = e.Search(desc, 0, 10, []uint8{1, 8}, "Seeders.Min:desc")
	if err != nil {
		t.Fatal(err)
	}
	if len(results.Hits) != 1 {
		t.Fatal("Wrong number of results")
	}
	if results.Hits[0].Id != t1.Hash {
		t.Fatal("Returned result from wrong category")
	}

	// Search with from statement
	results, err = e.Search(desc, 1, 10, []uint8{}, "Seeders.Min:desc")
	if err != nil {
		t.Fatal(err)
	}
	if len(results.Hits) != 1 {
		t.Fatal("Wrong number of results")
	}
	if results.Hits[0].Id != t2.Hash {
		t.Fatal("Returned wrong result")
	}

	// Search with no sort statement
	results, err = e.Search(desc, 0, 10, []uint8{}, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(results.Hits) != 2 {
		t.Fatal("Wrong number of results")
	}

	// Search with no query
	results, err = e.Search("", 0, 10, []uint8{}, "Seeders.Min:desc")
	if err != nil {
		t.Fatal(err)
	}
	if len(results.Hits) < 2 {
		t.Fatal("Wrong number of results")
	}
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
