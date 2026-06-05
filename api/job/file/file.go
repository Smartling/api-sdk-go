package jobfile

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"

	jobapi "github.com/Smartling/api-sdk-go/api/job"
	smclient "github.com/Smartling/api-sdk-go/helpers/sm_client"
	smerror "github.com/Smartling/api-sdk-go/helpers/sm_error"
)

const jobBasePath = "/jobs-api/v3/projects/"

// AddRequest is the payload for adding a file to a job.
type AddRequest struct {
	FileURI         string   `json:"fileUri"`
	TargetLocaleIDs []string `json:"targetLocaleIds,omitempty"`
}

// RemoveRequest is the payload for removing a file from a job.
type RemoveRequest struct {
	FileURI string `json:"fileUri"`
}

// Result reports how many strings (per locale) were affected by adding or
// removing a file. A file the API cannot match contributes to neither count.
type Result struct {
	SuccessCount int `json:"successCount"`
	FailCount    int `json:"failCount"`
}

// File is a single source file attached to a translation job.
type File struct {
	FileURI   string
	LocaleIDs []string
}

// ListResponse is a single page of a job's source files.
type ListResponse struct {
	Items      []File
	TotalCount int
}

// JobFile manages the files attached to a translation job. Add and Remove
// operate on a single fileUri; List returns a page of the job's files.
type JobFile interface {
	Add(ctx context.Context, projectID, translationJobUID string, req AddRequest) (Result, error)
	Remove(ctx context.Context, projectID, translationJobUID string, req RemoveRequest) (Result, error)
	List(ctx context.Context, projectID, translationJobUID string, limit, offset uint32) (ListResponse, error)
}

// NewJobFile returns new JobFile implementation
func NewJobFile(client *smclient.Client) JobFile {
	return newHttpJobFile(client)
}

// httpJobFile implements JobFile interface using HTTP client
type httpJobFile struct {
	client *smclient.Client
}

func newHttpJobFile(client *smclient.Client) httpJobFile {
	return httpJobFile{client: client}
}

// Add attaches a file (by URI) to a translation job.
func (h httpJobFile) Add(ctx context.Context, projectID, translationJobUID string, req AddRequest) (Result, error) {
	if err := requireFile(projectID, translationJobUID, req.FileURI); err != nil {
		return Result{}, err
	}
	payload, err := json.Marshal(req)
	if err != nil {
		return Result{}, fmt.Errorf("failed to marshal add file request: %w", err)
	}
	reqURL := path.Join(fileURL(projectID, translationJobUID), "add")
	var res Result
	if _, _, err := h.client.PostJSON(ctx, reqURL, payload, &res); err != nil {
		return Result{}, fmt.Errorf("failed to add file to job: %w", err)
	}
	return res, nil
}

// Remove detaches a file (by URI) from a translation job.
func (h httpJobFile) Remove(ctx context.Context, projectID, translationJobUID string, req RemoveRequest) (Result, error) {
	if err := requireFile(projectID, translationJobUID, req.FileURI); err != nil {
		return Result{}, err
	}
	payload, err := json.Marshal(req)
	if err != nil {
		return Result{}, fmt.Errorf("failed to marshal remove file request: %w", err)
	}
	reqURL := path.Join(fileURL(projectID, translationJobUID), "remove")
	var res Result
	if _, _, err := h.client.PostJSON(ctx, reqURL, payload, &res); err != nil {
		return Result{}, fmt.Errorf("failed to remove file from job: %w", err)
	}
	return res, nil
}

// List returns a single page of source files attached to a translation job.
func (h httpJobFile) List(ctx context.Context, projectID, translationJobUID string, limit, offset uint32) (ListResponse, error) {
	if err := requireIDs(projectID, translationJobUID); err != nil {
		return ListResponse{}, err
	}
	reqURL := path.Join(jobBasePath, url.PathEscape(projectID), "jobs",
		url.PathEscape(translationJobUID), "files")

	params := url.Values{}
	params.Set("limit", strconv.FormatUint(uint64(limit), 10))
	params.Set("offset", strconv.FormatUint(uint64(offset), 10))

	var page listFilesResponse
	_, code, err := h.client.GetJSON(ctx, reqURL, params, &page.Response.Data)
	if err != nil && code == http.StatusNotFound {
		return ListResponse{}, jobapi.ErrNotFound
	}
	if err != nil {
		return ListResponse{}, fmt.Errorf("failed to list job files: %w", err)
	}

	items := make([]File, 0, len(page.Response.Data.Items))
	for _, item := range page.Response.Data.Items {
		items = append(items, File{FileURI: item.URI, LocaleIDs: item.LocaleIDs})
	}
	return ListResponse{Items: items, TotalCount: page.Response.Data.TotalCount}, nil
}

// listFilesResponse is the raw decode shape for the job files listing.
type listFilesResponse struct {
	Response struct {
		Code string `json:"code"`
		Data struct {
			TotalCount int `json:"totalCount"`
			Items      []struct {
				URI       string   `json:"uri"`
				LocaleIDs []string `json:"localeIds"`
			} `json:"items"`
		} `json:"data"`
	} `json:"response"`
}

func fileURL(projectID, translationJobUID string) string {
	return path.Join(jobBasePath, url.PathEscape(projectID), "jobs",
		url.PathEscape(translationJobUID), "file")
}

func requireIDs(projectID, translationJobUID string) error {
	switch {
	case projectID == "":
		return smerror.ErrEmptyParam("projectID")
	case translationJobUID == "":
		return smerror.ErrEmptyParam("translationJobUID")
	}
	return nil
}

func requireFile(projectID, translationJobUID, fileURI string) error {
	if err := requireIDs(projectID, translationJobUID); err != nil {
		return err
	}
	if fileURI == "" {
		return smerror.ErrEmptyParam("fileUri")
	}
	return nil
}
