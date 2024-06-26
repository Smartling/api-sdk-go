// Copyright (c) 2020 Smartling
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// Package smartling is a client implementation of the Smartling Translation
// API v2 as documented at https://help.smartling.com/v1.0/reference
package smartling

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type ClientInterface interface {
	Authenticate() error
	DeleteFile(projectID string, uri string) error
	DownloadFile(projectID string, uri string) (io.Reader, error)
	DownloadTranslation(projectID string, localeID string, request FileDownloadRequest) (io.Reader, error)
	GetFileStatus(projectID string, fileURI string) (*FileStatus, error)
	Get(url string, params url.Values, options ...interface{}) (io.ReadCloser, int, error)
	GetJSON(url string, params url.Values, result interface{}, options ...interface{}) (json.RawMessage, int, error)
	GetProjectDetails(projectID string) (*ProjectDetails, error)
	Import(projectID string, localeID string, request ImportRequest) (*FileImportResult, error)
	LastModified(projectID string, request FileLastModifiedRequest) (*FileLastModifiedLocales, error)
	ListAllFiles(projectID string, request FilesListRequest) ([]File, error)
	ListFiles(projectID string, request FilesListRequest) (*FilesList, error)
	ListFileTypes(projectID string) ([]FileType, error)
	ListProjects(accountID string, request ProjectsListRequest) (*ProjectsList, error)
	Post(url string, payload []byte, result interface{}, options ...interface{}) (json.RawMessage, int, error)
	RenameFile(projectID string, oldURI string, newURI string) error
	UploadFile(projectID string, request FileUploadRequest) (*FileUploadResult, error)
}

type (
	// LogFunction represents abstract logger function interface which
	// can be used for setting up logging of library actions.
	LogFunction func(format string, args ...interface{})
)

var (
	// Version is a API SDK version, sent in User-Agent header.
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
func NewClient(userID string, tokenSecret string) *Client {
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
func (client *Client) SetInfoLogger(logger LogFunction) *Client {
	client.Logger.Infof = logger

	return client
}

// SetDebugLogger sets logger function which will be called for logging
// internal information like HTTP requests and their responses.
func (client *Client) SetDebugLogger(logger LogFunction) *Client {
	client.Logger.Debugf = logger

	return client
}
