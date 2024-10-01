package formatters

import (
	"fmt"
	"github.com/emanueldonalds/property-viewer/db"
	"strconv"
	"strings"
	"time"
)

func FormatPrice(value int) string {
	if value <= 0 {
		return ""
	}
	return FormatInt(value) + " €"
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
	formatted := t.Format("2 Jan")
	return formatted
}

func FormatFullDate(value string) string {
	if value == "" {
		return ""
	}
	t := parseTime(value)
	formatted := t.Format("2 Jan 2006")
	return formatted
}

func FormatDateTime(value string) string {
	if value == "" {
		return ""
	}
	t := parseTime(value)
	formatted := t.Format("2 Jan 15:04")
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
	return FormatPriceChange(lastPrice);
}

func FormatPriceChange(priceChange db.PriceChange) string {
	return strings.TrimSpace(fmt.Sprintf("%s € (%s)", FormatInt(priceChange.Price), FormatDate(priceChange.LastSeen)))
}
