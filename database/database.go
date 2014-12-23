package database

import (
	"errors"
	"sync"
)

type TorrentWriter interface {
	NewTorrent(t *Torrent)
}

type Database interface {
	Get(hash string) (*Torrent, error)
	Add(t *Torrent)
	List() []string
	AddClient(w TorrentWriter)
}
type TorrentDB struct {
	torrents map[string]*Torrent
	writers  []TorrentWriter
	lock     *sync.RWMutex
	r        chan struct{}
}

func NewTorrentDB() *TorrentDB {
	t := make(map[string]*Torrent)
	db := &TorrentDB{t, nil, &sync.RWMutex{}, make(chan struct{}, 10)}
	for i := 0; i < 10; i++ {
		db.r <- struct{}{}
	}
	return db
}

func (db *TorrentDB) NumTorrents() int {
	db.lock.RLock()
	defer db.lock.RUnlock()
	return len(db.torrents)
}

func (db *TorrentDB) Get(hash string) (*Torrent, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()
	torrent, ok := db.torrents[hash]
	if !ok {
		return nil, errors.New("No such hash stored")
	}
	return torrent, nil
}

func (db *TorrentDB) Add(t *Torrent) {
	db.lock.Lock()
	defer db.lock.Unlock()
	_, ok := db.torrents[t.Hash]
	if ok {
		return
	}
	db.torrents[t.Hash] = t
	for _, w := range db.writers {
		<-db.r
		go func() {
			w.NewTorrent(t)
			db.r <- struct{}{}
		}()
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

func (db *TorrentDB) AddClient(w TorrentWriter) {
	db.lock.Lock()
	defer db.lock.Unlock()
	db.writers = append(db.writers, w)
	for _, t := range db.torrents {
		go w.NewTorrent(t)
	}
}
