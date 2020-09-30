package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"

	"github.com/FreakyGranny/launchpad-api/internal/app"
	mockapp "github.com/FreakyGranny/launchpad-api/internal/app/mock"
	"github.com/FreakyGranny/launchpad-api/internal/models"
)

type UserSuite struct {
	suite.Suite
	mockAppCtl *gomock.Controller
	mockApp    *mockapp.MockApplication
}

func (s *UserSuite) SetupTest() {
	s.mockAppCtl = gomock.NewController(s.T())
	s.mockApp = mockapp.NewMockApplication(s.mockAppCtl)
}

func (s *UserSuite) TearDownTest() {
	s.mockAppCtl.Finish()
}

func (s *UserSuite) buildRequest() *http.Request {
	req := httptest.NewRequest(echo.GET, "/", bytes.NewBuffer(nil))
	req.Header.Set("Content-type", "application/json")

	return req
}

func (s *UserSuite) TestGetUserByID() {
	req := s.buildRequest()

	e := echo.New()
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/user/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	h := NewUserHandler(s.mockApp)
	user := &app.ExtendedUser{
		User: models.User{
			ID:        1,
			Username:  "X",
			FirstName: "Y",
			LastName:  "Z",
			Avatar:    "A",
			Email:     "E",
		},
		Participation: []models.Participation{
			{
				Cnt:           1,
				ProjectTypeID: 1,
			},
			{
				Cnt:           2,
				ProjectTypeID: 2,
			},
		},
	}

	s.mockApp.EXPECT().GetUser(1).Return(user, nil)

	s.Require().NoError(h.GetUser(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var tokenJSON = `{"id":1,"username":"X","first_name":"Y","last_name":"Z","avatar":"A","project_count":0,"success_rate":0,"participation":[{"count":1,"id":1},{"count":2,"id":2}]}`
	s.Require().Equal(tokenJSON, strings.Trim(rec.Body.String(), "\n"))
}

func (s *UserSuite) TestGetUserNotFound() {
	req := s.buildRequest()

	e := echo.New()
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/user/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	h := NewUserHandler(s.mockApp)

	s.mockApp.EXPECT().GetUser(1).Return(nil, app.ErrUserNotFound)

	s.Require().NoError(h.GetUser(c))
	s.Require().Equal(http.StatusNotFound, rec.Code)
}

func (s *UserSuite) TestGetCurrentUser() {
	req := s.buildRequest()

	e := echo.New()
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/user/")

	user := &app.ExtendedUser{
		User: models.User{
			ID:        1,
			Username:  "X",
			FirstName: "Y",
			LastName:  "Z",
			Avatar:    "A",
			Email:     "E",
		},
		Participation: []models.Participation{
			{
				Cnt:           1,
				ProjectTypeID: 1,
			},
			{
				Cnt:           2,
				ProjectTypeID: 2,
			},
		},
	}
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = float64(user.ID)

	c.Set("user", token)
	h := NewUserHandler(s.mockApp)

	s.mockApp.EXPECT().GetUser(1).Return(user, nil)
	s.Require().NoError(h.GetCurrentUser(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var tokenJSON = `{"id":1,"username":"X","first_name":"Y","last_name":"Z","avatar":"A","project_count":0,"success_rate":0,"participation":[{"count":1,"id":1},{"count":2,"id":2}]}`
	s.Require().Equal(tokenJSON, strings.Trim(rec.Body.String(), "\n"))
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserSuite))
}
