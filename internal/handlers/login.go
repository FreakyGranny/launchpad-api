package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"

	"github.com/FreakyGranny/launchpad-api/internal/app"
)

// TokenRequest - request for auth token
type TokenRequest struct {
	Code string `json:"code"`
}

// TokenResponse - response auth token
type TokenResponse struct {
	Token string `json:"token"`
}

// AuthHandler ...
type AuthHandler struct {
	app app.Application
}

// NewAuthHandler ...
func NewAuthHandler(a app.Application) *AuthHandler {
	return &AuthHandler{app: a}
}

// Login godoc
// @Summary Returns access token
// @Description get token for user
// @Tags auth
// @ID get-token
// @Accept json
// @Produce json
// @Param request body TokenRequest true "Request body"
// @Success 200 {object} TokenResponse
// @Router /login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	request := new(TokenRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	token, err := h.app.Authentificate(request.Code)
	switch err {
	case app.ErrGetAccessTokenFailed:
		log.Error(err)
		return c.JSON(http.StatusUnauthorized, errorResponse("unable to authentificate"))
	case app.ErrGetUserDataFailed:
		log.Error(err)
		return c.JSON(http.StatusUnauthorized, errorResponse("unable to authentificate"))
	case nil:
		return c.JSON(http.StatusOK, TokenResponse{Token: token})
	default:
		log.Error(err)
		return c.JSON(http.StatusInternalServerError, errorResponse("unexpected error"))
	}
}
