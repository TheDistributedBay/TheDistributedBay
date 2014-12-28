package frontend

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/TheDistributedBay/TheDistributedBay/core"
	"github.com/TheDistributedBay/TheDistributedBay/torrent"
)

type GetTorrentHandler struct {
	updater *torrent.StatsUpdater
	db      core.Database
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

	t, err := th.db.Get(hash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	th.updater.QueueTorrent(t)

	result := TorrentResult{
		Name:       t.Name,
		MagnetLink: t.MagnetLink(),
		Hash:       t.Hash,
		InfoHash:   t.NiceInfoHash(),
		Category:   t.Category(),
		CreatedAt:  t.CreatedAt,
		Size:       t.Size,
		Seeders:    t.Seeders,
		Leechers:   t.Leechers,
		Completed:  t.Completed,
		Files:      t.Files,
	}

	js, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
