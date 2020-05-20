package misc

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/gommon/log"
	"github.com/FreakyGranny/launchpad-api/db"
)

var recalcPipe chan uint
var updatePipe chan uint

// BackgroundInit background channels
func BackgroundInit() {
	recalcPipe = make(chan uint)
	updatePipe = make(chan uint)
}

// GetRecalcPipe returns recalc pipe
func GetRecalcPipe() chan uint {
	return recalcPipe
}

// GetUpdatePipe returns update pipe
func GetUpdatePipe() chan uint {
	return updatePipe
}

// RecalcProject update total for project
func RecalcProject() {
	defer close(recalcPipe)
	dbClient := db.GetDbClient()
	var project db.Project

	for {
		projectID := <-recalcPipe	
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

// UpdateUser update total for project
func UpdateUser() {
	defer close(updatePipe)
	dbClient := db.GetDbClient()
	var user db.User

	for {
		userID := <-updatePipe
		if err := dbClient.First(&user, userID).Error; gorm.IsRecordNotFoundError(err) {
			log.Error(err)
			return
		}
		var projectCount uint
		var closedCount uint
		var successCount uint
		var successRate float32 = 0
		
		dbClient.Model(&db.Project{}).Where("published = ? AND owner_id = ?", true, userID).Count(&projectCount)
		dbClient.Model(&db.Project{}).Where("closed = ? AND owner_id = ?", true, userID).Count(&closedCount)
		if closedCount > 0 {
			dbClient.Model(&db.Project{}).Where("closed = ? AND locked = ? AND owner_id = ?", true, true, userID).Count(&successCount)
			successRate = float32(successCount) / float32(closedCount)
		}

		dbClient.Model(&user).Updates(db.User{ProjectCount: projectCount, SuccessRate: successRate})
	}
}
