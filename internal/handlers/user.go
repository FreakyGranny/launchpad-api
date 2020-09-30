package handlers

import (
	"net/http"
	"strconv"

	"github.com/FreakyGranny/launchpad-api/internal/app"
	"github.com/labstack/echo/v4"
)

// UserHandler ...
type UserHandler struct {
	app app.Application
}

// NewUserHandler ...
func NewUserHandler(a app.Application) *UserHandler {
	return &UserHandler{app: a}
}

// GetCurrentUser godoc
// @Summary Show a current user
// @Description Returns user by ID from token
// @Tags user
// @ID get-user-by-token
// @Produce json
// @Success 200 {object} app.ExtendedUser
// @Security Bearer
// @Router /user [get]
func (h *UserHandler) GetCurrentUser(c echo.Context) error {
	userID, err := getUserIDFromToken(c.Get("user"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	user, err := h.app.GetUser(userID)
	switch err {
	case app.ErrUserNotFound:
		return c.JSON(http.StatusNotFound, nil)
	case nil:
		return c.JSON(http.StatusOK, user)
	default:
		return c.JSON(http.StatusInternalServerError, errorResponse("unexpected error"))
	}
}

// GetUser godoc
// @Summary Show a specific user
// @Description Returns user by ID
// @Tags user
// @ID get-user-by-id
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} app.ExtendedUser
// @Security Bearer
// @Router /user/{id} [get]
func (h *UserHandler) GetUser(c echo.Context) error {
	intID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse("wrong ID"))
	}
	user, err := h.app.GetUser(intID)
	switch err {
	case app.ErrUserNotFound:
		return c.JSON(http.StatusNotFound, nil)
	case nil:
		return c.JSON(http.StatusOK, user)
	default:
		return c.JSON(http.StatusInternalServerError, errorResponse("unexpected error"))
	}
}
