package job

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"

	smclient "github.com/Smartling/api-sdk-go/helpers/sm_client"
)

const jobBasePath = "/jobs-api/v3/projects/"

var ErrNotFound = errors.New("job not found")

// Job defines the job behaviour
type Job interface {
	GetJob(ctx context.Context, projectID, jobUID string) (GetJobResponse, error)
	SearchByName(ctx context.Context, projectID, name string) (jobs []GetJobResponse, err error)
	Progress(ctx context.Context, projectID string, jobUID string) (GetJobProgressResponse, error)
	ListFiles(ctx context.Context, projectID, jobUID string, limit, offset uint32) (ListJobFilesResponse, error)
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
func (h httpJob) GetJob(ctx context.Context, projectID, jobUID string) (GetJobResponse, error) {
	reqURL := path.Join(jobBasePath, url.PathEscape(projectID), "jobs", url.PathEscape(jobUID))

	var response getJobResponse
	_, code, err := h.client.GetJSON(ctx, reqURL, nil, &response.Response.Data)
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

	var res getJobsResponse
	_, _, err := h.client.GetJSON(ctx, reqURL, params, &res.Response.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to get jobs: %w", err)
	}
	jobs := toGetJobsResponse(res)
	return jobs, nil
}

// ListFiles returns a single page of source files attached to a translation job.
func (h httpJob) ListFiles(ctx context.Context, projectID, jobUID string, limit, offset uint32) (ListJobFilesResponse, error) {
	reqURL := path.Join(jobBasePath, url.PathEscape(projectID), "jobs", url.PathEscape(jobUID), "files")

	params := url.Values{}
	params.Set("limit", strconv.FormatUint(uint64(limit), 10))
	params.Set("offset", strconv.FormatUint(uint64(offset), 10))

	var page listJobFilesResponse
	_, code, err := h.client.GetJSON(ctx, reqURL, params, &page.Response.Data)
	if err != nil && code == http.StatusNotFound {
		return ListJobFilesResponse{}, ErrNotFound
	}
	if err != nil {
		return ListJobFilesResponse{}, fmt.Errorf("failed to list job files: %w", err)
	}

	items := make([]JobFile, 0, len(page.Response.Data.Items))
	for _, item := range page.Response.Data.Items {
		items = append(items, JobFile{
			FileURI:   item.URI,
			LocaleIDs: item.LocaleIDs,
		})
	}
	return ListJobFilesResponse{
		Items:      items,
		TotalCount: page.Response.Data.TotalCount,
	}, nil
}

// Progress returns a job related progress
func (h httpJob) Progress(ctx context.Context, projectID string, jobUID string) (GetJobProgressResponse, error) {
	reqURL := path.Join(jobBasePath, url.PathEscape(projectID), "jobs", url.PathEscape(jobUID), "progress")
	var response getJobProgressResponse
	_, code, err := h.client.GetJSON(ctx, reqURL, nil, &response.Response.Data)
	if err != nil && code == http.StatusNotFound {
		return GetJobProgressResponse{}, ErrNotFound
	}
	if err != nil {
		return GetJobProgressResponse{}, fmt.Errorf("failed to get job progress: %w", err)
	}
	return toGetJobProgressResponse(response, jobUID)
}
