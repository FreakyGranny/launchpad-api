package auth

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type VkProviderSuite struct {
	suite.Suite
	mockCtl  *gomock.Controller
	mockHTTP *MockHTTPClient
	vkClient *VkClient
}

func (s *VkProviderSuite) SetupTest() {
	s.mockCtl = gomock.NewController(s.T())
	s.mockHTTP = NewMockHTTPClient(s.mockCtl)
	s.vkClient = &VkClient{
		Client:      s.mockHTTP,
		AppID:       "AppID",
		AppSecret:   "ClientSecret",
		RedirectURI: "RedirectURI",
	}
}

func prepareTokenRequest() (*http.Request, error) {
	req, err := http.NewRequest("GET", "https://oauth.vk.com/access_token", nil)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = "client_id=AppID&client_secret=ClientSecret&code=magic_code&redirect_uri=RedirectURI"

	return req, nil
}

func prepareDataRequest() (*http.Request, error) {
	req, err := http.NewRequest("GET", "https://api.vk.com/method/users.get", nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("fields", userFields)
	q.Add("v", apiVersion)
	q.Add("user_id", "123")
	q.Add("access_token", "secret_token")

	req.URL.RawQuery = q.Encode()

	return req, nil
}

func (s *VkProviderSuite) TearDownTest() {
	s.mockCtl.Finish()
}

func (s *VkProviderSuite) TestGetTokenSuccess() {
	req, err := prepareTokenRequest()
	if err != nil {
		s.T().Fail()
	}
	json := `{"access_token":"Toooken!","expires_in":1980,"user_id":13,"email":"test@gmail.com"}`
	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))
	expected := &http.Response{
		StatusCode: 200,
		Body:       r,
	}

	s.mockHTTP.EXPECT().Do(req).Return(expected, nil)

	data, err := s.vkClient.GetAccessToken("magic_code")

	s.Require().NoError(err)
	s.Require().Equal("Toooken!", data.AccessToken)
	s.Require().Equal(uint(1980), data.Expires)
	s.Require().Equal(13, data.UserID)
	s.Require().Equal("test@gmail.com", data.Email)
}

func (s *VkProviderSuite) TestGetTokenError() {
	req, err := prepareTokenRequest()
	if err != nil {
		s.T().Fail()
	}
	s.mockHTTP.EXPECT().Do(req).Return(nil, errors.New("unexpected error"))

	data, err := s.vkClient.GetAccessToken("magic_code")

	s.Require().Error(err)
	s.Require().Nil(data)
}

func (s *VkProviderSuite) TestGetTokenWrongCode() {
	req, err := prepareTokenRequest()
	if err != nil {
		s.T().Fail()
	}
	expected := &http.Response{
		StatusCode: 404,
		Body:       nil,
	}
	s.mockHTTP.EXPECT().Do(req).Return(expected, nil)

	data, err := s.vkClient.GetAccessToken("magic_code")

	s.Require().Error(err)
	s.Require().Nil(data)
}

func (s *VkProviderSuite) TestGetTokenNotJSON() {
	req, err := prepareTokenRequest()
	if err != nil {
		s.T().Fail()
	}
	r := ioutil.NopCloser(bytes.NewReader([]byte("This is not json")))
	expected := &http.Response{
		StatusCode: 200,
		Body:       r,
	}

	s.mockHTTP.EXPECT().Do(req).Return(expected, nil)

	data, err := s.vkClient.GetAccessToken("magic_code")

	s.Require().Error(err)
	s.Require().Nil(data)
}

func (s *VkProviderSuite) TestGetDataSuccess() {
	req, err := prepareDataRequest()
	if err != nil {
		s.T().Fail()
	}

	json := `{"response":[{"screen_name":"Johnny86","first_name":"John","last_name":"Doe","photo_200":"https://avatar.com/john"}]}`
	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))
	expected := &http.Response{
		StatusCode: 200,
		Body:       r,
	}

	s.mockHTTP.EXPECT().Do(req).Return(expected, nil)

	data, err := s.vkClient.GetUserData(123, "secret_token")

	s.Require().NoError(err)
	s.Require().Equal("Johnny86", data.Username)
	s.Require().Equal("John", data.FirstName)
	s.Require().Equal("Doe", data.LastName)
	s.Require().Equal("https://avatar.com/john", data.Avatar)
}

func (s *VkProviderSuite) TestGetDataError() {
	req, err := prepareDataRequest()
	if err != nil {
		s.T().Fail()
	}
	s.mockHTTP.EXPECT().Do(req).Return(nil, errors.New("unexpected error"))

	data, err := s.vkClient.GetUserData(123, "secret_token")

	s.Require().Error(err)
	s.Require().Nil(data)
}

func (s *VkProviderSuite) TestGetDataWrongCode() {
	req, err := prepareDataRequest()
	if err != nil {
		s.T().Fail()
	}
	expected := &http.Response{
		StatusCode: 404,
		Body:       nil,
	}
	s.mockHTTP.EXPECT().Do(req).Return(expected, nil)

	data, err := s.vkClient.GetUserData(123, "secret_token")

	s.Require().Error(err)
	s.Require().Nil(data)
}

func (s *VkProviderSuite) TestGetDataNotJSON() {
	req, err := prepareDataRequest()
	if err != nil {
		s.T().Fail()
	}
	r := ioutil.NopCloser(bytes.NewReader([]byte("This is not json")))
	expected := &http.Response{
		StatusCode: 200,
		Body:       r,
	}
	s.mockHTTP.EXPECT().Do(req).Return(expected, nil)

	data, err := s.vkClient.GetUserData(123, "secret_token")

	s.Require().Error(err)
	s.Require().Nil(data)
}

func (s *VkProviderSuite) TestGetDataNoUsersAtResponse() {
	req, err := prepareDataRequest()
	if err != nil {
		s.T().Fail()
	}

	json := `{"response":[]}`
	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))
	expected := &http.Response{
		StatusCode: 200,
		Body:       r,
	}

	s.mockHTTP.EXPECT().Do(req).Return(expected, nil)

	data, err := s.vkClient.GetUserData(123, "secret_token")

	s.Require().Error(err)
	s.Require().Nil(data)
}

func TestProviderSuite(t *testing.T) {
	suite.Run(t, new(VkProviderSuite))
}
