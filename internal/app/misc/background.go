package misc

import (
	"github.com/FreakyGranny/launchpad-api/internal/app/models"
	"github.com/labstack/gommon/log"
)

// Background process
type Background struct {
	ProjectModel models.ProjectImpl
	RecalcChan   chan int
}

// UserUpdater process
type UserUpdater struct {
	UserModel  models.UserImpl
	UpdateChan chan int
}

// NewBackground return new background instance
func NewBackground(m models.ProjectImpl) *Background {
	return &Background{
		ProjectModel: m,
		RecalcChan:   make(chan int, 100),
	}
}

// NewUserUpdater returns new userUpdater instance
func NewUserUpdater(m models.UserImpl) *UserUpdater {
	return &UserUpdater{
		UserModel:  m,
		UpdateChan: make(chan int, 100),
	}
}

// GetRecalcPipe returns recalc pipe
func (b *Background) GetRecalcPipe() chan int {
	return b.RecalcChan
}

// GetUpdatePipe returns update pipe
func (uu *UserUpdater) GetUpdatePipe() chan int {
	return uu.UpdateChan
}

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
func (uu *UserUpdater) UpdateUser() {
	defer close(uu.UpdateChan)
	var user *models.User
	var ok bool
	var pGroups []models.ProjectGroup
	var err error
	var projectCount int
	var closedCount int
	var successCount int
	var successRate float32

	for {
		userID := <-uu.UpdateChan
		projectCount = 0
		closedCount = 0
		successCount = 0
		successRate = 0
	
		user, ok = uu.UserModel.Get(userID)
		if !ok {
			log.Errorf("user %d not found", userID)
			continue
		}
		pGroups, err = uu.UserModel.GetProjectsForRate(userID)
		if err != nil {
			log.Error("error while fetching project groups")
		}
		if len(pGroups) == 0 {
			continue
		}
		for _, group := range pGroups {
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
		user.ProjectCount = projectCount
		user.SuccessRate = successRate

		_, err = uu.UserModel.Update(user)
		if err != nil {
			log.Error("Error while trying update user")
		}
	}
}
