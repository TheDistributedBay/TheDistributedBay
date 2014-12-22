/*
This package mainly implements functions to encrypt/sign certain things and to smoosh stuff
into strings
*/
package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha512"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/TheDistributedBay/TheDistributedBay/database"
)

func HashTorrent(t *database.Torrent) string {
	h := sha512.New()
	binary.Write(h, binary.LittleEndian, t.PublicKey)
	binary.Write(h, binary.LittleEndian, t.MagnetLink)
	binary.Write(h, binary.LittleEndian, t.Name)
	binary.Write(h, binary.LittleEndian, t.Description)
	binary.Write(h, binary.LittleEndian, t.CategoryID)
	binary.Write(h, binary.LittleEndian, t.CreatedAt)
	for _, tag := range t.Tags {
		binary.Write(h, binary.LittleEndian, tag)
	}
	return hex.EncodeToString(h.Sum(nil))
}

type encodedKey struct {
	Curve string
	X, Y  *big.Int
}

func StringifyKey(k *ecdsa.PublicKey) (string, error) {
	if k.Curve != elliptic.P521() {
		panic("Incorrect curve in use")
	}

	e := encodedKey{"ecdsa:P521", k.X, k.Y}
	b, err := json.Marshal(e)
	if err != nil {
		return "", err
	}
	return string(b), err
}

func LoadKey(p string) (*ecdsa.PublicKey, error) {
	t := encodedKey{}
	err := json.Unmarshal([]byte(p), &t)
	if err != nil {
		return nil, err
	}
	if t.Curve != "ecdsa:P521" {
		return nil, errors.New("unrecognized key type :" + t.Curve)
	}

	k := &ecdsa.PublicKey{}
	k.Curve = elliptic.P521()
	k.X = t.X
	k.Y = t.Y
	return k, nil
}

type encodedSignature struct {
	R, S *big.Int
}

func Sign(data string, k *ecdsa.PrivateKey) (string, error) {
	sig := encodedSignature{}
	var err error
	sig.R, sig.S, err = ecdsa.Sign(rand.Reader, k, []byte(data))
	if err != nil {
		return "", err
	}

	b, err := json.Marshal(sig)
	return string(b), err
}

func Verify(data string, signature string, key string) error {
	sig := encodedSignature{}
	pk, err := LoadKey(key)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(signature), &sig)
	if err != nil {
		return err
	}

	ok := ecdsa.Verify(pk, []byte(data), sig.R, sig.S)
	if !ok {
		return errors.New("invalid signature")
	}
	return nil
}

func VerifyTorrent(t *database.Torrent) error {
	h := HashTorrent(t)
	if h != t.Hash {
		return errors.New(fmt.Sprintf("mutated hash %s vs %s", h, t.Hash))
	}

	return Verify(h, t.Signature, t.PublicKey)
}
