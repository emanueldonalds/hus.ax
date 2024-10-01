package web

import (
	"database/sql"
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

