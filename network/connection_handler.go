package network

import (
	"log"

	"github.com/TheDistributedBay/TheDistributedBay/database"
)

type ConnectionHandler struct {
	t  *Transcoder
	db database.Database
}

func NewConnectionHandler(t *Transcoder, db database.Database) *ConnectionHandler {
	h := &ConnectionHandler{t, db}
	go h.shovel()
	return h
}

// Shovels torrents into and out of the db
func (h *ConnectionHandler) shovel() {
	m := make(map[string]struct{})
	for _, h := range h.db.List() {
		m[h] = struct{}{}
	}
	go h.t.Write(Message{"TorrentList", m, nil})
	for {
		msg, err := h.t.Read()
		if err != nil {
			log.Printf("Error for %v : %v", h.t, err)
			return
		}
		switch msg.Code {
		case "TorrentList":
			torrents := make([]*database.Torrent, 0)
			for _, hash := range h.db.List() {
				if _, ok := msg.Torrents[hash]; !ok {
					t, err := h.db.Get(hash)
					if err != nil {
						log.Printf("Torrent dissapeared? : %v", hash)
						continue
					}
					torrents = append(torrents, t)
				}
			}
			go h.t.Write(Message{"Torrents", nil, torrents})
		case "Torrents":
			for _, t := range msg.Data {
				log.Print("ADDDING TORRENT WITHOUT VERIFICATION!!!!")
				h.db.Add(t)
			}
		}
	}
}

func (h *ConnectionHandler) Close() error {
	return h.t.Close()
}
