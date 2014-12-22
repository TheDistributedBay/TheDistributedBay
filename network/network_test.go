package network

import (
	"io"
	"testing"
	"time"

	"github.com/TheDistributedBay/TheDistributedBay/database"
)

var (
	testTorrent = &database.Torrent{"hash", "pk", "sig", "magnetlink", "name", "description", 1, 1, 1, time.Now(), []string{"tag1"}}
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

	time.Sleep(time.Second)

	if _, err := db2.Get("hash"); err != nil {
		t.Errorf("Expected torrent with hash in %v, error: %v", db2, err)
	}
	cm1.Close()
	cm2.Close()
}
