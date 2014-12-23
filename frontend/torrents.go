package frontend

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/TheDistributedBay/TheDistributedBay/database"
	"github.com/TheDistributedBay/TheDistributedBay/search"
)

type TorrentsHandler struct {
	s  *search.Searcher
	db database.Database
}

func (th TorrentsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	p := 0
	fmt.Sscan(r.URL.Query().Get("p"), &p)
	var results []*database.Torrent
	var err error
	if q != "" {
		results, _, err = th.s.Search(q, 35*p, 35)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	js, err := json.Marshal(results)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
