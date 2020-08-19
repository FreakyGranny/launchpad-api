package auth

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/FreakyGranny/launchpad-api/internal/app/config"
)

const apiVersion = "5.103"
const userFields = "photo_200,screen_name"

//go:generate mockgen -source=$GOFILE -destination=./vk_http_mock.go -package=auth . HTTPClient

// HTTPClient ...
type HTTPClient interface {
	// Do send http request
	Do(req *http.Request) (*http.Response, error)
}

// VKAuthData VK auth data response
type VKAuthData struct {
	AccessToken string `json:"access_token"`
	Expires     uint   `json:"expires_in"`
	UserID      int    `json:"user_id"`
	Email       string `json:"email"`
}

// UsersResponse users response
type UsersResponse struct {
	Response []VkUserData `json:"response"`
}

// VkUserData VK user data response
type VkUserData struct {
	Username  string `json:"screen_name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Avatar    string `json:"photo_200"`
}

//VkClient vk api client
type VkClient struct {
	Client      HTTPClient
	AppID       string
	AppSecret   string
	RedirectURI string
}

func (vk *VkClient) buildAuthRequest(code string) *http.Request {
	req, _ := http.NewRequest("GET", "https://oauth.vk.com/access_token", nil)

	q := req.URL.Query()
	q.Add("client_id", vk.AppID)
	q.Add("client_secret", vk.AppSecret)
	q.Add("redirect_uri", vk.RedirectURI)
	q.Add("code", code)

	req.URL.RawQuery = q.Encode()

	return req
}

func (vk *VkClient) buildUserDataRequest(userID int, token string) *http.Request {
	req, _ := http.NewRequest("GET", "https://api.vk.com/method/users.get", nil)

	q := req.URL.Query()
	q.Add("fields", userFields)
	q.Add("v", apiVersion)
	q.Add("user_id", strconv.Itoa(userID))
	q.Add("access_token", token)

	req.URL.RawQuery = q.Encode()

	return req
}

// GetAccessToken Get VK access token
func (vk *VkClient) GetAccessToken(code string) (*AccessData, error) {
	req := vk.buildAuthRequest(code)

	resp, err := vk.Client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("something went wrong")
	}

	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)

	response := VKAuthData{}
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return nil, err
	}

	return &AccessData{
		UserID:      response.UserID,
		AccessToken: response.AccessToken,
		Expires:     response.Expires,
		Email:       response.Email,
	}, nil
}

// GetUserData Get VK access token
func (vk *VkClient) GetUserData(userID int, token string) (*UserData, error) {
	req := vk.buildUserDataRequest(userID, token)
	resp, err := vk.Client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("something went wrong")
	}

	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)

	response := UsersResponse{}
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return nil, err
	}
	if len(response.Response) == 0 {
		return nil, errors.New("No valid user")
	}

	data := response.Response[0]

	return &UserData{
		Username:  data.Username,
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Avatar:    data.Avatar,
	}, nil
}

//NewVk initialize VK client
func NewVk(cfg config.VkAuth) *VkClient {
	return &VkClient{
		Client:  &http.Client{},
		AppID:       cfg.AppID,
		AppSecret:   cfg.ClientSecret,
		RedirectURI: cfg.RedirectURI,
	}
}
