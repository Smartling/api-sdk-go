package job

import (
	"encoding/json"
	"fmt"
	"io"

	smclient "github.com/Smartling/api-sdk-go/helpers/sm_client"
)

const jobBasePath = "/job-batches-api/v2/projects/"

// Job defines the job behaviour
type Job interface {
	GetJob(projectID string, translationJobUID string) (GetJobResponse, error)
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
func (h httpJob) GetJob(projectID string, translationJobUID string) (GetJobResponse, error) {
	url := jobBasePath + projectID + "/jobs/" + translationJobUID
	var response getJobResponse
	rawMessage, code, err := h.client.Get(url, nil)
	if err != nil {
		return GetJobResponse{}, err
	}
	if code != 200 {
		body, _ := io.ReadAll(rawMessage)
		h.client.Logger.Debugf("response body: %s\n", body)
		return GetJobResponse{}, fmt.Errorf("unexpected response code: %d", code)
	}
	body, err := io.ReadAll(rawMessage)
	if err != nil {
		body, _ := io.ReadAll(rawMessage)
		h.client.Logger.Debugf("response body: %s\n", body)
		return GetJobResponse{}, err
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return GetJobResponse{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return toGetJobResponse(response), nil
}
