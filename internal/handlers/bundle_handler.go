package handlers

import (
	"net/http"
	"be-ready-store/internal/models"
	"be-ready-store/internal/storage"
	"github.com/gin-gonic/gin"
)

func GetBundles(c *gin.Context) {
	bundles, err := storage.GetAllBundles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bundles)
}

func GetBundleByID(c *gin.Context) {
	bundle, err := storage.GetBundleByID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bundle)
}

func CreateBundle(c *gin.Context) {
	var newBundle models.Bundle
	if err := c.ShouldBindJSON(&newBundle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Assuming CreateBundle returns an error or success message
	err := storage.CreateBundle(newBundle)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, newBundle)
}