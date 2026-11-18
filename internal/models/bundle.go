package models

import "time"

type Item struct {
	ID                string    `json:"id"`
	BundleID          string    `json:"bundle_id"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	IndividualPrice   float64   `json:"individual_price"`
	ImageURL          string    `json:"image_url"`
	AmazonSearchQuery string    `json:"amazon_search_query"`
	CachedASIN        string    `json:"cached_asin"`
	AsinUpdatedAt     time.Time `json:"asin_updated_at"`
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