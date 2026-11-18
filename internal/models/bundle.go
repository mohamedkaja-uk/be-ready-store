package models

type Item struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	IndividualPrice float64 `json:"individual_price"`
	ImageURL        string  `json:"image_url"`
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