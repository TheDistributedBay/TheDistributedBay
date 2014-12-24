package database

import (
	"errors"
	"sync"
)

type TorrentWriter interface {
	NewTorrent(t *Torrent)
	NewSignature(s *Signature)
}

type Database interface {
	Get(hash string) (*Torrent, error)
	Add(t *Torrent)
	AddSignature(s *Signature)
	List() []string
	AddClient(w TorrentWriter)
}
type TorrentDB struct {
	torrents   map[string]*Torrent
	signatures map[string][]*Signature
	writers    []TorrentWriter
	lock       *sync.RWMutex
}

func NewTorrentDB() *TorrentDB {
	t := make(map[string]*Torrent)
	s := make(map[string][]*Signature)
	return &TorrentDB{t, s, nil, &sync.RWMutex{}}
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

	wg := &sync.WaitGroup{}
	wg.Add(len(db.writers))
	for _, w := range db.writers {
		go func() {
			w.NewTorrent(t)
			wg.Done()
		}()
	}
	wg.Wait()
}

func (db *TorrentDB) AddSignature(s *Signature) {
	db.lock.Lock()
	defer db.lock.Unlock()
	for _, t := range s.ListTorrents() {
		db.signatures[t] = append(db.signatures[t], s)
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(db.writers))
	for _, w := range db.writers {
		go func() {
			w.NewSignature(s)
			wg.Done()
		}()
	}
	wg.Wait()
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
