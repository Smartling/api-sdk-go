package job

import (
	"fmt"
	"time"

	"github.com/Smartling/api-sdk-go/helpers"
)

// JobSummary is a single row in a jobs listing. ProjectID and Priority are
// only populated by the account-level listing.
type JobSummary struct {
	TranslationJobUID string
	JobName           string
	JobNumber         string
	Description       string
	JobStatus         string
	ReferenceNumber   string
	TargetLocaleIDs   []string
	ProjectID         string
	Priority          int
	Dates             JobDates
}

// ListJobsResponse is a single page of a jobs listing.
type ListJobsResponse struct {
	Items      []JobSummary
	TotalCount int
}

// ListProjectJobsParams carries filters for listing jobs within a project.
type ListProjectJobsParams struct {
	JobName            string
	JobNumber          string
	TranslationJobUIDs []string
	JobStatus          []string
	Limit              uint32
	Offset             uint32
	SortBy             string
	SortDirection      string
}

// ListAccountJobsParams carries filters for listing jobs within an account.
type ListAccountJobsParams struct {
	JobName       string
	ProjectIDs    []string
	JobStatus     []string
	WithPriority  bool
	Limit         uint32
	Offset        uint32
	SortBy        string
	SortDirection string
}

// SearchJobsRequest is the body of the jobs search endpoint.
type SearchJobsRequest struct {
	FileURIs           []string
	Hashcodes          []string
	TranslationJobUIDs []string
}

// listJobsData is the shared wire shape of project, account, and search
// list responses.
type listJobsData struct {
	TotalCount int `json:"totalCount"`
	Items      []struct {
		TranslationJobUID string   `json:"translationJobUid"`
		JobName           string   `json:"jobName"`
		JobNumber         string   `json:"jobNumber"`
		Description       string   `json:"description"`
		JobStatus         string   `json:"jobStatus"`
		DueDate           string   `json:"dueDate"`
		CreatedDate       string   `json:"createdDate"`
		ReferenceNumber   string   `json:"referenceNumber"`
		TargetLocaleIDs   []string `json:"targetLocaleIds"`
		ProjectID         string   `json:"projectId"`
		Priority          int      `json:"priority"`
	} `json:"items"`
}

func toListJobsResponse(d listJobsData) (ListJobsResponse, error) {
	items := make([]JobSummary, 0, len(d.Items))
	for _, item := range d.Items {
		var dates JobDates
		var err error
		dates.Due, err = helpers.StringToTime(item.DueDate, time.RFC3339)
		if err != nil {
			return ListJobsResponse{}, fmt.Errorf("parse DueDate: %w", err)
		}
		dates.Created, err = helpers.StringToTime(item.CreatedDate, time.RFC3339)
		if err != nil {
			return ListJobsResponse{}, fmt.Errorf("parse CreatedDate: %w", err)
		}
		items = append(items, JobSummary{
			TranslationJobUID: item.TranslationJobUID,
			JobName:           item.JobName,
			JobNumber:         item.JobNumber,
			Description:       item.Description,
			JobStatus:         item.JobStatus,
			ReferenceNumber:   item.ReferenceNumber,
			TargetLocaleIDs:   item.TargetLocaleIDs,
			ProjectID:         item.ProjectID,
			Priority:          item.Priority,
			Dates:             dates,
		})
	}
	return ListJobsResponse{Items: items, TotalCount: d.TotalCount}, nil
}
