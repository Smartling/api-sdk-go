package mt

import (
	"fmt"

	"github.com/Smartling/api-sdk-go/helpers/sm_client"
	"github.com/Smartling/api-sdk-go/helpers/sm_file"
)

// Uploader defines uploader behaviour
type Uploader interface {
	UploadFile(accountUID AccountUID, projectID string, req smfile.FileUploadRequest) (UploadFileResponse, error)
}

// NewUploader returns new Uploader implementation
func NewUploader(client *smclient.Client) Uploader {
	return httpUploader{base: newBase(client)}
}

type httpUploader struct {
	base *base
}

// UploadFile uploads file
func (u httpUploader) UploadFile(accountUID AccountUID, projectID string, req smfile.FileUploadRequest) (UploadFileResponse, error) {
	filePath := buildUploadFilePath(accountUID)
	path := joinPath(mtBasePath, filePath)

	var response uploadFileResponse
	form, err := req.GetMTForm()
	if err != nil {
		return UploadFileResponse{}, fmt.Errorf("failed to create file upload form: %w", err)
	}

	err = form.Close()
	if err != nil {
		return UploadFileResponse{}, fmt.Errorf("failed to close upload file form: %w", err)
	}

	_, _, err = u.base.client.Post(
		path,
		form.Bytes(),
		&response,
		smclient.ContentTypeOption(form.GetContentType()),
	)
	if err != nil {
		return UploadFileResponse{}, fmt.Errorf(
			"failed to upload original file: %w",
			err,
		)
	}

	return UploadFileResponse{
		Code:    response.Response.Code,
		FileUID: FileUID(response.Response.Data.FileUID),
	}, nil
}

func buildUploadFilePath(accountUID AccountUID) string {
	return "/accounts/" + string(accountUID) + "/files"
}
