package locale

import (
	"context"
	"fmt"
	"net/url"
	"path"

	smclient "github.com/Smartling/api-sdk-go/helpers/sm_client"
	smerror "github.com/Smartling/api-sdk-go/helpers/sm_error"
)

const jobBasePath = "/jobs-api/v3/projects/"

// JobLocale manages target locales on a translation job.
type JobLocale interface {
	Add(ctx context.Context, projectID, jobUID, targetLocaleID string) error
	Remove(ctx context.Context, projectID, jobUID, targetLocaleID string) error
}

// NewJobLocale returns new JobLocale implementation
func NewJobLocale(client *smclient.Client) JobLocale {
	return newHttpJobLocale(client)
}

// httpJobLocale implements JobLocale interface using HTTP client
type httpJobLocale struct {
	client *smclient.Client
}

func newHttpJobLocale(client *smclient.Client) httpJobLocale {
	return httpJobLocale{client: client}
}

// Add assigns a target locale to a translation job.
func (h httpJobLocale) Add(ctx context.Context, projectID, jobUID, targetLocaleID string) error {
	if err := requireParams(projectID, jobUID, targetLocaleID); err != nil {
		return err
	}
	reqURL := localeURL(projectID, jobUID, targetLocaleID)
	if _, _, err := h.client.PostJSON(ctx, reqURL, nil, nil); err != nil {
		return fmt.Errorf("failed to add locale to job: %w", err)
	}
	return nil
}

// Remove unassigns a target locale from a translation job.
func (h httpJobLocale) Remove(ctx context.Context, projectID, jobUID, targetLocaleID string) error {
	if err := requireParams(projectID, jobUID, targetLocaleID); err != nil {
		return err
	}
	reqURL := localeURL(projectID, jobUID, targetLocaleID)
	if _, _, err := h.client.DeleteJSON(ctx, reqURL, nil); err != nil {
		return fmt.Errorf("failed to remove locale from job: %w", err)
	}
	return nil
}

func localeURL(projectID, translationJobUID, targetLocaleID string) string {
	return path.Join(jobBasePath, url.PathEscape(projectID), "jobs",
		url.PathEscape(translationJobUID), "locales", url.PathEscape(targetLocaleID))
}

func requireParams(projectID, translationJobUID, targetLocaleID string) error {
	switch {
	case projectID == "":
		return smerror.ErrEmptyParam("projectID")
	case translationJobUID == "":
		return smerror.ErrEmptyParam("translationJobUID")
	case targetLocaleID == "":
		return smerror.ErrEmptyParam("targetLocaleID")
	}
	return nil
}
