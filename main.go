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
	// database.ChangeVarcharToTextAtProducts()
	// products := database.AddColumnToTable(db, "products", "status", "VARCHAR", 1, "after", "type")
	// if products != nil {
	// 	// Tangani kesalahan
	// 	fmt.Println("Failed to add column:", products)
	// }

	r := routes.SetupRoutes(db)
	r.MaxMultipartMemory = 8 << 20
	r.Run(":9090")
	// defer db.Close()
	fmt.Println("hello world")
}
