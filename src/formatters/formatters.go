package formatters

import (
	"fmt"
	"github.com/emanueldonalds/property-viewer/db"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"strconv"
	"strings"
	"time"
)

var swedishMonths map[time.Month]string = map[time.Month]string{
	time.January:   "Jan",
	time.February:  "Feb",
	time.March:     "Mar",
	time.April:     "Apr",
	time.May:       "Maj",
	time.June:      "Jun",
	time.July:      "Jul",
	time.August:    "Aug",
	time.September: "Sep",
	time.October:   "Okt",
	time.November:  "Nov",
	time.December:  "Dec",
}

func FormatPrice(value int) string {
	if value <= 0 {
		return ""
	}
	var res = message.NewPrinter(language.Swedish).Sprintf("%d €", value)
	return res
}

func FormatInt(value int) string {
	if value <= 0 {
		return ""
	}
	return strconv.Itoa(value)
}

func FormaFloat(value float32) string {
	if value < 0 {
		return ""
	}
	return fmt.Sprintf("%.0f", value)
}

func FormatDate(value string) string {
	if value == "" {
		return ""
	}
	t := parseTime(value)
	formatted := fmt.Sprintf("%d %s %d", t.Day(), swedishMonths[t.Month()], t.Year())
	return formatted
}

func FormatFullDate(value string) string {
	if value == "" {
		return ""
	}
	t := parseTime(value)
	formatted := fmt.Sprintf("%d %s %d", t.Day(), swedishMonths[t.Month()], t.Year())
	return formatted
}

func FormatDateTime(value string) string {
	if value == "" {
		return ""
	}
	t := parseTime(value)
	formatted := fmt.Sprintf("%d %s %02d:%02d", t.Day(), swedishMonths[t.Month()], t.Hour(), t.Minute())
	return formatted
}

func FormatDateTimeRfc822(value string) string {
	if value == "" {
		return ""
	}
	t := parseTime(value)
	formatted := t.Format("Mon, 02 Jan 2006 15:04:05 -0700")
	return formatted
}

func parseTime(value string) time.Time {
	t, err := time.Parse("2006-01-02 15:04:05.999999", value)
	if err != nil {
		panic(err.Error())
	}
	return t.In(time.Local)
}

func FormatPrevPrice(priceHistory []db.PriceChange) string {
	if len(priceHistory) == 0 {
		return ""
	}
	var lastPrice = priceHistory[0]
	return FormatPriceChange(lastPrice)
}

func FormatPriceChange(priceChange db.PriceChange) string {
	return strings.TrimSpace(fmt.Sprintf("%s (%s)", FormatPrice(priceChange.Price), FormatDate(priceChange.LastSeen)))
}

func FormatBool(value bool) string {
  if (value == true) {
      return "Ja"
  } else {
      return "Nej"
  }
}
