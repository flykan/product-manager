package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string  `json:"name" gorm:"size:100;not null"`
	Description string  `json:"description" gorm:"size:500"`
	Price       float64 `json:"price" gorm:"type:decimal(10,2);not null"`
	Stock       int     `json:"stock" gorm:"not null;default:0"`
	Category    string  `json:"category" gorm:"size:50"`
}

// TableName 指定表名
func (Product) TableName() string {
	return "products"
}
