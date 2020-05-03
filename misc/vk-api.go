package misc

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"errors"

	"github.com/FreakyGranny/launchpad-api/config"
)

const apiVersion = "5.103"
const userFields = "photo_200,screen_name"

var vk *VkClient

// AuthData VK auth data response
type AuthData struct {
	AccessToken string `json:"access_token"`
	Expires     uint   `json:"expires_in"`
	UserID      uint   `json:"user_id"`
	Email       string `json:"email"`
}

// UsersResponse users response
type UsersResponse struct {
	Response  []VkUserData `json:"response"`
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
	HTTPClient *http.Client
	AppID string
	AppSecret string
	RedirectURI string
}

func (vk *VkClient) buildAuthRequest(code string) (*http.Request, error) {
	req, err := http.NewRequest("GET", "https://oauth.vk.com/access_token", nil)
    if err != nil {
        return nil, err
    }

	q := req.URL.Query()
	q.Add("client_id", vk.AppID)
	q.Add("client_secret", vk.AppSecret)
	q.Add("redirect_uri", vk.RedirectURI)
    q.Add("code", code)

	req.URL.RawQuery = q.Encode()

	return req, nil
}

func (vk *VkClient) buildUserDataRequest(userID uint, token string) (*http.Request, error) {
	req, err := http.NewRequest("GET", "https://api.vk.com/method/users.get", nil)
    if err != nil {
        return nil, err
    }

	q := req.URL.Query()
	q.Add("fields", userFields)
	q.Add("v", apiVersion)
	q.Add("user_id", fmt.Sprintf("%v", userID))
    q.Add("access_token", token)

	req.URL.RawQuery = q.Encode()
	
	return req, nil
}

// GetAccessToken Get VK access token
func (vk *VkClient) GetAccessToken(code string) (*AuthData, error) {
	req, err := vk.buildAuthRequest(code)
    if err != nil {
        return nil, err
	}
	
	resp, err := vk.HTTPClient.Do(req)
    if err != nil {
        return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("something went wrong")
	}

    defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)

	response := AuthData{}
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return nil, err
	}
	
	return &response, nil
}


// GetUserData Get VK access token
func (vk *VkClient) GetUserData(userID uint, token string) (*VkUserData, error) {
	req, err := vk.buildUserDataRequest(userID, token)
    if err != nil {
        return nil, err
	}
	resp, err := vk.HTTPClient.Do(req)
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
	
	return &response.Response[0], nil
}

//VkInit initialize VK client 
func VkInit (cfg config.VkAuth) {
	vk = &VkClient{
		HTTPClient: &http.Client{},
		AppID: cfg.AppID,
		AppSecret: cfg.ClientSecret,
		RedirectURI: cfg.RedirectURI,
	}
}

// GetVkClient returns vk client
func GetVkClient() *VkClient {
	return vk
}
