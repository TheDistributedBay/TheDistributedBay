package network

import (
	"io"
	"testing"
	"time"

	"github.com/TheDistributedBay/TheDistributedBay/database"
	"github.com/TheDistributedBay/TheDistributedBay/tls"
)

var (
	testTorrent = &database.Torrent{"hash", "pk", "sig", "magnetlink", "name", "description", 1, time.Now(), []string{"tag1"}}
)

type sewer struct {
	r io.ReadCloser
	w io.WriteCloser
}

func (s *sewer) Close() error {
	s.r.Close()
	s.w.Close()
	return nil
}

func (s *sewer) Read(b []byte) (n int, err error) {
	return s.r.Read(b)
}

func (s *sewer) Write(b []byte) (n int, err error) {
	return s.w.Write(b)
}

func (s *sewer) Protocol() string {
	return tls.Proto
}

type dummyListener chan *database.Torrent

func (d *dummyListener) NewTorrent(t *database.Torrent) {
	*d <- t
}

func testSewer() (*sewer, *sewer) {
	pr1, pw1 := io.Pipe()
	pr2, pw2 := io.Pipe()
	s1 := &sewer{pr1, pw2}
	s2 := &sewer{pr2, pw1}
	return s1, s2
}

func TestSingleHop(t *testing.T) {
	db1 := database.NewTorrentDB()
	db1.Add(testTorrent)
	cm1 := NewConnectionManager(db1)

	db2 := database.NewTorrentDB()
	cm2 := NewConnectionManager(db2)

	l, r := testSewer()
	go cm1.Handle(l)
	go cm2.Handle(r)

	c := make(chan *database.Torrent)
	listen := dummyListener(c)
	db2.AddClient(&listen)

	select {
	case <-listen:
	case <-time.After(time.Second):
	}

	recv, err := db2.Get("hash")
	if err != nil {
		t.Errorf("Expected torrent with hash in %v, error: %v", db2, err)
	}

	if recv.Hash != testTorrent.Hash {
		t.Errorf("Expected torrent with %v, got %v", testTorrent, recv)
	}

	t2 := &*testTorrent
	t2.Hash = "foobar"

	db1.Add(t2)

	select {
	case <-listen:
	case <-time.After(time.Second):
	}

	recv, err = db2.Get("foobar")
	if err != nil {
		t.Errorf("Expected torrent with hash in %v, error: %v", db2, err)
	}

	if recv.Hash != t2.Hash {
		t.Errorf("Expected torrent with %v, got %v", testTorrent, recv)
	}

	cm1.Close()
	cm2.Close()
}
