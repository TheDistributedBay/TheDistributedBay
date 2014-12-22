package frontend

import (
	"encoding/json"
	"log"
	"net/http"
)

type GetTorrentHandler struct {
	TorrentClient
}

func (th GetTorrentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Query().Get("hash")
	log.Println("HASH", hash)
	if hash == "" {
		http.Error(w, "Invalid request.", http.StatusBadRequest)
		return
	}

	torrent, err := th.TorrentClient.Database.Get(hash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(torrent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
