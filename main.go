package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/FreakyGranny/launchpad-api/db"
)
  

func main() {
	db.Init()
	client := db.GetDbClient()

	defer client.Close()
  
	fmt.Println(client)
  
	// Создание
	client.Create(&db.Product{Code: "L1212", Price: 1000})
    
	// Правка - обновление цены на 2000
	// client.Model(&product).Update("Price", 2000)
  
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		// Чтение
		var product db.Product
		client.First(&product, 1) // find product with id 1
		// client.First(&product, "code = ?", "L1212") // find product with code l1212

		return c.JSON(http.StatusOK, product)
	})
	e.Logger.Fatal(e.Start(":1323"))

	// Удаление - удалить продукт
	// client.Delete(&product)
  }
