package app

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"github.com/FreakyGranny/launchpad-api/internal/mocks"
	"github.com/FreakyGranny/launchpad-api/internal/models"
)

type DonationSuite struct {
	suite.Suite
	mockDonationCtl *gomock.Controller
	mockDonation    *mocks.MockDonationImpl
	recalcChan      chan int
	app             *App
}

func (s *DonationSuite) SetupTest() {
	s.mockDonationCtl = gomock.NewController(s.T())
	s.mockDonation = mocks.NewMockDonationImpl(s.mockDonationCtl)
	s.recalcChan = make(chan int, 1)
	s.app = New(nil, nil, nil, nil, s.mockDonation, nil, nil, "", s.recalcChan)
}

func (s *DonationSuite) TearDownTest() {
	s.mockDonationCtl.Finish()
	close(s.recalcChan)
}

func (s *DonationSuite) TestGetProjectDonations() {
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
	dons, err := s.app.GetProjectDonations(1)
	s.Require().NoError(err)
	s.Require().Equal(2, len(dons))
}

func (s *DonationSuite) TestGetUserDonations() {
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

	dons, err := s.app.GetUserDonations(111)
	s.Require().NoError(err)
	s.Require().Equal(donations, dons)
}

func (s *DonationSuite) TestCreateDonation() {
	donation := &models.Donation{
		Payment:   100,
		ProjectID: 10,
		UserID:    111,
	}
	s.mockDonation.EXPECT().Create(donation).Return(nil)
	newDon, err := s.app.CreateDonation(111, 10, 100)
	s.Require().NoError(err)
	s.Require().Equal(donation, newDon)

	select {
	case x := <-s.recalcChan:
		s.Require().Equal(10, x)
	default:
		s.T().Fail()
	}
}

func (s *DonationSuite) TestSetPayment() {
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

	newDon, err := s.app.UpdateDonation(1, 111, 200, false)
	s.Require().NoError(err)
	s.Require().Equal(donation, newDon)

	select {
	case x := <-s.recalcChan:
		s.Require().Equal(33, x)
	default:
		s.T().Fail()
	}
}

func (s *DonationSuite) TestSetPaymentWrongUser() {
	donation := &models.Donation{
		ID:        1,
		Payment:   100,
		UserID:    111,
		Paid:      false,
		Locked:    false,
		ProjectID: 33,
	}
	s.mockDonation.EXPECT().Get(1).Return(donation, true)

	newDon, err := s.app.UpdateDonation(1, 888, 1000, false)
	s.Require().Error(err)
	s.Require().Nil(newDon)
	s.Require().Equal(ErrDonationModifyNotAllowed, err)
}

func (s *DonationSuite) TestSetPaymentLocked() {
	donation := &models.Donation{
		ID:        1,
		Payment:   100,
		UserID:    111,
		Paid:      false,
		Locked:    true,
		ProjectID: 33,
	}
	s.mockDonation.EXPECT().Get(1).Return(donation, true)

	newDon, err := s.app.UpdateDonation(1, 111, 1000, false)
	s.Require().Error(err)
	s.Require().Nil(newDon)
	s.Require().Equal(ErrDonationModifyWrong, err)
}

func (s *DonationSuite) TestCheckPaid() {
	donation := &models.Donation{
		ID:        1,
		Payment:   100,
		UserID:    111,
		Paid:      false,
		Locked:    true,
		ProjectID: 33,
		Project: models.Project{
			OwnerID: 1212,
		},
	}
	s.mockDonation.EXPECT().Get(1).Return(donation, true)
	donation.Paid = true
	s.mockDonation.EXPECT().Update(donation).Return(nil)

	newDon, err := s.app.UpdateDonation(1, 1212, 0, true)
	s.Require().NoError(err)
	s.Require().Equal(donation, newDon)
}

func (s *DonationSuite) TestCheckPaidNotOwner() {
	donation := &models.Donation{
		ID:        1,
		Payment:   100,
		UserID:    111,
		Paid:      false,
		Locked:    true,
		ProjectID: 33,
		Project: models.Project{
			OwnerID: 1212,
		},
	}
	s.mockDonation.EXPECT().Get(1).Return(donation, true)

	newDon, err := s.app.UpdateDonation(1, 888, 0, true)
	s.Require().Error(err)
	s.Require().Nil(newDon)
	s.Require().Equal(ErrDonationModifyNotAllowed, err)
}

func (s *DonationSuite) TestCheckPaidNotLocked() {
	donation := &models.Donation{
		ID:        1,
		Payment:   100,
		UserID:    111,
		Paid:      false,
		Locked:    false,
		ProjectID: 33,
		Project: models.Project{
			OwnerID: 1212,
		},
	}
	s.mockDonation.EXPECT().Get(1).Return(donation, true)

	newDon, err := s.app.UpdateDonation(1, 1212, 0, true)
	s.Require().Error(err)
	s.Require().Nil(newDon)
	s.Require().Equal(ErrDonationModifyWrong, err)
}

func (s *DonationSuite) TestDeleteDonation() {
	expect := &models.Donation{
		ID:        1,
		UserID:    111,
		ProjectID: 44,
	}
	s.mockDonation.EXPECT().Get(1).Return(expect, true)
	s.mockDonation.EXPECT().Delete(expect).Return(nil)

	err := s.app.DeleteDonation(1, 111)
	s.Require().NoError(err)

	select {
	case x := <-s.recalcChan:
		s.Require().Equal(44, x)
	default:
		s.T().Fail()
	}
}

func TestDonationSuite(t *testing.T) {
	suite.Run(t, new(DonationSuite))
}
