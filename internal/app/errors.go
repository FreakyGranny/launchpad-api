package app

import "errors"

var (
	// ErrGetAccessTokenFailed error geting access token from provider.
	ErrGetAccessTokenFailed = errors.New("unable to get access token")
	// ErrGetUserDataFailed error geting user data from provider.
	ErrGetUserDataFailed = errors.New("unable to get user data")
)

var (
	// ErrUserNotFound user with given id not found.
	ErrUserNotFound = errors.New("user not found")
	// ErrGetUserParticipation error getting user participation
	ErrGetUserParticipation = errors.New("cant get user participation")
)
