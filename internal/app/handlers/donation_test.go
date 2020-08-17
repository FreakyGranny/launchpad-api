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

	"github.com/FreakyGranny/launchpad-api/internal/app/mocks"
	"github.com/FreakyGranny/launchpad-api/internal/app/models"
)

type DonationSuite struct {
	suite.Suite
	mockDonationCtl *gomock.Controller
	mockDonation    *mocks.MockDonationImpl
}

func (s *DonationSuite) SetupTest() {
	s.mockDonationCtl = gomock.NewController(s.T())
	s.mockDonation = mocks.NewMockDonationImpl(s.mockDonationCtl)
}

func (s *DonationSuite) TearDownTest() {
	s.mockDonationCtl.Finish()
}

func (s *DonationSuite) buildRequest() *http.Request {
	req := httptest.NewRequest(echo.GET, "/", bytes.NewBuffer(nil))
	req.Header.Set("Content-type", "application/json")

	return req
}

func (s *DonationSuite) TestGetProjectDonations() {
	req := s.buildRequest()

	e := echo.New()
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/donation/project/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	h := NewDonationHandler(s.mockDonation)

	donations := []models.Donation{
		{
			ID:      1,
			Payment: 100,
			Paid:    true,
			User: models.User{
				ID:        1,
				FirstName: "John",
				LastName:  "Doe",
			},
		},
		{
			ID:      2,
			Payment: 200,
			Paid:    true,
			User: models.User{
				ID:        2,
				FirstName: "Jane",
				LastName:  "Doe",
			},
		},
	}

	s.mockDonation.EXPECT().GetAllByProject(1).Return(donations, nil)

	s.Require().NoError(h.GetProjectDonations(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var pDonationsJSON = `[{"id":1,"user":{"id":1,"username":"","first_name":"John","last_name":"Doe","avatar":"","project_count":0,"success_rate":0},"locked":false,"paid":true},{"id":2,"user":{"id":2,"username":"","first_name":"Jane","last_name":"Doe","avatar":"","project_count":0,"success_rate":0},"locked":false,"paid":true}]`

	s.Require().Equal(pDonationsJSON, strings.Trim(rec.Body.String(), "\n"))
}

func (s *DonationSuite) TestGetUserDonations() {
	req := s.buildRequest()

	e := echo.New()
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/donation/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = float64(111)

	c.Set("user", token)

	h := NewDonationHandler(s.mockDonation)

	donations := []models.Donation{
		{
			ID:      1,
			Payment: 100,
			Paid:    true,
			ProjectID: 10,
		},
		{
			ID:      2,
			Payment: 200,
			Paid:    true,
			ProjectID: 20,
		},
	}

	s.mockDonation.EXPECT().GetAllByUser(111).Return(donations, nil)

	s.Require().NoError(h.GetUserDonations(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var pDonationsJSON = `[{"id":1,"payment":100,"locked":false,"paid":true,"project":10},{"id":2,"payment":200,"locked":false,"paid":true,"project":20}]`

	s.Require().Equal(pDonationsJSON, strings.Trim(rec.Body.String(), "\n"))
}

func TestDonationSuite(t *testing.T) {
	suite.Run(t, new(DonationSuite))
}