package misc

import (
	"github.com/FreakyGranny/launchpad-api/internal/app/models"
	"github.com/labstack/gommon/log"
)

// Background process
type Background struct {
	ProjectModel models.ProjectImpl
	RecalcChan chan int
}

// NewBackground return new background instance
func NewBackground(m models.ProjectImpl) *Background{
	return &Background{
		ProjectModel: m,
		RecalcChan: make(chan int, 10),
	}
}

// GetRecalcPipe returns recalc pipe
func (b *Background) GetRecalcPipe() chan int {
	return b.RecalcChan
}

// // GetUpdatePipe returns update pipe
// func GetUpdatePipe() chan uint {
// 	return updatePipe
// }

// // GetHarvestPipe returns harvest pipe
// func GetHarvestPipe() chan uint {
// 	return harvestPipe
// }

// RecalcProject update total for project
func (b *Background) RecalcProject() {
	defer close(b.RecalcChan)
	var project *models.Project
	var ok bool

	for {
		projectID := <-b.RecalcChan
		project, ok = b.ProjectModel.Get(projectID)
		if !ok {
			log.Errorf("project %d not found", projectID)
			continue
		}
		strategy, err := GetStrategy(&project.ProjectType, b.ProjectModel)
		if err != nil {
			log.Errorf("unable to get stategy for project %d", projectID)
		}
		strategy.Recalc(project)
		// strategy.CheckSearch(&project)
	}
}

// // HarvestCheck check all paid
// func HarvestCheck() {
// 	defer close(recalcPipe)
// 	dbClient := db.GetDbClient()
// 	var project db.Project

// 	for {
// 		projectID := <-harvestPipe
// 		if err := dbClient.Preload("ProjectType").First(&project, projectID).Error; gorm.IsRecordNotFoundError(err) {
// 			log.Error(err)
// 			return
// 		}

// 		strategy, err := GetStrategy(project.ProjectType)
// 		if err != nil {
// 			log.Error(err)
// 			return
// 		}
// 		strategy.CheckHarvest(&project)
// 	}
// }

// // UpdateUser update total for project
// func UpdateUser() {
// 	defer close(updatePipe)
// 	dbClient := db.GetDbClient()
// 	var user db.User

// 	for {
// 		userID := <-updatePipe
// 		if err := dbClient.First(&user, userID).Error; gorm.IsRecordNotFoundError(err) {
// 			log.Error(err)
// 			return
// 		}
// 		var projectCount uint
// 		var closedCount uint
// 		var successCount uint
// 		var successRate float32 = 0

// 		dbClient.Model(&db.Project{}).Where("published = ? AND owner_id = ?", true, userID).Count(&projectCount)
// 		dbClient.Model(&db.Project{}).Where("closed = ? AND owner_id = ?", true, userID).Count(&closedCount)
// 		if closedCount > 0 {
// 			dbClient.Model(&db.Project{}).Where("closed = ? AND locked = ? AND owner_id = ?", true, true, userID).Count(&successCount)
// 			successRate = float32(successCount) / float32(closedCount)
// 		}

// 		dbClient.Model(&user).Updates(db.User{ProjectCount: projectCount, SuccessRate: successRate})
// 	}
// }
