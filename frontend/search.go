package frontend

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/TheDistributedBay/TheDistributedBay/core"
	"github.com/TheDistributedBay/TheDistributedBay/search"
	"github.com/TheDistributedBay/TheDistributedBay/torrent"
)

type SearchHandler struct {
	s       *search.Searcher
	db      core.Database
	updater *torrent.StatsUpdater
}

type searchResult struct {
	Name, MagnetLink, Hash, Category string
	Size                             uint64
	CreatedAt                        time.Time
	Seeders, Leechers, Completed     core.Range
}

type TorrentBlob struct {
	Torrents []searchResult
	Pages    int
}

func (th SearchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	sort := r.URL.Query().Get("sort")
	size, err := strconv.Atoi(r.URL.Query().Get("size"))
	if err != nil {
		size = 35
	}
	categories := strings.Split(r.URL.Query().Get("category"), ",")
	categoryIds := make([]uint8, len(categories))
	for i, category := range categories {
		categoryIds[i] = core.CategoryToId(category)
	}
	p := 0
	fmt.Sscan(r.URL.Query().Get("p"), &p)

	results, count, err := th.s.Search(q, size*p, size, categoryIds, sort)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	b := TorrentBlob{}
	b.Torrents = make([]searchResult, len(results))
	for i, result := range results {
		th.updater.QueueTorrent(result)
		b.Torrents[i] = searchResult{
			Name:       result.Name,
			MagnetLink: result.MagnetLink(),
			Hash:       result.Hash,
			Category:   result.Category(),
			CreatedAt:  result.CreatedAt,
			Size:       result.Size,
			Seeders:    result.Seeders,
			Leechers:   result.Leechers,
			Completed:  result.Completed,
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
