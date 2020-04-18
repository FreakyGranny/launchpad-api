package db

import (
  "github.com/jinzhu/gorm"
)

// Product test model
type Product struct {
  gorm.Model
  Code string `json:"code"`
  Price uint  `json:"price"`
}
