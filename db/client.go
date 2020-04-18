package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB
var err error  


// Init db instance
func Init() {
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
	  panic("failed to connect database")
	}
	// Миграция схем
	db.AutoMigrate(&Product{})
	
}

// GetDbClient returns db client object 
func GetDbClient() *gorm.DB {
	return db
}
