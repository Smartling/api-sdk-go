package job

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"

	smclient "github.com/Smartling/api-sdk-go/helpers/sm_client"
)

const jobBasePath = "/jobs-api/v3/projects/"

// Job defines the job behaviour
type Job interface {
	Get(projectID string, translationJobUID string) (GetJobResponse, error)
	GetAllByName(projectID, name string) (jobs []GetJobResponse, err error)
	Progress(projectID string, translationJobUID string) (GetJobProgressResponse, error)
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

// Get gets a job related info
func (h httpJob) Get(projectID string, translationJobUID string) (GetJobResponse, error) {
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

// GetAllByName gets all jobs of a project by name
func (h httpJob) GetAllByName(projectID, name string) ([]GetJobResponse, error) {
	reqURL := jobBasePath + projectID + "/jobs"

	params := url.Values{}
	params.Set("jobName", name)

	rawMessage, code, err := h.client.Get(reqURL, params)
	if err != nil {
		return nil, err
	}
	if code != 200 {
		body, _ := io.ReadAll(rawMessage)
		h.client.Logger.Debugf("response body: %s\n", body)
		return nil, fmt.Errorf("unexpected response code: %d", code)
	}
	body, err := io.ReadAll(rawMessage)
	if err != nil {
		body, _ := io.ReadAll(rawMessage)
		h.client.Logger.Debugf("response body: %s\n", body)
		return nil, err
	}
	var res getJobsResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	jobs := toGetJobsResponse(res)

	return jobs, nil
}

// Progress returns a job related progress
func (h httpJob) Progress(projectID string, translationJobUID string) (GetJobProgressResponse, error) {
	url := jobBasePath + projectID + "/jobs/" + translationJobUID + "/progress"
	var response getJobProgressResponse
	rawMessage, code, err := h.client.Get(url, nil)
	if err != nil {
		return GetJobProgressResponse{}, err
	}
	if code != 200 {
		body, _ := io.ReadAll(rawMessage)
		h.client.Logger.Debugf("response body: %s\n", body)
		return GetJobProgressResponse{}, fmt.Errorf("unexpected response code: %d", code)
	}
	body, err := io.ReadAll(rawMessage)
	if err != nil {
		body, _ := io.ReadAll(rawMessage)
		h.client.Logger.Debugf("response body: %s\n", body)
		return GetJobProgressResponse{}, err
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return GetJobProgressResponse{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return toGetJobProgressResponse(response, translationJobUID)
}
