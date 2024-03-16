package models

import (
	"time"
)

type ProductType string

const (
	New  ProductType = "new"
	Sale ProductType = "sale"
)

type Product struct {
	ID           uint        `json:"id" gorm:"primary_key"`
	Name         string      `json:"name" binding:"required"`
	Slug         string      `json:"slug"`
	Price        float64     `json:"price" gorm:"type:decimal(10,2);not null"`
	PhotoProduct string      `json:"photo_product"`
	Description  string      `json:"description"`
	Type         ProductType `json:"type"`
	CreatedAt    time.Time   `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time   `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
