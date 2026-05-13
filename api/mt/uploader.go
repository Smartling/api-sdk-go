package mt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/textproto"

	"github.com/Smartling/api-sdk-go/helpers/sm_client"
)

// Uploader defines uploader behaviour
type Uploader interface {
	UploadFile(ctx context.Context, accountUID AccountUID, filename string, req UploadFileRequest) (UploadFileResponse, error)
}

// NewUploader returns new Uploader implementation
func NewUploader(client *smclient.Client) Uploader {
	return httpUploader{base: newBase(client)}
}

type httpUploader struct {
	base *base
}

// UploadFile uploads file
func (u httpUploader) UploadFile(ctx context.Context, accountUID AccountUID, filename string, req UploadFileRequest) (UploadFileResponse, error) {
	filePath := buildUploadFilePath(accountUID)
	path := joinPath(mtBasePath, filePath)

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	for _, locale := range req.LocalesToAuthorize {
		if locale == "" {
			continue
		}
		if err := writer.WriteField("localeIdsToAuthorize[]", locale); err != nil {
			return UploadFileResponse{}, err
		}
	}

	for directive, value := range req.Directives {
		if err := writer.WriteField("smartling."+directive, value); err != nil {
			return UploadFileResponse{}, err
		}
	}

	requestHeader := make(textproto.MIMEHeader)
	requestHeader.Set("Content-Disposition", `form-data; name="request"`)
	requestHeader.Set("Content-Type", "application/json")
	requestPart, err := writer.CreatePart(requestHeader)
	if err != nil {
		return UploadFileResponse{}, fmt.Errorf("failed to create request part: %w", err)
	}
	if err := json.NewEncoder(requestPart).Encode(struct {
		FileType Type `json:"fileType"`
	}{req.FileType}); err != nil {
		return UploadFileResponse{}, fmt.Errorf("failed to write request: %w", err)
	}

	fileHeader := make(textproto.MIMEHeader)
	fileHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, filename))
	fileHeader.Set("Content-Type", "text/plain")
	filePart, err := writer.CreatePart(fileHeader)
	if err != nil {
		return UploadFileResponse{}, fmt.Errorf("failed to create file part: %w", err)
	}
	if _, err := filePart.Write(req.File); err != nil {
		return UploadFileResponse{}, fmt.Errorf("failed to write file: %w", err)
	}
	if err := writer.Close(); err != nil {
		return UploadFileResponse{}, fmt.Errorf("failed to close writer: %w", err)
	}

	var response uploadFileResponse
	_, code, err := u.base.client.PostJSON(
		ctx,
		path,
		buf.Bytes(),
		&response.Response.Data,
		smclient.ContentTypeOption(writer.FormDataContentType()),
	)
	if err != nil {
		return UploadFileResponse{}, fmt.Errorf("failed to upload file: %w", err)
	}
	response.Response.Code = code
	return toUploadFileResponse(response), nil
}

func buildUploadFilePath(accountUID AccountUID) string {
	return "/accounts/" + string(accountUID) + "/files"
}
