package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
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

func BenchmarkCreate(b *testing.B) {
	c := time.Now()
	k, _ := NewKey()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CreateTorrent(k, "ml", "n", "d", "cid", c, []string{"foo"})
	}
}

func BenchmarkVerify(b *testing.B) {
	c := time.Now()
	k, _ := NewKey()
	t, _ := CreateTorrent(k, "ml", "n", "d", "cid", c, []string{"foo"})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		VerifyTorrent(t)
	}
}

func BenchmarkRawCreateP521(b *testing.B) {
	data := "foo"
	k, _ := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ecdsa.Sign(rand.Reader, k, []byte(data))
	}
}

func BenchmarkRawCreateP224(b *testing.B) {
	data := "foo"
	k, _ := ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ecdsa.Sign(rand.Reader, k, []byte(data))
	}
}

func BenchmarkRawVerifyP521(b *testing.B) {
	data := "foo"
	k, _ := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	R, S, _ := ecdsa.Sign(rand.Reader, k, []byte(data))
	pk := &k.PublicKey
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ecdsa.Verify(pk, []byte(data), R, S)
	}
}

func BenchmarkRawVerifyP224(b *testing.B) {
	data := "foo"
	k, _ := ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	R, S, _ := ecdsa.Sign(rand.Reader, k, []byte(data))
	pk := &k.PublicKey
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ecdsa.Verify(pk, []byte(data), R, S)
	}
}
