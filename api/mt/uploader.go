package mt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"

	"github.com/Smartling/api-sdk-go/helpers/sm_client"
)

// Uploader defines uploader behaviour
type Uploader interface {
	UploadFile(accountUID AccountUID, filename string, req UploadFileRequest) (UploadFileResponse, error)
}

// NewUploader returns new Uploader implementation
func NewUploader(client *smclient.Client) Uploader {
	return httpUploader{base: newBase(client)}
}

type httpUploader struct {
	base *base
}

// UploadFile uploads file
func (u httpUploader) UploadFile(accountUID AccountUID, filename string, req UploadFileRequest) (UploadFileResponse, error) {
	filePath := buildUploadFilePath(accountUID)
	path := joinPath(mtBasePath, filePath)

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	for _, locale := range req.LocalesToAuthorize {
		if locale == "" {
			continue
		}
		err := writer.WriteField("localeIdsToAuthorize[]", locale)
		if err != nil {
			return UploadFileResponse{}, err
		}
	}

	for directive, value := range req.Directives {
		err := writer.WriteField("smartling."+directive, value)
		if err != nil {
			return UploadFileResponse{}, err
		}
	}

	requestHeader := make(textproto.MIMEHeader)
	requestHeader.Set("Content-Disposition", `form-data; name="request"`)
	requestHeader.Set("Content-Type", "application/json")

	requestPart, err := writer.CreatePart(requestHeader)
	if err != nil {
		return UploadFileResponse{}, fmt.Errorf("failed to create request part: %v", err)
	}
	jsonRequest := fmt.Sprintf(`{"fileType":"%s"}`, req.FileType)
	_, err = requestPart.Write([]byte(jsonRequest))
	if err != nil {
		return UploadFileResponse{}, fmt.Errorf("failed to write request: %v", err)
	}

	fileHeader := make(textproto.MIMEHeader)
	fileHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, filename))
	fileHeader.Set("Content-Type", "text/plain")

	filePart, err := writer.CreatePart(fileHeader)
	if err != nil {
		return UploadFileResponse{}, fmt.Errorf("failed to create file part: %v", err)
	}
	_, err = filePart.Write(req.File)
	if err != nil {
		return UploadFileResponse{}, fmt.Errorf("failed to write: %v", err)
	}

	writer.Close()

	url := u.base.client.BaseURL + path
	u.base.client.Logger.Debugf(
		"<- %s %s [payload %d bytes]\n",
		"POST", url, buf.Len(),
	)
	request, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		return UploadFileResponse{}, fmt.Errorf("failed to create request: %v", err)
	}

	request.Header.Set("Authorization", "Bearer "+u.base.client.Credentials.AccessToken.Value)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := u.base.client.HTTP.Do(request)
	if err != nil {
		return UploadFileResponse{}, fmt.Errorf("request failed: %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			u.base.client.Logger.Debugf("failed to close response body: %v", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return UploadFileResponse{}, fmt.Errorf("failed to read response body: %v", err)
	}
	u.base.client.Logger.Debugf("response body: %s\n", body)

	var response uploadFileResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return UploadFileResponse{}, fmt.Errorf("failed to unmarshal: %v", err)
	}

	return UploadFileResponse{
		Code:    response.Response.Code,
		FileUID: FileUID(response.Response.Data.FileUID),
	}, nil
}

func buildUploadFilePath(accountUID AccountUID) string {
	return "/accounts/" + string(accountUID) + "/files"
}
