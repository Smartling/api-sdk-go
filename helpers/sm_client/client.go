package smclient

import "net/http"

type (
	// LogFunction represents abstract logger function interface which
	// can be used for setting up logging of library actions.
	LogFunction func(format string, args ...interface{})
)

var (
	// Version is an API SDK version, sent in User-Agent header.
	Version = "1.1"

	// DefaultBaseURL specifies base URL which will be used for calls unless
	// other is specified in the Client struct.
	DefaultBaseURL = "https://api.smartling.com"

	// DefaultHTTPClient specifies default HTTP client which will be used
	// for calls unless other is specified in the Client struct.
	DefaultHTTPClient = http.Client{}

	// DefaultUserAgent is a string that will be sent in User-Agent header.
	DefaultUserAgent = "smartling-api-sdk-go"
)

// Client represents Smartling API client.
type Client struct {
	BaseURL     string
	Credentials *Credentials

	HTTP *http.Client

	Logger struct {
		Infof  LogFunction
		Debugf LogFunction
	}

	UserAgent string
}

// NewClient returns new Smartling API client with specified authentication
// data.
func NewClient(userID, tokenSecret string) *Client {
	return &Client{
		BaseURL: DefaultBaseURL,
		Credentials: &Credentials{
			UserID: userID,
			Secret: tokenSecret,
		},

		HTTP: &DefaultHTTPClient,

		Logger: struct {
			Infof  LogFunction
			Debugf LogFunction
		}{
			Infof:  func(string, ...interface{}) {},
			Debugf: func(string, ...interface{}) {},
		},

		UserAgent: DefaultUserAgent + "/" + Version,
	}
}

// SetInfoLogger sets logger function which will be called for logging
// informational messages like progress of file download and so on.
func (c *Client) SetInfoLogger(logger LogFunction) *Client {
	c.Logger.Infof = logger

	return c
}

// SetDebugLogger sets logger function which will be called for logging
// internal information like HTTP requests and their responses.
func (c *Client) SetDebugLogger(logger LogFunction) *Client {
	c.Logger.Debugf = logger

	return c
}
