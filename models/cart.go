package models

import (
	"time"
)

type Cart struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	UserID    uint      `json:"user_id"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	ProductID uint      `json:"product_id"`
	Product   Product   `json:"product" gorm:"foreignKey:ProductID"`
	SizeID    uint      `json:"size_id"`
	Size      Size      `json:"size" gorm:"foreignKey:SizeID"`
	Quantity  uint      `json:"quantity"`
	Status    Status    `json:"status"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
