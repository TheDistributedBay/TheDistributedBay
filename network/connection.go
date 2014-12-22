package network

import (
	"encoding/json"
	"io"

	"github.com/TheDistributedBay/TheDistributedBay/database"
)

type Message struct {
	code     string
	torrents map[string]struct{}
	data     []*database.Torrent
}

type Connection interface {
	Read() (Message, error)
	Write(m Message) error
	Close() error
}

type Transcoder struct {
	c   io.ReadWriteCloser
	enc *json.Encoder
	dec *json.Decoder
}

// A transcoder takes messages and coverts them into a write format which is
// basically just JSON.
func NewTranscoder(c io.ReadWriteCloser) *Transcoder {
	t := &Transcoder{}
	t.c = c
	t.enc = json.NewEncoder(t.c)
	t.dec = json.NewDecoder(t.c)
	return t
}

func (t Transcoder) Write(m Message) error {
	return t.enc.Encode(m)
}

func (t Transcoder) Read() (Message, error) {
	var m Message
	err := t.dec.Decode(&m)
	return m, err
}

func (t Transcoder) Close() error {
	return t.c.Close()
}
