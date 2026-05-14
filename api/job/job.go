package job

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	smclient "github.com/Smartling/api-sdk-go/helpers/sm_client"
)

const jobBasePath = "/jobs-api/v3/projects/"

var ErrNotFound = errors.New("job not found")

// Job defines the job behaviour
type Job interface {
	GetJob(ctx context.Context, projectID, translationJobUID string) (GetJobResponse, error)
	SearchByName(ctx context.Context, projectID, name string) (jobs []GetJobResponse, err error)
	Progress(ctx context.Context, projectID string, translationJobUID string) (GetJobProgressResponse, error)
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
	if err != nil && code == http.StatusNotFound {
		return GetJobResponse{}, ErrNotFound

	}
	if err != nil {
		return GetJobResponse{}, fmt.Errorf("failed to get job: %w", err)
	}
	response.Response.Code = code
	return toGetJobResponse(response), nil
}

// SearchByName searches all jobs of a project by name
func (h httpJob) SearchByName(ctx context.Context, projectID, name string) ([]GetJobResponse, error) {
	reqURL := path.Join(jobBasePath, url.PathEscape(projectID), "jobs")

	params := url.Values{}
	params.Set("jobName", name)

	rawMessage, code, err := h.client.Get(ctx, reqURL, params)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rawMessage.Close(); err != nil {
			h.client.Logger.Debugf("failed to close response body: %v", err)
		}
	}()
	body, readErr := io.ReadAll(rawMessage)
	if code != http.StatusOK {
		if readErr != nil {
			return nil, fmt.Errorf("unexpected response code: %d, body: %s, readErr: %v", code, body, readErr)
		}
		return nil, fmt.Errorf("unexpected response code: %d, body: %s", code, body)
	}
	var res getJobsResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	jobs := toGetJobsResponse(res)

	return jobs, nil
}

// Progress returns a job related progress
func (h httpJob) Progress(ctx context.Context, projectID string, translationJobUID string) (GetJobProgressResponse, error) {
	reqURL := path.Join(jobBasePath, url.PathEscape(projectID), "jobs", url.PathEscape(translationJobUID), "progress")
	var response getJobProgressResponse
	rawMessage, code, err := h.client.Get(ctx, reqURL, nil)
	if err != nil {
		return GetJobProgressResponse{}, err
	}
	defer func() {
		if err := rawMessage.Close(); err != nil {
			h.client.Logger.Debugf("failed to close response body: %v", err)
		}
	}()
	body, err := io.ReadAll(rawMessage)
	if err != nil {
		return GetJobProgressResponse{}, fmt.Errorf("failed to read response body: %w", err)
	}
	if code != 200 {
		h.client.Logger.Debugf("response body: %s\n", body)
		return GetJobProgressResponse{}, fmt.Errorf("unexpected response code: %d", code)
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return GetJobProgressResponse{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return toGetJobProgressResponse(response, translationJobUID)
}
