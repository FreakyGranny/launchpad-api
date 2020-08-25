package misc

import (
	"errors"
	"time"

	"github.com/FreakyGranny/launchpad-api/internal/app/models"
)

// Strategy project strategy depends on type
type Strategy interface {
	Percent(p *models.Project) int
	Recalc(p *models.Project) error
	CheckSearch(p *models.Project) (bool, error)
	CheckHarvest(p *models.Project) (bool, error)
	CloseOutdated(p *models.Project) error
}

// MoneyStrategy simple money type
type MoneyStrategy struct {
	projectModel models.ProjectImpl
}

// NewMoneyStrategy Creates new money strategy
func NewMoneyStrategy(m models.ProjectImpl) *MoneyStrategy {
	return &MoneyStrategy{projectModel: m}
}

// Percent returns percent of completion
func (s *MoneyStrategy) Percent(p *models.Project) int {
	if p.GoalAmount == 0 {
		return 0
	}

	return int(float64(p.Total) / float64(p.GoalAmount) * 100)
}

// Recalc recalculate project
func (s *MoneyStrategy) Recalc(p *models.Project) error {
	return s.projectModel.UpdateTotalByPayment(p)
}

// CheckSearch check project for search stage ending
func (s *MoneyStrategy) CheckSearch(p *models.Project) (bool, error) {
	if s.Percent(p) >= 100 {
		return true, s.projectModel.Lock(p)
	}

	return false, nil
}

// CheckHarvest check project for harvest stage ending
func (s *MoneyStrategy) CheckHarvest(p *models.Project) (bool, error) {
	paid, err := s.projectModel.CheckForPaid(p.ID)
	if err != nil {
		return false, err
	}
	if !paid {
		return false, nil
	}

	return true, s.projectModel.Close(p)
}

// CloseOutdated check project is outdated
func (s *MoneyStrategy) CloseOutdated(p *models.Project) error {
	n := time.Now()
	d := p.ReleaseDate
	if n.Year() == d.Year() && n.Month() == d.Month() && n.Day() > d.Day() {
		return s.projectModel.Close(p)
	}

	return nil
}

// EventStrategy simple event type
type EventStrategy struct {
	projectModel models.ProjectImpl
}

// NewEventStrategy Creates new event strategy
func NewEventStrategy(m models.ProjectImpl) *EventStrategy {
	return &EventStrategy{projectModel: m}
}

// Percent returns percent of completion
func (s *EventStrategy) Percent(p *models.Project) int {
	if p.GoalPeople == 0 {
		return 0
	}

	return int(float64(p.Total) / float64(p.GoalPeople) * 100)
}

// Recalc recalculate project
func (s *EventStrategy) Recalc(p *models.Project) error {
	return s.projectModel.UpdateTotalByCount(p)
}

// CheckSearch check project for search stage ending
func (s *EventStrategy) CheckSearch(p *models.Project) (bool, error) {
	if s.Percent(p) >= 100 {
		return true, s.projectModel.Lock(p)
	}

	return false, nil
}

// CheckHarvest check project for harvest stage ending
func (s *EventStrategy) CheckHarvest(p *models.Project) (bool, error) {
	return true, s.projectModel.Close(p)
}

// CloseOutdated check project is outdated
func (s *EventStrategy) CloseOutdated(p *models.Project) error {
	n := time.Now()
	d := p.ReleaseDate
	if n.Year() == d.Year() && n.Month() == d.Month() && n.Day() > d.Day() {
		return s.projectModel.Close(p)
	}

	return nil
}

// EventDateStrategy event type with date
type EventDateStrategy struct {
	baseStrategy *EventStrategy
}

// NewEventDateStrategy ...
func NewEventDateStrategy(m models.ProjectImpl) *EventDateStrategy {
	return &EventDateStrategy{
		baseStrategy: NewEventStrategy(m),
	}
}

// Percent returns percent of completion
func (s *EventDateStrategy) Percent(p *models.Project) int {
	return s.baseStrategy.Percent(p)
}

// Recalc recalculate project
func (s *EventDateStrategy) Recalc(p *models.Project) error {
	return s.baseStrategy.Recalc(p)
}

// CheckSearch check project for search stage ending
func (s *EventDateStrategy) CheckSearch(p *models.Project) (bool, error) {
	n := time.Now()
	d := p.ReleaseDate
	if n.Year() == d.Year() && n.Month() == d.Month() && n.Day() == d.Day() {
		return s.baseStrategy.CheckSearch(p)
	}

	return false, nil
}

// CheckHarvest check project for harvest stage ending
func (s *EventDateStrategy) CheckHarvest(p *models.Project) (bool, error) {
	return s.baseStrategy.CheckHarvest(p)
}

// CloseOutdated check project is outdated
func (s *EventDateStrategy) CloseOutdated(p *models.Project) error {
	return s.baseStrategy.CloseOutdated(p)
}

// MoneyEqualStrategy money type with equal part splitting
type MoneyEqualStrategy struct {
	moneyStrategy *MoneyStrategy
	eventStrategy *EventDateStrategy
}

// NewMoneyEqualStrategy ...
func NewMoneyEqualStrategy(m models.ProjectImpl) *MoneyEqualStrategy {
	return &MoneyEqualStrategy{
		eventStrategy: NewEventDateStrategy(m),
		moneyStrategy: NewMoneyStrategy(m),
	}
}

// Percent returns percent of completion
func (s *MoneyEqualStrategy) Percent(p *models.Project) int {
	return s.eventStrategy.Percent(p)
}

// Recalc recalculate project
func (s *MoneyEqualStrategy) Recalc(p *models.Project) error {
	return s.eventStrategy.Recalc(p)
}

// CheckSearch check project for search stage ending
func (s *MoneyEqualStrategy) CheckSearch(p *models.Project) (bool, error) {
	evolved, err := s.eventStrategy.CheckSearch(p)
	if err != nil {
		return false, err
	}
	if !evolved {
		return evolved, err
	}
	err = s.moneyStrategy.projectModel.SetEqualDonation(p)

	return evolved, err
}

// CheckHarvest check project for harvest stage ending
func (s *MoneyEqualStrategy) CheckHarvest(p *models.Project) (bool, error) {
	return s.moneyStrategy.CheckHarvest(p)
}

// CloseOutdated check project is outdated
func (s *MoneyEqualStrategy) CloseOutdated(p *models.Project) error {
	return s.moneyStrategy.CloseOutdated(p)
}

// GetStrategy returns project strategy based on project type
func GetStrategy(pt *models.ProjectType, r models.ProjectImpl) (Strategy, error) {
	if pt.GoalByAmount && !pt.GoalByPeople {
		if pt.EndByGoalGain {
			return NewMoneyStrategy(r), nil
		}
	}
	if pt.GoalByPeople && !pt.GoalByAmount {
		if pt.EndByGoalGain {
			return NewEventStrategy(r), nil
		}
		return NewEventDateStrategy(r), nil
	}
	if pt.GoalByPeople && pt.GoalByAmount {
		if !pt.EndByGoalGain {
			return NewMoneyEqualStrategy(r), nil
		}
	}

	return nil, errors.New("no matched strategy")
}
