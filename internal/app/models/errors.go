package models

import (
	"errors"
)

// ErrDonationAlreadyExist donation to project for user already exist
var ErrDonationAlreadyExist = errors.New("donation to this project already exist")

// ErrDonationForbidden donation to project is not allowed
var ErrDonationForbidden = errors.New("donation to project is not allowed")

// ErrDonationModifyForbidden donation editing is not allowed
var ErrDonationModifyForbidden = errors.New("donation editing is not allowed")

// ErrUserNotFound user not found
var ErrUserNotFound = errors.New("user not found")

// ErrProjectModifyForbidden donation editing is not allowed
var ErrProjectModifyForbidden = errors.New("project editing is not allowed")
