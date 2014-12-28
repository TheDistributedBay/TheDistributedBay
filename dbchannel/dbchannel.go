/*
This pacakge provides a utility which will read from a database consistently,
dynamically switching from the slow channel to the fast channel when possible.
Note that the produced channel will contain duplicate hashes, that's just the
nature of the beast.
*/
package dbchannel

import (
	"github.com/TheDistributedBay/TheDistributedBay/core"
)

// This exists so that you can check if a torrent is needed before
// you force it to be loaded into memory by calling GetTorrent
type Torrent interface {
	GetHash() string
	GetTorrent() (*core.Torrent, error)
}

type wrapTorrent struct {
	t *core.Torrent
}

func (w wrapTorrent) GetHash() string { return w.t.Hash }
func (w wrapTorrent) GetTorrent() (*core.Torrent, error) {
	return w.t, nil
}

type wrapHash struct {
	h  string
	db core.Database
}

func (w wrapHash) GetHash() string {
	return w.h
}

func (w wrapHash) GetTorrent() (*core.Torrent, error) {
	return w.db.Get(w.h)
}

type DBChannel struct {
	// The channel to which information is dumped
	Output chan Torrent
	db     core.Database
	died   chan struct{}
}

func New(db core.Database) *DBChannel {
	o := make(chan Torrent, 2)
	d := make(chan struct{})
	c := &DBChannel{o, db, d}
	go c.fastLoop()
	go c.slowLoop()
	return c
}

func (d *DBChannel) fastLoop() {
	for {
		f := make(chan *core.Torrent, 10)
		go d.db.AddTorrentClient(f)
		for t := range f {
			d.Output <- wrapTorrent{t}
		}
		d.died <- struct{}{}
	}
}

func (d *DBChannel) slowLoop() {
	for {
		s := make(chan string)
		go d.db.GetTorrents(s)
		for t := range s {
			d.Output <- wrapHash{t, d.db}
		}
		<-d.died
	}
}
