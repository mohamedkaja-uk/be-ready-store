package amazon

import (
	"be-ready-store/internal/models"
	"be-ready-store/internal/storage"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

type tokenCache struct {
	mu          sync.Mutex
	accessToken string
	expiresAt   time.Time
}

var cache tokenCache

// getAccessToken fetches and caches an OAuth 2.0 Bearer token for Amazon Creators API
func getAccessToken(clientID, clientSecret string) (string, error) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	if cache.accessToken != "" && time.Now().Add(60*time.Second).Before(cache.expiresAt) {
		return cache.accessToken, nil
	}

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("scope", "creatorsapi::default")

	req, err := http.NewRequest("POST", "https://api.amazon.co.uk/auth/o2/token", bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil || res.AccessToken == "" {
		return "", fmt.Errorf("failed to obtain Creators API token")
	}

	cache.accessToken = res.AccessToken
	cache.expiresAt = time.Now().Add(time.Duration(res.ExpiresIn) * time.Second)
	return cache.accessToken, nil
}

// RefreshItemASIN queries Amazon Creators API (UK) for the top in-stock ASIN and caches it in Neon
func RefreshItemASIN(ctx context.Context, item *models.Item, partnerTag string) string {
	clientID := os.Getenv("AMAZON_CREATORS_CLIENT_ID")
	clientSecret := os.Getenv("AMAZON_CREATORS_CLIENT_SECRET")

	// If Creators API credentials aren't configured yet, return any existing cached ASIN
	if clientID == "" || clientSecret == "" {
		return item.CachedASIN
	}

	token, err := getAccessToken(clientID, clientSecret)
	if err != nil {
		return item.CachedASIN
	}

	payload := map[string]interface{}{
		"keywords":    item.AmazonSearchQuery,
		"searchIndex": "All",
		"itemCount":   1,
		"partnerTag":  partnerTag,
		"marketplace": "www.amazon.co.uk",
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, "POST", "https://creatorsapi.amazon/catalog/v1/searchItems", bytes.NewBuffer(body))
	if err != nil {
		return item.CachedASIN
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-marketplace", "www.amazon.co.uk")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return item.CachedASIN
	}
	defer resp.Body.Close()

	var result struct {
		SearchResult struct {
			Items []struct {
				ASIN string `json:"asin"`
			} `json:"items"`
		} `json:"searchResult"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err == nil && len(result.SearchResult.Items) > 0 {
		newASIN := result.SearchResult.Items[0].ASIN
		_, _ = storage.DB.Exec(ctx,
			"UPDATE items SET cached_asin = $1, asin_updated_at = NOW() WHERE id = $2",
			newASIN, item.ID,
		)
		item.CachedASIN = newASIN
	}

	return item.CachedASIN
}

// BuildAmazonCheckoutResponse generates a 1-click Amazon UK multi-item cart URL + individual item fallback links
func BuildAmazonCheckoutResponse(ctx context.Context, items []models.Item) (string, []map[string]string) {
	partnerTag := os.Getenv("AMAZON_ASSOCIATE_TAG")
	if partnerTag == "" {
		partnerTag = "bereadyuk-21" // Replace with your UK Associate tag in Render Env Vars
	}

	cartParams := url.Values{}
	cartParams.Set("AssociateTag", partnerTag)

	asinCount := 0
	var itemLinks []map[string]string

	for _, item := range items {
		asin := item.CachedASIN
		if asin == "" {
			asin = RefreshItemASIN(ctx, &item, partnerTag)
		}

		var itemURL string
		if asin != "" {
			asinCount++
			cartParams.Set(fmt.Sprintf("ASIN.%d", asinCount), asin)
			cartParams.Set(fmt.Sprintf("Quantity.%d", asinCount), "1")
			itemURL = fmt.Sprintf("https://www.amazon.co.uk/dp/%s?tag=%s", asin, partnerTag)
		} else {
			// Dynamic tagged search fallback when ASIN isn't cached yet
			itemURL = fmt.Sprintf("https://www.amazon.co.uk/s?k=%s&tag=%s", url.QueryEscape(item.AmazonSearchQuery), partnerTag)
		}

		itemLinks = append(itemLinks, map[string]string{
			"name": item.Name,
			"url":  itemURL,
		})
	}

	var multiCartURL string
	if asinCount > 0 {
		multiCartURL = "https://www.amazon.co.uk/gp/aws/cart/add.html?" + cartParams.Encode()
	}

	return multiCartURL, itemLinks
}
