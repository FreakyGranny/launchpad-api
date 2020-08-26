package models

import (
	"errors"
)

// ErrDonationAlreadyExist donation to project for user already exist
var ErrDonationAlreadyExist = errors.New("donation to this project already exist")

// ErrDonationForbidden donation to project is not allowed
var ErrDonationForbidden = errors.New("donation to project is not allowed")

// ErrUserNotFound user not found
var ErrUserNotFound = errors.New("user not found")
