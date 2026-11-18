package storage

import (
"be-ready-store/internal/models"
"context"
"errors"
"fmt"
)

func CreateBundle(b models.Bundle) models.Bundle {
	ctx := context.Background()
	_, err := DB.Exec(ctx, "INSERT INTO bundles (id, title, category, tagline, bundle_price, stock) VALUES ($1, $2, $3, $4, $5)",
		b.ID, b.Title, b.Category, b.Tagline, b.BundlePrice, b.Stock)
	if err != nil {
		fmt.Println("Error creating bundle:", err)
		return b
	}

	for _, item := range b.Items {
		_, _ = DB.Exec(ctx, "INSERT INTO items (id, bundle_id, name, description, individual_price, image_url) VALUES ($1, $2, $3, $4, $5, $6)",
			item.ID, b.ID, item.Name, item.Description, item.IndividualPrice, item.ImageURL)
	}

	return b
}

func GetAllBundles() []models.Bundle {
ctx := context.Background()
rows, err := DB.Query(ctx, "SELECT id, title, category, tagline, bundle_price, stock FROM bundles")
if err != nil {
fmt.Println("Error querying bundles:", err)
return []models.Bundle{}
}
defer rows.Close()

var bundles []models.Bundle
for rows.Next() {
var b models.Bundle
if err := rows.Scan(&b.ID, &b.Title, &b.Category, &b.Tagline, &b.BundlePrice, &b.Stock); err != nil {
continue
}
b.Items = GetItemsForBundle(b.ID)
bundles = append(bundles, b)
}
return bundles
}

func GetItemsForBundle(bundleID string) []models.Item {
ctx := context.Background()
rows, err := DB.Query(ctx, "SELECT id, name, description, individual_price, image_url FROM items WHERE bundle_id = $1", bundleID)
if err != nil {
return []models.Item{}
}
defer rows.Close()

var items []models.Item
for rows.Next() {
var i models.Item
if err := rows.Scan(&i.ID, &i.Name, &i.Description, &i.IndividualPrice, &i.ImageURL); err == nil {
items = append(items, i)
}
}
return items
}

func GetBundleByID(id string) (models.Bundle, error) {
ctx := context.Background()
var b models.Bundle
err := DB.QueryRow(ctx, "SELECT id, title, category, tagline, bundle_price, stock FROM bundles WHERE id = $1", id).
Scan(&b.ID, &b.Title, &b.Category, &b.Tagline, &b.BundlePrice, &b.Stock)
if err != nil {
return models.Bundle{}, errors.New("bundle not found")
}
b.Items = GetItemsForBundle(b.ID)
return b, nil
}

func ProcessOrderStock(items []models.OrderItem) error {
ctx := context.Background()
tx, err := DB.Begin(ctx)
if err != nil {
return err
}
defer tx.Rollback(ctx)

for _, item := range items {
var stock int
err := tx.QueryRow(ctx, "SELECT stock FROM bundles WHERE id = $1 FOR UPDATE", item.BundleID).Scan(&stock)
if err != nil {
return fmt.Errorf("bundle ID %s not found", item.BundleID)
}
if stock < item.Quantity {
return fmt.Errorf("insufficient stock (available: %d)", stock)
}

_, err = tx.Exec(ctx, "UPDATE bundles SET stock = stock - $1 WHERE id = $2", item.Quantity, item.BundleID)
if err != nil {
return err
}
}

return tx.Commit(ctx)
}

func SaveOrder(o models.Order) models.Order {
ctx := context.Background()
_, err := DB.Exec(ctx, "INSERT INTO orders (id, customer_name, customer_email, total_amount, status) VALUES ($1, $2, $3, $4, $5)",
o.ID, o.CustomerName, o.CustomerEmail, o.TotalAmount, o.Status)
if err != nil {
fmt.Println("Error saving order:", err)
return o
}

for _, item := range o.Items {
_, _ = DB.Exec(ctx, "INSERT INTO order_items (order_id, bundle_id, quantity) VALUES ($1, $2, $3)",
o.ID, item.BundleID, item.Quantity)
}

return o
}

func GetAllOrders() []models.Order {
ctx := context.Background()
rows, err := DB.Query(ctx, "SELECT id, customer_name, customer_email, total_amount, status, created_at FROM orders")
if err != nil {
return []models.Order{}
}
defer rows.Close()

var orders []models.Order
for rows.Next() {
var o models.Order
if err := rows.Scan(&o.ID, &o.CustomerName, &o.CustomerEmail, &o.TotalAmount, &o.Status, &o.CreatedAt); err == nil {
orders = append(orders, o)
}
}
return orders
}
