package main

import (
	"fmt"
	"github.com/emanueldonalds/husax/db"
	"github.com/emanueldonalds/husax/rss"
	"github.com/emanueldonalds/husax/web"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func main() {
	assetsDir := "./assets"

	_, err := os.Stat(assetsDir)

	if err != nil {
		panic("Could not stat assets directory. Make sure assets dir is in the working directory.")
	}

	router := mux.NewRouter()

	fs := http.FileServer(http.Dir(assetsDir))
	db := db.GetDb()

	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fs))

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { web.IndexHandler(w, r, db) })
	router.HandleFunc("/info/{id}", web.DetailsHandler(db))
	router.HandleFunc("/filter", func(w http.ResponseWriter, r *http.Request) { web.FilterHandler(w, r, db) })
	router.HandleFunc("/rss", func(w http.ResponseWriter, r *http.Request) { rss.RssHandler(w, r, db) })

	fmt.Println("Listening on :4932")
	log.Fatal(http.ListenAndServe(":4932", router))
}
