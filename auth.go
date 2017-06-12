// class responsible for maintaining user auth valid

package smartling

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// API endpoints
const (
	authApiAuth    = "/auth-api/v2/authenticate"
	authApiRefresh = "/auth-api/v2/authenticate/refresh"
)

const tokenExpirationSafetyDuration = time.Duration(30) * time.Second

// auth api call response data
type AuthApiResponse struct {
	AccessToken      string
	ExpiresIn        int32
	RefreshExpiresIn int32
	RefreshToken     string
}

type Token struct {
	Token          string
	ExpirationDate time.Time
}

func (t *Token) IsValid() bool {
	if len(t.Token) == 0 {
		return false
	}

	return time.Now().Before(t.ExpirationDate)
}

// true if token will last more then tokenExpirationSafetyDuration
func (t *Token) IsSafe() bool {
	if !t.IsValid() {
		return false
	}

	return time.Now().Add(tokenExpirationSafetyDuration).Before(t.ExpirationDate)
}

type Auth struct {
	userIdentifier string
	tokenSecret    string
	accessToken    Token
	reauthToken    Token
}

// get access token for request
func (a *Auth) AccessHeader(c *Client) (string, error) {
	// check if auth token is still valid
	if a.accessToken.IsSafe() {
		return fmt.Sprintf("Bearer %v", a.accessToken.Token), nil
	}

	// if reauth token is not valid as well - we need to relogin
	if !a.reauthToken.IsValid() {
		// do reauth call
		err := a.doAuthCall(c, false)
		if err != nil {
			return "", err
		}
	} else {
		// issue reauth call to refresh access token
		err := a.doAuthCall(c, true)
		if err != nil {
			return "", err
		}
	}

	return a.AccessHeader(c)
}

func (a *Auth) authData() ([]byte, error) {

	// validate input
	if len(a.userIdentifier) == 0 || len(a.tokenSecret) == 0 {
		return nil, fmt.Errorf("User credentials are not set")
	}

	// prepare the auth map
	authInfo := make(map[string]string)
	authInfo["userIdentifier"] = a.userIdentifier
	authInfo["userSecret"] = a.tokenSecret

	// marshall into bytes
	return json.Marshal(&authInfo)
}

func (a *Auth) reauthData() ([]byte, error) {
	// validate input
	if !a.reauthToken.IsValid() {
		return nil, fmt.Errorf("Reauth token is invalid")
	}

	// prepare the map
	authInfo := make(map[string]string)
	authInfo["refreshToken"] = a.reauthToken.Token

	// marshall into bytes
	return json.Marshal(&authInfo)
}

// actually performs auth call
func (a *Auth) doAuthCall(c *Client, isReauth bool) error {

	var authBytes []byte
	var err error = nil
	var apiCall string = ""
	if !isReauth {
		apiCall = authApiAuth
		authBytes, err = a.authData()
		if err != nil {
			return err
		}
	} else {
		log.Printf("REFRESH CALL")
		apiCall = authApiRefresh
		authBytes, err = a.reauthData()
		if err != nil {
			return err
		}
	}

	// use empty auth header
	bytes, statusCode, err := c.doPostRequest(c.baseUrl+apiCall, "", authBytes)
	if err != nil {
		return err
	}

	if statusCode != 200 {
		return fmt.Errorf("Auth call returned unexpected status code: %v", statusCode)
	}

	// unmarshal transport header
	apiResponse := SmartlingApiResponse{}
	err = json.Unmarshal(bytes, &apiResponse)
	if err != nil {
		return err
	}

	// check status
	if apiResponse.Response.Code != "SUCCESS" {
		return fmt.Errorf("Auth call returned unexpected response code: %v", apiResponse.Response.Code)
	}
	log.Printf("auth status %v", statusCode)

	// unmarshal auth body
	authResponse := AuthApiResponse{}
	err = json.Unmarshal(apiResponse.Response.Data, &authResponse)
	if err != nil {
		return err
	}

	// fill tokens
	a.accessToken.Token = authResponse.AccessToken
	a.accessToken.ExpirationDate = time.Now().Add(time.Duration(authResponse.ExpiresIn) * time.Second)

	a.reauthToken.Token = authResponse.RefreshToken
	a.reauthToken.ExpirationDate = time.Now().Add(time.Duration(authResponse.RefreshExpiresIn) * time.Second)

	log.Printf("%#v", authResponse)

	return nil
}
