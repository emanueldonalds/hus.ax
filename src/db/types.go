package db

type Listing struct {
	Id            string
	Agency        string
	Address       string
	Price         int
	Url           string
	Size          Size
	PriceOverArea int
	Rooms         int
	Year          int
	PriceHistory  []PriceChange
	FirstSeen     string
	LastSeen      string
	Deleted       bool
}

type Size struct {
	Value int
	Unit  string
}

type PriceChange struct {
	LastSeen  string
	Price     int
	ListingId string
}

type ScrapeEvent struct {
	Date        string
	Added       int
	Updated     int
	Deleted     int
	Undeleted   int
	TotalActive int
}
