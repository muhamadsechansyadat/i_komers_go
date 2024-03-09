package routes

import (
	"i_komers_go/controllers"
	// "i_komers_go/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func SetupRoutes(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	gin.SetMode(gin.DebugMode)
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
	})

	r.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		// Send a JSON response with the message "Hello, world"
		c.JSON(200, gin.H{"message": "Hello, world"})
	})
	r.POST("/register", controllers.RegisterHandler)
	r.POST("/login", controllers.LoginHandler)

	r.GET("/products", controllers.GetAllProductsHandler)
	r.POST("/product", controllers.CreateProductsHandler)
	return r
}
