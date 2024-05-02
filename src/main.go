package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	_ "github.com/go-sql-driver/mysql"
    "github.com/emanueldonalds/property-viewer/db"
    "github.com/emanueldonalds/property-viewer/handlers"
)

func main() {
	assetsDir := os.Getenv("PROPERTY_VIEWER_ASSETS_DIR")
	info, err := os.Stat(assetsDir)

	if err != nil {
		panic("Could not stat assets directory. PROPERTY_VIEWER_ASSETS_DIR must be set.")
	}
	if info.Mode().Perm()&0444 != 0444 {
		panic("Can not read assets")
	}

	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir(assetsDir))
	db := db.GetDb()

    mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { handlers.IndexHandler(w, r, db) })
	mux.HandleFunc("/filter", func(w http.ResponseWriter, r *http.Request) { handlers.FilterHandler(w, r, db) })

	fmt.Println("Listening on :4932")
	log.Fatal(http.ListenAndServe(":4932", mux))
}
