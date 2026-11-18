package storage

import (
"be-ready-store/internal/models"
"context"
"fmt"
"time"
)

// GetAllBundles returns just the slice to match single-value context in bundle_handler.go
func GetAllBundles() []models.Bundle {
bundles, err := GetBundles()
if err != nil {
return []models.Bundle{}
}
return bundles
}

func GetBundles() ([]models.Bundle, error) {
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

bundleQuery := `
SELECT id, title, category, tagline, bundle_price, stock
FROM bundles
ORDER BY bundle_price ASC
`
rows, err := DB.Query(ctx, bundleQuery)
if err != nil {
return nil, err
}
defer rows.Close()

var bundles []models.Bundle
for rows.Next() {
var b models.Bundle
if err := rows.Scan(&b.ID, &b.Title, &b.Category, &b.Tagline, &b.BundlePrice, &b.Stock); err != nil {
return nil, err
}

items, err := getItemsForBundle(ctx, b.ID)
if err != nil {
return nil, err
}
b.Items = items
bundles = append(bundles, b)
}

return bundles, nil
}

// GetBundleByID returns (models.Bundle, error) so `if err != nil` checks succeed in handlers
func GetBundleByID(bundleID string) (models.Bundle, error) {
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

bundleQuery := `
SELECT id, title, category, tagline, bundle_price, stock
FROM bundles
WHERE id = $1
`
var b models.Bundle
err := DB.QueryRow(ctx, bundleQuery, bundleID).Scan(
&b.ID, &b.Title, &b.Category, &b.Tagline, &b.BundlePrice, &b.Stock,
)
if err != nil {
return models.Bundle{}, err
}

items, err := getItemsForBundle(ctx, b.ID)
if err != nil {
return models.Bundle{}, err
}
b.Items = items

return b, nil
}

func CreateBundle(b models.Bundle) error {
return nil
}

func getItemsForBundle(ctx context.Context, bundleID string) ([]models.Item, error) {
itemQuery := `
SELECT i.id, i.name, i.description, i.individual_price, i.image_url, i.amazon_search_query, COALESCE(i.cached_asin, '')
FROM items i
JOIN bundle_items bi ON i.id = bi.item_id
WHERE bi.bundle_id = $1
ORDER BY i.id ASC
`
rows, err := DB.Query(ctx, itemQuery, bundleID)
if err != nil {
return nil, err
}
defer rows.Close()

var items []models.Item
for rows.Next() {
var it models.Item
it.BundleID = bundleID
if err := rows.Scan(
&it.ID,
&it.Name,
&it.Description,
&it.IndividualPrice,
&it.ImageURL,
&it.AmazonSearchQuery,
&it.CachedASIN,
); err != nil {
return nil, err
}
items = append(items, it)
}

return items, nil
}

func ProcessOrderStock(req models.CheckoutRequest) error {
_, err := ProcessOrder(req)
return err
}

func SaveOrder(order models.Order) error {
return nil
}

func GetAllOrders() []models.Order {
return []models.Order{}
}

func ProcessOrder(req models.CheckoutRequest) (models.Order, error) {
ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
defer cancel()

tx, err := DB.Begin(ctx)
if err != nil {
return models.Order{}, err
}
defer tx.Rollback(ctx)

var totalAmount float64
for _, cartItem := range req.Items {
var currentStock int
var price float64

lockQuery := `SELECT stock, bundle_price FROM bundles WHERE id = $1 FOR UPDATE`
if err := tx.QueryRow(ctx, lockQuery, cartItem.BundleID).Scan(&currentStock, &price); err != nil {
return models.Order{}, fmt.Errorf("bundle %s not found", cartItem.BundleID)
}

if currentStock < cartItem.Quantity {
return models.Order{}, fmt.Errorf("insufficient stock for bundle %s", cartItem.BundleID)
}

updateQuery := `UPDATE bundles SET stock = stock - $1 WHERE id = $2`
if _, err := tx.Exec(ctx, updateQuery, cartItem.Quantity, cartItem.BundleID); err != nil {
return models.Order{}, err
}

totalAmount += price * float64(cartItem.Quantity)
}

orderID := fmt.Sprintf("ORD-%d", time.Now().UnixNano()/1e6)
insertOrderQuery := `
INSERT INTO orders (id, customer_name, customer_email, shipping_address, total_amount, status, created_at)
VALUES ($1, $2, $3, $4, $5, $6, NOW())
`
if _, err := tx.Exec(ctx, insertOrderQuery,
orderID, req.CustomerName, req.CustomerEmail, req.ShippingAddress, totalAmount, "Confirmed",
); err != nil {
return models.Order{}, err
}

insertOrderItemQuery := `
INSERT INTO order_items (order_id, bundle_id, quantity)
VALUES ($1, $2, $3)
`
for _, cartItem := range req.Items {
if _, err := tx.Exec(ctx, insertOrderItemQuery, orderID, cartItem.BundleID, cartItem.Quantity); err != nil {
return models.Order{}, err
}
}

if err := tx.Commit(ctx); err != nil {
return models.Order{}, err
}

return models.Order{
ID:              orderID,
CustomerName:    req.CustomerName,
CustomerEmail:   req.CustomerEmail,
ShippingAddress: req.ShippingAddress,
Items:           req.Items,
TotalAmount:     totalAmount,
Status:          "Confirmed",
CreatedAt:       time.Now(),
}, nil
}
