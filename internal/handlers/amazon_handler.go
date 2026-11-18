package handlers

import (
	"be-ready-store/internal/amazon"
	"be-ready-store/internal/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AmazonCheckoutRequest struct {
	BundleID string `json:"bundle_id" binding:"required"`
}

func CreateAmazonCheckout(c *gin.Context) {
	var req AmazonCheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bundle_id is required"})
		return
	}

	bundle, found := storage.GetBundleByID(req.BundleID)
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bundle not found"})
		return
	}

	cartURL, itemLinks := amazon.BuildAmazonCheckoutResponse(c.Request.Context(), bundle.Items)

	c.JSON(http.StatusOK, gin.H{
		"bundle_title":   bundle.Title,
		"multi_cart_url": cartURL,   // Populated automatically when ASINs are cached/synced
		"item_links":     itemLinks, // Always populated (direct ASIN link or dynamic UK search link)
	})
}
