package job

import (
	"context"
	"fmt"

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
	_, code, err := h.client.GetJSON(ctx, url, nil, &response.Response.Data)
	if err != nil {
		return GetJobResponse{}, fmt.Errorf("failed to get job: %w", err)
	}
	response.Response.Code = code
	return toGetJobResponse(response), nil
}
