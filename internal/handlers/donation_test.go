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
	mockapp "github.com/FreakyGranny/launchpad-api/internal/app/mock"
	"github.com/FreakyGranny/launchpad-api/internal/models"
)

type DonationSuite struct {
	suite.Suite
	mockAppCtl *gomock.Controller
	mockApp    *mockapp.MockApplication
}

func (s *DonationSuite) SetupTest() {
	s.mockAppCtl = gomock.NewController(s.T())
	s.mockApp = mockapp.NewMockApplication(s.mockAppCtl)
}

func (s *DonationSuite) TearDownTest() {
	s.mockAppCtl.Finish()
}

func (s *DonationSuite) buildRequest() *http.Request {
	req := httptest.NewRequest(echo.GET, "/", bytes.NewBuffer(nil))
	req.Header.Set("Content-type", "application/json")

	return req
}

func (s *DonationSuite) TestGetProjectDonations() {
	e := echo.New()
	rec := httptest.NewRecorder()
	c := e.NewContext(s.buildRequest(), rec)
	c.SetPath("/donation/project/:id")
	c.SetParamNames("id")
	c.SetParamValues("1")

	h := NewDonationHandler(s.mockApp)

	donations := []app.ShortDonation{
		{
			ID:      1,
			Paid:    true,
			User: models.User{
				ID:        1,
				FirstName: "John",
				LastName:  "Doe",
			},
		},
		{
			ID:      2,
			Paid:    true,
			User: models.User{
				ID:        2,
				FirstName: "Jane",
				LastName:  "Doe",
			},
		},
	}
	s.mockApp.EXPECT().GetProjectDonations(1).Return(donations, nil)
	s.Require().NoError(h.GetProjectDonations(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var pDonationsJSON = `[{"id":1,"user":{"id":1,"username":"","first_name":"John","last_name":"Doe","avatar":"","project_count":0,"success_rate":0},"locked":false,"paid":true},{"id":2,"user":{"id":2,"username":"","first_name":"Jane","last_name":"Doe","avatar":"","project_count":0,"success_rate":0},"locked":false,"paid":true}]`

	s.Require().Equal(pDonationsJSON, strings.Trim(rec.Body.String(), "\n"))
}

func (s *DonationSuite) TestGetUserDonations() {
	e := echo.New()
	rec := httptest.NewRecorder()
	c := e.NewContext(s.buildRequest(), rec)
	c.SetPath("/donation")

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = float64(111)
	c.Set("user", token)

	h := NewDonationHandler(s.mockApp)
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
	s.mockApp.EXPECT().GetUserDonations(111).Return(donations, nil)
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
	e := echo.New()
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/donation")

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = float64(111)
	c.Set("user", token)
	expect := &models.Donation{
		ID: 111,
		ProjectID: reqStruct.ProjectID,
		Payment: reqStruct.Payment,
	}

	h := NewDonationHandler(s.mockApp)
	s.mockApp.EXPECT().CreateDonation(111, reqStruct.ProjectID, reqStruct.Payment).Return(expect, nil)
	s.Require().NoError(h.CreateDonation(c))
	s.Require().Equal(http.StatusCreated, rec.Code)

	var pDonationsJSON = `{"id":111,"payment":100,"locked":false,"paid":false,"project":10}`

	s.Require().Equal(pDonationsJSON, strings.Trim(rec.Body.String(), "\n"))
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

	h := NewDonationHandler(s.mockApp)
	donation := &models.Donation{
		ID:        1,
		Payment:   200,
		UserID:    111,
		Paid:      false,
		Locked:    false,
		ProjectID: 33,
	}
	s.mockApp.EXPECT().UpdateDonation(1, 111, 200, false).Return(donation, nil)
	s.Require().NoError(h.UpdateDonation(c))
	s.Require().Equal(http.StatusOK, rec.Code)

	var pDonationsJSON = `{"id":1,"payment":200,"locked":false,"paid":false,"project":33}`
	s.Require().Equal(pDonationsJSON, strings.Trim(rec.Body.String(), "\n"))
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

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = float64(111)
	c.Set("user", token)

	h := NewDonationHandler(s.mockApp)
	s.mockApp.EXPECT().DeleteDonation(1, 111).Return(nil)

	s.Require().NoError(h.DeleteDonation(c))
	s.Require().Equal(http.StatusNoContent, rec.Code)
}

func TestDonationSuite(t *testing.T) {
	suite.Run(t, new(DonationSuite))
}
