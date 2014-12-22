package database

import (
    "errors"
)

type Torrentdb struct {
	torrents map[string]*Torrent
}

func (t *Torrentdb) Get(hash string) (*Torrent, error) {
	torrent, ok := t.torrents[hash]
	if !ok {
		return nil, errors.New("No such hash stored")
	}
	return torrent, nil
}

func (t *Torrentdb) Add(r *Torrent) {
	t.torrents[r.Hash] = r
}

func (t *Torrentdb) List() []string {
	ts := make([]string, 0, len(t.torrents))
	for _, r := range t.torrents {
		ts = append(ts, r.Hash)
	}
	return ts
}
