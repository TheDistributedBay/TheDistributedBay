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
	"errors"
	"fmt"
	"io"
	"math/big"
	"time"

	"github.com/TheDistributedBay/TheDistributedBay/database"
)

func NewKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
}

func CreateTorrent(k *ecdsa.PrivateKey, magnetlink, name, description string, categoryid string, createdAt time.Time, tags []string) (*database.Torrent, error) {
	t := &database.Torrent{"", "", "", magnetlink, name, description, 1, createdAt, tags}
	t.Hash = hashTorrent(t)
	return t, nil
}

func hashTorrent(t *database.Torrent) string {
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

type EncodedKey struct {
	Curve string
	X, Y  *big.Int
}

func encodeKey(k *ecdsa.PublicKey) *EncodedKey {
	if k.Curve != elliptic.P224() {
		panic("Incorrect curve in use")
	}

	return &EncodedKey{"ecdsa:P224", k.X, k.Y}
}

func loadKey(e *EncodedKey) (*ecdsa.PublicKey, error) {
	if e.Curve != "ecdsa:P224" {
		return nil, errors.New("unrecognized key type :" + e.Curve)
	}
	k := &ecdsa.PublicKey{}
	k.Curve = elliptic.P224()
	k.X = e.X
	k.Y = e.Y
	return k, nil
}

func VerifyTorrent(t *database.Torrent) error {
	h := hashTorrent(t)
	if h != t.Hash {
		return errors.New(fmt.Sprintf("mutated hash %s vs %s", h, t.Hash))
	}
	return nil
}
