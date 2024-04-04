package models

type OrderItem struct {
	ID        uint    `json:"id" gorm:"primary_key"`
	OrderID   uint    `json:"order_id"`
	ProductID uint    `json:"product_id"`
	Product   Product `json:"product" gorm:"foreignKey:ProductID"`
	SizeID    uint    `json:"size_id"`
	Size      Size    `json:"size" gorm:"foreignKey:SizeID"`
	Quantity  uint    `json:"quantity"`
	Price     float64 `json:"price"`
}
