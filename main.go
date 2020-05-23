package main

import (
	"github.com/labstack/gommon/log"

	"github.com/FreakyGranny/launchpad-api/config"
	"github.com/FreakyGranny/launchpad-api/db"
	"github.com/FreakyGranny/launchpad-api/misc"
	"github.com/FreakyGranny/launchpad-api/route"
)
  

func main() {
	cfg := config.New()
	if cfg.DebugMode {
		log.SetLevel(log.DEBUG)
	}

	misc.VkInit(cfg.Vk)
	db.Init(cfg.Db)
	misc.BackgroundInit()

	go misc.RecalcProject()
	go misc.UpdateUser()
	go misc.HarvestCheck()

	client := db.GetDbClient()
	defer client.Close()
  	
	e := route.Init()

	log.Fatal(e.Start(":1323"))
}
