package app

import (
	"errors"
	"time"

	"github.com/FreakyGranny/launchpad-api/internal/auth"
	"github.com/FreakyGranny/launchpad-api/internal/models"
	"github.com/jonboulle/clockwork"
)

// Application business logic.
type Application interface {
	GetCategories() ([]models.Category, error)
	GetUser(id int) (*ExtendedUser, error)
	Authentificate(code string) (string, error)
	GetProjectTypes() ([]models.ProjectType, error)
	GetProject(id int) (*ExtendedProject, error)
	GetProjectsWithPagination(category, projectType, page, pageSize int, onlyOpen bool) ([]*ExtendedProject, int, bool, error)
	GetUserProjects(user int, onlyContributed, onlyOwned bool) ([]*ExtendedProject, error)
	CreateProject(user, goalPeople, goalAmount, category, projectType int, title, subtitle, descr, imageLink, instructions string, releaseDate, eventTime time.Time) (int, error)
	UpdateProject(id, user, goalPeople, goalAmount, category, projectType int, title, subtitle, descr, imageLink, instructions string, releaseDate, eventTime time.Time, published, dropEventDate bool) (*ExtendedProject, error)
	DeleteProject(iserID, projectID int) error
}

// App launchpad instance.
type App struct {
	categoryModel    models.CategoryImpl
	userModel        models.UserImpl
	projectModel     models.ProjectImpl
	projectTypeModel models.ProjectTypeImpl
	jwtSecret        string
	provider         auth.Provider
	clock            clockwork.Clock
}

// New returns new app.
func New(
	category models.CategoryImpl,
	user models.UserImpl,
	project models.ProjectImpl,
	projectType models.ProjectTypeImpl,
	donation models.DonationImpl,
	provider auth.Provider,
	clock clockwork.Clock,
	jwtSecret string,
) *App {
	return &App{
		categoryModel:    category,
		userModel:        user,
		projectModel:     project,
		projectTypeModel: projectType,
		jwtSecret:        jwtSecret,
		clock:            clock,
		provider:         provider,
	}
}

// GetCategories returns all existing category.
func (a *App) GetCategories() ([]models.Category, error) {
	return a.categoryModel.GetAll()
}

// GetUser returns all existing category.
func (a *App) GetUser(id int) (*ExtendedUser, error) {
	user, ok := a.userModel.Get(id)
	if !ok {
		return nil, ErrUserNotFound
	}
	pts, err := a.userModel.GetParticipation(id)
	if err != nil {
		return nil, ErrGetUserParticipation
	}

	return &ExtendedUser{User: *user, Participation: pts}, nil
}

// Authentificate authentificate user with given secure code.
func (a *App) Authentificate(code string) (string, error) {
	var token string
	data, err := a.provider.GetAccessToken(code)
	if err != nil {
		return token, ErrGetAccessTokenFailed
	}
	user, userExist := a.userModel.Get(data.UserID)
	user.ID = data.UserID
	user.Email = data.Email

	userData, err := a.provider.GetUserData(data.UserID, data.AccessToken)
	if err != nil {
		return token, ErrGetUserDataFailed
	}
	user.Username = userData.Username
	user.FirstName = userData.FirstName
	user.LastName = userData.LastName
	user.Avatar = userData.Avatar

	if !userExist {
		_, err = a.userModel.Create(user)
	} else {
		_, err = a.userModel.Update(user)
	}
	if err != nil {
		return token, errors.New("unable to create/update user")
	}
	token, err = auth.CreateToken(a.clock, a.jwtSecret, data.Expires, user)
	if err != nil {
		return "", err
	}

	return token, nil
}

// GetProjectTypes returns all existing project types
func (a *App) GetProjectTypes() ([]models.ProjectType, error) {
	return a.projectTypeModel.GetAll()
}

// GetProjectsWithPagination returns list of projects.
func (a *App) GetProjectsWithPagination(category, projectType, page, pageSize int, onlyOpen bool) ([]*ExtendedProject, int, bool, error) {
	var next int
	var hasNext bool

	paginator, err := a.projectModel.GetProjectsWithPagination(category, projectType, page, pageSize, onlyOpen)
	if err != nil {
		return nil, next, hasNext, ErrProjectRetrieve
	}
	next, hasNext = paginator.NextPage()
	projects, err := paginator.Retrieve()
	if err != nil {
		return nil, next, hasNext, ErrProjectRetrieve
	}

	projectList, err := a.extendProjectList(projects)
	if err != nil {
		return nil, next, hasNext, err
	}

	return projectList, next, hasNext, nil
}

