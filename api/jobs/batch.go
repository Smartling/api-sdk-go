package jobs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"

	smclient "github.com/Smartling/api-sdk-go/helpers/sm_client"
)

const jobBasePath = "/job-batches-api/v2/projects/"

type Batch interface {
	Create(ctx context.Context, projectID string, payload CreateBatchPayload) (CreateBatchResponse, error)
	CreateJob(ctx context.Context, projectID string, payload CreateJobPayload) (CreateJobResponse, error)
	UploadFile(ctx context.Context, projectID, batchUID string, payload UploadFilePayload) (UploadFileResponse, error)
	GetStatus(ctx context.Context, projectID, batchUID string) (GetStatusResponse, error)
}

func NewBatch(client *smclient.Client) Batch {
	return newHttpBatch(client)
}

type httpBatch struct {
	client *smclient.Client
}

func newHttpBatch(client *smclient.Client) httpBatch {
	return httpBatch{client: client}
}

func (h httpBatch) Create(ctx context.Context, projectID string, payload CreateBatchPayload) (CreateBatchResponse, error) {
	url := jobBasePath + projectID + "/batches"
	payloadB, err := json.Marshal(payload)
	if err != nil {
		return CreateBatchResponse{}, fmt.Errorf("unable to marshal: %w", err)
	}
	var response createBatchResponse
	resp, respCode, err := h.client.PostJSON(url, payloadB, &response)
	if err != nil {
		return CreateBatchResponse{}, fmt.Errorf("failed to create batch: %w", err)
	}
	if respCode != 200 {
		return CreateBatchResponse{}, fmt.Errorf("unexpected response code: %d, response: %s", respCode, resp)
	}
	return toCreateBatchResponse(response), nil
}

func (h httpBatch) CreateJob(ctx context.Context, projectID string, payload CreateJobPayload) (CreateJobResponse, error) {
	url := jobBasePath + projectID + "/jobs"
	payloadB, err := json.Marshal(payload)
	if err != nil {
		return CreateJobResponse{}, fmt.Errorf("unable to marshal: %w", err)
	}
	var response createJobResponse
	resp, respCode, err := h.client.PostJSON(url, payloadB, &response)
	if err != nil {
		return CreateJobResponse{}, fmt.Errorf("failed to create job: %w", err)
	}
	if respCode != 200 {
		return CreateJobResponse{}, fmt.Errorf("unexpected response code: %d, response: %s", respCode, resp)
	}
	return toCreateJobResponse(response), nil
}

func (h httpBatch) UploadFile(ctx context.Context, projectID, batchUID string, payload UploadFilePayload) (UploadFileResponse, error) {
	path := jobBasePath + projectID + "/batches/" + batchUID + "/file"
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	for _, locale := range payload.LocalesToAuthorize {
		if locale == "" {
			continue
		}
		err := writer.WriteField("localeIdsToAuthorize[]", locale)
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
	jsonRequest := fmt.Sprintf(`{"fileType":"%s"}`, payload.FileType)
	_, err = requestPart.Write([]byte(jsonRequest))
	if err != nil {
		return UploadFileResponse{}, fmt.Errorf("failed to write request: %v", err)
	}

	fileHeader := make(textproto.MIMEHeader)
	fileHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, payload.Filename))
	fileHeader.Set("Content-Type", "text/plain")

	filePart, err := writer.CreatePart(fileHeader)
	if err != nil {
		return UploadFileResponse{}, fmt.Errorf("failed to create file part: %v", err)
	}
	_, err = filePart.Write(payload.File)
	if err != nil {
		return UploadFileResponse{}, fmt.Errorf("failed to write: %v", err)
	}

	if err := writer.Close(); err != nil {
		return UploadFileResponse{}, fmt.Errorf("failed to close writer: %v", err)
	}

	url := h.client.BaseURL + path
	h.client.Logger.Debugf(
		"<- %s %s [payload %d bytes]\n",
		"POST", url, buf.Len(),
	)
	request, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		return UploadFileResponse{}, fmt.Errorf("failed to create request: %v", err)
	}

	request.Header.Set("Authorization", "Bearer "+h.client.Credentials.AccessToken.Value)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := h.client.HTTP.Do(request)
	if err != nil {
		return UploadFileResponse{}, fmt.Errorf("request failed: %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			h.client.Logger.Debugf("failed to close response body: %v", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return UploadFileResponse{}, fmt.Errorf("failed to read response body: %v", err)
	}
	h.client.Logger.Debugf("response body: %s\n", body)

	var response uploadFileResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return UploadFileResponse{}, fmt.Errorf("failed to unmarshal: %v", err)
	}

	return UploadFileResponse{
		Code: response.Response.Code,
	}, nil
}

func (h httpBatch) GetStatus(ctx context.Context, projectID, batchUID string) (GetStatusResponse, error) {
	url := jobBasePath + projectID + "/batches/" + batchUID
	var response getStatusResponse
	rawMessage, code, err := h.client.GetJSON(url, nil, &response)
	if err != nil {
		return GetStatusResponse{}, err
	}
	if code != 200 {
		return GetStatusResponse{}, fmt.Errorf("unexpected response code: %d with %s", code, rawMessage)
	}
	return toGetStatusResponse(response), nil
}
