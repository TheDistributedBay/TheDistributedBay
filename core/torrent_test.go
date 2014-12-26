package core

import (
	"encoding/hex"
	"testing"
	"time"
)

func TestTorrentCreationAndVerify(t *testing.T) {
	t1 := CreateTorrent([]byte("ml"), "n", "d", "cid", time.Now(), []string{"foo"}, 0, 0, 0, 0)
	err := t1.VerifyTorrent()
	if err != nil {
		t.Fatal(err)
	}
	t2 := CreateTorrent([]byte("ml"), "n", "d", "cid", time.Now(), []string{"2"}, 0, 0, 0, 0)

	if t1.Hash == t2.Hash {
		t.Fatal("Hashes are identical")
	}
}

func TestFailingVerification(t *testing.T) {
	t1 := CreateTorrent([]byte("ml"), "n", "d", "cid", time.Now(), []string{"foo"}, 0, 0, 0, 0)
	t1.Hash = "bad hash"
	err := t1.VerifyTorrent()
	if err == nil {
		t.Fatal("Verified bad Hash!")
	}
}

func TestCalculateHash(t *testing.T) {
	t1 := Torrent{
		InfoHash:    []byte("test"),
		Name:        "test name",
		Description: "test desc",
		CategoryID:  1,
		CreatedAt:   time.Unix(0, 0),
	}
	hash := hashTorrent(&t1)
	t1.CalculateHash()
	if hash != t1.Hash {
		t.Fatal("Calculated hash does not match set hash!", hash, t1.Hash)
	}
}

func TestHashTorrent(t *testing.T) {
	t1 := CreateTorrent([]byte("ml"), "n", "d", "cid", time.Unix(0, 0), []string{"foo"}, 0, 0, 0, 0)
	hash := hashTorrent(t1)
	if hash != "pbJJqGCrl_mcG_NabBNu5fUNZgj5MZDJqBZKfaR6aSA=" {
		t.Fatal("Hashing didn't match stored value.", hash)
	}
}

func TestTorrentCategory(t *testing.T) {
	m := map[string]uint8{
		"All":         0,
		"Anime":       1,
		"Software":    2,
		"Games":       3,
		"Adult":       4,
		"Movies":      5,
		"Music":       6,
		"Other":       7,
		"Series & TV": 8,
		"Books":       9,
	}
	for k, v := range m {
		tor := Torrent{CategoryID: v}
		if tor.Category() != k {
			t.Fatal("Torrent Category conversion failed!", k, v)
		}

		if CategoryToId(k) != v {
			t.Fatal("Failed to convert category to ID", k, v)
		}
	}
	// Test unknown edge cases
	tor := Torrent{CategoryID: 100}
	if tor.Category() != "Unknown" {
		t.Fatal("Torrent Category conversion failed! Unknown 100")
	}

	if CategoryToId("Blah") != 0 {
		t.Fatal("Failed to convert category Blah to ID 0")
	}
}

func TestTorrentNiceInfoHash(t *testing.T) {
	hashString := "0123456789abcdef"
	infoHash, _ := hex.DecodeString(hashString)
	t1 := CreateTorrent(infoHash, "test", "test", "Anime", time.Now(), []string{}, 0, 0, 0, 0)
	if t1.NiceInfoHash() != hashString {
		t.Fatal("Failed to export nice info hash.", t1.NiceInfoHash(), hashString)
	}
}
func TestTorrentMagnetLink(t *testing.T) {
	hashString := "0123456789abcdef"
	infoHash, _ := hex.DecodeString(hashString)
	t1 := CreateTorrent(infoHash, "test ;", "test", "Anime", time.Now(), []string{}, 0, 0, 0, 0)
	if t1.MagnetLink() != "magnet:?xt=urn:btih:0123456789abcdef&dn=test+%3B&tr=udp://open.demonii.com:1337/announce&tr=udp://tracker.publicbt.com:80/announce&tr=udp://tracker.istole.it:80/announce" {
		t.Fatal("Failed to export magnet link.", t1.MagnetLink(), hashString)
	}
}

func BenchmarkDefaultCreate(b *testing.B) {
	c := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CreateTorrent([]byte("ml"), "n", "d", "cid", c, []string{"foo"}, 0, 0, 0, 0)
	}
}

func BenchmarkDefaultVerify(b *testing.B) {
	c := time.Now()
	t := CreateTorrent([]byte("ml"), "n", "d", "cid", c, []string{"foo"}, 0, 0, 0, 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		t.VerifyTorrent()
	}
}

func BenchmarkDefaultMagnetLink(b *testing.B) {
	c := time.Now()
	t := CreateTorrent([]byte("ml"), "n", "d", "cid", c, []string{"foo"}, 0, 0, 0, 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		t.MagnetLink()
	}
}
