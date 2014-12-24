package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"
)

func TestUtil(t *testing.T) {
	k, err := NewKey()
	if err != nil {
		t.Fatal(err)
	}
	ek := EncodeKey(&k.PublicKey)
	lk, err := LoadKey(ek)
	if err != nil {
		t.Fatal(err)
	}
	if *lk != k.PublicKey {
		t.Fatal("Keys don't match")
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
