package mt

import (
	"fmt"

	"github.com/Smartling/api-sdk-go/helpers/sm_client"
	"github.com/Smartling/api-sdk-go/helpers/sm_file"
)

// Uploader defines uploader behaviour
type Uploader interface {
	UploadFile(accountUID AccountUID) (UploadFileResponse, error)
}

// NewUploader returns new Uploader implementation
func NewUploader(client *smclient.Client) Uploader {
	return httpUploader{base: newBase(client)}
}

type httpUploader struct {
	base *base
}

// UploadFile uploads file
func (u httpUploader) UploadFile(accountUID AccountUID) (UploadFileResponse, error) {
	filePath := buildUploadFilePath(accountUID)
	path := joinPath(mtBasePath, filePath)
	_, err := u.base.client.UploadFile(
		path,
		smfile.FileUploadRequest{},
	)
	if err != nil {
		return UploadFileResponse{}, fmt.Errorf("failed to download original file: %w", err)
	}

	// TODO set file id
	return UploadFileResponse{
		Code:    successResponseCode,
		FileUID: "",
	}, nil
}

func buildUploadFilePath(accountUID AccountUID) string {
	return "/accounts/" + string(accountUID) + "/files"
}
