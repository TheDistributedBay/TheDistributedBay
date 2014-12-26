package network

import (
	"io"
	"log"
	"os"
	"testing"
	"time"

	"github.com/TheDistributedBay/TheDistributedBay/core"
	"github.com/TheDistributedBay/TheDistributedBay/database"
	"github.com/TheDistributedBay/TheDistributedBay/tls"
)

func createDefaultTorrent(d string) *core.Torrent {
	return core.CreateTorrent([]byte("magnetlink"), "name", d, "category", time.Now(), nil, 0, 0, 0, 0)
}

type sewer struct {
	r io.ReadCloser
	w io.WriteCloser
}

func (s *sewer) Close() error {
	s.w.Close()
	return nil
}

func (s *sewer) Read(b []byte) (n int, err error) {
	return s.r.Read(b)
}

func (s *sewer) Write(b []byte) (n int, err error) {
	return s.w.Write(b)
}

func (s *sewer) Handshake() error {
	return nil
}

func (s *sewer) Protocol() string {
	return tls.Proto
}

func testSewer() (*sewer, *sewer) {
	pr1, pw1 := io.Pipe()
	pr2, pw2 := io.Pipe()
	s1 := &sewer{pr1, pw2}
	s2 := &sewer{pr2, pw1}
	return s1, s2
}

func TestSingleHop(t *testing.T) {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	t1 := createDefaultTorrent("test1")

	db1, err := database.NewTorrentDB("db1")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll("db1")
	db1.Add(t1)
	cm1 := NewConnectionManager(db1)

	db2, err := database.NewTorrentDB("db2")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll("db2")
	cm2 := NewConnectionManager(db2)

	l, r := testSewer()
	go cm1.Handle(l)
	go cm2.Handle(r)

	c := make(chan *core.Torrent, 2)
	db2.AddTorrentClient(c)

	select {
	case <-c:
	case <-time.After(time.Second):
	}

	recv, err := db2.Get(t1.Hash)
	if err != nil {
		t.Fatalf("Expected torrent %v, error: %v", t1, err)
	}

	if recv.Hash != t1.Hash {
		t.Errorf("Expected torrent %v, got %v", t1, recv)
	}

	t2 := createDefaultTorrent("test2")
	if t2.Hash == t1.Hash {
		t.Fatalf("identical hashes... wtf")
	}

	db1.Add(t2)

	select {
	case <-c:
	case <-time.After(time.Second):
	}

	recv, err = db2.Get(t2.Hash)
	if err != nil {
		t.Fatalf("Expected torrent with hash in %v, error: %v", db2, err)
	}

	if recv.Hash != t2.Hash {
		t.Errorf("Expected torrent with %v, got %v", t1, recv)
	}

	cm1.Close()
	cm2.Close()
}
