package controllers

import (
	"encoding/json"
	"fmt"
	"i_komers_go/helpers"
	"i_komers_go/models"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/gorm"
)

type CreateInput struct {
	Name        string  `json:"name" form:"name" binding:"required"`
	Price       float64 `json:"price" form:"price" gorm:"type:decimal(10,2);not null" binding:"required"`
	Description string  `json:"description" form:"description" binding:"required"`
	Type        string  `json:"type" form:"type" binding:"required,oneof=sale new"`
	// PhotoProduct []*string `json:"photo_product"`
}

func GetAllProductsHandler(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var product []models.Product
	if err := db.Find(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Product not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Successfully retrieve product",
		"data":    product,
	})
}

func CreateProductsHandler(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	name := c.PostForm("name")
	tableName := "products"
	form, err := c.MultipartForm()
	var input CreateInput

	// Binding form-data request and validating the required fields
	if err := c.ShouldBind(&input); err != nil {
		if validationErr, ok := err.(validator.ValidationErrors); ok {
			errors := helpers.ParseError(validationErr)
			helpers.HelperErrorWithDataJSON(c, errors)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}

	if err != nil {
		helpers.ErrorWithDataJSON(c, "errors", err.Error())
		return
	}

	// code upload files multiple
	files := form.File["photo_product[]"]
	if len(files) == 0 {
		helpers.ErrorJSON(c, "No files uploaded")
		return
	}

	// Ensure the directory exists before saving files
	if err := helpers.EnsureDir("./uploads/images/product"); err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %s", err.Error()))
		return
	}

	var filePaths []string
	for _, file := range files {
		// Generate unique filename
		newFileName := helpers.GenerateUniqueFileName(file.Filename)

		filePath := "./uploads/images/product/" + newFileName
		newPath := "/uploads/images/product/" + newFileName
		err := helpers.SaveUploadedFile(file, filePath)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %s", err.Error()))
			return
		}
		filePaths = append(filePaths, newPath)
	}

	slug := helpers.GenerateUniqueSlug(db, tableName, name)
	productType := models.ProductType(input.Type)
	photoPathsJSON, err := json.Marshal(filePaths)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %s", err.Error()))
		return
	}

	product := models.Product{
		Name:         input.Name,
		Price:        input.Price,
		Description:  input.Description,
		PhotoProduct: string(photoPathsJSON),
		Type:         productType,
		Slug:         slug,
	}

	if err := db.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product", "details": err.Error()})
		return
	}
	// code upload files multiple

	c.JSON(http.StatusOK, gin.H{"message": "Product created successfully", "product": product})
}

func GetOneData(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid product ID"})
		return
	}

	// Mengambil produk dari basis data berdasarkan ID
	var product models.Product
	if err := db.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Product not found"})
		return
	}

	// Mengirimkan respon JSON dengan data produk
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Successfully retrieve product",
		"data":    product,
	})
}

func LoadImageProduct(c *gin.Context) {
	filename := c.Param("filename")
	filePath := "./uploads/images/product/" + filename

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		defaultFilePath := "./uploads/images/default.jpg"
		c.File(defaultFilePath)
		return
	}

	c.File(filePath)
}

func UpdateProductHandler(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")
	name := c.PostForm("name")
	tableName := "products"
	form, err := c.MultipartForm()
	var input CreateInput

	// Mencari produk berdasarkan ID
	var product models.Product
	if err := db.Where("id = ?", id).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if err := c.ShouldBind(&input); err != nil {
		if validationErr, ok := err.(validator.ValidationErrors); ok {
			errors := helpers.ParseError(validationErr)
			helpers.HelperErrorWithDataJSON(c, errors)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		return
	}

	if err != nil {
		helpers.ErrorWithDataJSON(c, "errors", err.Error())
		return
	}

	// Mengambil foto produk dari form
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "Failed to process request"})
		return
	}
	files := form.File["photo_product[]"]
	var newFilePaths []string
	if len(files) > 0 {
		var photoPaths []string
		if err := json.Unmarshal([]byte(product.PhotoProduct), &photoPaths); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "Failed to parse photo paths", "details": err.Error()})
			return
		}

		for _, path := range photoPaths {
			// Periksa apakah file ada
			if _, err := os.Stat(path); os.IsNotExist(err) {
				continue // Langsung ke iterasi berikutnya jika file tidak ada
			}

			// Hapus file jika ada
			if err := os.Remove(path); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "Failed to remove old photo", "details": err.Error()})
				return
			}
		}

		// Simpan file baru ke dalam sistem file
		for _, file := range files {
			newFileName := helpers.GenerateUniqueFileName(file.Filename)
			newPath := "./uploads/images/product/" + newFileName
			if err := c.SaveUploadedFile(file, newPath); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "Failed to save new photo", "details": err.Error()})
				return
			}
			newFilePaths = append(newFilePaths, newPath)
		}
	}
	newFilePathsJSON, err := json.Marshal(newFilePaths)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "Failed to marshal new photo paths", "details": err.Error()})
		return
	}
	slug := helpers.GenerateUniqueSlug(db, tableName, name)

	// Memperbarui data produk
	product.Name = input.Name
	product.Price = input.Price
	product.Description = input.Description
	product.Type = models.ProductType(input.Type)
	product.PhotoProduct = string(newFilePathsJSON)
	product.Slug = slug

	if err := db.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "Failed to update product", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Product updated successfully", "data": product})
}

func DeleteProductHandler(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	id := c.Param("id")

	if err := db.Model(&models.Product{}).Where("id = ?", id).Update("status", "N").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": "Failed to mark product as deleted", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Product marked as deleted successfully"})
}

func GetProductWithSizesHandler(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	// Mencari produk berdasarkan ID
	var product models.Product
	if err := db.Where("id = ?", id).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "error": true, "message": "Product not found"})
		return
	}

	var sizes []models.Size
	if err := db.Where("product_id = ? AND (status = ? OR status = ? OR status IS NULL)", id, "", "Y").Find(&sizes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": true, "message": "Failed to retrieve sizes"})
		return
	}

	product.Sizes = sizes

	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Successfully retrieve product with sizes", "data": product})
}
