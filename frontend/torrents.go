package frontend

import (
	"encoding/json"
	"github.com/TheDistributedBay/TheDistributedBay/database"
	"log"
	"net/http"
)

type TorrentsHandler struct {
	TorrentClient
}

func (th TorrentsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	var results []*database.Torrent
	if q != "" {
		searchResults, err := th.TorrentClient.Search(q)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		results = make([]*database.Torrent, len(searchResults.Hits))
		for i, val := range searchResults.Hits {
			log.Println(i, val.ID)
			t, err := th.TorrentClient.Database.Get(val.ID)
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
