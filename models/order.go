package models

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
)

type OrderStatus string

const (
	Success OrderStatus = "Success"
	Failed  OrderStatus = "Failed"
	Expired OrderStatus = "Expired"
)

type Order struct {
	ID          uint        `json:"id" gorm:"primary_key"`
	UserID      uint        `json:"user_id"`
	User        User        `json:"user" gorm:"foreignKey:UserID"`
	OrderNumber string      `json:"order_number" gorm:"unique"`
	Items       []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
	Total       float64     `json:"total"`
	Status      OrderStatus `json:"status"`
	CreatedAt   time.Time   `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time   `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (o *Order) BeforeCreate(tx *gorm.DB) (err error) {
	o.OrderNumber = generateOrderNumber()
	return nil
}

func generateOrderNumber() string {
	return "Order-" + strconv.Itoa(rand.Intn(10000)) + strconv.Itoa(rand.Intn(10000))
}
