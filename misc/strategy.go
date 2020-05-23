package misc

import (
	"errors"

	"github.com/FreakyGranny/launchpad-api/db"
)

// Strategy project strategy depends on type
type Strategy interface{
	Percent(p *db.Project) uint
	Recalc(p *db.Project) uint
	CheckSearch(p *db.Project)
	CheckHarvest(p *db.Project)
}


// MoneyStrategy simple money type
type MoneyStrategy struct{}

// Percent returns percent of complition
func (s *MoneyStrategy) Percent(p *db.Project) uint {
	if p.GoalAmount == 0 {
		return 0
	}
	
	return uint(float64(p.Total) / float64(p.GoalAmount) * 100)
}

// Recalc recalculate project
func (s *MoneyStrategy) Recalc(p *db.Project) uint {
	dbClient := db.GetDbClient()
	var donations []db.Donation
	var total uint = 0

	dbClient.Where("project_id = ?", p.ID).Find(&donations)

	for _, d := range(donations) {
		total += d.Payment
	}

	return total
}

// CheckSearch check project for search stage ending
func (s *MoneyStrategy) CheckSearch(p *db.Project) {
	if s.Percent(p) >= 100 {
		p.Lock()
	}
}

// CheckHarvest check project for harvest stage ending
func (s *MoneyStrategy) CheckHarvest(p *db.Project) {
	dbClient := db.GetDbClient()
	var donations []db.Donation
	
	dbClient.Where("project_id = ?", p.ID).Find(&donations)

	for _, d := range(donations) {
		if !d.Paid {
			return
		}
	}
	p.Close()
}

// EventStrategy simple money type
type EventStrategy struct{}

// Percent returns percent of complition
func (s *EventStrategy) Percent(p *db.Project) uint {
	if p.GoalPeople == 0 {
		return 0
	}
	
	return uint(float64(p.Total) / float64(p.GoalPeople) * 100)
}

// Recalc recalculate project
func (s *EventStrategy) Recalc(p *db.Project) uint {
	dbClient := db.GetDbClient()
	var donationCount uint
	dbClient.Model(&db.Donation{}).Where("project_id = ?", p.ID).Count(&donationCount)

	return donationCount
}

// CheckSearch check project for search stage ending
func (s *EventStrategy) CheckSearch(p *db.Project) {
	if s.Percent(p) >= 100 {
		p.Lock()
		p.Close()
	}
}

// CheckHarvest check project for harvest stage ending
func (s *EventStrategy) CheckHarvest(p *db.Project) {}


// EventDateStrategy simple money type
type EventDateStrategy struct{
	baseStrategy EventStrategy
}

// Percent returns percent of complition
func (s *EventDateStrategy) Percent(p *db.Project) uint {
	return s.baseStrategy.Percent(p)
}

// Recalc recalculate project
func (s *EventDateStrategy) Recalc(p *db.Project) uint{
	return s.baseStrategy.Recalc(p)
}

// CheckSearch check project for search stage ending
func (s *EventDateStrategy) CheckSearch(p *db.Project) {
	// if day x
	s.baseStrategy.CheckSearch(p)
}

// CheckHarvest check project for harvest stage ending
func (s *EventDateStrategy) CheckHarvest(p *db.Project) {}


// MoneyEqualStrategy simple money type
type MoneyEqualStrategy struct{
	moneyStrategy MoneyStrategy
	eventStrategy EventStrategy
}

// Percent returns percent of complition
func (s *MoneyEqualStrategy) Percent(p *db.Project) uint {
	return s.moneyStrategy.Percent(p)
}

// Recalc recalculate project
func (s *MoneyEqualStrategy) Recalc(p *db.Project) uint{
	return s.eventStrategy.Recalc(p)
}

// CheckSearch check project for search stage ending
func (s *MoneyEqualStrategy) CheckSearch(p *db.Project) {
	// if day x
	s.eventStrategy.CheckSearch(p)
}

// CheckHarvest check project for harvest stage ending
func (s *MoneyEqualStrategy) CheckHarvest(p *db.Project) {
	s.moneyStrategy.CheckHarvest(p)
}


// GetStrategy returns project strategy based on project type
func GetStrategy(pt db.ProjectType) (Strategy, error){
    if pt.GoalByAmount && !pt.GoalByPeople {
        if pt.EndByGoalGain {
			return &MoneyStrategy{}, nil
		}            
	}
	if pt.GoalByPeople && !pt.GoalByAmount {
		if pt.EndByGoalGain {
			return &EventStrategy{}, nil
		}
		return &EventDateStrategy{baseStrategy: EventStrategy{}}, nil
	}
    if pt.GoalByPeople && pt.GoalByAmount {
        if !pt.EndByGoalGain {
			return &MoneyEqualStrategy{eventStrategy: EventStrategy{}, moneyStrategy: MoneyStrategy{}}, nil
		}            
	}

	return nil, errors.New("no matched strategy")
}
