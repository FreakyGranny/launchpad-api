package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"

	"github.com/FreakyGranny/launchpad-api/internal/app/mocks"
	"github.com/FreakyGranny/launchpad-api/internal/app/models"
)

type ProjectTypeSuite struct {
	suite.Suite
	mockProjectTypeCtl *gomock.Controller
	mockProjectType    *mocks.MockProjectTypeImpl
}

func (s *ProjectTypeSuite) SetupTest() {
	s.mockProjectTypeCtl = gomock.NewController(s.T())
	s.mockProjectType = mocks.NewMockProjectTypeImpl(s.mockProjectTypeCtl)
}

func (s *ProjectTypeSuite) TearDownTest() {
	s.mockProjectTypeCtl.Finish()
}

func (s *ProjectTypeSuite) buildRequest() *http.Request {
	req := httptest.NewRequest(echo.GET, "/", bytes.NewBuffer(nil))
	req.Header.Set("Content-type", "application/json")

	return req
}

func (s *ProjectTypeSuite) TestGetAllProjectTypes() {
	req := s.buildRequest()

	e := echo.New()
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/project_type")

	h := NewProjectTypeHandler(s.mockProjectType)

	projectTypes := []models.ProjectType{
		{
			ID:    1,
			Alias: "other",
			Name:  "Other",
			Options: []string{},
			GoalByAmount: false,
			GoalByPeople: true,
			EndByGoalGain: true,
		},
		{
			ID:    2,
			Alias: "some",
			Name:  "Some",
			Options: []string{},
			GoalByAmount: true,
			GoalByPeople: false,
			EndByGoalGain: true,
		},
	}

	s.mockProjectType.EXPECT().GetAll().Return(projectTypes, nil)

	s.Require().NoError(h.GetProjectTypes(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var ptJSON = "[{\"id\":1,\"alias\":\"other\",\"name\":\"Other\",\"options\":[],\"goal_by_people\":true,\"goal_by_amount\":false,\"end_by_goal_gain\":true},{\"id\":2,\"alias\":\"some\",\"name\":\"Some\",\"options\":[],\"goal_by_people\":false,\"goal_by_amount\":true,\"end_by_goal_gain\":true}]\n"

	s.Require().Equal(ptJSON, rec.Body.String())
}

func TestProjectTypeSuite(t *testing.T) {
	suite.Run(t, new(ProjectTypeSuite))
}
