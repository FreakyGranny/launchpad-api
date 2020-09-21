package app

import (
	"errors"

	"github.com/FreakyGranny/launchpad-api/internal/auth"
	"github.com/FreakyGranny/launchpad-api/internal/models"
	"github.com/jonboulle/clockwork"
)

// ExtendedUser user extended with participation.
type ExtendedUser struct {
	models.User
	Participation []models.Participation `json:"participation"`
}

// Application business logic.
type Application interface {
	GetCategories() ([]models.Category, error)
	GetUser(id int) (*ExtendedUser, error)
	Authentificate(code string) (string, error)
	GetProjectTypes() ([]models.ProjectType, error)
}

// App launchpad instance.
type App struct {
	categoryModel    models.CategoryImpl
	userModel        models.UserImpl
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
