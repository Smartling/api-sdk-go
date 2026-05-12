package job

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	smclient "github.com/Smartling/api-sdk-go/helpers/sm_client"
)

const jobBasePath = "/job-batches-api/v2/projects/"

// Job defines the job behaviour
type Job interface {
	GetJob(ctx context.Context, projectID, translationJobUID string) (GetJobResponse, error)
}

// NewJob returns new Job implementation
func NewJob(client *smclient.Client) Job {
	return newHttpJob(client)
}

// httpJob implements Job interface using HTTP client
type httpJob struct {
	client *smclient.Client
}

func newHttpJob(client *smclient.Client) httpJob {
	return httpJob{client: client}
}

// GetJob gets a job related info
func (h httpJob) GetJob(ctx context.Context, projectID, translationJobUID string) (GetJobResponse, error) {
	url := jobBasePath + projectID + "/jobs/" + translationJobUID
	var response getJobResponse
	rawMessage, code, err := h.client.Get(ctx, url, nil)
	if err != nil {
		return GetJobResponse{}, err
	}
	defer func() {
		if err := rawMessage.Close(); err != nil {
			h.client.Logger.Debugf("failed to close response body: %v", err)
		}
	}()
	body, err := io.ReadAll(rawMessage)
	if err != nil {
		return GetJobResponse{}, fmt.Errorf("failed to read response body: %w", err)
	}
	if code != 200 {
		return GetJobResponse{}, fmt.Errorf("unexpected response code: %d with %s", code, body)
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return GetJobResponse{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return toGetJobResponse(response), nil
}
