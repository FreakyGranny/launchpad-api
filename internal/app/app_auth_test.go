package app

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/suite"

	"github.com/FreakyGranny/launchpad-api/internal/auth"
	"github.com/FreakyGranny/launchpad-api/internal/mocks"
	"github.com/FreakyGranny/launchpad-api/internal/models"
)

type AuthSuite struct {
	suite.Suite
	mockUserCtl     *gomock.Controller
	mockUser        *mocks.MockUserImpl
	mockProviderCtl *gomock.Controller
	mockProvider    *mocks.MockProvider
	app             *App
}

func (s *AuthSuite) SetupTest() {
	s.mockUserCtl = gomock.NewController(s.T())
	s.mockUser = mocks.NewMockUserImpl(s.mockUserCtl)

	s.mockProviderCtl = gomock.NewController(s.T())
	s.mockProvider = mocks.NewMockProvider(s.mockProviderCtl)

	s.app = New(nil, s.mockUser, nil, nil, nil, s.mockProvider, clockwork.NewFakeClock(), "secret", nil)
}

func (s *AuthSuite) TearDownTest() {
	s.mockUserCtl.Finish()
	s.mockProviderCtl.Finish()
}

func (s *AuthSuite) TestWithCreateUser() {
	expectCode := "secret_code"
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
	s.mockUser.EXPECT().Get(a.UserID).Return(user, false)
	s.mockUser.EXPECT().Create(user).Return(user, nil)

	t, err := s.app.Authentificate(expectCode)
	s.Require().NoError(err)

	s.Require().Equal(user.Username, ud.Username)
	s.Require().Equal(user.FirstName, ud.FirstName)
	s.Require().Equal(user.LastName, ud.LastName)
	s.Require().Equal(user.Avatar, ud.Avatar)
	s.Require().Equal(user.ID, a.UserID)
	s.Require().Equal(user.Email, a.Email)

	s.Require().Equal("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6ZmFsc2UsImV4cCI6NDQ5ODg0OTIzLCJpZCI6MTN9.QDU7Og620wNdjDNFouk6jRmeBRqzhD9A6FgZa1x1gUs", t)
}

func (s *AuthSuite) TestWithUpdateUser() {
	expectCode := "secret_code"
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
	s.mockUser.EXPECT().Get(a.UserID).Return(user, true)
	s.mockUser.EXPECT().Update(user).Return(user, nil)

	t, err := s.app.Authentificate(expectCode)
	s.Require().NoError(err)

	s.Require().Equal(user.Username, ud.Username)
	s.Require().Equal(user.FirstName, ud.FirstName)
	s.Require().Equal(user.LastName, ud.LastName)
	s.Require().Equal(user.Avatar, ud.Avatar)
	s.Require().Equal(user.ID, a.UserID)
	s.Require().Equal(user.Email, a.Email)

	s.Require().Equal("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6ZmFsc2UsImV4cCI6NDQ5ODg0OTIzLCJpZCI6MTN9.QDU7Og620wNdjDNFouk6jRmeBRqzhD9A6FgZa1x1gUs", t)
}

func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(AuthSuite))
}
