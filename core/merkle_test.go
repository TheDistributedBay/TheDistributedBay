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
		r := CreateTorrent([]byte(fmt.Sprint(i)), "n", "d", "cat", time.Now(), []string{"tag"}, 0, 0, 0, 0)
		ts = append(ts, r)
	}

	s, err := SignTorrents(k, ts)
	if err != nil {
		t.Fatal(err)
	}

	s2, err := SignTorrents(k, ts[:2])
	if err != nil {
		t.Fatal(err)
	}

	if s2.Hash() == s.Hash() {
		t.Fatal("Hashes should be different %v v %v", s2.Hash(), s.Hash())
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

	// Make sure merkle verification fails.
	s.R.SetUint64(1)
	if s.VerifySignature() == nil {
		t.Fatal("Changing sig should cause failures")
	}

	s, err = SignTorrents(k, ts)
	if err != nil {
		t.Fatal(err)
	}

	s.Key.Curve = "badkey"
	if s.VerifySignature() == nil {
		t.Fatal("Changing key should cause failures")
	}

	s, err = SignTorrents(k, ts)
	if err != nil {
		t.Fatal(err)
	}

	s.M.l.r.hash = "badhash"
	if s.VerifySignature() == nil {
		t.Fatal("Changing key should cause failures")
	}

	s, err = SignTorrents(k, ts)
	if err != nil {
		t.Fatal(err)
	}

	s.M.r.l.hash = "badhash"
	if s.VerifySignature() == nil {
		t.Fatal("Changing key should cause failures")
	}
}
