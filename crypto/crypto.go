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
	"io"
	"math/big"
	"time"

	"github.com/TheDistributedBay/TheDistributedBay/database"
)

func NewKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
}

func CreateTorrent(k *ecdsa.PrivateKey, magnetlink, name, description string, categoryid string, createdAt time.Time, tags []string) (*database.Torrent, error) {
	t := &database.Torrent{"", "", "", magnetlink, name, description, 1, createdAt, tags}
	var err error
	t.PublicKey, err = StringifyKey(&k.PublicKey)
	if err != nil {
		return nil, err
	}
	t.Hash = HashTorrent(t)
	t.Signature, err = Sign(t.Hash, k)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func HashTorrent(t *database.Torrent) string {
	h := sha512.New()
	io.WriteString(h, t.PublicKey)
	io.WriteString(h, t.MagnetLink)
	io.WriteString(h, t.Name)
	io.WriteString(h, t.Description)
	binary.Write(h, binary.LittleEndian, t.CategoryID)
	binary.Write(h, binary.LittleEndian, t.CreatedAt.Unix())
	for _, tag := range t.Tags {
		io.WriteString(h, tag)
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
