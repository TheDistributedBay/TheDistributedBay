package tls

import (
	"io"
	"io/ioutil"
	"testing"
)

func TestTLS(t *testing.T) {
	l, err := Listen("localhost:8765")
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		c, err := l.Accept()
		if err != nil {
			t.Fatal(err)
		}
		_, err = io.WriteString(c, "hello")
		if err != nil {
			t.Fatal(err)
		}
		c.Close()
	}()
	c, err := Dial("localhost:8765")
	if err != nil {
		t.Fatal(err)
	}

	b, err := ioutil.ReadAll(c)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "hello" {
		t.Fatal("unexpected response:", string(b))
	}
}
