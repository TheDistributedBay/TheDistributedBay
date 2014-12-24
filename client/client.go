// Implements a buffered torrent client
package client

import (
	"sync"

	"github.com/TheDistributedBay/TheDistributedBay/core"
)

type BufferedClient struct {
	c  core.TorrentWriter
	t  []*core.Torrent
	s  []*core.Signature
	lt *sync.Mutex
	ls *sync.Mutex
	ct *sync.Cond
	cs *sync.Cond
}

func New(c core.TorrentWriter) *BufferedClient {
	lt := &sync.Mutex{}
	ls := &sync.Mutex{}
	ct := sync.NewCond(lt)
	cs := sync.NewCond(ls)
	b := &BufferedClient{c, nil, nil, lt, ls, ct, cs}
	go b.torrents()
	go b.sigs()
	return b
}

func (b *BufferedClient) torrents() {
	b.lt.Lock()
	for {
		for len(b.t) > 0 {
			t := b.t[0]
			b.t = b.t[1:]
			b.lt.Unlock()
			b.c.NewTorrent(t)
			b.lt.Lock()
		}
		b.ct.Wait()
	}
}

func (b *BufferedClient) sigs() {
	b.ls.Lock()
	for {
		for len(b.s) > 0 {
			s := b.s[0]
			b.s = b.s[1:]
			b.ls.Unlock()
			b.c.NewSignature(s)
			b.ls.Lock()
		}
		b.cs.Wait()
	}
}

func (b *BufferedClient) NewTorrent(t *core.Torrent) {
	b.lt.Lock()
	b.t = append(b.t, t)
	b.lt.Unlock()
	b.ct.Broadcast()
}

func (b *BufferedClient) NewSignature(s *core.Signature) {
	b.ls.Lock()
	b.s = append(b.s, s)
	b.ls.Unlock()
	b.cs.Broadcast()
}
