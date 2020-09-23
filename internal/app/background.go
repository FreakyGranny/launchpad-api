package app

import (
	"context"
	"math"
	"sync"
	"time"

	"github.com/FreakyGranny/launchpad-api/internal/models"
	"github.com/labstack/gommon/log"
)

// Background process
type Background struct {
	systemModel  models.SystemImpl
	projectModel models.ProjectImpl
	userModel    models.UserImpl
	recalcChan   chan int
	updateChan   chan int
	searchChan   chan *models.Project
	harverstChan chan *models.Project
	wg           *sync.WaitGroup
}

// NewBackground return new background instance
func NewBackground(ms models.SystemImpl, mp models.ProjectImpl, mu models.UserImpl) *Background {
	return &Background{
		systemModel:  ms,
		projectModel: mp,
		userModel:    mu,
		recalcChan:   make(chan int, 100),
		updateChan:   make(chan int, 100),
		searchChan:   make(chan *models.Project, 10),
		harverstChan: make(chan *models.Project, 10),
		wg:           &sync.WaitGroup{},
	}
}

// GetRecalcPipe returns recalc pipe
func (b *Background) GetRecalcPipe() chan int {
	return b.recalcChan
}

// Start starts background pipeline.
func (b *Background) Start(ctx context.Context) {
	go b.PeriodicCheck(ctx, b.wg)
	go b.RecalcProject(b.wg)
	go b.CheckSearch(b.wg)
	go b.HarvestCheck(b.wg)
	go b.UpdateUser(b.wg)
	b.wg.Add(5)

}

// Wait waits background pipeline.
func (b *Background) Wait() {
	b.wg.Wait()
}

// PeriodicCheck returns recalc pipe
func (b *Background) PeriodicCheck(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(b.recalcChan)
	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()
	for {
		select {
		case t := <-ticker.C:
			system, err := b.systemModel.Get()
			if err != nil {
				log.Error("unable to get system settings")
			}
			system.LastCheck = system.LastCheck.Add(24 * time.Hour)
			if t.Before(system.LastCheck) {
				continue
			}
			log.Info("checking active projects")
			err = b.systemModel.Update(system)
			if err != nil {
				log.Error(err)
				continue
			}
			projects, err := b.projectModel.GetActiveProjects()
			if err != nil {
				log.Error(err)
				continue
			}
			for _, project := range *projects {
				b.recalcChan <- project.ID
			}
		case <-ctx.Done():
			log.Info("stop periodic check")
			return
		}
	}
}

// RecalcProject update total for project
func (b *Background) RecalcProject(wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(b.searchChan)
	var project *models.Project
	var ok bool

	for {
		projectID, open := <-b.recalcChan
		if projectID == 0 && !open {
			log.Info("stop recalc")
			return
		}
		project, ok = b.projectModel.Get(projectID)
		if !ok {
			log.Errorf("project %d not found", projectID)
			continue
		}
		if project.Locked {
			b.searchChan <- project
			continue
		}
		strategy, err := GetStrategy(&project.ProjectType, b.projectModel)
		if err != nil {
			log.Errorf("unable to get stategy for project %d", projectID)
			continue
		}
		err = strategy.Recalc(project)
		if err != nil {
			log.Error(err)
			log.Errorf("unable to recalc project %d", projectID)
			continue
		}
		b.searchChan <- project
	}
}

// CheckSearch check project for search stage
func (b *Background) CheckSearch(wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(b.harverstChan)
	for {
		project, open := <-b.searchChan
		if project == nil && !open {
			log.Info("stop checking search")
			return
		}
		if project.Locked {
			b.harverstChan <- project
			continue
		}
		strategy, err := GetStrategy(&project.ProjectType, b.projectModel)
		if err != nil {
			log.Errorf("unable to get stategy for project %d", project.ID)
			continue
		}
		_, err = strategy.CheckSearch(project)
		if err != nil {
			log.Errorf("unable to check search for project %d", project.ID)
			continue
		}
		b.harverstChan <- project
	}
}

// HarvestCheck check project for harvest stage
func (b *Background) HarvestCheck(wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(b.updateChan)
	for {
		project, open := <-b.harverstChan
		if project == nil && !open {
			log.Info("stop checking harvest")
			return
		}
		strategy, err := GetStrategy(&project.ProjectType, b.projectModel)
		if err != nil {
			log.Errorf("unable to get stategy for project %d", project.ID)
			continue
		}
		if !project.Locked {
			closed, err := strategy.CloseOutdated(project)
			if err != nil {
				log.Errorf("unable to check outdate for project %d", project.ID)
			}
			if closed {
				b.updateChan <- project.OwnerID
			}

			continue
		}
		evolved, err := strategy.CheckHarvest(project)
		if err != nil {
			log.Errorf("unable to check harvest for project %d", project.ID)
			continue
		}
		if evolved {
			b.updateChan <- project.OwnerID
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
		userID, open := <-b.updateChan
		if userID == 0 && !open {
			log.Info("stop update users")
			return
		}
		user, ok = b.userModel.Get(userID)
		if !ok {
			log.Errorf("user %d not found", userID)
			continue
		}
		pGroups, err = b.userModel.GetProjectsForRate(userID)
		if err != nil {
			log.Error("error while fetching project groups")
		}
		user.ProjectCount, user.SuccessRate = getStats(pGroups)

		_, err = b.userModel.Update(user)
		if err != nil {
			log.Error("Error while trying update user")
		}
	}
}

func getStats(groups []models.ProjectGroup) (int, float64) {
	var projectCount int
	var closedCount int
	var successCount int
	var successRate float64

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
		successRate = math.Round(float64(successCount)*100/float64(closedCount)) / 100
	}

	return projectCount, successRate
}
