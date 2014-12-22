package database

import (
	"errors"
)

type TorrentDB struct {
	torrents map[string]*Torrent
}

func NewTorrentDB() *TorrentDB {
	t := make(map[string]*Torrent)
	return &TorrentDB{t}
}

func (t *TorrentDB) Get(hash string) (*Torrent, error) {
	torrent, ok := t.torrents[hash]
	if !ok {
		return nil, errors.New("No such hash stored")
	}
	return torrent, nil
}

func (t *TorrentDB) Add(r *Torrent) {
	t.torrents[r.Hash] = r
}

func (t *TorrentDB) List() []string {
	ts := make([]string, 0, len(t.torrents))
	for _, r := range t.torrents {
		ts = append(ts, r.Hash)
	}
	return ts
}
