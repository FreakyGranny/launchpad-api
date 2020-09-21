package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"

	"github.com/FreakyGranny/launchpad-api/internal/app"
	"github.com/FreakyGranny/launchpad-api/internal/mocks"
	"github.com/FreakyGranny/launchpad-api/internal/models"
)

type CategorySuite struct {
	suite.Suite
	mockCategoryCtl *gomock.Controller
	mockCategory    *mocks.MockCategoryImpl
}

func (s *CategorySuite) SetupTest() {
	s.mockCategoryCtl = gomock.NewController(s.T())
	s.mockCategory = mocks.NewMockCategoryImpl(s.mockCategoryCtl)
}

func (s *CategorySuite) TearDownTest() {
	s.mockCategoryCtl.Finish()
}

func (s *CategorySuite) buildRequest() *http.Request {
	req := httptest.NewRequest(echo.GET, "/", bytes.NewBuffer(nil))
	req.Header.Set("Content-type", "application/json")

	return req
}

func (s *CategorySuite) TestGetAllCategories() {
	req := s.buildRequest()

	e := echo.New()
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/category")

	app := app.New(s.mockCategory, nil, nil, nil, nil, nil, nil, "")
	h := NewCategoryHandler(app)

	categories := []models.Category{
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

	s.mockCategory.EXPECT().GetAll().Return(categories, nil)

	s.Require().NoError(h.GetCategories(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var categoriesJSON = `[{"id":1,"alias":"other","name":"Other"},{"id":2,"alias":"some","name":"Some"}]`

	s.Require().Equal(categoriesJSON, strings.Trim(rec.Body.String(), "\n"))
}

func (s *CategorySuite) TestNoCategories() {
	req := s.buildRequest()

	e := echo.New()
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/category")

	app := app.New(s.mockCategory, nil, nil, nil, nil, nil, nil, "")
	h := NewCategoryHandler(app)
	s.mockCategory.EXPECT().GetAll().Return([]models.Category{}, nil)

	s.Require().NoError(h.GetCategories(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var emptyJSON = "[]\n"

	s.Require().Equal(emptyJSON, rec.Body.String())
}

func TestCategorySuite(t *testing.T) {
	suite.Run(t, new(CategorySuite))
}
