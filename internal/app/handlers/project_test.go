package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"

	"github.com/FreakyGranny/launchpad-api/internal/app/mocks"
	"github.com/FreakyGranny/launchpad-api/internal/app/models"
)

type ProjectSuite struct {
	suite.Suite
	mockProjectCtl   *gomock.Controller
	mockProject      *mocks.MockProjectImpl
	mockPaginatorCtl *gomock.Controller
	mockPaginator    *mocks.MockProjectPaginatorImpl
}

func (s *ProjectSuite) SetupTest() {
	s.mockProjectCtl = gomock.NewController(s.T())
	s.mockProject = mocks.NewMockProjectImpl(s.mockProjectCtl)
	s.mockPaginatorCtl = gomock.NewController(s.T())
	s.mockPaginator = mocks.NewMockProjectPaginatorImpl(s.mockPaginatorCtl)
}

func (s *ProjectSuite) TearDownTest() {
	s.mockProjectCtl.Finish()
	s.mockPaginatorCtl.Finish()
}

func (s *ProjectSuite) buildRequest() *http.Request {
	req := httptest.NewRequest(echo.GET, "/", bytes.NewBuffer(nil))
	req.Header.Set("Content-type", "application/json")

	return req
}

func (s *ProjectSuite) makeProjectList() *[]models.Project {
	return &[]models.Project{
		{
			ID:        1,
			Title:     "Title",
			SubTitle:  "Subtitle",
			Locked:    true,
			Published: true,
			Closed:    true,
			Category: models.Category{
				ID: 1,
			},
			ProjectType: models.ProjectType{
				ID: 1,
				GoalByAmount: true,
				EndByGoalGain: true,
			},
		},
		{
			ID:        2,
			Title:     "Second Project",
			SubTitle:  "2 Subtitle",
			Locked:    true,
			Published: true,
			Closed:    true,
			Category: models.Category{
				ID: 2,
			},
			ProjectType: models.ProjectType{
				ID: 2,
				GoalByPeople: true,
				EndByGoalGain: true,
			},
		},
	}
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
		Owner: models.User{
			ID:        1,
			FirstName: "John",
			LastName:  "Doe",
		},
		Category: models.Category{
			ID: 1,
		},
		ProjectType: models.ProjectType{
			ID: 1,
			GoalByAmount: true,
			EndByGoalGain: true,
		},
		Total: 344,
		GoalAmount: 1000,
	}

	s.mockProject.EXPECT().Get(1).Return(project, true)

	s.Require().NoError(h.GetSingleProject(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var pJSON = `{"id":1,"title":"Title","subtitle":"Subtitle","status":"search","release_date":"0001-01-01","event_date":null,"image_link":"","total":344,"percent":34,"category":{"id":1,"alias":"","name":""},"project_type":{"id":1,"alias":"","name":"","options":null,"goal_by_people":false,"goal_by_amount":true,"end_by_goal_gain":true},"goal_people":0,"goal_amount":1000,"description":"","instructions":"","owner":{"id":1,"username":"","first_name":"John","last_name":"Doe","avatar":"","project_count":0,"success_rate":0}}`

	s.Require().Equal(pJSON, strings.Trim(rec.Body.String(), "\n"))
}

func (s *ProjectSuite) TestGetProjectsWithPagination() {
	req := s.buildRequest()

	e := echo.New()
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/project")
	page := 1
	pageSize := 2
	c.QueryParams().Add("category", "1")
	c.QueryParams().Add("type", "2")
	c.QueryParams().Add("open", "true")

	c.QueryParams().Add("page", strconv.Itoa(page))
	c.QueryParams().Add("page_size", strconv.Itoa(pageSize))

	h := NewProjectHandler(s.mockProject)

	filter := models.ProjectListFilter{
		Category:    1,
		ProjectType: 2,
		OnlyOpen:    true,
		PageSize:    pageSize,
		Page:        page,
	}
	s.mockProject.EXPECT().GetProjectsWithPagination(&filter).Return(s.mockPaginator, nil)

	s.mockPaginator.EXPECT().NextPage().Return(0, false)
	s.mockPaginator.EXPECT().Retrieve().Return(s.makeProjectList(), nil)

	s.Require().NoError(h.GetProjects(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var pJSON = `{"results":[{"id":1,"title":"Title","subtitle":"Subtitle","status":"success","release_date":"0001-01-01","event_date":null,"image_link":"","total":0,"percent":0,"category":{"id":1,"alias":"","name":""},"project_type":{"id":1,"alias":"","name":"","options":null,"goal_by_people":false,"goal_by_amount":true,"end_by_goal_gain":true}},{"id":2,"title":"Second Project","subtitle":"2 Subtitle","status":"success","release_date":"0001-01-01","event_date":null,"image_link":"","total":0,"percent":0,"category":{"id":2,"alias":"","name":""},"project_type":{"id":2,"alias":"","name":"","options":null,"goal_by_people":true,"goal_by_amount":false,"end_by_goal_gain":true}}],"next":0,"has_next":false}`

	s.Require().Equal(pJSON, strings.Trim(rec.Body.String(), "\n"))
}

func (s *ProjectSuite) TestGetUserProjects() {
	req := s.buildRequest()

	e := echo.New()
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/project/user/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")
	c.QueryParams().Add("owned", "true")
	c.QueryParams().Add("contributed", "true")

	h := NewProjectHandler(s.mockProject)
	filter := models.ProjectUserFilter{
		UserID:      1,
		Owned:       true,
		Contributed: true,
	}
	s.mockProject.EXPECT().GetUserProjects(&filter).Return(s.makeProjectList(), nil)

	s.Require().NoError(h.GetUserProjects(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var pJSON = `[{"id":1,"title":"Title","subtitle":"Subtitle","status":"success","release_date":"0001-01-01","event_date":null,"image_link":"","total":0,"percent":0,"category":{"id":1,"alias":"","name":""},"project_type":{"id":1,"alias":"","name":"","options":null,"goal_by_people":false,"goal_by_amount":true,"end_by_goal_gain":true}},{"id":2,"title":"Second Project","subtitle":"2 Subtitle","status":"success","release_date":"0001-01-01","event_date":null,"image_link":"","total":0,"percent":0,"category":{"id":2,"alias":"","name":""},"project_type":{"id":2,"alias":"","name":"","options":null,"goal_by_people":true,"goal_by_amount":false,"end_by_goal_gain":true}}]`

	s.Require().Equal(pJSON, strings.Trim(rec.Body.String(), "\n"))
}

func TestProjectSuite(t *testing.T) {
	suite.Run(t, new(ProjectSuite))
}
