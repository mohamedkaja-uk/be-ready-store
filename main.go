package main

import (
	"be-ready-store/internal/handlers"
	"be-ready-store/internal/storage"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize PostgreSQL Connection Pool
	storage.InitDB()

	r := gin.Default()

	// Serve the storefront UI at the root URL
	r.StaticFile("/", "./static/index.html")

	// API routes
	v1 := r.Group("/api/v1")
	{
		v1.GET("/bundles", handlers.GetBundles)
		v1.GET("/bundles/:id", handlers.GetBundleByID)

		v1.POST("/checkout", handlers.CreateCheckout)
		v1.GET("/orders", handlers.GetOrders)
	}

	// Use dynamic port for cloud deployment, fallback to 8080 locally
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}