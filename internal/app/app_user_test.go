package app

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"github.com/FreakyGranny/launchpad-api/internal/mocks"
	"github.com/FreakyGranny/launchpad-api/internal/models"
)

type UserSuite struct {
	suite.Suite
	mockUserCtl *gomock.Controller
	mockUser    *mocks.MockUserImpl
	app         *App
}

func (s *UserSuite) SetupTest() {
	s.mockUserCtl = gomock.NewController(s.T())
	s.mockUser = mocks.NewMockUserImpl(s.mockUserCtl)
	s.app = New(nil, s.mockUser, nil, nil, nil, nil, nil, "", nil)
}

func (s *UserSuite) TearDownTest() {
	s.mockUserCtl.Finish()
}

func (s *UserSuite) TestGetUserByID() {
	user := &models.User{
		ID:        1,
		Username:  "X",
		FirstName: "Y",
		LastName:  "Z",
		Avatar:    "A",
		Email:     "E",
	}
	pts := []models.Participation{
		{
			Cnt:           1,
			ProjectTypeID: 1,
		},
		{
			Cnt:           2,
			ProjectTypeID: 2,
		},
	}

	s.mockUser.EXPECT().Get(1).Return(user, true)
	s.mockUser.EXPECT().GetParticipation(1).Return(pts, nil)

	extUser, err := s.app.GetUser(1)
	s.Require().NoError(err)
	s.Require().NotNil(extUser)
}

func (s *UserSuite) TestGetParticipationErr() {
	s.mockUser.EXPECT().Get(1).Return(&models.User{}, true)
	s.mockUser.EXPECT().GetParticipation(1).Return(nil, errors.New("unexpected error"))

	extUser, err := s.app.GetUser(1)
	s.Require().Error(err)
	s.Require().Nil(extUser)
	s.Require().Equal(ErrGetUserParticipation, err)
}

func (s *UserSuite) TestGetUserNotFound() {
	s.mockUser.EXPECT().Get(1).Return(&models.User{}, false)

	extUser, err := s.app.GetUser(1)
	s.Require().Error(err)
	s.Require().Nil(extUser)
	s.Require().Equal(ErrUserNotFound, err)
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserSuite))
}
