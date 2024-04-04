package controllers

import (
	"i_komers_go/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func AddOrderFromSelectedItemsHandler(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Ambil user dari konteks
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "error": true, "message": "User data not found in context"})
		return
	}

	// Ambil ID pengguna dari data pengguna
	userData, ok := user.(*models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": true, "message": "Failed to parse user data"})
		return
	}

	// Ambil daftar item yang dipilih dari request JSON
	var selectedItems []uint
	if err := c.BindJSON(&selectedItems); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": true, "message": "Invalid request body"})
		return
	}

	// Ambil item keranjang pengguna yang dipilih
	var carts []models.Cart
	if err := db.Where("user_id = ? AND (status = '' OR status IS NULL) AND id IN (?)", userData.ID, selectedItems).Find(&carts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": true, "message": "Failed to retrieve user's selected carts"})
		return
	}

	// Buat pesanan baru dan salin detail dari item keranjang yang dipilih
	order := models.Order{
		UserID: userData.ID,
		Items:  make([]models.OrderItem, len(carts)),
	}

	for i, cart := range carts {
		order.Items[i] = models.OrderItem{
			ProductID: cart.ProductID,
			SizeID:    cart.SizeID,
			Quantity:  cart.Quantity,
			Price:     calculatePrice(cart.ProductID, cart.SizeID, cart.Quantity),
		}
		order.Total += order.Items[i].Price
	}

	// Simpan pesanan ke dalam basis data
	if err := db.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": true, "message": "Failed to create order"})
		return
	}

	// Set status item keranjang yang dipilih menjadi "completed" atau hapus item keranjang dari basis data
	if err := db.Model(&models.Cart{}).Where("user_id = ? AND status = ? AND id IN (?)", userData.ID, "active", selectedItems).Update("status", "N").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": true, "message": "Failed to update selected carts status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Order successfully created from selected items",
		"data":    order,
	})
}

func calculatePrice(productID, sizeID, quantity uint) float64 {
	// Implementasi logika perhitungan harga di sini
	return 10.00 // Contoh sederhana
}
