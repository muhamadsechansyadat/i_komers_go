package controllers

import (
	"i_komers_go/helpers"
	"i_komers_go/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/gorm"
)

type CreateSizeInput struct {
	Name      string `json:"name" form:"name" binding:"required"`
	Quantity  uint   `json:"quantity" form:"quantity" binding:"required"`
	ProductID uint   `json:"product_id" form:"product_id" binding:"required"`
}

func GetAllSizesHandler(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var size []models.Size
	if err := db.Find(&size).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Size not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Successfully retrieve size",
		"data":    size,
	})
}

func CreateSizeHandler(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var input CreateSizeInput
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
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "Id Product is not valid", "message": err.Error()})
		return
	}

	if input.Quantity < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": true, "message": "Quantity cannot be less than zero"})
		return
	}

	size := models.Size{
		Name:      input.Name,
		Quantity:  input.Quantity,
		ProductID: input.ProductID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Menyimpan data ke database
	if err := db.Create(&size).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": true, "message": "Failed to create size"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "ok", "message": "Size created successfully", "size": size})
}

func GetSizeData(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": true, "message": "Invalid size ID"})
		return
	}

	var size models.Size
	if err := db.First(&size, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "error": true, "message": "Size not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Successfully retrieve size",
		"data":    size,
	})
}

func UpdateSizeHandler(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")
	var input CreateSizeInput

	var size models.Size
	if err := db.Where("id = ?", id).First(&size).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "error": true, "message": "Size not found"})
		return
	}

	if err := c.ShouldBind(&input); err != nil {
		if validationErr, ok := err.(validator.ValidationErrors); ok {
			errors := helpers.ParseError(validationErr)
			helpers.HelperErrorWithDataJSON(c, errors)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": true, "message": err.Error()})
		return
	}

	if input.Quantity < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": true, "message": "Quantity cannot be less than zero"})
		return
	}

	var product models.Product
	if err := db.Where("id = ?", input.ProductID).First(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "Failed to fetch product", "message": err.Error()})
		return
	}

	// Update data ukuran
	size.Name = input.Name
	size.Quantity = input.Quantity
	size.ProductID = input.ProductID
	size.UpdatedAt = time.Now()

	if err := db.Save(&size).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": true, "message": "Failed to update size"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Size updated successfully", "size": size})
}

func DeleteSizeHandler(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	id := c.Param("id")

	var size models.Size
	if err := db.Where("id = ?", id).First(&size).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "error": true, "message": "Size not found"})
		return
	}

	if err := db.Model(&size).Where("id = ?", id).Update("status", "N").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "Failed to mark size as deleted", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Size marked as deleted successfully"})
}
