package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jonboulle/clockwork"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"

	"github.com/FreakyGranny/launchpad-api/internal/app/auth"
	"github.com/FreakyGranny/launchpad-api/internal/app/mocks"
	"github.com/FreakyGranny/launchpad-api/internal/app/models"
)

type LoginSuite struct {
	suite.Suite
	mockUserCtl     *gomock.Controller
	mockUser        *mocks.MockUserImpl
	mockProviderCtl *gomock.Controller
	mockProvider    *mocks.MockProvider
}

func (s *LoginSuite) SetupTest() {
	s.mockUserCtl = gomock.NewController(s.T())
	s.mockUser = mocks.NewMockUserImpl(s.mockUserCtl)

	s.mockProviderCtl = gomock.NewController(s.T())
	s.mockProvider = mocks.NewMockProvider(s.mockProviderCtl)
}

func (s *LoginSuite) TearDownTest() {
	s.mockUserCtl.Finish()
	s.mockProviderCtl.Finish()
}

func (s *LoginSuite) buildRequest(code string) (*http.Request, error) {
	body, err := json.Marshal(TokenRequest{Code: code})
	if err != nil {
		return nil, err
	}

	req := httptest.NewRequest(echo.POST, "/", bytes.NewBuffer(body))
	req.Header.Set("Content-type", "application/json")

	return req, nil
}

func (s *LoginSuite) TestWithCreateUser() {
	expectCode := "secret_code"
	req, err := s.buildRequest(expectCode)
	if err != nil {
		s.T().Fail()
	}

	e := echo.New()
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/login")

	h := NewAuthHandler(
		"secret",
		s.mockUser,
		s.mockProvider,
		clockwork.NewFakeClock(),
	)
	a := &auth.AccessData{
		AccessToken: "token",
		Expires:     123,
		UserID:      13,
		Email:       "some",
	}
	ud := &auth.UserData{
		Username:  "1",
		FirstName: "2",
		LastName:  "3",
		Avatar:    "4",
	}
	user := &models.User{}

	s.mockProvider.EXPECT().GetAccessToken(expectCode).Return(a, nil)
	s.mockProvider.EXPECT().GetUserData(a.UserID, a.AccessToken).Return(ud, nil)
	s.mockUser.EXPECT().FindByID(a.UserID).Return(user, false)
	s.mockUser.EXPECT().Create(user).Return(user, nil)

	s.Require().NoError(h.Login(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	s.Require().Equal(user.Username, ud.Username)
	s.Require().Equal(user.FirstName, ud.FirstName)
	s.Require().Equal(user.LastName, ud.LastName)
	s.Require().Equal(user.Avatar, ud.Avatar)
	s.Require().Equal(user.ID, a.UserID)
	s.Require().Equal(user.Email, a.Email)

	var tokenJSON = `{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6ZmFsc2UsImV4cCI6NDQ5ODg0OTIzLCJpZCI6MTN9.QDU7Og620wNdjDNFouk6jRmeBRqzhD9A6FgZa1x1gUs"}`

	s.Require().Equal(tokenJSON, strings.Trim(rec.Body.String(), "\n"))
}

func (s *LoginSuite) TestWithUpdateUser() {
	expectCode := "secret_code"
	req, err := s.buildRequest(expectCode)
	if err != nil {
		s.T().Fail()
	}

	e := echo.New()
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/login")

	h := NewAuthHandler(
		"secret",
		s.mockUser,
		s.mockProvider,
		clockwork.NewFakeClock(),
	)
	a := &auth.AccessData{
		AccessToken: "token",
		Expires:     123,
		UserID:      13,
		Email:       "some",
	}
	ud := &auth.UserData{
		Username:  "1",
		FirstName: "2",
		LastName:  "3",
		Avatar:    "4",
	}
	user := &models.User{
		ID:        13,
		Username:  "not updated",
		FirstName: "not updated",
		LastName:  "not updated",
		Avatar:    "not updated",
		Email:     "not updated",
	}

	s.mockProvider.EXPECT().GetAccessToken(expectCode).Return(a, nil)
	s.mockProvider.EXPECT().GetUserData(a.UserID, a.AccessToken).Return(ud, nil)
	s.mockUser.EXPECT().FindByID(a.UserID).Return(user, true)
	s.mockUser.EXPECT().Update(user).Return(user, nil)

	s.Require().NoError(h.Login(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	s.Require().Equal(user.Username, ud.Username)
	s.Require().Equal(user.FirstName, ud.FirstName)
	s.Require().Equal(user.LastName, ud.LastName)
	s.Require().Equal(user.Avatar, ud.Avatar)
	s.Require().Equal(user.ID, a.UserID)
	s.Require().Equal(user.Email, a.Email)

	var tokenJSON = `{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6ZmFsc2UsImV4cCI6NDQ5ODg0OTIzLCJpZCI6MTN9.QDU7Og620wNdjDNFouk6jRmeBRqzhD9A6FgZa1x1gUs"}`
	s.Require().Equal(tokenJSON, strings.Trim(rec.Body.String(), "\n"))
}

func TestLoginSuite(t *testing.T) {
	suite.Run(t, new(LoginSuite))
}
