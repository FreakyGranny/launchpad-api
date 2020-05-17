package misc

import (
	"errors"
	"github.com/FreakyGranny/launchpad-api/db"
)

// Strategy project strategy depends on type
type Strategy interface{
	Percent(p *db.Project) uint
	Recalc(projectID uint) uint
}


// MoneyFastStrategy simple money type
type MoneyFastStrategy struct{}

// Percent returns percent of complition
func (s *MoneyFastStrategy) Percent(p *db.Project) uint {
	if p.GoalAmount == 0 {
		return 0
	}
	
	return uint(float64(p.Total) / float64(p.GoalAmount) * 100)
}

// Recalc recalculate project
func (s *MoneyFastStrategy) Recalc(projectID uint) uint {
	dbClient := db.GetDbClient()
	var donations []db.Donation
	var total uint = 0

	dbClient.Where("project_id = ?", projectID).Find(&donations)

	for _, d := range(donations) {
		total += d.Payment
	}

	return total
}


// MoneyEqualStrategy simple money type
type MoneyEqualStrategy struct{}

// Percent returns percent of complition
func (s *MoneyEqualStrategy) Percent(p *db.Project) uint {
	if p.GoalPeople == 0 {
		return 0
	}
	
	return uint(float64(p.Total) / float64(p.GoalPeople) * 100)
}

// Recalc recalculate project
func (s *MoneyEqualStrategy) Recalc(projectID uint) uint{
	dbClient := db.GetDbClient()
	var donationCount uint
	dbClient.Model(&db.Donation{}).Where("project_id = ?", projectID).Count(&donationCount)

	return donationCount
}


// EventDateStrategy simple money type
type EventDateStrategy struct{}

// Percent returns percent of complition
func (s *EventDateStrategy) Percent(p *db.Project) uint {
	if p.GoalPeople == 0 {
		return 0
	}
	
	return uint(float64(p.Total) / float64(p.GoalPeople) * 100)
}

// Recalc recalculate project
func (s *EventDateStrategy) Recalc(projectID uint) uint{
	dbClient := db.GetDbClient()
	var donationCount uint
	dbClient.Model(&db.Donation{}).Where("project_id = ?", projectID).Count(&donationCount)

	return donationCount
}


// EventFastStrategy simple money type
type EventFastStrategy struct{}

// Percent returns percent of complition
func (s *EventFastStrategy) Percent(p *db.Project) uint {
	if p.GoalPeople == 0 {
		return 0
	}
	
	return uint(float64(p.Total) / float64(p.GoalPeople) * 100)
}

// Recalc recalculate project
func (s *EventFastStrategy) Recalc(projectID uint) uint {
	dbClient := db.GetDbClient()
	var donationCount uint
	dbClient.Model(&db.Donation{}).Where("project_id = ?", projectID).Count(&donationCount)

	return donationCount
}


// GetStrategy returns project strategy based on project type
func GetStrategy(pt db.ProjectType) (Strategy, error){
    if pt.GoalByAmount && !pt.GoalByPeople {
        if pt.EndByGoalGain {
			return &MoneyFastStrategy{}, nil
		}            
	}
	if pt.GoalByPeople && !pt.GoalByAmount {
		if pt.EndByGoalGain {
			return &EventFastStrategy{}, nil
		}
		return &EventDateStrategy{}, nil
	}
    if pt.GoalByPeople && pt.GoalByAmount {
        if !pt.EndByGoalGain {
			return &MoneyEqualStrategy{}, nil
		}            
	}

	return nil, errors.New("no matched strategy")
}
