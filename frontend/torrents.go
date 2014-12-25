package frontend

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/TheDistributedBay/TheDistributedBay/core"
	"github.com/TheDistributedBay/TheDistributedBay/search"
)

type TorrentsHandler struct {
	s  *search.Searcher
	db core.Database
}

type searchResult struct {
	Name, MagnetLink, Hash, Category string
	CreatedAt                        time.Time
}

type TorrentBlob struct {
	Torrents []searchResult
	Pages    int
}

func (th TorrentsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	category := core.CategoryToId(r.URL.Query().Get("category"))
	p := 0
	fmt.Sscan(r.URL.Query().Get("p"), &p)
	b := TorrentBlob{}
	results, count, err := th.s.Search(q, 35*p, 35, category)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	b.Torrents = make([]searchResult, len(results))
	for i, result := range results {
		b.Torrents[i] = searchResult{
			Name:       result.Name,
			MagnetLink: result.MagnetLink(),
			Hash:       result.Hash,
			Category:   result.Category(),
			CreatedAt:  result.CreatedAt,
		}
	}
	b.Pages = count / 35
	if count%35 > 0 {
		b.Pages += 1
	}

	js, err := json.Marshal(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
