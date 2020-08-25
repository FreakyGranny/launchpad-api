package misc

import (
	"sync"
	"time"

	"github.com/FreakyGranny/launchpad-api/internal/app/models"
	"github.com/labstack/gommon/log"
)

// Background process
type Background struct {
	SystemModel  models.SystemImpl
	ProjectModel models.ProjectImpl
	UserModel    models.UserImpl
	DoneChan     chan struct{}
	RecalcChan   chan int
	UpdateChan   chan int
	SearchChan   chan *models.Project
	HarverstChan chan *models.Project
}

// NewBackground return new background instance
func NewBackground(ms models.SystemImpl, mp models.ProjectImpl, mu models.UserImpl) *Background {
	return &Background{
		SystemModel:  ms,
		ProjectModel: mp,
		UserModel:    mu,
		DoneChan:     make(chan struct{}, 1),
		RecalcChan:   make(chan int, 100),
		UpdateChan:   make(chan int, 100),
		SearchChan:   make(chan *models.Project, 10),
		HarverstChan: make(chan *models.Project, 10),
	}
}

// Terminate stops all channels
func (b *Background) Terminate() {
	close(b.DoneChan)
}

// GetRecalcPipe returns recalc pipe
func (b *Background) GetRecalcPipe() chan int {
	return b.RecalcChan
}

// PeriodicCheck returns recalc pipe
func (b *Background) PeriodicCheck(wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(b.RecalcChan)
	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()
	for {
		select {
		case t := <-ticker.C:
			system, err := b.SystemModel.Get()
			if err != nil {
				log.Error("unable to get system settings")
			}
			system.LastCheck = system.LastCheck.Add(24 * time.Hour)
			if t.Before(system.LastCheck) {
				continue
			}
			log.Info("checking active projects")
			err = b.SystemModel.Update(system)
			if err != nil {
				log.Error(err)
				continue
			}
			projects, err := b.ProjectModel.GetActiveProjects()
			if err != nil {
				log.Error(err)
				continue
			}
			for _, project := range *projects {
				b.RecalcChan <- project.ID
			}
		case <-b.DoneChan:
			log.Info("stop periodic check")
			return
		}
	}
}

// RecalcProject update total for project
func (b *Background) RecalcProject(wg *sync.WaitGroup) {
	defer wg.Done()
	var project *models.Project
	var ok bool

	for {
		projectID, open := <-b.RecalcChan
		if projectID == 0 && !open {
			log.Info("stop recalc")
			close(b.SearchChan)
			return
		}
		project, ok = b.ProjectModel.Get(projectID)
		if !ok {
			log.Errorf("project %d not found", projectID)
			continue
		}
		if project.Locked {
			b.SearchChan <- project
			continue
		}
		strategy, err := GetStrategy(&project.ProjectType, b.ProjectModel)
		if err != nil {
			log.Errorf("unable to get stategy for project %d", projectID)
			continue
		}
		err = strategy.Recalc(project)
		if err != nil {
			log.Errorf("unable to recalc project %d", projectID)
			continue
		}
		b.SearchChan <- project
	}
}

// CheckSearch check project for search stage
func (b *Background) CheckSearch(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		project, open := <-b.SearchChan
		if project == nil && !open {
			log.Info("stop checking search")
			close(b.HarverstChan)
			return
		}
		if project.Locked {
			b.HarverstChan <- project
			continue
		}
		strategy, err := GetStrategy(&project.ProjectType, b.ProjectModel)
		if err != nil {
			log.Errorf("unable to get stategy for project %d", project.ID)
			continue
		}
		evolved, err := strategy.CheckSearch(project)
		if err != nil {
			log.Errorf("unable to check search for project %d", project.ID)
			continue
		}
		if evolved {
			b.HarverstChan <- project
		} else {
			err = strategy.CloseOutdated(project)
			if err != nil {
				log.Errorf("unable to check outdate for project %d", project.ID)
			}
		}
	}
}

// HarvestCheck check project for harvest stage
func (b *Background) HarvestCheck(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		project, open := <-b.HarverstChan
		if project == nil && !open {
			log.Info("stop checking harvest")
			close(b.UpdateChan)
			return
		}
		strategy, err := GetStrategy(&project.ProjectType, b.ProjectModel)
		if err != nil {
			log.Errorf("unable to get stategy for project %d", project.ID)
			continue
		}
		evolved, err := strategy.CheckHarvest(project)
		if err != nil {
			log.Errorf("unable to check search for project %d", project.ID)
			continue
		}
		if evolved {
			b.UpdateChan <- project.OwnerID
		}
	}
}

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
