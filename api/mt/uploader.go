package mt

import (
	smclient "github.com/Smartling/api-sdk-go/helpers/sm_client"
)

type Uploader interface {
	UploadFile() error
}

func NewUploader(client *smclient.Client) Uploader {
	return httpUploader{client: client}
}

type httpUploader struct {
	client *smclient.Client
}

func (h httpUploader) UploadFile() error {
	//TODO implement me
	panic("implement me")
}
