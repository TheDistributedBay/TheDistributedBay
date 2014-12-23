package crypto

import (
	"testing"
	"time"
)

func TestTorrentCreationAndVerify(t *testing.T) {
	k, err := NewKey()
	if err != nil {
		t.Fatal(err)
	}

	t1, err := CreateTorrent(k, "ml", "n", "d", "cid", time.Now(), []string{"foo"})
	if err != nil {
		t.Fatal(err)
	}

	err = VerifyTorrent(t1)
	if err != nil {
		t.Fatal(err)
	}
	t2, err := CreateTorrent(k, "ml", "n", "d", "cid", time.Now(), []string{"2"})
	if err != nil {
		t.Fatal(err)
	}

	if t1.Hash == t2.Hash {
		t.Fatal("Hashes are identical")
	}
}
