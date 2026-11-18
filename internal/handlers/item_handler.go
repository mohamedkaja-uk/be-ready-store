package handlers

import (
"be-ready-store/internal/storage"
"net/http"

"github.com/gin-gonic/gin"
)

func ListItems(c *gin.Context) {
items, err := storage.GetAllItems()
if err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch store items"})
return
}
c.JSON(http.StatusOK, items)
}
