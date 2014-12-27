package frontend

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/TheDistributedBay/TheDistributedBay/core"
)

type GetTorrentHandler struct {
	db core.Database
}

type TorrentResult struct {
	Name, MagnetLink, Hash, InfoHash, Category string
	CreatedAt                                  time.Time
	Seeders, Leechers, Completed               core.Range
	Size                                       uint64
	Files                                      uint
}

func (th GetTorrentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Query().Get("hash")
	if hash == "" {
		http.Error(w, "Invalid request.", http.StatusBadRequest)
		return
	}

	torrent, err := th.db.Get(hash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	result := TorrentResult{
		Name:       torrent.Name,
		MagnetLink: torrent.MagnetLink(),
		Hash:       torrent.Hash,
		InfoHash:   torrent.NiceInfoHash(),
		Category:   torrent.Category(),
		CreatedAt:  torrent.CreatedAt,
		Size:       torrent.Size,
		Seeders:    torrent.Seeders,
		Leechers:   torrent.Leechers,
		Completed:  torrent.Completed,
		Files:      torrent.Files,
	}

	js, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
