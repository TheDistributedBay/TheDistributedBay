package frontend

import (
	"encoding/json"
	"github.com/TheDistributedBay/TheDistributedBay/database"
	"net/http"
)

func TorrentsHandler(w http.ResponseWriter, r *http.Request) {
	profile := [1]database.Torrent{database.Torrent{
		Name:        "Banana",
		MagnetLink:  "magnet:blahblahblah",
		Description: "This is a description!"}}

	js, err := json.Marshal(profile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
