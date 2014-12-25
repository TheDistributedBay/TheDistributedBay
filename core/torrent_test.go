package core

import (
	"testing"
	"time"
)

func TestTorrentCreationAndVerify(t *testing.T) {
	t1 := CreateTorrent([]byte("ml"), "n", "d", "cid", time.Now(), []string{"foo"})
	err := t1.VerifyTorrent()
	if err != nil {
		t.Fatal(err)
	}
	t2 := CreateTorrent([]byte("ml"), "n", "d", "cid", time.Now(), []string{"2"})

	if t1.Hash == t2.Hash {
		t.Fatal("Hashes are identical")
	}
}

func BenchmarkDefaultCreate(b *testing.B) {
	c := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CreateTorrent([]byte("ml"), "n", "d", "cid", c, []string{"foo"})
	}
}

func BenchmarkDefaultVerify(b *testing.B) {
	c := time.Now()
	t := CreateTorrent([]byte("ml"), "n", "d", "cid", c, []string{"foo"})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		t.VerifyTorrent()
	}
}
