package search

import (
	"fmt"
	"testing"
	"time"

	"github.com/TheDistributedBay/TheDistributedBay/core"
)

func testTorrent(i string, c time.Time) *core.Torrent {
	return simpleTorrent(i, i, i, c)
}

func simpleTorrent(hash, name, description string, c time.Time) *core.Torrent {
	return &core.Torrent{hash, []byte("magnetlink"), name, description, 1, c, []string{"tags"}}
}

func TestBleve(t *testing.T) {
	bleve, err := NewBleve("test.bleve")
	if err != nil {
		t.Fatal(err)
	}
	bleve.NewTorrent(simpleTorrent("t1", "foo", "", time.Now()))
	r, err := bleve.Search("foo", 0, 10)
	if err != nil {
		t.Fatal(err)
	}
	if r.Total != 1 {
		t.Fatal("Unable to find foo")
	}
	bleve.NewBatchedTorrent(simpleTorrent("t2", "bar", "", time.Now()))
	bleve.Flush()
	r, err = bleve.Search("bar", 0, 10)
	if err != nil {
		t.Fatal(err)
	}
	if r.Total != 1 {
		t.Fatal("Unable to find bar")
	}
	// Double index see what happens
	bleve.NewBatchedTorrent(simpleTorrent("t2", "bar", "", time.Now()))
	bleve.Flush()
	r, err = bleve.Search("bar", 0, 10)
	if err != nil {
		t.Fatal(err)
	}
	if r.Total != 1 {
		t.Fatal("Double indexed")
	}
	if err := bleve.Exists("t2"); err != nil {
		t.Fatal("should have returned nil", err)
	}
}

func BenchmarkTorrentCreation(b *testing.B) {
	c := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testTorrent(fmt.Sprint(i), c)
	}
}

func BenchmarkNormalBleve(b *testing.B) {
	bleve, _ := NewBleve("test.bleve")
	c := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bleve.NewTorrent(testTorrent(fmt.Sprint(i), c))
	}
}

func BenchmarkBatchBleve10(b *testing.B) {
	bleve, _ := NewBleve("test.bleve")
	bleve.BatchSize = 10
	c := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bleve.NewBatchedTorrent(testTorrent(fmt.Sprint(i), c))
	}
	bleve.Flush()
}

func BenchmarkBatchBleve50(b *testing.B) {
	bleve, _ := NewBleve("test.bleve")
	bleve.BatchSize = 50
	c := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bleve.NewBatchedTorrent(testTorrent(fmt.Sprint(i), c))
	}
	bleve.Flush()
}

func BenchmarkBatchBleve100(b *testing.B) {
	bleve, _ := NewBleve("test.bleve")
	bleve.BatchSize = 100
	c := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bleve.NewBatchedTorrent(testTorrent(fmt.Sprint(i), c))
	}
	bleve.Flush()
}

func BenchmarkBatchBleve200(b *testing.B) {
	bleve, _ := NewBleve("test.bleve")
	bleve.BatchSize = 200
	c := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bleve.NewBatchedTorrent(testTorrent(fmt.Sprint(i), c))
	}
	bleve.Flush()
}

func BenchmarkBatchBleve500(b *testing.B) {
	bleve, _ := NewBleve("test.bleve")
	bleve.BatchSize = 500
	c := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bleve.NewBatchedTorrent(testTorrent(fmt.Sprint(i), c))
	}
	bleve.Flush()
}

func BenchmarkBatchBleve1000(b *testing.B) {
	bleve, _ := NewBleve("test.bleve")
	bleve.BatchSize = 1000
	c := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bleve.NewBatchedTorrent(testTorrent(fmt.Sprint(i), c))
	}
	bleve.Flush()
}
