package controllers

import (
	"i_komers_go/helpers"
	"i_komers_go/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type CreateInput struct {
	Name string `json:"name" form:"name" binding:"required"`
	// End  time.Time `json:"end" form:"end" binding:"required_without=Start,omitempty,gt|gtfield=Start"`
}

func GetAllProductsHandler(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var products []models.Product
	db.Find(&products)
	helpers.SuccessJSON(c, "products", products)
}

func CreateProductsHandler(c *gin.Context) {
	// db := c.MustGet("db").(*gorm.DB)
	// name := c.PostForm("name")
	// tableName := "products"

	var input CreateInput

	// Binding form-data request and validating the required fields
	if err := c.ShouldBind(&input); err != nil {
		errors := helpers.ParseError(err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	// If the slug is not provided, generate it based on the name
	// if product.Slug == "" {
	// 	product.Slug = helpers.GenerateUniqueSlug(db, tableName, name)
	// }

	// helpers.ErrorWithDataJSON(c, "product", product)

	// Proses file-file yang diunggah
	// files := c.Request.MultipartForm.File["photo_product"]
	// var photoURLs []string

	// for _, file := range files {
	// 	// Validasi jenis berkas
	// 	if err := validateFile(file); err != nil {
	// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 		return
	// 	}

	// 	// Simpan berkas di server atau lakukan tugas lain yang diperlukan
	// 	err := c.SaveUploadedFile(file, filepath.Join("./uploads", file.Filename))
	// 	if err != nil {
	// 		c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %s", err.Error()))
	// 		return
	// 	}

	// 	// Simpan URL foto ke dalam slice
	// 	photoURLs = append(photoURLs, "/uploads/"+file.Filename)
	// }

	// // Simpan data ke dalam database
	// product := models.Product{
	// 	Name:         name,
	// 	PhotoProduct: toJSON(photoURLs),
	// }
	// db.Create(&product)

	// // Kirim respons yang mencakup data produk yang baru saja diunggah
	// c.JSON(http.StatusOK, gin.H{
	// 	"message": "Product uploaded successfully",
	// 	"product": gin.H{
	// 		"id":           product.ID,
	// 		"name":         product.Name,
	// 		"photoProduct": parseJSON(product.PhotoProduct),
	// 	},
	// })
}
