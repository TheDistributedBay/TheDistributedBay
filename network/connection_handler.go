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
}

func NewConnectionHandler(t *Transcoder, db core.Database) *ConnectionHandler {
	h := &ConnectionHandler{t, db, make(map[string]struct{}), &sync.RWMutex{}}

	m := make(map[string]struct{})
	c := make(chan string)
	go h.db.GetTorrents(c)
	for h := range c {
		m[h] = struct{}{}
	}
	go h.t.Write(Message{"TorrentList", m, nil})
	go h.shovel()
	return h
}

func (h *ConnectionHandler) NewTorrent(t *core.Torrent) {
	h.lock.RLock()
	defer h.lock.RUnlock()
	if _, seen := h.torrentList[t.Hash]; !seen {
		h.t.Write(Message{"Torrents", nil, []*core.Torrent{t}})
	}
}

func (h *ConnectionHandler) NewSignature(s *core.Signature) {}

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
			go h.db.AddClient(h)
		case "Torrents":
			for _, t := range msg.Data {
				h.lock.Lock()
				h.torrentList[t.Hash] = struct{}{}
				h.lock.Unlock()
				h.db.Add(t)
			}
		}
	}
}

func (h *ConnectionHandler) Close() error {
	return h.t.Close()
}
