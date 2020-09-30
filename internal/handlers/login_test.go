package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"

	"github.com/FreakyGranny/launchpad-api/internal/app"
	mockapp "github.com/FreakyGranny/launchpad-api/internal/app/mock"
)

type LoginSuite struct {
	suite.Suite
	mockAppCtl *gomock.Controller
	mockApp    *mockapp.MockApplication
}

func (s *LoginSuite) SetupTest() {
	s.mockAppCtl = gomock.NewController(s.T())
	s.mockApp = mockapp.NewMockApplication(s.mockAppCtl)
}

func (s *LoginSuite) TearDownTest() {
	s.mockAppCtl.Finish()
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

func (s *LoginSuite) TestSuccessLogin() {
	expectCode := "secret_code"
	req, err := s.buildRequest(expectCode)
	if err != nil {
		s.T().Fail()
	}
	e := echo.New()
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/login")

	s.mockApp.EXPECT().Authentificate(expectCode).Return("MOCKED_TOKEN", nil)

	h := NewAuthHandler(s.mockApp)
	s.Require().NoError(h.Login(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	s.Require().Equal(`{"token":"MOCKED_TOKEN"}`, strings.Trim(rec.Body.String(), "\n"))
}

func (s *LoginSuite) TestBadRequest() {
	req := httptest.NewRequest(echo.POST, "/", bytes.NewBuffer([]byte("this is JSON body?")))
	req.Header.Set("Content-type", "application/json")

	e := echo.New()
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/login")

	h := NewAuthHandler(s.mockApp)
	s.Require().NoError(h.Login(c))
	s.Require().Equal(http.StatusBadRequest, rec.Code)
}

func (s *LoginSuite) TestGetTokenFail() {
	expectCode := "secret_code"
	req, err := s.buildRequest(expectCode)
	if err != nil {
		s.T().Fail()
	}
	e := echo.New()
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/login")

	s.mockApp.EXPECT().Authentificate(expectCode).Return("", app.ErrGetAccessTokenFailed)

	h := NewAuthHandler(s.mockApp)
	s.Require().NoError(h.Login(c))
	s.Require().Equal(http.StatusUnauthorized, rec.Code)

	s.Require().Equal(`{"error":"unable to authentificate"}`, strings.Trim(rec.Body.String(), "\n"))
}

func (s *LoginSuite) TestGetUserDataFail() {
	expectCode := "secret_code"
	req, err := s.buildRequest(expectCode)
	if err != nil {
		s.T().Fail()
	}
	e := echo.New()
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/login")

	s.mockApp.EXPECT().Authentificate(expectCode).Return("", app.ErrGetUserDataFailed)

	h := NewAuthHandler(s.mockApp)
	s.Require().NoError(h.Login(c))
	s.Require().Equal(http.StatusUnauthorized, rec.Code)

	s.Require().Equal(`{"error":"unable to authentificate"}`, strings.Trim(rec.Body.String(), "\n"))
}

func (s *LoginSuite) TestLoginWithUnexpectedError() {
	expectCode := "secret_code"
	req, err := s.buildRequest(expectCode)
	if err != nil {
		s.T().Fail()
	}
	e := echo.New()
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/login")

	s.mockApp.EXPECT().Authentificate(expectCode).Return("", errors.New("terrible error"))

	h := NewAuthHandler(s.mockApp)
	s.Require().NoError(h.Login(c))
	s.Require().Equal(http.StatusInternalServerError, rec.Code)

	s.Require().Equal(`{"error":"unexpected error"}`, strings.Trim(rec.Body.String(), "\n"))
}

func TestLoginSuite(t *testing.T) {
	suite.Run(t, new(LoginSuite))
}
