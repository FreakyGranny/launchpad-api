package main

import (
	"github.com/FreakyGranny/launchpad-api/db"
	"github.com/FreakyGranny/launchpad-api/config"
	"github.com/FreakyGranny/launchpad-api/route"
)
  

func main() {
	cfg := config.New()

	db.Init(cfg.Db)
	client := db.GetDbClient()

	defer client.Close()
  
	// Создание
	client.Create(&db.Product{Code: "L1212", Price: 1000})
    
	// Правка - обновление цены на 2000
	// client.Model(&product).Update("Price", 2000)
  
	e := route.Init()
	e.Logger.Fatal(e.Start(":1323"))

	// Удаление - удалить продукт
	// client.Delete(&product)
  }
