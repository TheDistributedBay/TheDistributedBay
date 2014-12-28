package dbchannel

import (
	"testing"

	"github.com/TheDistributedBay/TheDistributedBay/core"
)

type testdb struct {
	slow chan string
	fast chan *core.Torrent
}

func (t testdb) Get(hash string) (*core.Torrent, error) {
	s := &core.Torrent{}
	s.Hash = hash
	return s, nil
}
func (t testdb) Add(*core.Torrent) error {
	return nil
}
func (t testdb) AddSignature(s *core.Signature) {}
func (t testdb) AddTorrentClient(c chan *core.Torrent) {
	for s := range t.fast {
		c <- s
	}
	close(c)
}
func (t testdb) GetTorrents(c chan string) {
	for s := range t.slow {
		c <- s
	}
	close(c)
}

func (t *testdb) CloseSlow() {
	o := t.slow
	t.slow = make(chan string)
	close(o)
}

func (t *testdb) CloseFast() {
	o := t.fast
	t.fast = make(chan *core.Torrent)
	close(o)
}

func TestDBChannel(t *testing.T) {
	db := &testdb{make(chan string), make(chan *core.Torrent)}
	c := New(db)

	go func() {
		db.slow <- "foo"
	}()
	r := <-c.Output
	if r.GetHash() != "foo" {
		t.Fatal("Expected to recieve torrent foo received", r.GetHash())
	}
	a, _ := r.GetTorrent()
	h := a.Hash
	if h != "foo" {
		t.Fatal("Expected to recieve torrent foo received", h)
	}

	go func() {
		s := &core.Torrent{}
		s.Hash = "foo"
		db.fast <- s
	}()
	r = <-c.Output
	if r.GetHash() != "foo" {
		t.Fatal("Expected to recieve torrent foo received", r.GetHash())
	}
	a, _ = r.GetTorrent()
	h = a.Hash
	if h != "foo" {
		t.Fatal("Expected to recieve torrent foo received", h)
	}

	db.CloseSlow()
	go func() {
		s := &core.Torrent{}
		s.Hash = "foo"
		db.fast <- s
	}()
	r = <-c.Output
	if r.GetHash() != "foo" {
		t.Fatal("Expected to recieve torrent foo received", r.GetHash())
	}
	a, _ = r.GetTorrent()
	h = a.Hash
	if h != "foo" {
		t.Fatal("Expected to recieve torrent foo received", h)
	}

	db.CloseFast()
	go func() {
		s := &core.Torrent{}
		s.Hash = "foo"
		db.fast <- s
	}()
	r = <-c.Output
	if r.GetHash() != "foo" {
		t.Fatal("Expected to recieve torrent foo received", r.GetHash())
	}
	a, _ = r.GetTorrent()
	h = a.Hash
	if h != "foo" {
		t.Fatal("Expected to recieve torrent foo received", h)
	}

	go func() {
		db.slow <- "foo"
	}()
	r = <-c.Output
	if r.GetHash() != "foo" {
		t.Fatal("Expected to recieve torrent foo received", r.GetHash())
	}
	a, _ = r.GetTorrent()
	h = a.Hash
	if h != "foo" {
		t.Fatal("Expected to recieve torrent foo received", h)
	}
}
