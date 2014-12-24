package database

import (
	"errors"
	"sync"

	"github.com/TheDistributedBay/TheDistributedBay/client"
	"github.com/TheDistributedBay/TheDistributedBay/core"
)

type TorrentDB struct {
	torrents   map[string]*core.Torrent
	signatures map[string][]*core.Signature
	writers    []core.TorrentWriter
	lock       *sync.RWMutex
}

func NewTorrentDB() *TorrentDB {
	t := make(map[string]*core.Torrent)
	s := make(map[string][]*core.Signature)
	return &TorrentDB{t, s, nil, &sync.RWMutex{}}
}

func (db *TorrentDB) NumTorrents() int {
	db.lock.RLock()
	defer db.lock.RUnlock()
	return len(db.torrents)
}

func (db *TorrentDB) Get(hash string) (*core.Torrent, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()
	torrent, ok := db.torrents[hash]
	if !ok {
		return nil, errors.New("No such hash stored")
	}
	return torrent, nil
}

func (db *TorrentDB) Add(t *core.Torrent) {
	db.lock.Lock()
	defer db.lock.Unlock()
	_, ok := db.torrents[t.Hash]
	if ok {
		return
	}
	db.torrents[t.Hash] = t

	for _, w := range db.writers {
		w.NewTorrent(t)
	}
}

func (db *TorrentDB) AddSignature(s *core.Signature) {
	db.lock.Lock()
	defer db.lock.Unlock()
	for _, t := range s.ListTorrents() {
		db.signatures[t] = append(db.signatures[t], s)
	}

	for _, w := range db.writers {
		w.NewSignature(s)
	}
}

func (db *TorrentDB) List() []string {
	db.lock.RLock()
	defer db.lock.RUnlock()
	ts := make([]string, 0, len(db.torrents))
	for _, r := range db.torrents {
		ts = append(ts, r.Hash)
	}
	return ts
}

func (db *TorrentDB) AddClient(w core.TorrentWriter) {
	db.lock.Lock()
	defer db.lock.Unlock()
	ww := client.New(w)
	db.writers = append(db.writers, ww)
	for _, t := range db.torrents {
		ww.NewTorrent(t)
	}
}
