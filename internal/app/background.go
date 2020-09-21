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
	SystemModel  models.SystemImpl
	ProjectModel models.ProjectImpl
	UserModel    models.UserImpl
	RecalcChan   chan int
	UpdateChan   chan int
	SearchChan   chan *models.Project
	HarverstChan chan *models.Project
	wg           *sync.WaitGroup
}

// NewBackground return new background instance
func NewBackground(ms models.SystemImpl, mp models.ProjectImpl, mu models.UserImpl) *Background {
	return &Background{
		SystemModel:  ms,
		ProjectModel: mp,
		UserModel:    mu,
		RecalcChan:   make(chan int, 100),
		UpdateChan:   make(chan int, 100),
		SearchChan:   make(chan *models.Project, 10),
		HarverstChan: make(chan *models.Project, 10),
		wg:           &sync.WaitGroup{},
	}
}

// GetRecalcPipe returns recalc pipe
func (b *Background) GetRecalcPipe() chan int {
	return b.RecalcChan
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
		case <-ctx.Done():
			log.Info("stop periodic check")
			return
		}
	}
}

// RecalcProject update total for project
func (b *Background) RecalcProject(wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(b.SearchChan)
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
			log.Error(err)
			log.Errorf("unable to recalc project %d", projectID)
			continue
		}
		b.SearchChan <- project
	}
}

// CheckSearch check project for search stage
func (b *Background) CheckSearch(wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(b.HarverstChan)
	for {
		project, open := <-b.SearchChan
		if project == nil && !open {
			log.Info("stop checking search")
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
		_, err = strategy.CheckSearch(project)
		if err != nil {
			log.Errorf("unable to check search for project %d", project.ID)
			continue
		}
		b.HarverstChan <- project
	}
}

// HarvestCheck check project for harvest stage
func (b *Background) HarvestCheck(wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(b.UpdateChan)
	for {
		project, open := <-b.HarverstChan
		if project == nil && !open {
			log.Info("stop checking harvest")
			return
		}
		strategy, err := GetStrategy(&project.ProjectType, b.ProjectModel)
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
				b.UpdateChan <- project.OwnerID
			}

			continue
		}
		evolved, err := strategy.CheckHarvest(project)
		if err != nil {
			log.Errorf("unable to check harvest for project %d", project.ID)
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
