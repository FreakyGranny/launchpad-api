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

type ProjectSuite struct {
	suite.Suite
	mockProjectCtl *gomock.Controller
	mockProject    *mocks.MockProjectImpl
}

func (s *ProjectSuite) SetupTest() {
	s.mockProjectCtl = gomock.NewController(s.T())
	s.mockProject = mocks.NewMockProjectImpl(s.mockProjectCtl)
}

func (s *ProjectSuite) TearDownTest() {
	s.mockProjectCtl.Finish()
}

func (s *ProjectSuite) buildRequest() *http.Request {
	req := httptest.NewRequest(echo.GET, "/", bytes.NewBuffer(nil))
	req.Header.Set("Content-type", "application/json")

	return req
}

func (s *ProjectSuite) TestGetSingleProject() {
	req := s.buildRequest()

	e := echo.New()
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/project/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	h := NewProjectHandler(s.mockProject)

	project := &models.Project{
		ID:        1,
		Title:     "Title",
		SubTitle:  "Subtitle",
		Locked:    false,
		Published: true,
		Closed:    false,
	}

	s.mockProject.EXPECT().Get(1).Return(project, true)

	s.Require().NoError(h.GetSingleProject(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	// var pJSON = `{"id":1,"title":"title","subtitle":"Subtitle","status":"","release_date":"0001-01-01","event_date":null,"image_link":"","total":0,"percent":0,"category":{"id":0,"alias":"","name":""},"project_type":{"id":0,"alias":"","name":"","options":null,"goal_by_people":false,"goal_by_amount":false,"end_by_goal_gain":false},"goal_people":0,"goal_amount":0,"description":"","instructions":"","owner":{"id":0,"username":"","first_name":"","last_name":"","avatar":"","project_count":0,"success_rate":0}}`

	// s.Require().Equal(pJSON, rec.Body.String())
}

func TestProjectSuite(t *testing.T) {
	suite.Run(t, new(ProjectSuite))
}
