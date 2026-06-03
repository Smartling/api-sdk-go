package job

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"

	smclient "github.com/Smartling/api-sdk-go/helpers/sm_client"
)

const (
	jobBasePath     = "/jobs-api/v3/projects/"
	accountBasePath = "/jobs-api/v3/accounts/"
)

var ErrNotFound = errors.New("job not found")

// Job defines the job behaviour
type Job interface {
	GetJob(ctx context.Context, projectID, jobUID string) (GetJobResponse, error)
	ListProjectJobs(ctx context.Context, projectID string, params ListProjectJobsParams) (ListJobsResponse, error)
	ListAccountJobs(ctx context.Context, accountUID string, params ListAccountJobsParams) (ListJobsResponse, error)
	SearchJobs(ctx context.Context, projectID string, req SearchJobsRequest) (ListJobsResponse, error)
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
	return toGetJobResponse(response)
}

// ListProjectJobs lists jobs within a project, applying the given filters.
func (h httpJob) ListProjectJobs(ctx context.Context, projectID string, params ListProjectJobsParams) (ListJobsResponse, error) {
	reqURL := path.Join(jobBasePath, url.PathEscape(projectID), "jobs")

	q := url.Values{}
	if params.JobName != "" {
		q.Set("jobName", params.JobName)
	}
	if params.JobNumber != "" {
		q.Set("jobNumber", params.JobNumber)
	}
	for _, uid := range params.TranslationJobUIDs {
		q.Add("translationJobUids", uid)
	}
	for _, st := range params.JobStatus {
		q.Add("translationJobStatus", st)
	}
	if params.Limit > 0 {
		q.Set("limit", strconv.FormatUint(uint64(params.Limit), 10))
	}
	if params.Offset > 0 {
		q.Set("offset", strconv.FormatUint(uint64(params.Offset), 10))
	}
	if params.SortBy != "" {
		q.Set("sortBy", params.SortBy)
	}
	if params.SortDirection != "" {
		q.Set("sortDirection", params.SortDirection)
	}

	var data listJobsData
	_, code, err := h.client.GetJSON(ctx, reqURL, q, &data)
	if err != nil && code == http.StatusNotFound {
		return ListJobsResponse{}, ErrNotFound
	}
	if err != nil {
		return ListJobsResponse{}, fmt.Errorf("failed to list jobs: %w", err)
	}
	return toListJobsResponse(data)
}

// ListAccountJobs lists jobs within an account, applying the given filters.
func (h httpJob) ListAccountJobs(ctx context.Context, accountUID string, params ListAccountJobsParams) (ListJobsResponse, error) {
	reqURL := path.Join(accountBasePath, url.PathEscape(accountUID), "jobs")

	q := url.Values{}
	if params.JobName != "" {
		q.Set("jobName", params.JobName)
	}
	for _, pid := range params.ProjectIDs {
		q.Add("projectIds", pid)
	}
	for _, st := range params.JobStatus {
		q.Add("translationJobStatus", st)
	}
	if params.WithPriority {
		q.Set("withPriority", "true")
	}
	if params.Limit > 0 {
		q.Set("limit", strconv.FormatUint(uint64(params.Limit), 10))
	}
	if params.Offset > 0 {
		q.Set("offset", strconv.FormatUint(uint64(params.Offset), 10))
	}
	if params.SortBy != "" {
		q.Set("sortBy", params.SortBy)
	}
	if params.SortDirection != "" {
		q.Set("sortDirection", params.SortDirection)
	}

	var data listJobsData
	_, code, err := h.client.GetJSON(ctx, reqURL, q, &data)
	if err != nil && code == http.StatusNotFound {
		return ListJobsResponse{}, ErrNotFound
	}
	if err != nil {
		return ListJobsResponse{}, fmt.Errorf("failed to list account jobs: %w", err)
	}
	return toListJobsResponse(data)
}

// SearchJobs finds jobs by file URIs, hashcodes, or job UIDs.
func (h httpJob) SearchJobs(ctx context.Context, projectID string, req SearchJobsRequest) (ListJobsResponse, error) {
	reqURL := path.Join(jobBasePath, url.PathEscape(projectID), "jobs", "search")

	body, err := json.Marshal(map[string][]string{
		"fileUris":           req.FileURIs,
		"hashcodes":          req.Hashcodes,
		"translationJobUids": req.TranslationJobUIDs,
	})
	if err != nil {
		return ListJobsResponse{}, fmt.Errorf("failed to encode search request: %w", err)
	}

	var data listJobsData
	_, code, err := h.client.PostJSON(ctx, reqURL, body, &data)
	if err != nil && code == http.StatusNotFound {
		return ListJobsResponse{}, ErrNotFound
	}
	if err != nil {
		return ListJobsResponse{}, fmt.Errorf("failed to search jobs: %w", err)
	}
	return toListJobsResponse(data)
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
