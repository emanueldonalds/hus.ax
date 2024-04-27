package main

import (
	"net/http"
	"testing"
)

var db = GetDb()
var req, err = http.NewRequest(http.MethodGet, "/filter?order_by=price&sort_order=desc&price_min=&price_max=&build_year_min=&build_year_max=&size_value_min=&size_value_max=&price_over_area_min=&price_over_area_max=&agency=&include_deleted=true", nil)
var listings = GetListings(req, db)

func BenchmarkGetListings(B *testing.B) {

	GetListings(req, db)
	//l := GetListings(req, db)
	//	println(len(l))
}

func BenchmarkGetLastScrape(B *testing.B) {
	GetLastScrape(db)
}

func BenchmarkGetPriceChanges(B *testing.B) {
	GetPriceChanges(listings, db)
}
