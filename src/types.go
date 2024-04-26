package main

type Listing struct {
    id string
    agency string
    address string
    price int
    url string
    size Size
    priceOverArea int
    rooms int
    year int
    priceHistory []PriceChange
    firstSeen string
    lastSeen string
}

type Size struct {
    value int
    unit string
}

type PriceChange struct {
    lastSeen string
    price int
    listingId string
}

type ScrapeEvent struct {
    date string
    added int
    updated int
    deleted int
    undeleted int
    totalActive int
}
