package test

import (
	"net/http"
	"os"
	"testing"

	"github.com/emanueldonalds/husax/db"
)

func TestQueryReturnsCorrectAmountOfListings(t *testing.T) {

	os.Setenv("HUSAX_DB_HOST", "0.0.0.0")
	os.Setenv("HUSAX_DB_PASSWORD", "abc123")

	sqldb := db.GetDb()
	req, err := http.NewRequest("GET", "http://test.com/?deleted=false", nil)
	if err != nil {
		panic(err)
	}

	listings := db.GetListings(req, sqldb)
    actual := len(listings)
    expected := 298

	if actual != expected {
		t.Fatalf(`Expected %d listings but was %d`, expected, actual)
	}
}

func TestItemWithPriceHistoryIsNotDuplicated(t *testing.T) {

	os.Setenv("HUSAX_DB_HOST", "0.0.0.0")
	os.Setenv("HUSAX_DB_PASSWORD", "abc123")

	sqldb := db.GetDb()
	req, err := http.NewRequest("GET", "http://test.com/?deleted=false", nil)
	if err != nil {
		panic(err)
	}

	listings := db.GetListings(req, sqldb)
    var items []db.Listing
    
    for i := 0; i < len(listings); i++ {
        if listings[i].Id == "322"  {
            l := listings[i]
            items = append(items, l)
        }
    }

    actual := len(items)
    expected := 1

	if actual != expected {
		t.Fatalf(`Expected %d listings but was %d`, expected, actual)
	}
}
