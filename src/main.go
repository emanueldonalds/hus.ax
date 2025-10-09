package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/emanueldonalds/husax/db"
	"github.com/emanueldonalds/husax/rss"
	"github.com/emanueldonalds/husax/web"
	"github.com/gorilla/mux"
)

// Same etag for all files, generates a new every time server restarts
var etag string = "W/\"" + fmt.Sprint(time.Now().UTC().Unix()) + "\""

func main() {
	assetsDir := "./assets"

	_, err := os.Stat(assetsDir)

	if err != nil {
		panic("Could not stat assets directory. Make sure assets dir is in the working directory.")
	}

	mux := mux.NewRouter()
	db := db.GetDb()

	assetsHandler := http.StripPrefix("/assets/", http.FileServer(http.Dir(assetsDir)))
	mux.PathPrefix("/assets/").Handler(cacheControl(assetsHandler))

	mux.Handle("/", cacheControl(web.IndexHandler(db)))
	mux.Handle("/info/{id}", cacheControl(web.DetailsHandler(db)))
	mux.Handle("/stats", cacheControl(web.StatisticsHandler(db)))
	mux.Handle("/filter", cacheControl(web.FilterHandler(db)))
	mux.Handle("/rss", cacheControl(rss.RssHandler(db)))

	fmt.Println("Listening on :4932")
	log.Fatal(http.ListenAndServe(":4932", mux))
}

func cacheControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ifNoneMatch := r.Header.Get("If-None-Match")

		if ifNoneMatch == etag {
			w.WriteHeader(304)
			return
		}

		w.Header().Set("ETag", etag)
		w.Header().Set("Cache-Control", "no-cache, must-revalidate")
		h.ServeHTTP(w, r)
	})
}
