package handlers

import (
"be-ready-store/internal/models"
"be-ready-store/internal/storage"
"fmt"
"math"
"net/http"
"time"

"github.com/gin-gonic/gin"
)

func CreateCheckout(c *gin.Context) {
var req models.CheckoutRequest
if err := c.ShouldBindJSON(&req); err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid checkout request: " + err.Error()})
return
}

var totalAmount float64
for _, item := range req.Items {
if item.Quantity <= 0 {
c.JSON(http.StatusBadRequest, gin.H{"error": "Item quantity must be at least 1"})
return
}
bundle, err := storage.GetBundleByID(item.BundleID)
if err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Bundle ID %s not found", item.BundleID)})
return
}
totalAmount += bundle.BundlePrice * float64(item.Quantity)
}

// Check and decrement stock atomically
if err := storage.ProcessOrderStock(req); err != nil {
c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
return
}

totalAmount = math.Round(totalAmount*100) / 100

order := models.Order{
ID:            fmt.Sprintf("UK-ORD-%d", time.Now().Unix()),
CustomerName:  req.CustomerName,
CustomerEmail: req.CustomerEmail,
Items:         req.Items,
TotalAmount:   totalAmount,
Status:        "CONFIRMED",
CreatedAt:     time.Now(),
}

savedOrder := storage.SaveOrder(order)
c.JSON(http.StatusCreated, savedOrder)
}

func GetOrders(c *gin.Context) {
c.JSON(http.StatusOK, storage.GetAllOrders())
}
