package frontend

//go:generate bash -c "cd angular; npm install"
//go:generate bash -c "cd angular; node_modules/bower/bin/bower install"
//go:generate bash -c "cd angular; node_modules/grunt-cli/bin/grunt build"

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
        "os"
)

func Serve() {
        r := mux.NewRouter()
	r.PathPrefix("/api/").Handler(ApiRouter())
        if (os.Getenv("DEBUG") != "") {
          fmt.Println("Debug mode is on.")
          r.PathPrefix("/styles/").Handler(http.FileServer(http.Dir("frontend/angular/.tmp/")))
          r.PathPrefix("/bower_components/").Handler(http.FileServer(http.Dir("frontend/angular/")))
          r.PathPrefix("/").Handler(NotFoundHook{
            http.FileServer(http.Dir("frontend/angular/app/")),
            "frontend/angular/app/index.html" })
        } else {
          r.PathPrefix("/").Handler(NotFoundHook{
            http.FileServer(http.Dir("frontend/angular/dist/")),
            "frontend/angular/dist/index.html" })
        }
        http.Handle("/", r)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}

func ApiRouter() *mux.Router {
    r := mux.NewRouter()
    r.HandleFunc("/api/torrents", TorrentsHandler)

    return r
}

type hookedResponseWriter struct {
    http.ResponseWriter
    *http.Request
    file string
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
  return hrw.ResponseWriter.Header();
}

type NotFoundHook struct {
    h http.Handler
    file string
}

func (nfh NotFoundHook) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  nfh.h.ServeHTTP(&hookedResponseWriter{ResponseWriter: w, Request: r, file: nfh.file}, r)
}
