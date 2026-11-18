package models

//import "time"

type Item struct {
	ID                string  `json:"id"`
	Name              string  `json:"name"`
	Category          string  `json:"category"` // <-- Add this field
	Description       string  `json:"description"`
	IndividualPrice   float64 `json:"individual_price"`
	ImageURL          string  `json:"image_url"`
	AmazonSearchQuery string  `json:"amazon_search_query"`
	CachedASIN        string  `json:"cached_asin"`
}

type Bundle struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Category    string  `json:"category"`
	Tagline     string  `json:"tagline"`
	BundlePrice float64 `json:"bundle_price"`
	Stock       int     `json:"stock"`
	Items       []Item  `json:"items"`
}
