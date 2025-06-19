package mt

import (
	smclient "github.com/Smartling/api-sdk-go/helpers/sm_client"
)

// Uploader defines uploader behaviour
type Uploader interface {
	UploadFile() (UploadFileResponse, error)
}

// NewUploader returns new Uploader implementation
func NewUploader(client *smclient.Client) Uploader {
	return httpUploader{base: newBase(client)}
}

type httpUploader struct {
	base *base
}

// UploadFile uploads file
func (h httpUploader) UploadFile() (UploadFileResponse, error) {
	//TODO implement me
	panic("implement me")
}
