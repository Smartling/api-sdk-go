package job

import (
	"fmt"
	"time"

	"github.com/Smartling/api-sdk-go/helpers"
)

// FindJobsByStringsRequest is the body of the find-jobs-by-strings endpoint.
type FindJobsByStringsRequest struct {
	Hashcodes []string
	LocaleIDs []string
}

// JobHashcodesByLocale lists the matching hashcodes for one locale within a job.
type JobHashcodesByLocale struct {
	LocaleID  string
	Hashcodes []string
}

// JobWithStrings is a job that contains some of the searched strings, grouped by
// locale.
type JobWithStrings struct {
	TranslationJobUID string
	JobName           string
	DueDate           time.Time
	HashcodesByLocale []JobHashcodesByLocale
}

// FindJobsByStringsResponse is the result of the find-jobs-by-strings endpoint.
type FindJobsByStringsResponse struct {
	TotalCount int
	Items      []JobWithStrings
}

// findJobsByStringsData is the wire shape of the find-jobs-by-strings response.
type findJobsByStringsData struct {
	TotalCount int `json:"totalCount"`
	Items      []struct {
		TranslationJobUID string `json:"translationJobUid"`
		JobName           string `json:"jobName"`
		DueDate           string `json:"dueDate"`
		HashcodesByLocale []struct {
			LocaleID  string   `json:"localeId"`
			Hashcodes []string `json:"hashcodes"`
		} `json:"hashcodesByLocale"`
	} `json:"items"`
}

func toFindJobsByStringsResponse(d findJobsByStringsData) (FindJobsByStringsResponse, error) {
	items := make([]JobWithStrings, 0, len(d.Items))
	for _, item := range d.Items {
		dueDate, err := helpers.StringToTime(item.DueDate, time.RFC3339)
		if err != nil {
			return FindJobsByStringsResponse{}, fmt.Errorf("parse DueDate: %w", err)
		}
		byLocale := make([]JobHashcodesByLocale, 0, len(item.HashcodesByLocale))
		for _, hl := range item.HashcodesByLocale {
			byLocale = append(byLocale, JobHashcodesByLocale{
				LocaleID:  hl.LocaleID,
				Hashcodes: hl.Hashcodes,
			})
		}
		items = append(items, JobWithStrings{
			TranslationJobUID: item.TranslationJobUID,
			JobName:           item.JobName,
			DueDate:           dueDate,
			HashcodesByLocale: byLocale,
		})
	}
	return FindJobsByStringsResponse{TotalCount: d.TotalCount, Items: items}, nil
}
