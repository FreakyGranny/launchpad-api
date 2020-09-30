package app

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"github.com/FreakyGranny/launchpad-api/internal/mocks"
	"github.com/FreakyGranny/launchpad-api/internal/models"
)

type ProjectTypeSuite struct {
	suite.Suite
	mockProjectTypeCtl *gomock.Controller
	mockProjectType    *mocks.MockProjectTypeImpl
	app                *App
}

func (s *ProjectTypeSuite) SetupTest() {
	s.mockProjectTypeCtl = gomock.NewController(s.T())
	s.mockProjectType = mocks.NewMockProjectTypeImpl(s.mockProjectTypeCtl)
	s.app = New(nil, nil, nil, s.mockProjectType, nil, nil, nil, "", nil)
}

func (s *ProjectTypeSuite) TearDownTest() {
	s.mockProjectTypeCtl.Finish()
}

func (s *ProjectTypeSuite) TestGetAllProjectTypes() {
	projectTypes := []models.ProjectType{
		{
			ID:            1,
			Alias:         "other",
			Name:          "Other",
			Options:       []string{},
			GoalByAmount:  false,
			GoalByPeople:  true,
			EndByGoalGain: true,
		},
		{
			ID:            2,
			Alias:         "some",
			Name:          "Some",
			Options:       []string{},
			GoalByAmount:  true,
			GoalByPeople:  false,
			EndByGoalGain: true,
		},
	}
	s.mockProjectType.EXPECT().GetAll().Return(projectTypes, nil)
	pts, err := s.app.GetProjectTypes()
	s.Require().NoError(err)
	s.Require().Equal(projectTypes, pts)
}

func TestProjectTypeSuite(t *testing.T) {
	suite.Run(t, new(ProjectTypeSuite))
}