// GetUserProjects returns list of projects for user.
func (a *App) GetUserProjects(user int, onlyContributed, onlyOwned bool) ([]*ExtendedProject, error) {
	projects, err := a.projectModel.GetUserProjects(user, onlyContributed, onlyOwned)
	if err != nil {
		return nil, ErrProjectRetrieve
	}

	return a.extendProjectList(projects)
}

// GetProject returns project by given id.
func (a *App) GetProject(id int) (*ExtendedProject, error) {
	project, ok := a.projectModel.Get(id)
	if !ok {
		return nil, ErrProjectNotFound
	}

	return a.extendProject(project)
}

func (a *App) extendProjectList(projects *[]models.Project) ([]*ExtendedProject, error) {
	result := make([]*ExtendedProject, 0)
	for _, project := range *projects {
		entry, err := a.extendProject(&project)
		if err != nil {
			return nil, err
		}
		result = append(result, entry)
	}

	return result, nil
}

func (a *App) extendProject(project *models.Project) (*ExtendedProject, error) {
	strategy, err := GetStrategy(&project.ProjectType, a.projectModel)
	if err != nil {
		return nil, err
	}
	extended := &ExtendedProject{
		ID:           project.ID,
		Title:        project.Title,
		SubTitle:     project.SubTitle,
		Status:       project.Status(),
		ReleaseDate:  project.ReleaseDate.Format(DateLayout),
		ImageLink:    project.ImageLink,
		Total:        project.Total,
		Percent:      strategy.Percent(project),
		Category:     project.Category,
		ProjectType:  project.ProjectType,
		GoalPeople:   project.GoalPeople,
		GoalAmount:   project.GoalAmount,
		Description:  project.Description,
		Instructions: project.Instructions,
		Owner:        project.Owner,
	}

	if !project.EventDate.IsZero() {
		ed := project.EventDate.Format(DateTimeLayout)
		extended.EventDate = &ed
	}

	return extended, nil
}

// CreateProject creates new prject.
func (a *App) CreateProject(user, goalPeople, goalAmount, category, projectType int, title, subtitle, descr, imageLink, instructions string, releaseDate, eventTime time.Time) (int, error) {
	newProject := models.Project{
		OwnerID:       user,
		Title:         title,
		SubTitle:      subtitle,
		ReleaseDate:   releaseDate,
		EventDate:     eventTime,
		GoalPeople:    goalPeople,
		GoalAmount:    goalAmount,
		Description:   descr,
		ImageLink:     imageLink,
		Instructions:  instructions,
		CategoryID:    category,
		ProjectTypeID: projectType,
		Closed:        false,
		Locked:        false,
		Published:     false,
		Total:         0,
	}

	return newProject.ID, a.projectModel.Create(&newProject)
}

// UpdateProject updates prject.
func (a *App) UpdateProject(id, user, goalPeople, goalAmount, category, projectType int, title, subtitle, descr, imageLink, instructions string, releaseDate, eventTime time.Time, published, dropEventDate bool) (*ExtendedProject, error) {
	project, ok := a.projectModel.Get(id)
	if !ok {
		return nil, ErrProjectNotFound
	}
	if project.Published || project.OwnerID != user {
		return nil, ErrProjectModifyNotAllowed
	}

	if dropEventDate {
		err := a.projectModel.DropEventDate(project)
		if err != nil {
			return nil, err
		}
		return a.extendProject(project)
	}

	project.Title = title
	project.SubTitle = subtitle
	project.Instructions = instructions
	project.Description = descr
	project.ImageLink = imageLink
	project.CategoryID = category
	project.ProjectTypeID = projectType
	project.GoalAmount = goalAmount
	project.GoalPeople = goalPeople
	project.ReleaseDate = releaseDate
	project.EventDate = eventTime
	project.Published = published

	err := a.projectModel.Update(project)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return a.extendProject(project)
}

// DeleteProject deletes project with given id.
func (a *App) DeleteProject(userID, projectID int) error {
	project, ok := a.projectModel.Get(projectID)
	if !ok {
		return ErrProjectNotFound
	}
	if project.Published || project.OwnerID != userID {
		return ErrProjectModifyNotAllowed
	}

	return a.projectModel.Delete(project)
}
