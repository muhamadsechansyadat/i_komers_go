package models

import "time"

type User struct {
	ID           uint      `json:"id" gorm:"primary_key"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	Password     string    `json:"-"`
	PhotoProfile *string   `json:"photo_profile"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Expiration   string    `json:"expiration"`
}
