package models

import (
	"time"
)

type Status string

const (
	Yes Status = "Y"
	No  Status = "N"
)

type Size struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	ProductID uint      `json:"product_id"`
	Name      string    `json:"name"`
	Quantity  uint      `json:"quantity"`
	Status    Status    `json:"status"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
