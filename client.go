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
	"io"

	"github.com/Smartling/api-sdk-go/helpers/sm_client"
	"github.com/Smartling/api-sdk-go/helpers/sm_file"
)

type APIClient interface {
	Authenticate() error
	DeleteFile(projectID string, uri string) error
	DownloadFile(projectID string, uri string) (io.Reader, error)
	DownloadTranslation(projectID string, localeID string, request FileDownloadRequest) (io.Reader, error)
	GetFileStatus(projectID string, fileURI string) (*smfile.FileStatus, error)
	GetProjectDetails(projectID string) (*ProjectDetails, error)
	Import(projectID string, localeID string, request smfile.ImportRequest) (*FileImportResult, error)
	ListAllFiles(projectID string, request smfile.FilesListRequest) ([]smfile.File, error)
	ListProjects(accountID string, request smfile.ProjectsListRequest) (*ProjectsList, error)
	RenameFile(projectID string, oldURI string, newURI string) error
	UploadFile(projectID string, request smfile.FileUploadRequest) (*smfile.FileUploadResult, error)
}

// Client represents Smartling API client.
type httpAPIClient struct {
	Client *smclient.Client
}

// NewAPIClient returns new Smartling API client with specified authentication
// data.
func NewAPIClient(userID, tokenSecret string) APIClient {
	smclient := smclient.NewClient(userID, tokenSecret)
	return &httpAPIClient{
		Client: smclient,
	}
}
