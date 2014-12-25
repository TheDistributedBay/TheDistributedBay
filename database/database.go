package database

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/jmhodges/levigo"

	"github.com/TheDistributedBay/TheDistributedBay/core"
)

type TorrentDB struct {
	db      *levigo.DB
	writers []chan *core.Torrent
	lock    *sync.RWMutex
}

func NewTorrentDB(dir string) (*TorrentDB, error) {
	opts := levigo.NewOptions()
	opts.SetCache(levigo.NewLRUCache(10 << 20))
	opts.SetCreateIfMissing(true)
	defer opts.Close()
	db, err := levigo.Open(dir, opts)
	if err != nil {
		return nil, err
	}
	return &TorrentDB{db, nil, &sync.RWMutex{}}, nil
}

// Slow path designed for use by clients which can't keep up or to resync initially.
// Note that this channel will be closed after all torrents in a database are read.
func (db *TorrentDB) GetTorrents(c chan string) {
	ro := levigo.NewReadOptions()
	ro.SetFillCache(false)
	defer ro.Close()
	it := db.db.NewIterator(ro)
	for it.SeekToFirst(); it.Valid(); it.Next() {
		k := it.Key()
		if k[0] == 't' {
			c <- string(k[1:])
		}
	}
	close(c)
}

func (db *TorrentDB) Get(hash string) (*core.Torrent, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()
	ro := levigo.NewReadOptions()
	defer ro.Close()
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
	defer wo.Close()
	err = db.db.Put(wo, []byte("t"+t.Hash), data)
	if err != nil {
		return err
	}

	bad := make([]int, 0)
	for i, w := range db.writers {
		select {
		case w <- t:
		default:
			close(w)
			bad = append(bad, i)
		}
	}

	for c, i := range bad {
		i = i - c
		db.writers = append(db.writers[:i], db.writers[i+1:]...)
	}
	return nil
}

func (db *TorrentDB) AddSignature(s *core.Signature) {
	db.lock.Lock()
	defer db.lock.Unlock()
	data, err := json.Marshal(s)
	wo := levigo.NewWriteOptions()
	defer wo.Close()
	err = db.db.Put(wo, []byte("s"+s.Hash()), data)
	if err != nil {
		log.Print(err)
	}
}

// Fast path for getting torrents, if the channel blocks we close it.
func (db *TorrentDB) AddTorrentClient(c chan *core.Torrent) {
	db.lock.Lock()
	defer db.lock.Unlock()
	db.writers = append(db.writers, c)
}
