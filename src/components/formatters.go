package components

import (
	"fmt"
	"github.com/emanueldonalds/property-viewer/db"
	"strconv"
	"strings"
	"time"
)

func formatPrice(value int) string {
	if value <= 0 {
		return ""
	}
	return formatInt(value) + " €"
}

func formatInt(value int) string {
	if value <= 0 {
		return ""
	}
	return strconv.Itoa(value)
}

func formatFloat(value float32) string {
	if value < 0 {
		return ""
	}
	return fmt.Sprintf("%.0f", value)
}

func formatDate(value string) string {
	if value == "" {
		return ""
	}
	t := parseTime(value)
	formatted := t.Format("2 Jan")
	return formatted
}

func formatDateTime(value string) string {
	if value == "" {
		return ""
	}
	t := parseTime(value)
	formatted := t.Format("2 Jan 15:04")
	return formatted
}

func parseTime(value string) time.Time {
	t, err := time.Parse("2006-01-02 15:04:05.999999", value)
	if err != nil {
		panic(err.Error())
	}
	return t.In(time.Local)
}

func formatPrevPrice(priceHistory []db.PriceChange) string {
	if len(priceHistory) == 0 {
		return ""
	}
	var lastPrice = priceHistory[0]
	return strings.TrimSpace(fmt.Sprintf("%s € (%s)", formatInt(lastPrice.Price), formatDate(lastPrice.LastSeen)))
}
