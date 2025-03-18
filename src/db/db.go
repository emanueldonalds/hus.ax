package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"os"
	"slices"
	"strings"
)

func GetDb() *sql.DB {
	dbHost := os.Getenv("PROPERTY_VIEWER_DB_HOST")
	dbPass := os.Getenv("PROPERTY_VIEWER_DB_PASSWORD")

	if dbHost == "" {
		panic("PROPERTY_VIEWER_DB_HOST must be set.")
	}
	if dbPass == "" {
		panic("PROPERTY_VIEWER_DB_PASSWORD must be set.")
	}

	connString := fmt.Sprintf("property-viewer:%s@tcp(%s:3306)/property_api", dbPass, dbHost)

	db, err := sql.Open("mysql", connString)
	if err != nil {
		panic(err.Error())
	}

	db.SetConnMaxLifetime(180)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db
}

func GetListing(id int, sqldb *sql.DB) Listing {

	query, err := sqldb.Query(
		"SELECT "+
			"listing.id, "+
			"IFNULL(address, \"\"), "+
			"IFNULL(listing.price, -1), "+
			"IFNULL(build_year, -1), "+
			"FLOOR(IFNULL(size_value, -1)), "+
			"IFNULL(size_name, \"\"),"+
			"IFNULL(FLOOR(listing.price/size_value), -1) as price_over_area, "+
			"IFNULL(rooms, -1), "+
			"first_seen, "+
			"listing.last_seen, "+
			"IFNULL(listing.last_updated, listing.last_seen), "+
			"agency, "+
			"url, "+
			"deleted = 1 "+
			"FROM listing "+
			"WHERE id = ?",
		id,
	)

	if err != nil {
		panic(err.Error())
	}

	query.Next()

    var listing Listing
    err = query.Scan(
        &listing.Id,
        &listing.Address,
        &listing.Price,
        &listing.Year,
        &listing.Size.Value,
        &listing.Size.Unit,
        &listing.PriceOverArea,
        &listing.Rooms,
        &listing.FirstSeen,
        &listing.LastSeen,
        &listing.LastUpdated,
        &listing.Agency,
        &listing.Url,
        &listing.Deleted)

	query.Close()

	if err != nil {
		panic(err.Error())
	}


	priceChanges := GetPriceChanges([]Listing{listing}, sqldb)

	for _, priceChange := range priceChanges {
        listing.PriceHistory = append(listing.PriceHistory, priceChange)
	}

	return listing
}

func GetListings(r *http.Request, sqldb *sql.DB) []Listing {

}

func GetPriceChanges(listings []Listing, sqldb *sql.DB) []PriceChange {
	if len(listings) == 0 {
		return nil
	}
	listingsIds := []string{}

	for _, listing := range listings {
		listingsIds = append(listingsIds, listing.Id)
	}

	joinedIds := strings.Join(listingsIds, ", ")

	query, err := sqldb.Query("SELECT IFNULL(price, 0), COALESCE(last_seen, ''), listing_id FROM price_change WHERE listing_id IN (" + joinedIds + ") ORDER BY last_seen DESC")

	if err != nil {
		panic(err.Error())
	}

	priceChanges := []PriceChange{}

	for query.Next() {
		var rowPriceChange PriceChange
		err := query.Scan(&rowPriceChange.Price, &rowPriceChange.LastSeen, &rowPriceChange.ListingId)

		if err != nil {
			panic(err.Error())
		}
		priceChanges = append(priceChanges, rowPriceChange)
	}

	query.Close()

	return priceChanges
}

func GetLastScrape(sqldb *sql.DB) ScrapeEvent {
	query, qErr := sqldb.Query("SELECT date, added, updated, deleted, undeleted, total_active from scrape_event ORDER BY date DESC LIMIT 1")

	if qErr != nil {
		panic(qErr.Error())
	}

	query.Next()

	var scrapeEvent ScrapeEvent
	sErr := query.Scan(&scrapeEvent.Date, &scrapeEvent.Added, &scrapeEvent.Updated, &scrapeEvent.Deleted, &scrapeEvent.Undeleted, &scrapeEvent.TotalActive)

	query.Close()

	if sErr != nil {
		panic(sErr.Error())
	}

	return scrapeEvent
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
			"build_year",
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
