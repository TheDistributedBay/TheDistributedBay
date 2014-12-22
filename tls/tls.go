// A simple package to handle all the wierd tls tasks required to use TLS without caring about certs
package tls

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"math/big"
	"net"
	"time"
)

const (
	Proto = "distributedbay/1/json"
)

func GenerateEmptyConfig() (*tls.Config, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		return nil, err
	}
	number, err := rand.Int(rand.Reader, big.NewInt(0).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, err
	}
	cert := &x509.Certificate{
		SerialNumber:          number,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
	}
	certbytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &priv.PublicKey, priv)
	if err != nil {
		return nil, err
	}

	ct := tls.Certificate{[][]byte{certbytes}, priv, nil, cert}
	c := &tls.Config{InsecureSkipVerify: true}
	c.Certificates = append(c.Certificates, ct)
	c.NextProtos = []string{Proto}
	return c, nil
}

func Dial(addr string) (*wrap, error) {
	co, err := GenerateEmptyConfig()
	if err != nil {
		return nil, err
	}
	c, err := tls.Dial("tcp", addr, co)
	return Wrap(c), err
}
func Listen(addr string) (net.Listener, error) {
	co, err := GenerateEmptyConfig()
	if err != nil {
		return nil, err
	}
	c, err := tls.Listen("tcp", addr, co)
	return c, err
}

type wrap struct {
	*tls.Conn
}

func (w *wrap) Protocol() string {
	return w.ConnectionState().NegotiatedProtocol
}

func Wrap(c net.Conn) *wrap {
	return &wrap{c.(*tls.Conn)}
}
