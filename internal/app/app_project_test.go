package app

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"github.com/FreakyGranny/launchpad-api/internal/mocks"
	"github.com/FreakyGranny/launchpad-api/internal/models"
)

type ProjectSuite struct {
	suite.Suite
	mockProjectCtl   *gomock.Controller
	mockProject      *mocks.MockProjectImpl
	mockPaginatorCtl *gomock.Controller
	mockPaginator    *mocks.MockProjectPaginatorImpl
	app              *App
}

func (s *ProjectSuite) SetupTest() {
	s.mockProjectCtl = gomock.NewController(s.T())
	s.mockProject = mocks.NewMockProjectImpl(s.mockProjectCtl)
	s.mockPaginatorCtl = gomock.NewController(s.T())
	s.mockPaginator = mocks.NewMockProjectPaginatorImpl(s.mockPaginatorCtl)
	s.app = New(nil, nil, s.mockProject, nil, nil, nil, nil, "", nil)
}

func (s *ProjectSuite) TearDownTest() {
	s.mockProjectCtl.Finish()
	s.mockPaginatorCtl.Finish()
}

func (s *ProjectSuite) TestGetSingleProject() {
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
			ID:            1,
			GoalByAmount:  true,
			EndByGoalGain: true,
		},
		Total:      344,
		GoalAmount: 1000,
	}
	expect := &ExtendedProject{
		ID:           project.ID,
		Title:        project.Title,
		SubTitle:     project.SubTitle,
		Status:       "search",
		ReleaseDate:  "0001-01-01",
		ImageLink:    project.ImageLink,
		Total:        project.Total,
		Percent:      34,
		Category:     project.Category,
		ProjectType:  project.ProjectType,
		GoalPeople:   project.GoalPeople,
		GoalAmount:   project.GoalAmount,
		Description:  project.Description,
		Instructions: project.Instructions,
		Owner:        project.Owner,
	}

	s.mockProject.EXPECT().Get(1).Return(project, true)
	pr, err := s.app.GetProject(1)
	s.Require().NoError(err)
	s.Require().Equal(expect, pr)
}

func (s *ProjectSuite) TestCreateProject() {
	title := "project"
	subtitle := "some subtitle"
	releaseDate := time.Date(2020, 8, 20, 0, 0, 0, 0, time.UTC)
	eventTime := time.Time{}
	category := 1
	goalPeople := 0
	goalAmount := 1000
	imageLink := "https://avatar.com"
	instructions := "instructions"
	descr := "description"
	projectType := 1
	userID := 113

	expect := models.Project{
		Title:         title,
		SubTitle:      subtitle,
		ReleaseDate:   releaseDate,
		EventDate:     eventTime,
		GoalPeople:    goalPeople,
		GoalAmount:    goalAmount,
		Description:   descr,
		ImageLink:     imageLink,
		Instructions:  instructions,
		OwnerID:       userID,
		CategoryID:    category,
		ProjectTypeID: projectType,
	}
	s.mockProject.EXPECT().Create(&expect).Return(nil)
	id, err := s.app.CreateProject(
		userID, 
		goalPeople, 
		goalAmount, 
		category, 
		projectType, 
		title, 
		subtitle, 
		descr, 
		imageLink, 
		instructions, 
		releaseDate, 
		eventTime,
	)
	s.Require().NoError(err)
	s.Require().Equal(0, id)
}

func (s *ProjectSuite) TestUpdateProject() {
	expect := &models.Project{
		ID:    17,
		Title: "before_Title",
		ProjectType: models.ProjectType{
			GoalByAmount:  true,
			EndByGoalGain: true,
		},
		OwnerID: 42,
	}
	s.mockProject.EXPECT().Get(17).Return(expect, true)
	s.mockProject.EXPECT().Update(expect).Return(nil)
	eProject, err := s.app.UpdateProject(17, 42, 0, 0, 0, 0, "ChangeProject", "", "", "", "", time.Time{}, time.Time{}, false, false)
	s.Require().NoError(err)
	s.Require().Equal("ChangeProject", eProject.Title)
}

func (s *ProjectSuite) TestDeleteProject() {
	expect := &models.Project{
		ID:      1,
		OwnerID: 111,
	}
	s.mockProject.EXPECT().Get(1).Return(expect, true)
	s.mockProject.EXPECT().Delete(expect).Return(nil)
	s.Require().NoError(s.app.DeleteProject(111, 1))
}

func (s *ProjectSuite) TestGetProjectsWithPagination() {
	category := 1
	projectType := 2
	page := 1
	pageSize := 2

	s.mockProject.EXPECT().GetProjectsWithPagination(category, projectType, page, pageSize, false).Return(s.mockPaginator, nil)
	s.mockPaginator.EXPECT().NextPage().Return(0, false)
	s.mockPaginator.EXPECT().Retrieve().Return(s.makeProjectList(), nil)

	list, next, hasNext, err := s.app.GetProjectsWithPagination(category, projectType, page, pageSize, false)
	s.Require().NoError(err)
	s.Require().Equal(2, len(list))
	s.Require().Equal(0, next)
	s.Require().False(hasNext)
}

func (s *ProjectSuite) makeProjectList() *[]models.Project {
	return &[]models.Project{
		{
			ProjectType: models.ProjectType{
				ID:            1,
				GoalByAmount:  true,
				EndByGoalGain: true,
			},
		},
		{
			ProjectType: models.ProjectType{
				ID:            2,
				GoalByPeople:  true,
				EndByGoalGain: true,
			},
		},
	}
}

func (s *ProjectSuite) TestGetUserProjects() {
	s.mockProject.EXPECT().GetUserProjects(1, false, false).Return(s.makeProjectList(), nil)

	list, err := s.app.GetUserProjects(1, false, false)
	s.Require().NoError(err)
	s.Require().Equal(2, len(list))
}

func TestProjectSuite(t *testing.T) {
	suite.Run(t, new(ProjectSuite))
}
