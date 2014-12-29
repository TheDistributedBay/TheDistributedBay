package frontend

//go:generate bash -c "cd angular; npm install"
//go:generate bash -c "cd angular; node_modules/bower/bin/bower install"
//go:generate bash -c "cd angular; node_modules/grunt-cli/bin/grunt build"

import (
	"log"
	"mime"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/TheDistributedBay/TheDistributedBay/core"
	"github.com/TheDistributedBay/TheDistributedBay/search"
	"github.com/TheDistributedBay/TheDistributedBay/torrent"
)

func Serve(httpAddress *string, db core.Database, s *search.Searcher, devAssets bool, updater *torrent.StatsUpdater) {

	// Add SVG to mime directory
	mime.AddExtensionType(".svg", "image/svg+xml")

	r := mux.NewRouter()

	apirouter := ApiRouter(db, s, updater)

	r.PathPrefix("/api/").Handler(apirouter)
	if devAssets {
		log.Println("Debug mode is on. Serving development assets from angular/app.")
		r.PathPrefix("/styles/").Handler(http.FileServer(http.Dir("frontend/angular/.tmp/")))
		r.PathPrefix("/bower_components/").Handler(http.FileServer(http.Dir("frontend/angular/")))
		r.PathPrefix("/").Handler(NotFoundHook{
			http.FileServer(http.Dir("frontend/angular/app/")),
			"frontend/angular/app/index.html"})
	} else {
		r.PathPrefix("/").Handler(NotFoundHook{
			http.FileServer(http.Dir("frontend/angular/dist/")),
			"frontend/angular/dist/index.html"})
	}
	http.Handle("/", r)
	log.Println("Web server listening on", *httpAddress)
	err := http.ListenAndServe(*httpAddress, nil)
	if err != nil {
		log.Println(err)
	}
}

func ApiRouter(db core.Database, s *search.Searcher, updater *torrent.StatsUpdater) *mux.Router {
	r := mux.NewRouter()
	r.Methods("GET").Path("/api/search").Handler(SearchHandler{s, db, updater})
	r.Methods("POST").Path("/api/add_torrent").Handler(AddTorrentHandler{db})
	r.Methods("GET").Path("/api/torrent").Handler(GetTorrentHandler{s, updater, db})

	return r
}

type hookedResponseWriter struct {
	http.ResponseWriter
	*http.Request
	file   string
	ignore bool
}

func (hrw *hookedResponseWriter) WriteHeader(status int) {
	if status == 404 {
		hrw.ignore = true
		// Write custom error here to hrw.ResponseWriter
		hrw.ResponseWriter.Header().Set("Content-Type", "text/html")
		http.ServeFile(hrw.ResponseWriter, hrw.Request, hrw.file)
	} else {

		hrw.ResponseWriter.WriteHeader(status)
	}
}

func (hrw *hookedResponseWriter) Write(p []byte) (int, error) {
	if hrw.ignore {
		return len(p), nil
	}
	return hrw.ResponseWriter.Write(p)
}

func (hrw *hookedResponseWriter) Header() http.Header {
	return hrw.ResponseWriter.Header()
}

type NotFoundHook struct {
	h    http.Handler
	file string
}

func (nfh NotFoundHook) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	nfh.h.ServeHTTP(&hookedResponseWriter{ResponseWriter: w, Request: r, file: nfh.file}, r)
}
