/*
This package mainly implements functions to encrypt/sign certain things and to smoosh stuff
into strings
*/
package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"math/big"
)

func NewKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
}

type EncodedKey struct {
	Curve string
	X, Y  *big.Int
}

func EncodeKey(k *ecdsa.PublicKey) *EncodedKey {
	if k.Curve != elliptic.P521() {
		panic("Incorrect curve in use")
	}

	return &EncodedKey{"ecdsa:P521", k.X, k.Y}
}

func LoadKey(e *EncodedKey) (*ecdsa.PublicKey, error) {
	if e.Curve != "ecdsa:P521" {
		return nil, errors.New("unrecognized key type :" + e.Curve)
	}
	k := &ecdsa.PublicKey{}
	k.Curve = elliptic.P521()
	k.X = e.X
	k.Y = e.Y
	return k, nil
}
