package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"

	"github.com/FreakyGranny/launchpad-api/internal/app/mocks"
	"github.com/FreakyGranny/launchpad-api/internal/app/models"
)

type UserSuite struct {
	suite.Suite
	mockUserCtl *gomock.Controller
	mockUser    *mocks.MockUserImpl
}

func (s *UserSuite) SetupTest() {
	s.mockUserCtl = gomock.NewController(s.T())
	s.mockUser = mocks.NewMockUserImpl(s.mockUserCtl)
}

func (s *UserSuite) TearDownTest() {
	s.mockUserCtl.Finish()
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

	h := NewUserHandler(s.mockUser)

	user := &models.User{
		ID:        1,
		Username:  "X",
		FirstName: "Y",
		LastName:  "Z",
		Avatar:    "A",
		Email:     "E",
	}

	s.mockUser.EXPECT().FindByID(1).Return(user, true)

	s.Require().NoError(h.GetUser(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var tokenJSON = "{\"id\":1,\"username\":\"X\",\"first_name\":\"Y\",\"last_name\":\"Z\",\"avatar\":\"A\",\"is_admin\":false,\"project_count\":0,\"success_rate\":0}\n"

	s.Require().Equal(tokenJSON, rec.Body.String())
}

func (s *UserSuite) TestGetUserNotFound() {
	req := s.buildRequest()

	e := echo.New()
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/user/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	h := NewUserHandler(s.mockUser)

	s.mockUser.EXPECT().FindByID(1).Return(&models.User{}, false)

	s.Require().NoError(h.GetUser(c))
	s.Require().Equal(http.StatusNotFound, rec.Code)
}

func (s *UserSuite) TestGetCurrentUser() {
	req := s.buildRequest()

	e := echo.New()
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/user/")

	user := &models.User{
		ID:        1,
		Username:  "X",
		FirstName: "Y",
		LastName:  "Z",
		Avatar:    "A",
		Email:     "E",
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = float64(user.ID)

	c.Set("user", token)

	h := NewUserHandler(s.mockUser)

	s.mockUser.EXPECT().FindByID(1).Return(user, true)

	s.Require().NoError(h.GetCurrentUser(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var tokenJSON = "{\"id\":1,\"username\":\"X\",\"first_name\":\"Y\",\"last_name\":\"Z\",\"avatar\":\"A\",\"is_admin\":false,\"project_count\":0,\"success_rate\":0}\n"

	s.Require().Equal(tokenJSON, rec.Body.String())
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserSuite))
}
