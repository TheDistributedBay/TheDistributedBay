package frontend

import (
	"encoding/json"
	"log"
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
	var results []*database.Torrent
	if q != "" {
		searchResults, err := th.s.Search(q)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		results = make([]*database.Torrent, len(searchResults.Hits))
		for i, val := range searchResults.Hits {
			log.Println(i, val.ID)
			t, err := th.db.Get(val.ID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			results[i] = t
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
