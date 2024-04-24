package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		IndexHandler(w, r)
	})

	fmt.Println("Listening on :3000")
	log.Fatal(http.ListenAndServe(":3000", mux))
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	index := index(GetListings(r))
	index.Render(r.Context(), w)
}

func GetListings(r *http.Request) []Listing {
	db := GetDb()
	defer db.Close()

	agency := r.URL.Query().Get("agency")
	qPriceMin := r.URL.Query().Get("price_min")
	qPriceMax := r.URL.Query().Get("price_max")
	qYearMin := r.URL.Query().Get("year_min")
	qYearMax := r.URL.Query().Get("year_max")
	qSizeMin := r.URL.Query().Get("size_min")
	qSizeMax := r.URL.Query().Get("size_max")
	qPriceOverAreaMin := r.URL.Query().Get("price_over_area_min")
	qPriceOverAreaMax := r.URL.Query().Get("price_over_area_max")
	qRoomsMin := r.URL.Query().Get("rooms_min")
	qRoomsMax := r.URL.Query().Get("rooms_max")
	qFirstSeenMin := r.URL.Query().Get("first_seen_min")
	qLastSeen := r.URL.Query().Get("last_seen")

	qOrderBy := r.URL.Query().Get("order_by")
	qSortOrder := r.URL.Query().Get("sort_order")
	qIncludeDeleted := r.URL.Query().Get("include_deleted")

	query, err := db.Query(
		"SELECT "+
			"IFNULL(address, \"\"), "+
			"IFNULL(listing.price, -1), "+
			"IFNULL(year, -1), "+
			"FLOOR(IFNULL(size_value, -1)), "+
			"IFNULL(size_name, \"\"),"+
            "IFNULL(FLOOR(listing.price/size_value), -1) as price_over_area, "+
			"IFNULL(rooms, -1), "+
			"IFNULL(first_seen, \"\"), "+
			"IFNULL(listing.last_seen, \"\"), "+
			"agency, "+
			"url, "+
            "IFNULL(price_change.price, -1), "+
            "IFNULL(price_change.last_seen, \"\") "+
			"FROM listing "+
            //"LEFT JOIN (SELECT price_change.price, price_change.last_seen FROM price_change ORDER BY price_change.last_seen DESC LIMIT 1)"+
			"LEFT JOIN price_change on price_change.listing_id = listing.id "+
			"WHERE deleted = ? "+
			"AND agency = COALESCE(NULLIF(?, ''), agency) "+
			"AND (listing.price IS NULL OR listing.price >= COALESCE(NULLIF(?, ''), listing.price)) "+
			"AND (listing.price IS NULL OR listing.price <= COALESCE(NULLIF(?, ''), listing.price)) "+
			"AND (year IS NULL OR year >= COALESCE(NULLIF(?, ''), year)) "+
			"AND (year IS NULL OR year <= COALESCE(NULLIF(?, ''), year)) "+
			"AND (size_value IS NULL OR size_value >= COALESCE(NULLIF(?, ''), size_value)) "+
			"AND (size_value IS NULL OR size_value <= COALESCE(NULLIF(?, ''), size_value)) "+
			"AND (rooms IS NULL OR rooms >= COALESCE(NULLIF(?, ''), rooms)) "+
			"AND (rooms IS NULL OR rooms <= COALESCE(NULLIF(?, ''), rooms)) "+
			"AND first_seen >= COALESCE(NULLIF(?, ''), first_seen) "+
			"AND listing.last_seen <= COALESCE(NULLIF(?, ''), listing.last_seen ) "+
			"HAVING (price_over_area IS NULL OR price_over_area >= COALESCE(NULLIF(?, ''), price_over_area)) "+
			"AND (price_over_area IS NULL OR price_over_area <= COALESCE(NULLIF(?, ''), price_over_area)) "+
			ResolveOrder(qOrderBy, qSortOrder),
		ResolveDeleted(qIncludeDeleted),
		agency,
		qPriceMin,
		qPriceMax,
		qYearMin,
		qYearMax,
		qSizeMin,
		qSizeMax,
		qRoomsMin,
		qRoomsMax,
		qFirstSeenMin,
		qLastSeen,
		qPriceOverAreaMin,
		qPriceOverAreaMax,
	)

	if err != nil {
		panic(err.Error())
	}

	//var priceChanges  = make(map[int][]PriceChange)
	listings := []Listing{}

	for query.Next() {
		var rowListing Listing
        var rowPriceChange PriceChange
		err := query.Scan(
			&rowListing.address,
			&rowListing.price,
			&rowListing.year,
			&rowListing.size.value,
			&rowListing.size.unit,
            &rowListing.priceOverArea,
			&rowListing.rooms,
			&rowListing.firstSeen,
			&rowListing.lastSeen,
			&rowListing.agency,
			&rowListing.url,
            &rowPriceChange.price,
            &rowPriceChange.lastSeen)
		if err != nil {
			panic(err.Error())
		}

		listings = append(listings, rowListing)
	}
	return listings
}

func GetDb() *sql.DB {
	dbPassEnv := os.Getenv("PROPERTY_VIEWER_DB_PASSWORD")
	connString := fmt.Sprintf("property-viewer:%s@tcp(10.0.1.12:3306)/property_api", dbPassEnv)

	db, err := sql.Open("mysql", connString)
	if err != nil {
		panic(err.Error())
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db
}

func ResolveOrder(qOrderBy string, qSortOrder string) string {
	if qOrderBy == "" {
		return ""
	}

	if !slices.Contains(
		[]string{
			"agency",
			"address",
			"price",
			"url",
			"size_value",
			"price_over_area",
			"rooms",
			"year",
			"first_seen",
			"last_seen",
		},
		strings.ToLower(qOrderBy)) {
		panic("Invalid order by value " + qOrderBy)
	}

	if !slices.Contains([]string{"", "asc", "desc"}, strings.ToLower(qSortOrder)) {
		panic("Invalid sort order " + qSortOrder)
	}

	return fmt.Sprintf("ORDER BY %s %s", qOrderBy, qSortOrder)
}

func ResolveDeleted(qIncludeDeleted string) string {
	if qIncludeDeleted == "" {
		return "false"
	}
	if !slices.Contains([]string{"true", "false"}, strings.ToLower(qIncludeDeleted)) {
		panic("Invalid include deleted value " + qIncludeDeleted)
	}

	return qIncludeDeleted
}
