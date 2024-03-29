package models

import (
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

func SetupDB() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	USER := os.Getenv("DB_USERNAME")
	PASS := os.Getenv("DB_PASSWORD")
	HOST := os.Getenv("DB_HOST")
	PORT := os.Getenv("DB_PORT")
	DBNAME := os.Getenv("DB_DATABASE")
	DBCONNECTION := os.Getenv("DB_CONNECTION")

	URL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", USER, PASS, HOST, PORT, DBNAME)
	// URL := USER + ":" + PASS + "@tcp(" + HOST + ":" + PORT + ")/" + DBNAME
	// db, err := gorm.Open("mysql", "root:1234567890@tcp(localhost:33066)/i_komers_db?charset=utf8&parseTime=True&loc=Local&allowPublicKeyRetrieval=true")
	// db, err := gorm.Open("mysql", "root:1234567890@tcp(localhost:33066)/i_komers_db?charset=utf8&parseTime=True&loc=Local&allowPublicKeyRetrieval=true")

	db, err := gorm.Open(DBCONNECTION, URL)
	if err != nil {
		panic(err.Error())
	}

	return db
}
