package network

import (
	"log"
	"sync"

	"github.com/TheDistributedBay/TheDistributedBay/crypto"
	"github.com/TheDistributedBay/TheDistributedBay/database"
)

type ConnectionHandler struct {
	t           *Transcoder
	db          database.Database
	torrentList map[string]struct{}
	lock        *sync.RWMutex
}

func NewConnectionHandler(t *Transcoder, db database.Database) *ConnectionHandler {
	h := &ConnectionHandler{t, db, nil, &sync.RWMutex{}}

	m := make(map[string]struct{})
	for _, h := range h.db.List() {
		m[h] = struct{}{}
	}
	go h.t.Write(Message{"TorrentList", m, nil})
	go h.shovel()
	return h
}

func (h *ConnectionHandler) NewTorrent(t *database.Torrent) {
	h.lock.RLock()
	defer h.lock.RUnlock()
	if _, seen := h.torrentList[t.Hash]; !seen {
		h.t.Write(Message{"Torrents", nil, []*database.Torrent{t}})
	}
}

// Shovels torrents into the db
func (h *ConnectionHandler) shovel() {
	for {
		msg, err := h.t.Read()
		if err != nil {
			log.Printf("Error for %v : %v", h.t, err)
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
				if err := crypto.VerifyTorrent(t); err != nil {
					log.Printf("Invalid torrent %s recieved : %v", t, err)
				} else {
					h.lock.Lock()
					h.torrentList[t.Hash] = struct{}{}
					h.lock.Unlock()
					h.db.Add(t)
				}
			}
		}
	}
}

func (h *ConnectionHandler) Close() error {
	return h.t.Close()
}
