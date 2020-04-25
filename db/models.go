package db

import (
  "time"
)

// Product test model
type Product struct {
	ID        uint `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
  DeletedAt *time.Time `sql:"index" json:"-"`
  
  Code string `json:"code"`
  Price uint  `json:"price"`
}
