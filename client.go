// Package smartling is a client implementation of the Smartling Translation
// API v2 as documented at https://help.smartling.com/v1.0/reference
package smartling

import (
	"net/http"
	"time"
)

var (
	// DefaultBaseURL specifies base URL which will be used for calls unless
	// other is specified in the Client struct.
	DefaultBaseURL = "https://api.smartling.com"

	// DefaultHTTPClient specifies default HTTP client which will be used
	// for calls unless other is specified in the Client struct.
	DefaultHTTPClient = http.Client{Timeout: 60 * time.Second}
)

// Client represents Smartling API client.
type Client struct {
	BaseURL     string
	Credentials *Credentials

	HTTP *http.Client
}

// NewClient returns new Smartling API client with specified authentication
// data.
func NewClient(userID string, tokenSecret string) *Client {
	return &Client{
		BaseURL: DefaultBaseURL,
		Credentials: &Credentials{
			UserID: userID,
			Secret: tokenSecret,
		},

		HTTP: &DefaultHTTPClient,
	}
}
