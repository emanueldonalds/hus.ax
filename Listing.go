package main

type Listing struct {
    agency string
    address string
    price int
    url string
    size Size
    rooms int
    year int
    priceHistory []PriceChange
    firstSeen string
    lastSeen string
}

type Size struct {
    value float32
    unit string
}

type PriceChange struct {
    lastSeen string
    price int
}
