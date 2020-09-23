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

var (
	// ErrProjectRetrieve error getting projects.
	ErrProjectRetrieve = errors.New("unable to get projects")
	// ErrProjectNotFound project with given id not found.
	ErrProjectNotFound = errors.New("project not found")
	// ErrProjectModifyNotAllowed project modifying not allowed.
	ErrProjectModifyNotAllowed = errors.New("modifying forbidden")
)

var (
	// ErrDonationNotFound donation with given id not found.
	ErrDonationNotFound = errors.New("donation not found")
	// ErrDonationModifyNotAllowed donation locked.
	ErrDonationModifyNotAllowed = errors.New("modifying forbidden")
	// ErrDonationModifyWrong modifying params are wrong.
	ErrDonationModifyWrong = errors.New("wrong modifying params")

)

var (
	// ErrNoStrategy no mathed strategy for project type.
	ErrNoStrategy = errors.New("no matched strategy")
)
