package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
    //println("main");
	mux := http.NewServeMux()

	assetsDir := os.Getenv("PROPERTY_VIEWER_ASSETS_DIR")

    info, err := os.Stat(assetsDir);
    if err != nil {
        panic(err.Error())
    }
    if info.Mode().Perm()&0444 != 0444 {
        panic("Can not read assets")
    }


	db := GetDb()
	fs := http.FileServer(http.Dir(assetsDir))

    //println("before assets");
    mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

    //println("before handlers");
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { IndexHandler(w, r, db) })
	mux.HandleFunc("/filter", func(w http.ResponseWriter, r *http.Request) { FilterHandler(w, r, db) })

	fmt.Println("Listening on :4932")
	log.Fatal(http.ListenAndServe(":4932", mux))
}

func IndexHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    //println("index handler");
	listings := GetListings(r, db)
	lastScrape := GetLastScrape(db)
	index := index(listings, lastScrape)
	index.Render(r.Context(), w)
}

func FilterHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    //println("filter handler");
	listings := GetListings(r, db)
	index := result(listings)
	index.Render(r.Context(), w)
}

func GetListings(r *http.Request, db *sql.DB) []Listing {
    //println("get listings");
	agency := r.URL.Query().Get("agency")
	qPriceMin := r.URL.Query().Get("price_min")
	qPriceMax := r.URL.Query().Get("price_max")
	qYearMin := r.URL.Query().Get("year_min")
	qYearMax := r.URL.Query().Get("year_max")
	qSizeMin := r.URL.Query().Get("size_value_min")
	qSizeMax := r.URL.Query().Get("size_value_max")
	qPriceOverAreaMin := r.URL.Query().Get("price_over_area_min")
	qPriceOverAreaMax := r.URL.Query().Get("price_over_area_max")
	qRoomsMin := r.URL.Query().Get("rooms_min")
	qRoomsMax := r.URL.Query().Get("rooms_max")
	qFirstSeenMin := r.URL.Query().Get("first_seen_min")
	qLastSeen := r.URL.Query().Get("last_seen")

	qOrderBy := r.URL.Query().Get("order_by")
	qSortOrder := r.URL.Query().Get("sort_order")
	qIncludeDeleted := r.URL.Query().Get("include_deleted")


    //println("before query");
	query, err := db.Query(
		"SELECT "+
			"listing.id, "+
			"IFNULL(address, \"\"), "+
			"IFNULL(listing.price, -1), "+
			"IFNULL(year, -1), "+
			"FLOOR(IFNULL(size_value, -1)), "+
			"IFNULL(size_name, \"\"),"+
			"IFNULL(FLOOR(listing.price/size_value), -1) as price_over_area, "+
			"IFNULL(rooms, -1), "+
			"first_seen, "+
			"listing.last_seen, "+
			"agency, "+
			"url, "+
			"IFNULL(price_change.price, -1), "+
			"IFNULL(price_change.last_seen, \"\") "+
			"FROM listing "+
			"LEFT JOIN price_change on price_change.listing_id = listing.id "+
			"WHERE "+
            ResolveDeleted(qIncludeDeleted)+
			"AND agency = COALESCE(NULLIF(?, ''), agency) "+
			"AND (listing.price IS NULL OR listing.price >= COALESCE(NULLIF(?, ''), listing.price-1)) "+
			"AND (listing.price IS NULL OR listing.price <= COALESCE(NULLIF(?, ''), listing.price+1)) "+
			"AND (year IS NULL OR year >= COALESCE(NULLIF(?, ''), year-1)) "+
			"AND (year IS NULL OR year <= COALESCE(NULLIF(?, ''), year+1)) "+
			"AND (size_value IS NULL OR size_value >= COALESCE(NULLIF(?, ''), size_value-1)) "+
			"AND (size_value IS NULL OR size_value <= COALESCE(NULLIF(?, ''), size_value+1)) "+
			"AND (rooms IS NULL OR rooms >= COALESCE(NULLIF(?, ''), rooms-1)) "+
			"AND (rooms IS NULL OR rooms <= COALESCE(NULLIF(?, ''), rooms+1)) "+
			"AND first_seen >= COALESCE(NULLIF(?, ''), first_seen) "+
			"AND listing.last_seen <= COALESCE(NULLIF(?, ''), listing.last_seen ) "+
			"HAVING (price_over_area IS NULL OR price_over_area >= COALESCE(NULLIF(?, ''), price_over_area-1)) "+
			"AND (price_over_area IS NULL OR price_over_area <= COALESCE(NULLIF(?, ''), price_over_area+1)) "+
			ResolveOrder(qOrderBy, qSortOrder),
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

    //println("after query");

	if err != nil {
		panic(err.Error())
	}

	listings := []Listing{}

    //println("before next");

	for query.Next() {
		var rowListing Listing
		var rowPriceChange PriceChange
		err := query.Scan(
			&rowListing.id,
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

    query.Close()

    //println("after next");

	priceChanges := GetPriceChanges(listings, db)

	// Add price changes to listings
	for _, priceChange := range priceChanges {
		for i, listing := range listings {
			if priceChange.listingId == listing.id {
				listings[i].priceHistory = append(listing.priceHistory, priceChange)
			}
		}
	}

    //println("return listings")
	return listings

}

func GetPriceChanges(listings []Listing, db *sql.DB) []PriceChange {
    //println("get price changes");

	if len(listings) == 0 {
		return nil
	}
	listingsIds := []string{}

	for _, listing := range listings {
		listingsIds = append(listingsIds, listing.id)
	}

	joinedIds := strings.Join(listingsIds, ", ")
    
	query, err := db.Query("SELECT price, COALESCE(last_seen, ''), listing_id FROM price_change WHERE listing_id IN (" + joinedIds + ") ORDER BY last_seen DESC")

	if err != nil {
		panic(err.Error())
	}

	priceChanges := []PriceChange{}

	for query.Next() {
		var rowPriceChange PriceChange
		err := query.Scan(&rowPriceChange.price, &rowPriceChange.lastSeen, &rowPriceChange.listingId)

		if err != nil {
			panic(err.Error())
		}
		priceChanges = append(priceChanges, rowPriceChange)
	}

    query.Close()

	return priceChanges
}

func GetLastScrape(db *sql.DB) ScrapeEvent {
	query, qErr := db.Query("SELECT date, added, updated, deleted, undeleted, total_active from scrape_event ORDER BY date DESC LIMIT 1")

	if qErr != nil {
		panic(qErr.Error())
	}

	query.Next()

	var scrapeEvent ScrapeEvent
	sErr := query.Scan(&scrapeEvent.date, &scrapeEvent.added, &scrapeEvent.updated, &scrapeEvent.deleted, &scrapeEvent.undeleted, &scrapeEvent.totalActive)

    query.Close();

	if sErr != nil {
		panic(sErr.Error())
	}


    return scrapeEvent
}

func GetDb() *sql.DB {
    //println("get db");
	dbPassEnv := os.Getenv("PROPERTY_VIEWER_DB_PASSWORD")
	connString := fmt.Sprintf("property-viewer:%s@tcp(10.0.1.12:3306)/property_api", dbPassEnv)

	db, err := sql.Open("mysql", connString)
	if err != nil {
		panic(err.Error())
	}
	db.SetConnMaxLifetime(180)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
    
    //println("return db");
	return db
}

func ResolveOrder(qOrderBy string, qSortOrder string) string {
	if qOrderBy == "" {
		return "ORDER BY first_seen desc"
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

	orderBy := "listing." + qOrderBy
	if qOrderBy == "price_over_area" {
		orderBy = qOrderBy
	}

	return fmt.Sprintf("ORDER BY %s %s", orderBy, qSortOrder)
}

func ResolveDeleted(qIncludeDeleted string) string {
	if qIncludeDeleted == "" {
		return "deleted = false "
	}
	if slices.Contains([]string{"true", "false"}, strings.ToLower(qIncludeDeleted)) {
		return "deleted IN (true, false) "
	}

    panic("Invalid include deleted value " + qIncludeDeleted)
}

