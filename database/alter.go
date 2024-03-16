package database

import (
	"fmt"
	"i_komers_go/models"

	"github.com/jinzhu/gorm"
)

// ExecuteMigration adalah fungsi yang digunakan untuk mengeksekusi pernyataan SQL migrasi.
func ChangeVarcharToTextAtProducts() {
	db := models.SetupDB()
	db.Exec("ALTER TABLE products MODIFY photo_product TEXT;")
	db.Exec("ALTER TABLE products MODIFY description TEXT;")
}

func AddColumnToTable(db *gorm.DB, tableName string, columnName string, dataType string, length int, position string, targetColumn string) error {
	// Cek apakah posisi yang diminta valid (before atau after)
	if position != "before" && position != "after" {
		return fmt.Errorf("Invalid position. Must be 'before' or 'after'")
	}

	// Buat string untuk tipe data dengan panjangnya
	dataTypeWithLength := dataType
	if length > 0 {
		dataTypeWithLength = fmt.Sprintf("%s(%d)", dataType, length)
	}

	// Query untuk menambahkan kolom
	var query string
	if position == "before" {
		query = fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s BEFORE %s", tableName, columnName, dataTypeWithLength, targetColumn)
	} else {
		query = fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s AFTER %s", tableName, columnName, dataTypeWithLength, targetColumn)
	}

	// Eksekusi query
	if err := db.Exec(query).Error; err != nil {
		return err
	}

	return nil
}
