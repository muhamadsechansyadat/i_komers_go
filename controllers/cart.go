package controllers

import (
	"i_komers_go/helpers"
	"i_komers_go/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/gorm"
)

type AddToCartRequest struct {
	Quantity  uint `json:"quantity" form:"quantity" binding:"required"`
	ProductID uint `json:"product_id" form:"product_id" binding:"required"`
	SizeID    uint `json:"size_id" form:"size_id" binding:"required"`
}

func AddToCartHandler(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User data not found in context"})
		return
	}

	userData, ok := user.(*models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user data"})
		return
	}

	var input AddToCartRequest
	if err := c.ShouldBind(&input); err != nil {
		if validationErr, ok := err.(validator.ValidationErrors); ok {
			errors := helpers.ParseError(validationErr)
			helpers.HelperErrorWithDataJSON(c, errors)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}

	var product models.Product
	if err := db.Where("id = ?", input.ProductID).First(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": true, "message": "Id Product not found"})
		return
	}

	var size models.Size
	if err := db.Where("id = ? AND product_id = ?", input.SizeID, input.ProductID).First(&size).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": true, "message": "Product With Size not found"})
		return
	}

	if input.Quantity < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": true, "message": "Quantity cannot be less than zero"})
		return
	}

	var existCart models.Cart
	if err := db.Where("user_id = ? AND product_id = ? AND size_id = ? AND (status = '' OR status IS NULL)", userData.ID, input.ProductID, input.SizeID).First(&existCart).Error; err != nil {
		cart := models.Cart{
			UserID:    userData.ID,
			Quantity:  input.Quantity,
			ProductID: input.ProductID,
			SizeID:    input.SizeID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Menyimpan data ke database
		if err := db.Create(&cart).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": true, "message": "Failed to create cart"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"status": "ok", "message": "Cart created successfully", "cart": cart})
	} else {
		existCart.UserID = userData.ID
		existCart.Quantity = (existCart.Quantity + input.Quantity)
		existCart.ProductID = input.ProductID
		existCart.SizeID = input.SizeID
		existCart.UpdatedAt = time.Now()

		if err := db.Save(&existCart).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": true, "message": "Failed to update cart"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Cart updated successfully", "cart": existCart})
	}
}

func GetAllCartsHandler(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "error": true, "message": "User data not found in context"})
		return
	}

	userData, ok := user.(*models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": true, "message": "Failed to parse user data"})
		return
	}

	var cart []models.Cart
	if err := db.Preload("User").Preload("Product").Preload("Size").Where("user_id = ? AND (status = '' OR status IS NULL)", userData.ID).Find(&cart).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "error": true, "message": "Cart not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Successfully retrieve cart",
		"data":    cart,
	})
}

func CountAllCartsHandler(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "error": true, "message": "User data not found in context"})
		return
	}

	userData, ok := user.(*models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": true, "message": "Failed to parse user data"})
		return
	}

	var count int64
	if err := db.Model(&models.Cart{}).Where("user_id = ? AND (status = '' OR status IS NULL)", userData.ID).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": true, "message": "Failed to count carts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Successfully retrieved cart count",
		"data":    count,
	})
}

func DeleteCartsHandler(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	id := c.Param("id")

	var cart models.Cart
	if err := db.Where("id = ? AND (status = '' OR status IS NULL)", id).First(&cart).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "error": true, "message": "Cart not found"})
		return
	}

	if err := db.Model(&cart).Where("id = ?", id).Update("status", "N").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "Failed to mark Cart as deleted", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Cart marked as deleted successfully"})
}
