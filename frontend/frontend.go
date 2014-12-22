package frontend

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func Serve() {
	r := mux.NewRouter()
	r.HandleFunc("/", IndexHandler)

	//r.HandleFunc("/torrents", TorrentsHandler)
	//r.HandleFunc("/torrent/{torrent}", TorrentHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("frontend/static/"))))
	http.Handle("/", r)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}
