package database

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/jmhodges/levigo"

	"github.com/TheDistributedBay/TheDistributedBay/client"
	"github.com/TheDistributedBay/TheDistributedBay/core"
)

type TorrentDB struct {
	db      *levigo.DB
	writers []core.TorrentWriter
	lock    *sync.RWMutex
}

func NewTorrentDB(dir string) (*TorrentDB, error) {
	opts := levigo.NewOptions()
	opts.SetCache(levigo.NewLRUCache(10 << 20))
	opts.SetCreateIfMissing(true)
	db, err := levigo.Open(dir, opts)
	if err != nil {
		return nil, err
	}
	return &TorrentDB{db, nil, &sync.RWMutex{}}, nil
}

func (db *TorrentDB) NumTorrents() int {
	db.lock.RLock()
	defer db.lock.RUnlock()
	ro := levigo.NewReadOptions()
	ro.SetFillCache(false)
	it := db.db.NewIterator(ro)
	count := 0
	for it.SeekToFirst(); it.Valid(); it.Next() {
		count += 1
	}
	return count
}

func (db *TorrentDB) Get(hash string) (*core.Torrent, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()
	ro := levigo.NewReadOptions()
	data, err := db.db.Get(ro, []byte("t"+hash))
	if err != nil {
		return nil, err
	}

	t := core.Torrent{}
	err = json.Unmarshal(data, &t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (db *TorrentDB) Add(t *core.Torrent) error {
	db.lock.Lock()
	defer db.lock.Unlock()
	data, err := json.Marshal(t)
	wo := levigo.NewWriteOptions()
	err = db.db.Put(wo, []byte("t"+t.Hash), data)
	if err != nil {
		return err
	}

	for _, w := range db.writers {
		w.NewTorrent(t)
	}
	return nil
}

func (db *TorrentDB) AddSignature(s *core.Signature) {
	db.lock.Lock()
	defer db.lock.Unlock()
	data, err := json.Marshal(s)
	wo := levigo.NewWriteOptions()
	err = db.db.Put(wo, []byte("s"+s.Hash()), data)
	if err != nil {
		log.Print(err)
	}

	for _, w := range db.writers {
		w.NewSignature(s)
	}
}

func (db *TorrentDB) List() []string {
	db.lock.RLock()
	defer db.lock.RUnlock()
	ro := levigo.NewReadOptions()
	ro.SetFillCache(false)
	it := db.db.NewIterator(ro)
	ts := make([]string, 0)
	for it.SeekToFirst(); it.Valid(); it.Next() {
		k := it.Key()
		if k[0] == 't' {
			ts = append(ts, string(it.Key()))
		}
	}
	return ts
}

func (db *TorrentDB) AddClient(w core.TorrentWriter) {
	db.lock.Lock()
	defer db.lock.Unlock()
	ww := client.New(w)
	db.writers = append(db.writers, ww)
	ro := levigo.NewReadOptions()
	ro.SetFillCache(false)
	it := db.db.NewIterator(ro)
	for it.SeekToFirst(); it.Valid(); it.Next() {
		k := it.Key()
		if k[0] == 't' {
			var t core.Torrent
			err := json.Unmarshal(it.Value(), &t)
			if err != nil {
				log.Print(err)
				continue
			}
			w.NewTorrent(&t)
		} else if k[0] == 's' {
			var s core.Signature
			err := json.Unmarshal(it.Value(), &s)
			if err != nil {
				log.Print(err)
				continue
			}
			w.NewSignature(&s)
		}
	}
}
