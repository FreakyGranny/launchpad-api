package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"

	"github.com/FreakyGranny/launchpad-api/internal/app"
	"github.com/FreakyGranny/launchpad-api/internal/mocks"
	"github.com/FreakyGranny/launchpad-api/internal/models"
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
	recalcChan := make(chan int, 1)
	defer close(recalcChan)

	e := echo.New()
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/donation/project/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	app := app.New(nil, nil, nil, nil, s.mockDonation, nil, nil, "", recalcChan)
	h := NewDonationHandler(app)

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
	recalcChan := make(chan int, 1)
	defer close(recalcChan)

	e := echo.New()
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/donation")

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = float64(111)

	c.Set("user", token)

	app := app.New(nil, nil, nil, nil, s.mockDonation, nil, nil, "", recalcChan)
	h := NewDonationHandler(app)

	donations := []models.Donation{
		{
			ID:        1,
			Payment:   100,
			Paid:      true,
			ProjectID: 10,
		},
		{
			ID:        2,
			Payment:   200,
			Paid:      true,
			ProjectID: 20,
		},
	}
	s.mockDonation.EXPECT().GetAllByUser(111).Return(donations, nil)
	s.Require().NoError(h.GetUserDonations(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var pDonationsJSON = `[{"id":1,"payment":100,"locked":false,"paid":true,"project":10},{"id":2,"payment":200,"locked":false,"paid":true,"project":20}]`

	s.Require().Equal(pDonationsJSON, strings.Trim(rec.Body.String(), "\n"))
}

func (s *DonationSuite) TestCreateDonation() {
	reqStruct := DonationCreateRequest{
		ProjectID: 10,
		Payment:   100,
	}
	body, err := json.Marshal(reqStruct)
	if err != nil {
		s.T().Fail()
	}
	req := httptest.NewRequest(echo.POST, "/", bytes.NewBuffer(body))
	req.Header.Set("Content-type", "application/json")

	recalcChan := make(chan int, 1)
	defer close(recalcChan)

	e := echo.New()
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/donation")

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = float64(111)
	c.Set("user", token)

	app := app.New(nil, nil, nil, nil, s.mockDonation, nil, nil, "", recalcChan)
	h := NewDonationHandler(app)

	donation := models.Donation{
		Payment:   reqStruct.Payment,
		ProjectID: reqStruct.ProjectID,
		UserID:    111,
	}
	s.mockDonation.EXPECT().Create(&donation).Return(nil)
	s.Require().NoError(h.CreateDonation(c))
	s.Require().Equal(http.StatusCreated, rec.Code)

	var pDonationsJSON = `{"id":0,"payment":100,"locked":false,"paid":false,"project":10}`

	s.Require().Equal(pDonationsJSON, strings.Trim(rec.Body.String(), "\n"))
	s.Require().Equal(10, <-recalcChan)
}

func (s *DonationSuite) TestUpdateDonation() {
	body, err := json.Marshal(DonationUpdateRequest{
		Payment: 200,
		Paid:    false,
	})
	if err != nil {
		s.T().Fail()
	}
	req := httptest.NewRequest(echo.PATCH, "/", bytes.NewBuffer(body))
	req.Header.Set("Content-type", "application/json")

	recalcChan := make(chan int, 1)
	defer close(recalcChan)

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

	app := app.New(nil, nil, nil, nil, s.mockDonation, nil, nil, "", recalcChan)
	h := NewDonationHandler(app)

	donation := &models.Donation{
		ID:        1,
		Payment:   100,
		UserID:    111,
		Paid:      false,
		Locked:    false,
		ProjectID: 33,
	}
	s.mockDonation.EXPECT().Get(1).Return(donation, true)
	donation.Payment = 200
	s.mockDonation.EXPECT().Update(donation).Return(nil)
	s.Require().NoError(h.UpdateDonation(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var pDonationsJSON = `{"id":1,"payment":200,"locked":false,"paid":false,"project":33}`

	s.Require().Equal(pDonationsJSON, strings.Trim(rec.Body.String(), "\n"))
	s.Require().Equal(33, <-recalcChan)
}

func (s *DonationSuite) TestDeleteDonation() {
	req := httptest.NewRequest(echo.DELETE, "/", bytes.NewBuffer(nil))
	req.Header.Set("Content-type", "application/json")

	recalcChan := make(chan int, 1)
	defer close(recalcChan)

	e := echo.New()
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/donation/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	expect := &models.Donation{
		ID:        1,
		UserID:    111,
		ProjectID: 44,
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = float64(111)
	c.Set("user", token)

	app := app.New(nil, nil, nil, nil, s.mockDonation, nil, nil, "", recalcChan)
	h := NewDonationHandler(app)

	s.mockDonation.EXPECT().Get(1).Return(expect, true)
	s.mockDonation.EXPECT().Delete(expect).Return(nil)
	s.Require().NoError(h.DeleteDonation(c))
	s.Require().Equal(http.StatusNoContent, rec.Code)

	s.Require().Equal("", rec.Body.String())
	s.Require().Equal(44, <-recalcChan)
}

func TestDonationSuite(t *testing.T) {
	suite.Run(t, new(DonationSuite))
}
