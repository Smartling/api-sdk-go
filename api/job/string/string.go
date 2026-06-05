package jobstring

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"strconv"

	smclient "github.com/Smartling/api-sdk-go/helpers/sm_client"
	smerror "github.com/Smartling/api-sdk-go/helpers/sm_error"
)

const jobBasePath = "/jobs-api/v3/projects/"

// AddRequest is the payload for adding strings to a job.
type AddRequest struct {
	Hashcodes       []string `json:"hashcodes"`
	TargetLocaleIDs []string `json:"targetLocaleIds,omitempty"`
	MoveEnabled     bool     `json:"moveEnabled,omitempty"`
}

// RemoveRequest is the payload for removing strings from a job.
type RemoveRequest struct {
	Hashcodes []string `json:"hashcodes"`
	LocaleIDs []string `json:"localeIds,omitempty"`
}

// ListParams carries the filters for listing a job's strings.
type ListParams struct {
	TargetLocaleID string
	Limit          uint32
	Offset         uint32
}

// StringHashcode is a single string entry in a job.
type StringHashcode struct {
	TargetLocaleID string `json:"targetLocaleId"`
	Hashcode       string `json:"hashcode"`
}

// ListResponse is the result of listing a job's strings.
type ListResponse struct {
	TotalCount uint32           `json:"totalCount"`
	Items      []StringHashcode `json:"items"`
}

// Result reports how many strings (per locale) were affected by an
// add/remove operation. Hashcodes that don't exist in the project are silently
// ignored by the API and counted in neither field.
type Result struct {
	SuccessCount int `json:"successCount"`
	FailCount    int `json:"failCount"`
}

// JobString manages strings on a translation job.
type JobString interface {
	Add(ctx context.Context, projectID, translationJobUID string, req AddRequest) (Result, error)
	Remove(ctx context.Context, projectID, translationJobUID string, req RemoveRequest) (Result, error)
	List(ctx context.Context, projectID, translationJobUID string, params ListParams) (ListResponse, error)
}

// NewJobString returns new JobString implementation
func NewJobString(client *smclient.Client) JobString {
	return newHttpJobString(client)
}

// httpJobString implements JobString interface using HTTP client
type httpJobString struct {
	client *smclient.Client
}

func newHttpJobString(client *smclient.Client) httpJobString {
	return httpJobString{client: client}
}

// Add assigns strings (by hashcode) to a translation job. The returned Result
// reports how many strings were actually added; nonexistent hashcodes are
// ignored by the API and counted in neither successCount nor failCount.
func (h httpJobString) Add(ctx context.Context, projectID, translationJobUID string, req AddRequest) (Result, error) {
	if err := requireStrings(projectID, translationJobUID, req.Hashcodes); err != nil {
		return Result{}, err
	}
	payload, err := json.Marshal(req)
	if err != nil {
		return Result{}, fmt.Errorf("failed to marshal add strings request: %w", err)
	}
	reqURL := path.Join(stringsURL(projectID, translationJobUID), "add")
	var res Result
	if _, _, err := h.client.PostJSON(ctx, reqURL, payload, &res); err != nil {
		return Result{}, fmt.Errorf("failed to add strings to job: %w", err)
	}
	return res, nil
}

// Remove unassigns strings (by hashcode) from a translation job. The returned
// Result reports how many strings were actually removed.
func (h httpJobString) Remove(ctx context.Context, projectID, translationJobUID string, req RemoveRequest) (Result, error) {
	if err := requireStrings(projectID, translationJobUID, req.Hashcodes); err != nil {
		return Result{}, err
	}
	payload, err := json.Marshal(req)
	if err != nil {
		return Result{}, fmt.Errorf("failed to marshal remove strings request: %w", err)
	}
	reqURL := path.Join(stringsURL(projectID, translationJobUID), "remove")
	var res Result
	if _, _, err := h.client.PostJSON(ctx, reqURL, payload, &res); err != nil {
		return Result{}, fmt.Errorf("failed to remove strings from job: %w", err)
	}
	return res, nil
}

// List returns the strings assigned to a translation job.
func (h httpJobString) List(ctx context.Context, projectID, translationJobUID string, params ListParams) (ListResponse, error) {
	if err := requireIDs(projectID, translationJobUID); err != nil {
		return ListResponse{}, err
	}
	q := url.Values{}
	if params.TargetLocaleID != "" {
		q.Set("targetLocaleId", params.TargetLocaleID)
	}
	if params.Limit > 0 {
		q.Set("limit", strconv.FormatUint(uint64(params.Limit), 10))
	}
	if params.Offset > 0 {
		q.Set("offset", strconv.FormatUint(uint64(params.Offset), 10))
	}

	var resp ListResponse
	if _, _, err := h.client.GetJSON(ctx, stringsURL(projectID, translationJobUID), q, &resp); err != nil {
		return ListResponse{}, fmt.Errorf("failed to list job strings: %w", err)
	}
	return resp, nil
}

func stringsURL(projectID, translationJobUID string) string {
	return path.Join(jobBasePath, url.PathEscape(projectID), "jobs",
		url.PathEscape(translationJobUID), "strings")
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

func requireStrings(projectID, translationJobUID string, hashcodes []string) error {
	if err := requireIDs(projectID, translationJobUID); err != nil {
		return err
	}
	if len(hashcodes) == 0 {
		return smerror.ErrEmptyParam("hashcodes")
	}
	return nil
}
