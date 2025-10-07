package web

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/emanueldonalds/property-viewer/db"
)

func IndexHandler(w http.ResponseWriter, r *http.Request, sqldb *sql.DB) {
	orderBy := r.URL.Query().Get("order_by")
	sortOrder := r.URL.Query().Get("sort_order")

	listings := db.GetListings(r, sqldb)
	lastScrape := db.GetLastScrape(sqldb)
	index := Index(listings, lastScrape, orderBy, sortOrder)
	index.Render(r.Context(), w)
}

func FilterHandler(w http.ResponseWriter, r *http.Request, sqldb *sql.DB) {
	orderBy := r.URL.Query().Get("order_by")
	sortOrder := r.URL.Query().Get("sort_order")

	listings := db.GetListings(r, sqldb)
	lastScrape := db.GetLastScrape(sqldb)
	index := Listings(listings, lastScrape, orderBy, sortOrder)
	index.Render(r.Context(), w)
}

func DetailsHandler(sqldb *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		idParam := params["id"]

		id, err := strconv.Atoi(idParam)

		if err != nil {
			fmt.Printf("Could not convert ID param [%s] to int", idParam)
		}

		listing := db.GetListing(id, sqldb)
		listingHistory := db.GetListingHistory(id, sqldb)

		for i := range listingHistory {
			historicListing := &listingHistory[i]
			historicListing.InfoUrl = "https://hus.ax/info/" + historicListing.Id
		}

		listingPage := Listing(listing, listingHistory)
		listingPage.Render(r.Context(), w)
	}
}
