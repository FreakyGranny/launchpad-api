package misc

import (
	"testing"

	"github.com/FreakyGranny/launchpad-api/internal/app/mocks"
	"github.com/FreakyGranny/launchpad-api/internal/app/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

func isMoneyStrategy(t interface{}) bool {
	_, ok := t.(*MoneyStrategy)

	return ok
}

func isEventStrategy(t interface{}) bool {
	_, ok := t.(*EventStrategy)

	return ok
}

func isEventDateStrategy(t interface{}) bool {
	_, ok := t.(*EventDateStrategy)

	return ok
}

func isMoneyEqualStrategy(t interface{}) bool {
	_, ok := t.(*MoneyEqualStrategy)

	return ok
}

type StrategySuite struct {
	suite.Suite
	mockProjectCtl *gomock.Controller
	mockProject    *mocks.MockProjectImpl
}

func (s *StrategySuite) SetupTest() {
	s.mockProjectCtl = gomock.NewController(s.T())
	s.mockProject = mocks.NewMockProjectImpl(s.mockProjectCtl)
}

func (s *StrategySuite) TearDownTest() {
	s.mockProjectCtl.Finish()
}

func (s *StrategySuite) TestGetStrategyMoney() {
	pt := &models.ProjectType{
		GoalByAmount: true,
		EndByGoalGain: true,
	}
	st, err := GetStrategy(pt, s.mockProject)
	s.Require().NoError(err)
	s.Require().True(isMoneyStrategy(st))
}

func (s *StrategySuite) TestGetStrategyEvent() {
	pt := &models.ProjectType{
		GoalByPeople: true,
		EndByGoalGain: true,
	}
	st, err := GetStrategy(pt, s.mockProject)
	s.Require().NoError(err)
	s.Require().True(isEventStrategy(st))
}

func (s *StrategySuite) TestGetStrategyEventDate() {
	pt := &models.ProjectType{
		GoalByPeople: true,
		EndByGoalGain: false,
	}
	st, err := GetStrategy(pt, s.mockProject)
	s.Require().NoError(err)
	s.Require().True(isEventDateStrategy(st))
}

func (s *StrategySuite) TestGetStrategyMoneyEqual() {
	pt := &models.ProjectType{
		GoalByAmount: true,
		GoalByPeople: true,
		EndByGoalGain: false,
	}
	st, err := GetStrategy(pt, s.mockProject)
	s.Require().NoError(err)
	s.Require().True(isMoneyEqualStrategy(st))
}

func (s *StrategySuite) TestGetStrategyUnknown() {
	pt := &models.ProjectType{}
	st, err := GetStrategy(pt, s.mockProject)
	s.Require().Error(err)
	s.Require().Nil(st)
}

func (s *StrategySuite) TestMoneyPercent() {
	st := MoneyStrategy{s.mockProject}
	proj := &models.Project{
		Total:      344,
		GoalAmount: 1000,
	}
	s.Require().Equal(34, st.Percent(proj))
}

func (s *StrategySuite) TestMoneyPercentZero() {
	st := MoneyStrategy{s.mockProject}
	proj := &models.Project{
		Total:      344,
		GoalAmount: 0,
	}
	s.Require().Equal(0, st.Percent(proj))
}

func (s *StrategySuite) TestEventPercent() {
	st := EventStrategy{s.mockProject}
	proj := &models.Project{
		Total:      3,
		GoalPeople: 9,
	}
	s.Require().Equal(33, st.Percent(proj))
}

func (s *StrategySuite) TestEventPercentZero() {
	st := EventStrategy{s.mockProject}
	proj := &models.Project{
		Total:      4,
		GoalPeople: 0,
	}
	s.Require().Equal(0, st.Percent(proj))
}

func (s *StrategySuite) TestEventDatePercent() {
	st := EventDateStrategy{
		baseStrategy: EventStrategy{s.mockProject},
	}
	proj := &models.Project{
		Total:      4,
		GoalPeople: 9,
	}
	s.Require().Equal(44, st.Percent(proj))
}

func (s *StrategySuite) TestMoneyEqualPercent() {
	st := MoneyEqualStrategy{
		moneyStrategy: MoneyStrategy{s.mockProject},
		eventStrategy: EventStrategy{s.mockProject},
	}
	proj := &models.Project{
		Total:      2,
		GoalPeople: 7,
	}
	s.Require().Equal(28, st.Percent(proj))
}

func TestStrategySuite(t *testing.T) {
	suite.Run(t, new(StrategySuite))
}
