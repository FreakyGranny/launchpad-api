package misc

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/gommon/log"
	"github.com/FreakyGranny/launchpad-api/db"
)

var recalcPipe chan uint

// BackgroundInit background channels
func BackgroundInit() {
	recalcPipe = make(chan uint)
}

// GetRecalcPipe returns recalc pipe
func GetRecalcPipe() chan uint {
	return recalcPipe
}

// RecalcProject update total for project
func RecalcProject(ch chan uint) {
	defer close(ch)
	dbClient := db.GetDbClient()
	var project db.Project

	for {
		projectID := <-ch	
		if err := dbClient.Preload("ProjectType").First(&project, projectID).Error; gorm.IsRecordNotFoundError(err) {
			log.Error(err)
			return
		}
	
		strategy, err := GetStrategy(project.ProjectType)
		if err != nil {
			log.Error(err)
			return
		}
		total := strategy.Recalc(project.ID)
		dbClient.Model(&project).Update("total", total)
	}
}
