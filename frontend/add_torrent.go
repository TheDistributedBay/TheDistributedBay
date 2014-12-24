package frontend

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"github.com/TheDistributedBay/TheDistributedBay/core"
)

type AddTorrentHandler struct {
	db core.Database
}

func (th AddTorrentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var t core.Torrent
	err := decoder.Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	t.CreatedAt = time.Now()
	t.CalculateHash()

	th.db.Add(&t)

	js, err := json.Marshal(t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
func randomStr() string {
	size := 32 // change the length of the generated random string here

	rb := make([]byte, size)
	rand.Read(rb)

	return base64.URLEncoding.EncodeToString(rb)
}
