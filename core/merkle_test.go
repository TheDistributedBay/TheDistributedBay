package core

import (
	"fmt"
	"testing"
	"time"

	"github.com/TheDistributedBay/TheDistributedBay/crypto"
)

func TestMerkle(t *testing.T) {
	k, err := crypto.NewKey()
	if err != nil {
		t.Fatal(err)
	}
	ts := make([]*Torrent, 0)

	for i := 0; i < 5; i++ {
		r := CreateTorrent([]byte(fmt.Sprint(i)), "n", "d", "cat", time.Now(), []string{"tag"})
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
