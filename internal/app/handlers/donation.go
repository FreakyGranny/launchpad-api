package handlers

import (
	"net/http"
	"strconv"

	"github.com/FreakyGranny/launchpad-api/internal/app/models"
	"github.com/labstack/echo/v4"
)

// DonationHandler ...
type DonationHandler struct {
	DonationModel models.DonationImpl
}

// NewDonationHandler ...
func NewDonationHandler(m models.DonationImpl) *DonationHandler {
	return &DonationHandler{DonationModel: m}
}

// DonationCreateRequest ...
type DonationCreateRequest struct {
	ProjectID int `json:"project"`
	Payment   int `json:"payment"`
}

// DonationUpdateRequest ...
type DonationUpdateRequest struct {
	Paid    bool `json:"paid,omitempty"`
	Payment int  `json:"payment,omitempty"`
}

// ProjectDonation for project donations response
type ProjectDonation struct {
	ID     int         `json:"id"`
	User   models.User `json:"user"`
	Locked bool        `json:"locked"`
	Paid   bool        `json:"paid"`
}

// GetUserDonations godoc
// @Summary Returns list of user's donations
// @Description Returns list of user's donations
// @Tags donation
// @ID get-user-donations
// @Produce json
// @Success 200 {object} []models.Donation
// @Security Bearer
// @Router /donation [get]
func (h *DonationHandler) GetUserDonations(c echo.Context) error {
	userID, err := getUserIDFromToken(c.Get("user"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	donations, err := h.DonationModel.GetAllByUser(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, donations)
}

// GetProjectDonations godoc
// @Summary Returns list of project donations
// @Description Returns list of project donations
// @Tags donation
// @ID get-project-donations
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {object} []ProjectDonation
// @Security Bearer
// @Router /donation/project/{id} [get]
func (h *DonationHandler) GetProjectDonations(c echo.Context) error {
	intID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("wrong ID"))
	}
	donations, err := h.DonationModel.GetAllByProject(intID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	projectDonations := make([]ProjectDonation, 0, len(donations))

	for _, donation := range donations {
		projectDonations = append(projectDonations, ProjectDonation{
			ID:     donation.ID,
			User:   donation.User,
			Locked: donation.Locked,
			Paid:   donation.Paid,
		})
	}
	return c.JSON(http.StatusOK, projectDonations)
}

// CreateDonation godoc
// @Summary Create donation
// @Description Create new donation
// @Tags donation
// @ID post-donation
// @Accept json
// @Produce json
// @Param request body DonationCreateRequest true "Request body"
// @Success 200 {object} models.Donation
// @Security Bearer
// @Router /donation [post]
func (h *DonationHandler) CreateDonation(c echo.Context) error {
	request := new(DonationCreateRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	userID, err := getUserIDFromToken(c.Get("user"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	donation := &models.Donation{
		UserID:    userID,
		ProjectID: request.ProjectID,
		Payment:   request.Payment,
	}
	err = h.DonationModel.Create(donation)

	switch err {
	case nil:
		return c.JSON(http.StatusCreated, donation)
	case models.ErrDonationAlreadyExist:
		return c.JSON(http.StatusForbidden, err)
	case models.ErrDonationForbidden:
		return c.JSON(http.StatusForbidden, err)
	case models.ErrUserNotFound:
		return c.JSON(http.StatusBadRequest, err)
	default:
		return c.JSON(http.StatusInternalServerError, err)
	}
}

// DeleteDonation godoc
// @Summary Delete not locked donation
// @Description Delete not locked donation
// @Tags donation
// @ID delete-donation
// @Param id path int true "Donation ID"
// @Success 204
// @Security Bearer
// @Router /donation/{id} [delete]
func (h *DonationHandler) DeleteDonation(c echo.Context) error {
	donationID, _ := strconv.Atoi(c.Param("id"))
	userID, err := getUserIDFromToken(c.Get("user"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("wrong ID"))
	}
	err = h.DonationModel.Delete(donationID, userID)
	switch err {
	case nil:
		return c.NoContent(http.StatusNoContent)
	case models.ErrDonationModifyForbidden:
		return c.JSON(http.StatusForbidden, err)
	default:
		return c.JSON(http.StatusInternalServerError, err)
	}
}

// UpdateDonation godoc
// @Summary Update not locked donation
// @Description Update not locked donation
// @Tags donation
// @ID update-donation
// @Accept json
// @Produce json
// @Param request body DonationUpdateRequest true "Request body"
// @Param id path int true "Donation ID"
// @Success 200 {object} models.Donation
// @Security Bearer
// @Router /donation/{id} [patch]
func (h *DonationHandler) UpdateDonation(c echo.Context) error {
	request := new(DonationUpdateRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	userID, err := getUserIDFromToken(c.Get("user"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	donationID, _ := strconv.Atoi(c.Param("id"))

	donation := &models.Donation{
		ID:      donationID,
		UserID:  userID,
		Paid:    request.Paid,
		Payment: request.Payment,
	}
	err = h.DonationModel.Update(donation)
	switch err {
	case nil:
		return c.JSON(http.StatusOK, donation)
	case models.ErrDonationModifyForbidden:
		return c.JSON(http.StatusForbidden, err)
	default:
		return c.JSON(http.StatusInternalServerError, err)
	}
}
