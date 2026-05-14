package batches

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/textproto"

	smclient "github.com/Smartling/api-sdk-go/helpers/sm_client"
)

const jobBasePath = "/job-batches-api/v2/projects/"

var ErrNotFound = errors.New("batch not found")

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
	_, code, err := h.client.PostJSON(ctx, url, payloadB, &response.Response.Data)
	if err != nil {
		return CreateBatchResponse{}, fmt.Errorf("failed to create batch: %w", err)
	}
	response.Response.Code = code
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
	_, code, err := h.client.PostJSON(ctx, url, payloadB, &response.Response.Data)
	if err != nil {
		return CreateJobResponse{}, fmt.Errorf("failed to create job: %w", err)
	}
	response.Response.Code = code
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
		if err := writer.WriteField("localeIdsToAuthorize[]", locale); err != nil {
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
		if err := writer.WriteField(directive, value); err != nil {
			return UploadFileResponse{}, err
		}
	}

	fileHeader := make(textproto.MIMEHeader)
	fileHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, payload.Filename))
	fileHeader.Set("Content-Type", "text/plain")
	filePart, err := writer.CreatePart(fileHeader)
	if err != nil {
		return UploadFileResponse{}, fmt.Errorf("failed to create file part: %w", err)
	}
	if _, err := filePart.Write(payload.File); err != nil {
		return UploadFileResponse{}, fmt.Errorf("failed to write file: %w", err)
	}
	if err := writer.Close(); err != nil {
		return UploadFileResponse{}, fmt.Errorf("failed to close writer: %w", err)
	}

	_, code, err := h.client.PostJSON(
		ctx,
		path,
		buf.Bytes(),
		nil,
		smclient.ContentTypeOption(writer.FormDataContentType()),
	)
	if err != nil {
		return UploadFileResponse{}, fmt.Errorf("failed to upload file: %w", err)
	}
	return UploadFileResponse{Code: code}, nil
}

// GetStatus retrieves the status of a batch in the specified project
func (h httpBatch) GetStatus(ctx context.Context, projectID, batchUID string) (GetStatusResponse, error) {
	url := jobBasePath + projectID + "/batches/" + batchUID

	var response getStatusResponse
	_, code, err := h.client.GetJSON(ctx, url, nil, &response.Response.Data)
	if err != nil {
		return GetStatusResponse{}, fmt.Errorf("failed to get batch status: %w", err)
	}
	response.Response.Code = code
	return toGetStatusResponse(response), nil
}
