package network

import (
	"log"

	"github.com/TheDistributedBay/TheDistributedBay/database"
)

type Message struct {
	code     string
	torrents map[string]struct{}
	data     []*database.Torrent
}

type Connection interface {
	Read() (Message, error)
	Write(Message) error
}

type ConnectionHandler struct {
	c    Connection
	db   database.Database
	quit chan struct{}
}

func NewConnectionHandler(c Connection, db database.Database) *ConnectionHandler {
	h := &ConnectionHandler{c, db, nil}
	go h.shovel()
	return h
}

// Shovels torrents into and out of the db
func (h *ConnectionHandler) shovel() {
	m := make(map[string]struct{})
	for _, h := range h.db.List() {
		m[h] = struct{}{}
	}
	h.c.Write(Message{"TorrentList", m, nil})
	for {
		msg, err := h.c.Read()
		if err != nil {
			log.Printf("Error for %v : %v", h.c, err)
			return
		}
		switch msg.code {
		case "TorrentList":
			torrents := make([]*database.Torrent, 0, len(msg.torrents))
			for _, hash := range h.db.List() {
				if _, ok := msg.torrents[hash]; !ok {
					t, err := h.db.Get(hash)
					if err != nil {
						log.Printf("Torrent dissapeared? : %v", hash)
						continue
					}
					torrents = append(torrents, t)
				}
			}
			h.c.Write(Message{"Torrents", nil, torrents})
		case "Torrents":
			for _, t := range msg.data {
				log.Print("ADDDING TORRENT WITHOUT VERIFICATION!!!!")
				h.db.Add(t)
			}
		}
	}
}
