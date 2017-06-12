// Package smartling is a client implementation of the Smartling Translation API v2 as documented at
// http://docs.smartling.com/pages/API/v2/
package smartling

import (
	"net/http"
	"time"
)


// API base url
const apiBaseUrl  = "https://api.smartling.com"


// Smartling SDK clinet
type Client struct {
	baseUrl         string
	auth            *Auth
	httpClient      *http.Client
}


// Smartling client initialization
func NewClient(userIdentifier string, tokenSecret string) *Client {
	return &Client{
		baseUrl:    apiBaseUrl,
		auth:       &Auth { userIdentifier, tokenSecret, Token{}, Token{} },
		httpClient: &http.Client{ Timeout: 60 * time.Second },
	}
}


// custom initialization with overrided base url
func NewClientWithBaseUrl(userIdentifier string, tokenSecret string, baseUrl string) *Client {
	client := NewClient(userIdentifier, tokenSecret)
	client.baseUrl = baseUrl
	return client
}

// attempts to authenticate with smartling
// not a required call
func (c *Client) AuthenticationTest() error {
	return c.auth.doAuthCall(c, false)
}

func (c *Client) SetHttpTimeout(t time.Duration) {
	c.httpClient.Timeout = t
}
