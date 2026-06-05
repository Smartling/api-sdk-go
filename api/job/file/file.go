package jobfile

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"path"

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

// JobFile manages the files attached to a translation job. Each call operates on
// a single fileUri.
type JobFile interface {
	Add(ctx context.Context, projectID, translationJobUID string, req AddRequest) (Result, error)
	Remove(ctx context.Context, projectID, translationJobUID string, req RemoveRequest) (Result, error)
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

func fileURL(projectID, translationJobUID string) string {
	return path.Join(jobBasePath, url.PathEscape(projectID), "jobs",
		url.PathEscape(translationJobUID), "file")
}

func requireFile(projectID, translationJobUID, fileURI string) error {
	switch {
	case projectID == "":
		return smerror.ErrEmptyParam("projectID")
	case translationJobUID == "":
		return smerror.ErrEmptyParam("translationJobUID")
	case fileURI == "":
		return smerror.ErrEmptyParam("fileUri")
	}
	return nil
}
