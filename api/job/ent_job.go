package job

import (
	"fmt"
	"time"

	"github.com/Smartling/api-sdk-go/helpers"
)

// JobSourceFile is a source file attached to a job.
type JobSourceFile struct {
	Name    string
	URI     string
	FileUID string
}

// JobCustomField is a custom field value on a job.
type JobCustomField struct {
	FieldUID   string
	FieldName  string
	FieldValue string
}

// JobIssues holds the issue counts on a job.
type JobIssues struct {
	SourceIssuesCount      int
	TranslationIssuesCount int
}

// JobDates groups the timestamps on a job. Fields the API omits (e.g.
// completion dates on an in-progress job) are the zero time.Time.
type JobDates struct {
	Due             time.Time
	Modified        time.Time
	Created         time.Time
	FirstAuthorized time.Time
	LastAuthorized  time.Time
	FirstCompleted  time.Time
	LastCompleted   time.Time
}

// GetJobResponse defines the full job detail response.
type GetJobResponse struct {
	Code              int
	TranslationJobUID string
	JobName           string
	JobNumber         string
	Description       string
	ReferenceNumber   string
	JobStatus         string
	Priority          int
	TargetLocaleIDs   []string
	SourceFiles       []JobSourceFile
	CustomFields      []JobCustomField
	Issues            JobIssues
	CreatedByUserUID  string
	ModifiedByUserUID string
	Dates             JobDates
}

// FindFirstJobByName finds the first job by name from the list of jobs.
func FindFirstJobByName(jobs []JobSummary, name string) (JobSummary, bool) {
	for _, job := range jobs {
		if job.JobName == name {
			return job, true
		}
	}
	return JobSummary{}, false
}

type getJobResponse struct {
	Response struct {
		Code int
		Data struct {
			TranslationJobUID   string   `json:"translationJobUid"`
			JobName             string   `json:"jobName"`
			JobNumber           string   `json:"jobNumber"`
			Description         string   `json:"description"`
			ReferenceNumber     string   `json:"referenceNumber"`
			JobStatus           string   `json:"jobStatus"`
			DueDate             string   `json:"dueDate"`
			ModifiedDate        string   `json:"modifiedDate"`
			CreatedDate         string   `json:"createdDate"`
			ModifiedByUserUID   string   `json:"modifiedByUserUid"`
			CreatedByUserUID    string   `json:"createdByUserUid"`
			FirstAuthorizedDate string   `json:"firstAuthorizedDate"`
			LastAuthorizedDate  string   `json:"lastAuthorizedDate"`
			FirstCompletedDate  string   `json:"firstCompletedDate"`
			LastCompletedDate   string   `json:"lastCompletedDate"`
			Priority            int      `json:"priority"`
			TargetLocaleIDs     []string `json:"targetLocaleIds"`
			SourceFiles         []struct {
				Name    string `json:"name"`
				URI     string `json:"uri"`
				FileUID string `json:"fileUid"`
			} `json:"sourceFiles"`
			CustomFields []struct {
				FieldUID   string `json:"fieldUid"`
				FieldName  string `json:"fieldName"`
				FieldValue string `json:"fieldValue"`
			} `json:"customFields"`
			Issues struct {
				SourceIssuesCount      int `json:"sourceIssuesCount"`
				TranslationIssuesCount int `json:"translationIssuesCount"`
			} `json:"issues"`
		} `json:"data"`
	} `json:"response"`
}

func toGetJobResponse(r getJobResponse) (GetJobResponse, error) {
	data := r.Response.Data
	sourceFiles := make([]JobSourceFile, 0, len(data.SourceFiles))
	for _, f := range data.SourceFiles {
		sourceFiles = append(sourceFiles, JobSourceFile{Name: f.Name, URI: f.URI, FileUID: f.FileUID})
	}
	customFields := make([]JobCustomField, 0, len(data.CustomFields))
	for _, c := range data.CustomFields {
		customFields = append(customFields, JobCustomField{FieldUID: c.FieldUID, FieldName: c.FieldName, FieldValue: c.FieldValue})
	}

	dates, err := toJobDates(r)
	if err != nil {
		return GetJobResponse{}, err
	}

	return GetJobResponse{
		Code:              r.Response.Code,
		TranslationJobUID: data.TranslationJobUID,
		JobName:           data.JobName,
		JobNumber:         data.JobNumber,
		Description:       data.Description,
		ReferenceNumber:   data.ReferenceNumber,
		JobStatus:         data.JobStatus,
		Priority:          data.Priority,
		TargetLocaleIDs:   data.TargetLocaleIDs,
		SourceFiles:       sourceFiles,
		CustomFields:      customFields,
		Issues: JobIssues{
			SourceIssuesCount:      data.Issues.SourceIssuesCount,
			TranslationIssuesCount: data.Issues.TranslationIssuesCount,
		},
		CreatedByUserUID:  data.CreatedByUserUID,
		ModifiedByUserUID: data.ModifiedByUserUID,
		Dates:             dates,
	}, nil
}

// toJobDates converts the raw timestamps to JobDates, failing fast on the
// first unparseable value so the caller can decide how to handle it.
func toJobDates(r getJobResponse) (JobDates, error) {
	var res JobDates
	var err error
	res.Due, err = helpers.StringToTime(r.Response.Data.DueDate, time.RFC3339)
	if err != nil {
		return JobDates{}, fmt.Errorf("parse DueDate: %w", err)
	}
	res.Modified, err = helpers.StringToTime(r.Response.Data.ModifiedDate, time.RFC3339)
	if err != nil {
		return JobDates{}, fmt.Errorf("parse ModifiedDate: %w", err)
	}
	res.Created, err = helpers.StringToTime(r.Response.Data.CreatedDate, time.RFC3339)
	if err != nil {
		return JobDates{}, fmt.Errorf("parse CreatedDate: %w", err)
	}
	res.FirstAuthorized, err = helpers.StringToTime(r.Response.Data.FirstAuthorizedDate, time.RFC3339)
	if err != nil {
		return JobDates{}, fmt.Errorf("parse FirstAuthorizedDate: %w", err)
	}
	res.LastAuthorized, err = helpers.StringToTime(r.Response.Data.LastAuthorizedDate, time.RFC3339)
	if err != nil {
		return JobDates{}, fmt.Errorf("parse LastAuthorizedDate: %w", err)
	}
	res.FirstCompleted, err = helpers.StringToTime(r.Response.Data.FirstCompletedDate, time.RFC3339)
	if err != nil {
		return JobDates{}, fmt.Errorf("parse FirstCompletedDate: %w", err)
	}
	res.LastCompleted, err = helpers.StringToTime(r.Response.Data.LastCompletedDate, time.RFC3339)
	if err != nil {
		return JobDates{}, fmt.Errorf("parse LastCompletedDate: %w", err)
	}
	return res, nil
}
