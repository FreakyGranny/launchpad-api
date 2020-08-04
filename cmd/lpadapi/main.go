package main

import (
	"github.com/jonboulle/clockwork"
	"github.com/labstack/gommon/log"

	"github.com/FreakyGranny/launchpad-api/internal/app/auth"
	"github.com/FreakyGranny/launchpad-api/internal/app/config"
	"github.com/FreakyGranny/launchpad-api/internal/app/db"
	"github.com/FreakyGranny/launchpad-api/internal/app/handlers"
	"github.com/FreakyGranny/launchpad-api/internal/app/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := config.New()
	if cfg.DebugMode {
		log.SetLevel(log.DEBUG)
	}

	e := echo.New()

	d, err := db.Connect(cfg.Db)
	if err != nil {
		log.Fatal(err)
	}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{AllowOrigins: []string{"*"}}))

	ha := handlers.NewAuthHandler(
		cfg.JWTSecret,
		models.NewUserModel(d),
		auth.NewVk(cfg.Vk),
		clockwork.NewRealClock(),
	)

	e.POST("/login", ha.Login)
	e.OPTIONS("/login", ha.Login)

	hu := handlers.NewUserHandler(models.NewUserModel(d))

	u := e.Group("/user")
	// u.Use(middleware.JWT([]byte("secret")))
	// u.GET("", hu.GetUser)
	u.GET("/:id", hu.GetUser)

	// misc.VkInit(cfg.Vk)
	// db.Init(cfg.Db)
	// misc.BackgroundInit()

	// go misc.RecalcProject()
	// go misc.UpdateUser()
	// go misc.HarvestCheck()

	// client := db.GetDbClient()
	// defer client.Close()

	// e := route.Init()

	log.Fatal(e.Start(":1323"))
}
