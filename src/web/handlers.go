package web

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"

	"github.com/emanueldonalds/property-viewer/db"
)

func IndexHandler(w http.ResponseWriter, r *http.Request, sqldb *sql.DB) {
	listings := db.GetListings(r, sqldb)
	lastScrape := db.GetLastScrape(sqldb)
	index := Index(listings, lastScrape)
	index.Render(r.Context(), w)
}

func FilterHandler(w http.ResponseWriter, r *http.Request, sqldb *sql.DB) {
	listings := db.GetListings(r, sqldb)
	lastScrape := db.GetLastScrape(sqldb)
	index := Listings(listings, lastScrape)
	index.Render(r.Context(), w)
}

func StatsHandler(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
        id := params["id"]
		fmt.Println("Param id is ")
		fmt.Println(id)


        listing := db.GetListing(id);


		fmt.Fprintf(w, "Stats\n%s", params["id"])
	}
}
