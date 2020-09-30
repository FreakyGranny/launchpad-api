package app

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"github.com/FreakyGranny/launchpad-api/internal/mocks"
	"github.com/FreakyGranny/launchpad-api/internal/models"
)

type CategorySuite struct {
	suite.Suite
	mockCategoryCtl *gomock.Controller
	mockCategory    *mocks.MockCategoryImpl
	app             *App
}

func (s *CategorySuite) SetupTest() {
	s.mockCategoryCtl = gomock.NewController(s.T())
	s.mockCategory = mocks.NewMockCategoryImpl(s.mockCategoryCtl)
	s.app = New(s.mockCategory, nil, nil, nil, nil, nil, nil, "", nil)
}

func (s *CategorySuite) TearDownTest() {
	s.mockCategoryCtl.Finish()
}

func (s *CategorySuite) TestGetAllCategories() {
	expect := []models.Category{
		{
			ID:    1,
			Alias: "other",
			Name:  "Other",
		},
		{
			ID:    2,
			Alias: "some",
			Name:  "Some",
		},
	}

	s.mockCategory.EXPECT().GetAll().Return(expect, nil)

	cat, err := s.app.GetCategories()
	s.Require().NoError(err)
	s.Require().Equal(expect, cat)
}

func TestCategorySuite(t *testing.T) {
	suite.Run(t, new(CategorySuite))
}
