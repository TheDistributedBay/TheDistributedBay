package frontend

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"github.com/TheDistributedBay/TheDistributedBay/database"
	"net/http"
	"time"
)

type AddTorrentHandler struct {
	TorrentClient
}

func (th AddTorrentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var t database.Torrent
	err := decoder.Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	t.Size = 0
	t.FileCount = 0
	t.CreatedAt = time.Now()

	// TODO set public key, signature, hash
	t.PublicKey = "wub a dub a blub blub"
	t.Signature = "sign me please!"
	t.Hash = randomStr()

	th.TorrentClient.Database.Add(&t)

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
