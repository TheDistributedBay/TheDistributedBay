package crypto

import (
	"fmt"
	"testing"
	"time"

	"github.com/TheDistributedBay/TheDistributedBay/database"
)

func TestMerkle(t *testing.T) {
	k, err := NewKey()
	if err != nil {
		t.Fatal(err)
	}
	ts := make([]*database.Torrent, 0)

	for i := 0; i < 5; i++ {
		r, err := CreateTorrent(k, fmt.Sprint(i), "n", "d", "cat", time.Now(), []string{"tag"})
		if err != nil {
			t.Fatal(err)
		}
		ts = append(ts, r)
	}

	s, err := SignTorrents(k, ts)
	if err != nil {
		t.Fatal(err)
	}

	err = s.VerifySignature()
	if err != nil {
		t.Fatal(err)
	}

	for _, r := range ts {
		found := false
		for _, h := range s.ListTorrents() {
			if h == r.Hash {
				found = true
			}
		}
		if !found {
			t.Fatal("Couldn't find torrent %s in %v", r.Hash, s.ListTorrents())
		}
	}
}
