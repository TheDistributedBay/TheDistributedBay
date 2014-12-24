package core

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha512"
	"errors"
	"io"
	"math/big"

	"github.com/TheDistributedBay/TheDistributedBay/crypto"
)

type Signature struct {
	Key  *crypto.EncodedKey
	R, S *big.Int
	M    *MerkleNode
}

type MerkleNode struct {
	hash string
	l    *MerkleNode
	r    *MerkleNode
}

func SignTorrents(k *ecdsa.PrivateKey, ts []*Torrent) (*Signature, error) {
	m := buildMerkle(ts)
	R, S, err := ecdsa.Sign(rand.Reader, k, []byte(m.hash))
	if err != nil {
		return nil, err
	}
	ek := crypto.EncodeKey(&k.PublicKey)
	return &Signature{ek, R, S, m}, nil
}

func (s *Signature) VerifySignature() error {
	err := verifyMerkle(s.M)
	if err != nil {
		return err
	}
	pk, err := crypto.LoadKey(s.Key)
	if err != nil {
		return err
	}
	ok := ecdsa.Verify(pk, []byte(s.M.hash), s.R, s.S)
	if !ok {
		return errors.New("Invalid signature")
	}
	return nil
}

func (s *Signature) ListTorrents() []string {
	t := make([]string, 0)
	listMerkle(s.M, &t)
	return t
}

func hash(a, b string) string {
	h := sha512.New()
	io.WriteString(h, a)
	io.WriteString(h, b)
	return string(h.Sum(nil))
}

func buildMerkle(ts []*Torrent) *MerkleNode {
	if len(ts) == 0 {
		panic("This should never happen")
	}
	if len(ts) == 1 {
		return &MerkleNode{ts[0].Hash, nil, nil}
	}
	if len(ts) >= 2 {
		mid := len(ts) / 2
		l := buildMerkle(ts[:mid])
		r := buildMerkle(ts[mid:])
		h := hash(l.hash, r.hash)
		return &MerkleNode{h, l, r}
	}
	panic("Should be impossible")
}

func verifyMerkle(m *MerkleNode) error {
	if m.r == nil || m.l == nil {
		return nil
	}
	if hash(m.l.hash, m.r.hash) != m.hash {
		return errors.New("Invalid signature")
	}
	err := verifyMerkle(m.l)
	if err != nil {
		return err
	}
	err = verifyMerkle(m.r)
	if err != nil {
		return err
	}
	return nil
}

func listMerkle(m *MerkleNode, r *[]string) {
	leaf := true
	if m.r != nil {
		listMerkle(m.r, r)
		leaf = false
	}
	if m.l != nil {
		listMerkle(m.l, r)
		leaf = false
	}
	if leaf {
		*r = append(*r, m.hash)
	}
}
