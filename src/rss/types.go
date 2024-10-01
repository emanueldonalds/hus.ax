package rss

type Feed struct {
	Title       string
	Description string
	PubDate     string
	WebMaster   string
	Items       []Item
}

type Item struct {
	Title       string
	Link        string
	Description string
	PubDate     string
	Guid        string
}
