package rss

import (
	"bytes"
	"database/sql"
	"fmt"
	"html"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/emanueldonalds/husax/db"
	"github.com/emanueldonalds/husax/formatters"
)

var rssTemplate = readFile("rss/template.xml")

func RssHandler(w http.ResponseWriter, r *http.Request, sqldb *sql.DB) {
	listings := db.GetListings(r, sqldb)
	lastScrape := db.GetLastScrape(sqldb)

	items := []Item{}
	for i := 0; i < len(listings); i++ {
		listing := listings[i]

		priceChanges := []string{}

		for y := 0; y < len(listing.PriceHistory); y++ {
			var hist db.PriceChange = listing.PriceHistory[y]
			priceChanges = append(priceChanges, formatters.FormatPriceChange(hist))
		}

		items = append(items, Item{
			Title: listing.Address,
			Link:  listing.Url,
			Description: fmt.Sprintf(
				html.EscapeString("<ul><li>Price: %s</li><li>Size: %s %s</li><li>Year: %s</li><li>Rooms: %s</li><li>Agency: %s</li><li>Address: %s</li><li>Price over area: %s €/m2</li><li>Price history: %s</li><li>First seen: %s</li></ul>"),
				formatters.FormatPrice(listing.Price),
				formatters.FormatInt(listing.Size.Value),
				listing.Size.Unit,
				formatters.FormatInt(listing.Year),
				formatters.FormatInt(listing.Rooms),
				listing.Agency,
				listing.Address,
				formatters.FormatInt(listing.PriceOverArea),
				strings.Join(priceChanges, ", "),
				formatters.FormatFullDate(listing.FirstSeen),
			),
			PubDate: formatters.FormatDateTimeRfc822(listing.FirstSeen),
			Guid:    listing.Id,
		})
	}

	data := Feed{
		Title:       "Hus.ax",
		Description: "Property listings on Åland RSS feed",
		PubDate:     formatters.FormatDateTimeRfc822(lastScrape.Date),
		WebMaster:   "husax@protonmail.com",
		Items:       items,
	}

	tmpl, err := template.New("feed").Parse(rssTemplate)
	if err != nil {
		panic(err)
	}

	var res bytes.Buffer
	tmpl.Execute(&res, data)

	w.Header().Set("Content-Type", "application/rss+xml")
	w.Write([]byte(res.Bytes()))
}

func readFile(filename string) string {
	fmt.Println("Loading " + filename)
	fileBytes, err := os.ReadFile(filename)
	if err != nil {
		pwd, _ := os.Getwd()
		fmt.Println("Error reading file at " + pwd + "/" + filename)
		panic(err)
	}
	content := string(fileBytes)
	return content
}
