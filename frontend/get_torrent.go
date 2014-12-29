package frontend

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/TheDistributedBay/TheDistributedBay/core"
	"github.com/TheDistributedBay/TheDistributedBay/search"
	"github.com/TheDistributedBay/TheDistributedBay/torrent"
)

type GetTorrentHandler struct {
	s       *search.Searcher
	updater *torrent.StatsUpdater
	db      core.Database
}

type TorrentResult struct {
	Name, MagnetLink, Hash, InfoHash, Category, Description string
	CreatedAt                                               time.Time
	Seeders, Leechers, Completed                            core.Range
	Size                                                    uint64
	Files                                                   uint
	MoreLikeThis                                            []searchResult
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
		Name:        t.Name,
		MagnetLink:  t.MagnetLink(),
		Hash:        t.Hash,
		InfoHash:    t.NiceInfoHash(),
		Category:    t.Category(),
		CreatedAt:   t.CreatedAt,
		Size:        t.Size,
		Seeders:     t.Seeders,
		Leechers:    t.Leechers,
		Completed:   t.Completed,
		Files:       t.Files,
		Description: t.Description,
	}

	mlt, err := th.s.MoreLikeThis(t.Hash)
	if err != nil {
		log.Println("Search MLT ERR", err)
	}
	mlt_length := len(mlt)
	if mlt_length > 5 {
		mlt_length = 5
	}
	result.MoreLikeThis = make([]searchResult, mlt_length)
	for i, res := range mlt[:mlt_length] {
		th.updater.QueueTorrent(res)
		result.MoreLikeThis[i] = searchResult{
			Name:       res.Name,
			MagnetLink: res.MagnetLink(),
			Hash:       res.Hash,
			Category:   res.Category(),
			CreatedAt:  res.CreatedAt,
			Size:       res.Size,
			Seeders:    res.Seeders,
			Leechers:   res.Leechers,
			Completed:  res.Completed,
		}
	}

	js, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
