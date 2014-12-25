package network

import (
	"io"
	"log"
	"sync"

	"github.com/TheDistributedBay/TheDistributedBay/core"
)

type ConnectionHandler struct {
	t           *Transcoder
	db          core.Database
	torrentList map[string]struct{}
	lock        *sync.RWMutex
	cond        *sync.Cond
}

func NewConnectionHandler(t *Transcoder, db core.Database) *ConnectionHandler {
	h := &ConnectionHandler{t, db, make(map[string]struct{}), &sync.RWMutex{}, nil}
	h.cond = sync.NewCond(h.lock)

	go h.shovel()
	go h.shovelout()
	return h
}

func (h *ConnectionHandler) SeenTorrent(hash string) bool {
	h.lock.RLock()
	defer h.lock.RUnlock()
	_, seen := h.torrentList[hash]
	return seen
}

func (h *ConnectionHandler) RecordTorrent(hash string) {
	h.lock.Lock()
	defer h.lock.Unlock()
	h.torrentList[hash] = struct{}{}
}

func (h *ConnectionHandler) shovelout() {
	log.Print("Building torrent list")
	m := make(map[string]struct{})
	c := make(chan string)
	go h.db.GetTorrents(c)
	for h := range c {
		m[h] = struct{}{}
	}
	log.Print("Done building")
	go h.t.Write(Message{"TorrentList", m, nil})

	h.lock.Lock()
	log.Print("Waiting for other sides torrent list")
	h.cond.Wait()
	log.Print("Received")
	h.lock.Unlock()

	// Shove slow torrents and fast torrents
	f := make(chan *core.Torrent, 2)
	s := make(chan string)
	go h.db.AddTorrentClient(f)
	go h.db.GetTorrents(s)
	fOpen := true
	sOpen := true
	for {
		if !sOpen && !fOpen {
			log.Print("restarting slow & fast sync")
			go h.db.AddTorrentClient(f)
			go h.db.GetTorrents(s)
		}
		var t *core.Torrent = nil
		var err error
		var hash string
		select {
		case t, fOpen = <-f:
			if !fOpen {
				log.Print("too slow to keep fast path open")
				f = make(chan *core.Torrent, 2)
				continue
			}
			if h.SeenTorrent(t.Hash) {
				continue
			}
		case hash, sOpen = <-s:
			if !sOpen {
				log.Print("slow sync finished")
				s = make(chan string)
				continue
			}
			if h.SeenTorrent(hash) {
				continue
			}
			t, err = h.db.Get(hash)
			if err != nil {
				log.Printf("Error fetching %v : %v", hash, err)
				continue
			}
		}
		log.Print("sending torrent")
		h.RecordTorrent(t.Hash)
		h.t.Write(Message{"Torrents", nil, []*core.Torrent{t}})
	}
}

// Shovels torrents into the db
func (h *ConnectionHandler) shovel() {
	for {
		msg, err := h.t.Read()
		if err != nil {
			if err != io.EOF {
				log.Printf("Error for %v : %v", h.t, err)
			}
			return
		}
		switch msg.Code {
		case "TorrentList":
			// Once we know what the other side doesn't have, start dumping
			h.lock.Lock()
			h.torrentList = msg.Torrents
			h.lock.Unlock()
			h.cond.Broadcast()
		case "Torrents":
			for _, t := range msg.Data {
				log.Print("torrent recieved")
				err := t.VerifyTorrent()
				if err != nil {
					log.Print(err)
					continue
				}
				log.Print("torrent verified")
				h.RecordTorrent(t.Hash)
				err = h.db.Add(t)
				if err != nil {
					log.Print(err)
				}
			}
		}
	}
}

func (h *ConnectionHandler) Close() error {
	return h.t.Close()
}
