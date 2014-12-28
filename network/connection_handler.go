package network

import (
	"io"
	"log"
	"sync"

	"github.com/TheDistributedBay/TheDistributedBay/core"
	"github.com/TheDistributedBay/TheDistributedBay/dbchannel"
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

	dc := dbchannel.New(h.db)
	log.Print("streaming torrents")
	for t := range dc.Output {
		if h.SeenTorrent(t.GetHash()) {
			continue
		}
		tr, err := t.GetTorrent()
		if err != nil {
			log.Printf("Error fetching %v : %v", t.GetHash(), err)
			continue
		}
		h.RecordTorrent(tr.Hash)
		h.t.Write(Message{"Torrents", nil, []*core.Torrent{tr}})
	}
	log.Print("dbchannel closed, should be impossible")
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
