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

	JWTmiddleware := middleware.JWT([]byte(cfg.JWTSecret))

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
	u.Use(JWTmiddleware)
	u.GET("", hu.GetCurrentUser)
	u.GET("/:id", hu.GetUser)

	hc := handlers.NewCategoryHandler(models.NewCategoryModel(d))
	c := e.Group("/category")
	c.Use(JWTmiddleware)
	c.GET("", hc.GetCategories)

	// misc.BackgroundInit()

	// go misc.RecalcProject()
	// go misc.UpdateUser()
	// go misc.HarvestCheck()

	log.Fatal(e.Start(":1323"))
}
