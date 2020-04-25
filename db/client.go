package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/FreakyGranny/launchpad-api/config"
)

var db *gorm.DB
var err error  


func sslMode(sslEnable bool) string {
	if sslEnable {
		return "enable"
	}

	return "disable"
}

// Init db instance
func Init(cfg config.PgConnection) {
	connectString := fmt.Sprintf(
		"host=%s port=%v user=%s dbname=%s password=%s sslmode=%s", 
		cfg.Host,
		cfg.Port,
		cfg.Username,
		cfg.DbName,
		cfg.Password,
		sslMode(cfg.SslEnable),
	)
	db, err = gorm.Open("postgres", connectString)
	if err != nil {
		fmt.Println(err)
		panic("DB Connection Error")
	}

	db.AutoMigrate(&Product{})
}

// GetDbClient returns db client object 
func GetDbClient() *gorm.DB {
	return db
}
