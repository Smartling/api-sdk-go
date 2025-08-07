package batches

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
	smerror "github.com/Smartling/api-sdk-go/helpers/sm_error"
)

const jobBasePath = "/job-batches-api/v2/projects/"

// Batch defines the batch behaviour
type Batch interface {
	Create(ctx context.Context, projectID string, payload CreateBatchPayload) (CreateBatchResponse, error)
	CreateJob(ctx context.Context, projectID string, payload CreateJobPayload) (CreateJobResponse, error)
	UploadFile(ctx context.Context, projectID, batchUID string, payload UploadFilePayload) (UploadFileResponse, error)
	GetStatus(ctx context.Context, projectID, batchUID string) (GetStatusResponse, error)
}

// NewBatch returns new Batch implementation
func NewBatch(client *smclient.Client) Batch {
	return newHttpBatch(client)
}

// httpBatch implements Batch interface using HTTP client
type httpBatch struct {
	client *smclient.Client
}

func newHttpBatch(client *smclient.Client) httpBatch {
	return httpBatch{client: client}
}

// Create creates a new batch in the specified project
func (h httpBatch) Create(ctx context.Context, projectID string, payload CreateBatchPayload) (CreateBatchResponse, error) {
	url := jobBasePath + projectID + "/batches"
	payloadB, err := json.Marshal(payload)
	if err != nil {
		return CreateBatchResponse{}, fmt.Errorf("unable to marshal: %w", err)
	}

	var response createBatchResponse
	resp, err := h.client.Post(url, payloadB, &response)
	if err != nil {
		return CreateBatchResponse{}, fmt.Errorf("failed to create batch: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return CreateBatchResponse{}, fmt.Errorf("unexpected response code: %d, response: %s", resp.StatusCode, resp)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return CreateBatchResponse{}, smerror.APIError{
			Cause:   err,
			URL:     url,
			Payload: payloadB,
		}
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return CreateBatchResponse{}, err
	}
	return toCreateBatchResponse(response), nil
}

// CreateJob creates a new job in the specified project
func (h httpBatch) CreateJob(ctx context.Context, projectID string, payload CreateJobPayload) (CreateJobResponse, error) {
	url := jobBasePath + projectID + "/jobs"
	payloadB, err := json.Marshal(payload)
	if err != nil {
		return CreateJobResponse{}, fmt.Errorf("unable to marshal: %w", err)
	}
	var response createJobResponse
	resp, err := h.client.Post(url, payloadB, &response)
	if err != nil {
		return CreateJobResponse{}, fmt.Errorf("failed to create job: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return CreateJobResponse{}, fmt.Errorf("unexpected response code: %d, response: %s", resp.StatusCode, resp)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return CreateJobResponse{}, smerror.APIError{
			Cause:   err,
			URL:     url,
			Payload: payloadB,
		}
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return CreateJobResponse{}, err
	}
	return toCreateJobResponse(response), nil
}

// UploadFile uploads a file to the specified batch in the project
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

	fileUri := payload.Filename
	if fileUri == "" {
		fileUri = payload.FileUri
	}
	if err := writer.WriteField("fileUri", fileUri); err != nil {
		return UploadFileResponse{}, err
	}
	if err := writer.WriteField("fileType", payload.FileType.String()); err != nil {
		return UploadFileResponse{}, err
	}
	for directive, value := range payload.Directives {
		if directive == "" {
			continue
		}
		if err := writer.WriteField("smartling."+directive, value); err != nil {
			return UploadFileResponse{}, err
		}
	}

	requestHeader := make(textproto.MIMEHeader)
	requestHeader.Set("Content-Disposition", `form-data; name="request"`)
	requestHeader.Set("Content-Type", "application/json")

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
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &buf)
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

// GetStatus retrieves the status of a batch in the specified project
func (h httpBatch) GetStatus(ctx context.Context, projectID, batchUID string) (GetStatusResponse, error) {
	url := jobBasePath + projectID + "/batches/" + batchUID
	var response getStatusResponse
	rawMessage, code, err := h.client.Get(url, nil)
	if err != nil {
		return GetStatusResponse{}, err
	}
	if code != 200 {
		return GetStatusResponse{}, fmt.Errorf("unexpected response code: %d with %s", code, rawMessage)
	}
	body, err := io.ReadAll(rawMessage)
	if err != nil {
		return GetStatusResponse{}, smerror.APIError{
			Cause:   err,
			URL:     url,
			Payload: []byte(fmt.Sprintf("%v", rawMessage)),
		}
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return GetStatusResponse{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return toGetStatusResponse(response), nil
}
