package job

import (
	"encoding/json"
	"fmt"
)

// GetJobProgressResponse defines get job progress response
type GetJobProgressResponse struct {
	TranslationJobUID string
	TotalWordCount    uint32
	PercentComplete   uint32
	Json              []byte
}
type getJobProgressResponse struct {
	Response struct {
		Code string `json:"code"`
		Data struct {
			ContentProgressReport []struct {
				Progress struct {
					PercentComplete int `json:"percentComplete"`
					TotalWordCount  int `json:"totalWordCount"`
				} `json:"progress"`
				TargetLocaleDescription    string `json:"targetLocaleDescription"`
				TargetLocaleId             string `json:"targetLocaleId"`
				UnauthorizedProgressReport struct {
					StringCount int `json:"stringCount"`
					WordCount   int `json:"wordCount"`
				} `json:"unauthorizedProgressReport"`
				WorkflowProgressReportList []struct {
					WorkflowName                      string `json:"workflowName"`
					WorkflowStepSummaryReportItemList []struct {
						StringCount      int    `json:"stringCount"`
						WordCount        int    `json:"wordCount"`
						WorkflowStepName string `json:"workflowStepName"`
						WorkflowStepType string `json:"workflowStepType"`
						WorkflowStepUid  string `json:"workflowStepUid"`
					} `json:"workflowStepSummaryReportItemList"`
					WorkflowUid string `json:"workflowUid"`
				} `json:"workflowProgressReportList"`
			} `json:"contentProgressReport"`
			Progress struct {
				PercentComplete int `json:"percentComplete"`
				TotalWordCount  int `json:"totalWordCount"`
			} `json:"progress"`
			SummaryReport []struct {
				StringCount      int    `json:"stringCount"`
				WordCount        int    `json:"wordCount"`
				WorkflowStepName string `json:"workflowStepName"`
			} `json:"summaryReport"`
		} `json:"data"`
	} `json:"response"`
}

func toGetJobProgressResponse(r getJobProgressResponse, translationJobUID string) (GetJobProgressResponse, error) {
	data, err := json.Marshal(r.Response.Data)
	if err != nil {
		return GetJobProgressResponse{}, fmt.Errorf("failed to marshal job progress response: %w", err)
	}
	return GetJobProgressResponse{
		TranslationJobUID: translationJobUID,
		TotalWordCount:    uint32(r.Response.Data.Progress.TotalWordCount),
		PercentComplete:   uint32(r.Response.Data.Progress.PercentComplete),
		Json:              data,
	}, nil
}
