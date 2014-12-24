package frontend

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/TheDistributedBay/TheDistributedBay/core"
	"github.com/TheDistributedBay/TheDistributedBay/search"
)

type TorrentsHandler struct {
	s  *search.Searcher
	db core.Database
}

type TorrentBlob struct {
	Torrents []*core.Torrent
	Pages    uint64
}

func (th TorrentsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	p := 0
	fmt.Sscan(r.URL.Query().Get("p"), &p)
	b := TorrentBlob{}
	if q != "" {
		results, count, err := th.s.Search(q, 35*p, 35)
		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		b.Torrents = results
		b.Pages = count / 35
		if count%35 > 0 {
			b.Pages += 1
		}
	}

	js, err := json.Marshal(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
