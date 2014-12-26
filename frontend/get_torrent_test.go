package frontend

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

func TestGetTorrent(t *testing.T) {
	db, _, endpoint := setupHttpTest(t)
	setupHttpTest(t)

	// Valid Request
	res, err := http.Get(endpoint + "/api/torrent?hash=test%20hash")
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != 200 {
		t.Errorf("Success expected: %d", res.StatusCode)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
	t1 := &TorrentResult{}
	err = json.Unmarshal(body, t1)
	if err != nil {
		t.Error(err)
	}
	t2, err := db.Get("test hash")
	if err != nil {
		t.Error(err)
	}
	log.Println(t1)
	if t1.Name != t2.Name || t1.Hash != t2.Hash || t1.MagnetLink != t2.MagnetLink() || t1.InfoHash != t2.NiceInfoHash() || t1.CreatedAt != t2.CreatedAt || t1.Category != t2.Category() || t1.Seeders != t2.Seeders || t1.Leechers != t2.Leechers || t1.Completed != t2.Completed || t1.Files != t2.Files || t1.Size != t2.Size {
		t.Error("The two hashes are not the same.")
	}

	// Request with no hash
	res, err = http.Get(endpoint + "/api/torrent")
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != 400 {
		t.Errorf("Failure expected: %d", res.StatusCode)
	}
	defer res.Body.Close()

	// Request with non-existant hash
	res, err = http.Get(endpoint + "/api/torrent?hash=bad%20hash")
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != 404 {
		t.Errorf("Failure expected: %d", res.StatusCode)
	}
	defer res.Body.Close()
}
