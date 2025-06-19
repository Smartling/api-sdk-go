package mt

import (
	"context"
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
func (h httpUploader) UploadFile(ctx context.Context, fileType FileType) (UploadFileResponse, error) {
	//TODO implement me
	panic("implement me")
	url := "https://api.smartling.com/file-translations-api/v2/accounts/$smartlingAccountId/files"
}
