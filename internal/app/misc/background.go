package misc

import (
	"sync"

	"github.com/FreakyGranny/launchpad-api/internal/app/models"
	"github.com/labstack/gommon/log"
)

// Background process
type Background struct {
	ProjectModel models.ProjectImpl
	UserModel    models.UserImpl
	RecalcChan   chan int
	UpdateChan   chan int
}

// NewBackground return new background instance
func NewBackground(mp models.ProjectImpl, mu models.UserImpl) *Background {
	return &Background{
		ProjectModel: mp,
		UserModel:    mu,
		RecalcChan:   make(chan int, 100),
		UpdateChan:   make(chan int, 100),
	}
}

// GetRecalcPipe returns recalc pipe
func (b *Background) GetRecalcPipe() chan int {
	return b.RecalcChan
}

// GetUpdatePipe returns update pipe
func (b *Background) GetUpdatePipe() chan int {
	return b.UpdateChan
}

// // GetHarvestPipe returns harvest pipe
// func GetHarvestPipe() chan uint {
// 	return harvestPipe
// }

// RecalcProject update total for project
func (b *Background) RecalcProject(wg *sync.WaitGroup) {
	defer wg.Done()
	var project *models.Project
	var ok bool

	for {
		projectID, open := <-b.RecalcChan
		if projectID == 0 && !open {
			log.Info("stop recalc")
			return
		}
		project, ok = b.ProjectModel.Get(projectID)
		if !ok {
			log.Errorf("project %d not found", projectID)
			continue
		}
		strategy, err := GetStrategy(&project.ProjectType, b.ProjectModel)
		if err != nil {
			log.Errorf("unable to get stategy for project %d", projectID)
			continue
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

// UpdateUser update user's rate
func (b *Background) UpdateUser(wg *sync.WaitGroup) {
	defer wg.Done()
	var user *models.User
	var ok bool
	var pGroups []models.ProjectGroup
	var err error

	for {
		userID, open := <-b.UpdateChan
		if userID == 0 && !open {
			log.Info("stop update users")
			return
		}
		user, ok = b.UserModel.Get(userID)
		if !ok {
			log.Errorf("user %d not found", userID)
			continue
		}
		pGroups, err = b.UserModel.GetProjectsForRate(userID)
		if err != nil {
			log.Error("error while fetching project groups")
		}
		user.ProjectCount, user.SuccessRate = getStats(pGroups)

		_, err = b.UserModel.Update(user)
		if err != nil {
			log.Error("Error while trying update user")
		}
	}
}

func getStats(groups []models.ProjectGroup) (int, float32) {
	var projectCount int
	var closedCount int
	var successCount int
	var successRate float32

	for _, group := range groups {
		projectCount += group.Cnt
		if group.Closed {
			closedCount += group.Cnt
		}
		if group.Closed && group.Locked {
			successCount += group.Cnt
		}
	}
	if closedCount > 0 {
		successRate = float32(successCount) / float32(closedCount)
	}

	return projectCount, successRate
}
