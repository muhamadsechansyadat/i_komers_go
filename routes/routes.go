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

	// r.GET("/products", controllers.GetAllProductsHandler)
	// r.POST("/product", controllers.CreateProductsHandler)
	// r.GET("/product/:id", controllers.GetOneData)
	// r.GET("/uploads/images/product/:filename", controllers.LoadImageProduct)
	v1 := r.Group("/v1")
	{
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/register", controllers.RegisterHandler)
			authGroup.POST("/login", controllers.LoginHandler)
		}

		productsGroup := v1.Group("/products")
		{
			productsGroup.GET("", controllers.GetAllProductsHandler)
			productsGroup.POST("", controllers.CreateProductsHandler)
			productsGroup.GET("/:id", controllers.GetOneData)
			productsGroup.PATCH("/:id", controllers.UpdateProductHandler)
			productsGroup.DELETE("/:id", controllers.DeleteProductHandler)
		}

		uploadsGroup := v1.Group("/uploads")
		{
			uploadsGroup.GET("/images/product/:filename", controllers.LoadImageProduct)
		}
	}
	return r
}
