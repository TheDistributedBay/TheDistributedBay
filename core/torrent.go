package core

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"time"
)

type Torrent struct {
	// Hash of everything in this struct
	Hash string
	// Torrent information
	MagnetLink  string
	Name        string
	Description string
	CategoryID  uint64
	CreatedAt   time.Time
	Tags        []string
}

func CreateTorrent(magnetlink, name, description string, categoryid string, createdAt time.Time, tags []string) *Torrent {
	t := &Torrent{"", magnetlink, name, description, 1, createdAt, tags}
	t.CalculateHash()
	return t
}

func (t *Torrent) CalculateHash() {
	t.Hash = hashTorrent(t)
}

func hashTorrent(t *Torrent) string {
	h := sha256.New()
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

func (t *Torrent) VerifyTorrent() error {
	h := hashTorrent(t)
	if h != t.Hash {
		return errors.New(fmt.Sprintf("mutated hash %s vs %s", h, t.Hash))
	}
	return nil
}
