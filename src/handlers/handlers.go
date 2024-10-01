package handlers

import (
	"database/sql"
	"fmt"
	"github.com/emanueldonalds/property-viewer/components"
	"github.com/emanueldonalds/property-viewer/db"
	"net/http"
	"slices"
	"strings"
)

func IndexHandler(w http.ResponseWriter, r *http.Request, sqldb *sql.DB) {
	listings := GetListings(r, sqldb)
	lastScrape := GetLastScrape(sqldb)
	index := components.Index(listings, lastScrape)
	index.Render(r.Context(), w)
}

func FilterHandler(w http.ResponseWriter, r *http.Request, sqldb *sql.DB) {
	listings := GetListings(r, sqldb)
	lastScrape := GetLastScrape(sqldb)
	index := components.Listings(listings, lastScrape)
	index.Render(r.Context(), w)
}

func RssHandler(w http.ResponseWriter, r *http.Request, sqldb *sql.DB, mux *http.ServeMux) {

    listings := GetListings(r, sqldb)

	rssPage := 
`<?xml version="1.0"?>
  <rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
  <channel>
  <title>NASA Space Station News</title>
  <link>http://www.nasa.gov/</link>
  <description>A RSS news feed containing the latest NASA press releases on the International Space Station.</description>
  <language>en-us</language>
  <pubDate>Tue, 10 Jun 2003 04:00:00 GMT</pubDate>
  <lastBuildDate>Fri, 21 Jul 2023 09:04 EDT</lastBuildDate>
  <docs>https://www.rssboard.org/rss-specification</docs>
  <generator>Blosxom 2.1.2</generator>
  <managingEditor>neil.armstrong@example.com (Neil Armstrong)</managingEditor>
  <webMaster>sally.ride@example.com (Sally Ride)</webMaster>
  <atom:link href="https://www.rssboard.org/files/sample-rss-2.xml" rel="self" type="application/rss+xml" />
`;

    //Add content
    for i := 0; i < len(listings); i++ {
        rssPage +=
`  <item>
    <title>Louisiana Students to Hear from NASA Astronauts Aboard Space Station</title>
    <link>http://www.nasa.gov/press-release/louisiana-students-to-hear-from-nasa-astronauts-aboard-space-station</link>
    <description>As part of the state's first Earth-to-space call, students from Louisiana will have an opportunity soon to hear from NASA astronauts aboard the International Space Station.</description>
    <pubDate>Fri, 21 Jul 2023 09:04 EDT</pubDate>
    <guid>http://www.nasa.gov/press-release/louisiana-students-to-hear-from-nasa-astronauts-aboard-space-station</guid>
  </item>
`
    }


   rssPage +=
`  </channel>
</rss>
`;
	w.Header().Set("Content-Type", "application/rss+xml")
	w.Write([]byte(rssPage))

}

func GetListings(r *http.Request, sqldb *sql.DB) []db.Listing {
	agency := r.URL.Query().Get("agency")
	qPriceMin := r.URL.Query().Get("price_min")
	qPriceMax := r.URL.Query().Get("price_max")
	qYearMin := r.URL.Query().Get("build_year_min")
	qYearMax := r.URL.Query().Get("build_year_max")
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
			"agency, "+
			"url, "+
			"deleted = 1, "+
			"IFNULL(price_change.price, -1), "+
			"IFNULL(price_change.last_seen, \"\") "+
			"FROM listing "+
			"LEFT JOIN price_change on price_change.listing_id = listing.id "+
			"WHERE "+
			ResolveDeleted(qIncludeDeleted)+
			"AND agency = COALESCE(NULLIF(?, ''), agency) "+
			"AND (listing.price IS NULL OR listing.price >= COALESCE(NULLIF(?, ''), listing.price-1)) "+
			"AND (listing.price IS NULL OR listing.price <= COALESCE(NULLIF(?, ''), listing.price+1)) "+
			"AND (build_year IS NULL OR build_year >= COALESCE(NULLIF(?, ''), build_year-1)) "+
			"AND (build_year IS NULL OR build_year <= COALESCE(NULLIF(?, ''), build_year+1)) "+
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

	if err != nil {
		panic(err.Error())
	}

	listings := []db.Listing{}

	for query.Next() {
		var rowListing db.Listing
		var rowPriceChange db.PriceChange
		err := query.Scan(
			&rowListing.Id,
			&rowListing.Address,
			&rowListing.Price,
			&rowListing.Year,
			&rowListing.Size.Value,
			&rowListing.Size.Unit,
			&rowListing.PriceOverArea,
			&rowListing.Rooms,
			&rowListing.FirstSeen,
			&rowListing.LastSeen,
			&rowListing.Agency,
			&rowListing.Url,
			&rowListing.Deleted,
			&rowPriceChange.Price,
			&rowPriceChange.LastSeen)

		if err != nil {
			panic(err.Error())
		}

		listings = append(listings, rowListing)
	}

	query.Close()

	priceChanges := GetPriceChanges(listings, sqldb)

	// Add price changes to listings
	for _, priceChange := range priceChanges {
		for i, listing := range listings {
			if priceChange.ListingId == listing.Id {
				listings[i].PriceHistory = append(listing.PriceHistory, priceChange)
			}
		}
	}

	return listings

}

func GetPriceChanges(listings []db.Listing, sqldb *sql.DB) []db.PriceChange {
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

	priceChanges := []db.PriceChange{}

	for query.Next() {
		var rowPriceChange db.PriceChange
		err := query.Scan(&rowPriceChange.Price, &rowPriceChange.LastSeen, &rowPriceChange.ListingId)

		if err != nil {
			panic(err.Error())
		}
		priceChanges = append(priceChanges, rowPriceChange)
	}

	query.Close()

	return priceChanges
}

func GetLastScrape(sqldb *sql.DB) db.ScrapeEvent {
	query, qErr := sqldb.Query("SELECT date, added, updated, deleted, undeleted, total_active from scrape_event ORDER BY date DESC LIMIT 1")

	if qErr != nil {
		panic(qErr.Error())
	}

	query.Next()

	var scrapeEvent db.ScrapeEvent
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
