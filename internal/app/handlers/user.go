package handlers

import (
	"net/http"
	"strconv"

	"github.com/FreakyGranny/launchpad-api/internal/app/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type extendedUser struct {
	models.User
	Participation []models.Participation `json:"participation"`
}

// UserHandler ...
type UserHandler struct {
	UserModel models.UserImpl
}

// NewUserHandler ...
func NewUserHandler(u models.UserImpl) *UserHandler {
	return &UserHandler{UserModel: u}
}

// GetCurrentUser godoc
// @Summary Show a current user
// @Description Returns user by ID from token
// @Tags user
// @ID get-user-by-token
// @Produce json
// @Success 200 {object} extendedUser
// @Security Bearer
// @Router /user [get]
func (h *UserHandler) GetCurrentUser(c echo.Context) error {
	userID, err := getUserIDFromToken(c.Get("user"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	user, ok := h.UserModel.FindByID(userID)
	if !ok {
		return c.JSON(http.StatusNotFound, nil)
	}

	pts, err := h.UserModel.GetParticipation(userID)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, errorResponse("cant get participation"))
	}

	return c.JSON(http.StatusOK, extendedUser{
		User:          *user,
		Participation: pts,
	})
}

// GetUser godoc
// @Summary Show a specific user
// @Description Returns user by ID
// @Tags user
// @ID get-user-by-id
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} extendedUser
// @Security Bearer
// @Router /user/{id} [get]
func (h *UserHandler) GetUser(c echo.Context) error {
	intID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("wrong ID"))
	}
	user, ok := h.UserModel.FindByID(intID)
	if !ok {
		return c.JSON(http.StatusNotFound, nil)
	}
	pts, err := h.UserModel.GetParticipation(intID)
	if err != nil {
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, errorResponse("cant get participation"))
	}

	return c.JSON(http.StatusOK, extendedUser{
		User:          *user,
		Participation: pts,
	})
}
