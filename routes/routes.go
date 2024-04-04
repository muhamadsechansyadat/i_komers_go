package routes

import (
	"i_komers_go/controllers"
	"i_komers_go/middleware"

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
		c.JSON(200, gin.H{"message": "Hello, world"})
	})

	v1 := r.Group("/v1")
	{
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/register", controllers.RegisterHandler)
			authGroup.POST("/login", controllers.LoginHandler)
		}

		productsGroup := v1.Group("/products")
		productsGroup.Use(middleware.AuthMiddleware())
		{
			productsGroup.GET("", controllers.GetAllProductsHandler)
			productsGroup.POST("", controllers.CreateProductsHandler)
			productsGroup.GET("/:id", controllers.GetProductWithSizesHandler)
			productsGroup.PATCH("/:id", controllers.UpdateProductHandler)
			productsGroup.DELETE("/:id", controllers.DeleteProductHandler)
		}

		sizesGroup := v1.Group("/sizes")
		sizesGroup.Use(middleware.AuthMiddleware())
		{
			sizesGroup.GET("", controllers.GetAllSizesHandler)
			sizesGroup.POST("", controllers.CreateSizeHandler)
			sizesGroup.GET("/:id", controllers.GetSizeData)
			sizesGroup.PATCH("/:id", controllers.UpdateSizeHandler)
			sizesGroup.DELETE("/:id", controllers.DeleteSizeHandler)
		}

		cartsGroup := v1.Group("/carts")
		cartsGroup.Use(middleware.AuthMiddleware())
		{
			cartsGroup.GET("", controllers.GetAllCartsHandler)
			cartsGroup.POST("", controllers.AddToCartHandler)
			cartsGroup.GET("count", controllers.CountAllCartsHandler)
			cartsGroup.DELETE("/:id", controllers.DeleteCartsHandler)
		}

		ordersGroup := v1.Group("/orders")
		ordersGroup.Use(middleware.AuthMiddleware())
		{
			ordersGroup.POST("", controllers.AddOrderFromSelectedItemsHandler)
		}

		uploadsGroup := v1.Group("/uploads")
		{
			uploadsGroup.GET("/images/product/:filename", controllers.LoadImageProduct)
		}
	}
	return r
}
