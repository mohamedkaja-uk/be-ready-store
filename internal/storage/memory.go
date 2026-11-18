package storage

import (
	"be-ready-store/internal/models"
	"context"
	"fmt"
	"time"
)

// GetAllBundles returns all bundles with their associated items using pgx
func GetAllBundles() ([]models.Bundle, error) {
	if DB == nil {
		return nil, fmt.Errorf("database connection not initialized")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id, title, category, tagline, bundle_price, stock, image_url FROM bundles`
	rows, err := DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bundles []models.Bundle
	for rows.Next() {
		var b models.Bundle
		var imageURL *string
		if err := rows.Scan(&b.ID, &b.Title, &b.Category, &b.Tagline, &b.BundlePrice, &b.Stock, &imageURL); err != nil {
			return nil, err
		}
		if imageURL != nil {
			b.ImageURL = *imageURL
		}

		itemQuery := `
			SELECT i.id, i.name, i.category, i.description, i.individual_price, i.image_url, i.amazon_search_query, COALESCE(i.cached_asin, '')
			FROM items i
			JOIN bundle_items bi ON i.id = bi.item_id
			WHERE bi.bundle_id = $1`
		
		itemRows, err := DB.Query(ctx, itemQuery, b.ID)
		if err != nil {
			return nil, err
		}
		
		var items []models.Item
		for itemRows.Next() {
			var item models.Item
			var itemImg, searchQuery, cachedASIN *string
			if err := itemRows.Scan(&item.ID, &item.Name, &item.Category, &item.Description, &item.IndividualPrice, &itemImg, &searchQuery, &cachedASIN); err != nil {
				itemRows.Close()
				return nil, err
			}
			if itemImg != nil { item.ImageURL = *itemImg }
			if searchQuery != nil { item.AmazonSearchQuery = *searchQuery }
			if cachedASIN != nil { item.CachedASIN = *cachedASIN }

			items = append(items, item)
		}
		itemRows.Close()

		b.Items = items
		bundles = append(bundles, b)
	}

	if bundles == nil {
		bundles = []models.Bundle{}
	}
	return bundles, nil
}

// GetBundleByID returns a single bundle by ID using pgx
func GetBundleByID(bundleID string) (models.Bundle, error) {
	if DB == nil {
		return models.Bundle{}, fmt.Errorf("database connection not initialized")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	bundleQuery := `
	SELECT id, title, category, tagline, bundle_price, stock, image_url
	FROM bundles
	WHERE id = $1
	`
	var b models.Bundle
	var imageURL *string
	err := DB.QueryRow(ctx, bundleQuery, bundleID).Scan(
		&b.ID, &b.Title, &b.Category, &b.Tagline, &b.BundlePrice, &b.Stock, &imageURL,
	)
	if err != nil {
		return models.Bundle{}, err
	}
	if imageURL != nil {
		b.ImageURL = *imageURL
	}

	items, err := getItemsForBundle(ctx, b.ID)
	if err != nil {
		return models.Bundle{}, err
	}
	b.Items = items

	return b, nil
}

// GetAllItems returns all items using pgx
func GetAllItems() ([]models.Item, error) {
	if DB == nil {
		return nil, fmt.Errorf("database connection not initialized")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id, name, category, description, individual_price, image_url, amazon_search_query, COALESCE(cached_asin, '') FROM items`
	rows, err := DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var item models.Item
		var itemImg, searchQuery, cachedASIN *string
		if err := rows.Scan(&item.ID, &item.Name, &item.Category, &item.Description, &item.IndividualPrice, &itemImg, &searchQuery, &cachedASIN); err != nil {
			return nil, err
		}
		if itemImg != nil { item.ImageURL = *itemImg }
		if searchQuery != nil { item.AmazonSearchQuery = *searchQuery }
		if cachedASIN != nil { item.CachedASIN = *cachedASIN }
		items = append(items, item)
	}
	if items == nil {
		items = []models.Item{}
	}
	return items, nil
}

func CreateBundle(b models.Bundle) error {
	return nil
}

func getItemsForBundle(ctx context.Context, bundleID string) ([]models.Item, error) {
	itemQuery := `
	SELECT i.id, i.name, i.category, i.description, i.individual_price, i.image_url, i.amazon_search_query, COALESCE(i.cached_asin, '')
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
		var itemImg, searchQuery, cachedASIN *string
		if err := rows.Scan(
			&it.ID,
			&it.Name,
			&it.Category,
			&it.Description,
			&it.IndividualPrice,
			&itemImg,
			&searchQuery,
			&cachedASIN,
		); err != nil {
			return nil, err
		}
		if itemImg != nil { it.ImageURL = *itemImg }
		if searchQuery != nil { it.AmazonSearchQuery = *searchQuery }
		if cachedASIN != nil { it.CachedASIN = *cachedASIN }
		
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
	if DB == nil {
		return models.Order{}, fmt.Errorf("database connection not initialized")
	}
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