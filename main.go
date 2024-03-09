package main

import (
	"fmt"
	"i_komers_go/models"
	"i_komers_go/routes"
)

func main() {
	db := models.SetupDB()
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Product{})

	r := routes.SetupRoutes(db)
	r.Run(":9090")
	// defer db.Close()
	fmt.Println("hello world")
}
