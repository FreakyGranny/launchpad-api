package auth

import (
	"time"

	"github.com/FreakyGranny/launchpad-api/internal/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/jonboulle/clockwork"
)

// CreateToken creates token for user.
func CreateToken(clock clockwork.Clock, secret string, expires uint, user *models.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["admin"] = user.IsAdmin
	claims["exp"] = clock.Now().Add(time.Second * time.Duration(expires)).Unix()

	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return t, nil
}
