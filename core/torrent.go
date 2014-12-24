package core

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"html"
	"io"
	"time"
)

var Trackers = [...]string{
	"udp://open.demonii.com:1337/announce",
	"udp://tracker.publicbt.com:80/announce",
	"udp://tracker.istole.it:80/announce",
}

type Torrent struct {
	// Hash of everything in this struct
	Hash string
	// Torrent information
	InfoHash    []byte
	Name        string
	Description string
	CategoryID  uint8
	CreatedAt   time.Time
	Tags        []string
}

func (t Torrent) Category() string {
	switch t.CategoryID {
	case 0:
		return "All"
	case 1:
		return "Anime"
	case 2:
		return "Software"
	case 3:
		return "Games"
	case 4:
		return "Adult"
	case 5:
		return "Movies"
	case 6:
		return "Music"
	case 7:
		return "Other"
	case 8:
		return "Series & TV"
	case 9:
		return "Books"
	}
	return "Unknown"
}
func (t Torrent) NiceInfoHash() string {
	return hex.EncodeToString(t.InfoHash)
}
func (t Torrent) MagnetLink() string {
	infoHash := t.NiceInfoHash()
	name := html.EscapeString(t.Name)
	magnet := fmt.Sprintf("magnet:?xt=urn:btih:%s&dn=%s", infoHash, name)

	for _, tracker := range Trackers {
		magnet += "&tr=" + tracker
	}

	return magnet
}

func CreateTorrent(infoHash []byte, name, description string, categoryid string, createdAt time.Time, tags []string) *Torrent {
	t := &Torrent{"", infoHash, name, description, 1, createdAt, tags}
	t.CalculateHash()
	return t
}

func (t *Torrent) CalculateHash() {
	t.Hash = hashTorrent(t)
}

func hashTorrent(t *Torrent) string {
	h := sha256.New()
	io.WriteString(h, (string)(t.InfoHash))
	io.WriteString(h, t.Name)
	io.WriteString(h, t.Description)
	binary.Write(h, binary.LittleEndian, t.CategoryID)
	binary.Write(h, binary.LittleEndian, t.CreatedAt.Unix())
	for _, tag := range t.Tags {
		io.WriteString(h, tag)
	}
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

func (t *Torrent) VerifyTorrent() error {
	h := hashTorrent(t)
	if h != t.Hash {
		return errors.New(fmt.Sprintf("mutated hash %s vs %s", h, t.Hash))
	}
	return nil
}
