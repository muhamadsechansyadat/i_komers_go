package helpers

import (
	"i_komers_go/models"
	"strconv"

	"github.com/gosimple/slug"
	"github.com/jinzhu/gorm"
)

func GenerateUniqueSlug(db *gorm.DB, tableName, name string) string {
	baseSlug := slug.Make(name)
	proposedSlug := baseSlug
	counter := 1

	for {
		// Check if the proposedSlug already exists in the database
		var existingProduct models.Product
		if err := db.Table(tableName).Where("slug = ?", proposedSlug).First(&existingProduct).Error; err != nil {
			// Slug doesn't exist, use it
			return proposedSlug
		}

		// If the slug exists, append a counter and try again
		proposedSlug = baseSlug + "-" + slug.Make(strconv.Itoa(counter))
		counter++
	}
}
