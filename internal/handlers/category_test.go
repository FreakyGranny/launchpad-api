package handlers

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"

	mockapp "github.com/FreakyGranny/launchpad-api/internal/app/mock"
	"github.com/FreakyGranny/launchpad-api/internal/models"
)

type CategorySuite struct {
	suite.Suite
	mockAppCtl *gomock.Controller
	mockApp    *mockapp.MockApplication
}

func (s *CategorySuite) SetupTest() {
	s.mockAppCtl = gomock.NewController(s.T())
	s.mockApp = mockapp.NewMockApplication(s.mockAppCtl)
}

func (s *CategorySuite) TearDownTest() {
	s.mockAppCtl.Finish()
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

	h := NewCategoryHandler(s.mockApp)

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

	s.mockApp.EXPECT().GetCategories().Return(categories, nil)

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

	h := NewCategoryHandler(s.mockApp)
	s.mockApp.EXPECT().GetCategories().Return([]models.Category{}, nil)

	s.Require().NoError(h.GetCategories(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var emptyJSON = "[]\n"

	s.Require().Equal(emptyJSON, rec.Body.String())
}

func (s *CategorySuite) TestError() {
	req := s.buildRequest()

	e := echo.New()
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/category")

	h := NewCategoryHandler(s.mockApp)
	s.mockApp.EXPECT().GetCategories().Return(nil, errors.New("some error"))

	s.Require().NoError(h.GetCategories(c))
	s.Require().Equal(http.StatusInternalServerError, rec.Code)
}

func TestCategorySuite(t *testing.T) {
	suite.Run(t, new(CategorySuite))
}
