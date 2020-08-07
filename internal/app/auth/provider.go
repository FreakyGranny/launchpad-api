package auth

// AccessData VK auth data response
type AccessData struct {
	AccessToken string
	Expires     uint
	UserID      int
	Email       string
}

// UserData VK user data response
type UserData struct {
	Username  string
	FirstName string
	LastName  string
	Avatar    string
}

//go:generate mockgen -destination=../mocks/auth_provider_mock.go -package=mocks . Provider

// Provider ...
type Provider interface {
	// GetAccessToken access token
	GetAccessToken(code string) (*AccessData, error)

	// GetUserData Get user data
	GetUserData(userID int, token string) (*UserData, error)
}
