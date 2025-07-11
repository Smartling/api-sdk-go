package jobs

import (
	"strings"
	"time"

	"github.com/Smartling/api-sdk-go/helpers/file"
)

type CreateBatchPayload struct {
	Authorize         bool             `json:"authorize"`
	TranslationJobUID string           `json:"translationJobUid"`
	FileUris          []string         `json:"fileUris"`
	LocaleWorkflows   []LocaleWorkflow `json:"localeWorkflows"`
}
type LocaleWorkflow struct {
	TargetLocaleID string `json:"targetLocaleId"`
	WorkflowUid    string `json:"workflowUid"`
}
type CreateBatchResponse struct {
	Code     string
	BatchUID string
}
type createBatchResponse struct {
	Response struct {
		Code string `json:"code"`
		Data struct {
			BatchUID string `json:"batchUid"`
		} `json:"data"`
	} `json:"response"`
}

func toCreateBatchResponse(r createBatchResponse) CreateBatchResponse {
	return CreateBatchResponse{
		Code:     r.Response.Code,
		BatchUID: r.Response.Data.BatchUID,
	}
}

// UploadFilePayload defines upload file payload
type UploadFilePayload struct {
	Filename           string
	File               []byte
	FileType           file.Type
	LocalesToAuthorize []string
}

type UploadFileResponse struct {
	Code string
}
type uploadFileResponse struct {
	Response struct {
		Code string      `json:"code"`
		Data interface{} `json:"data"`
	} `json:"response"`
}

func toUploadFileResponse(r uploadFileResponse) UploadFileResponse {
	return UploadFileResponse{
		Code: r.Response.Code,
	}
}

type GetStatusResponse struct {
	Code              string
	Authorized        bool
	GeneralErrors     string
	ProjectID         string
	Status            string
	TranslationJobUID string
	UpdatedDate       time.Time
	Files             []GetStatusFile
}

type GetStatusFile struct {
	Errors        string
	FileUri       string
	Status        string
	TargetLocales []TargetLocale
	UpdatedDate   time.Time
}

type TargetLocale struct {
	LocaleID     string
	StringsAdded int
}

type getStatusResponse struct {
	Response struct {
		Code string `json:"code"`
		Data struct {
			Authorized string `json:"authorized"`
			Files      []struct {
				Errors        string `json:"errors"`
				FileUri       string `json:"fileUri"`
				Status        string `json:"status"`
				TargetLocales []struct {
					LocaleID     string `json:"localeId"`
					StringsAdded int    `json:"stringsAdded"`
				} `json:"targetLocales"`
				UpdatedDate time.Time `json:"updatedDate"`
			} `json:"files"`
			GeneralErrors     string    `json:"generalErrors"`
			ProjectID         string    `json:"projectId"`
			Status            string    `json:"status"`
			TranslationJobUID string    `json:"translationJobUid"`
			UpdatedDate       time.Time `json:"updatedDate"`
		} `json:"data"`
	} `json:"response"`
}

func toGetStatusResponse(r getStatusResponse) GetStatusResponse {
	res := GetStatusResponse{
		Code:              r.Response.Code,
		Authorized:        strings.ToLower(r.Response.Data.Authorized) == "true",
		GeneralErrors:     r.Response.Data.GeneralErrors,
		ProjectID:         r.Response.Data.ProjectID,
		Status:            r.Response.Data.Status,
		TranslationJobUID: r.Response.Data.TranslationJobUID,
		UpdatedDate:       r.Response.Data.UpdatedDate,
	}
	res.Files = make([]GetStatusFile, len(r.Response.Data.Files))
	for i, file := range r.Response.Data.Files {
		res.Files[i] = GetStatusFile{
			Errors:      file.Errors,
			FileUri:     file.FileUri,
			Status:      file.Status,
			UpdatedDate: file.UpdatedDate,
		}
		res.Files[i].TargetLocales = make([]TargetLocale, len(file.TargetLocales))
		for j, locale := range file.TargetLocales {
			res.Files[i].TargetLocales[j] = TargetLocale{
				LocaleID:     locale.LocaleID,
				StringsAdded: locale.StringsAdded,
			}
		}
	}
	return res
}

type CreateJobPayload struct {
	NameTemplate    string   `json:"nameTemplate"`
	Description     string   `json:"description"`
	TargetLocaleIds []string `json:"targetLocaleIds"`
	Mode            string   `json:"mode"`
	Salt            string   `json:"salt"`
	TimeZoneName    string   `json:"timeZoneName"`
}

type CreateJobResponse struct {
	Code              string
	TranslationJobUID string
	JobName           string
	JobNumber         string
	TargetLocaleIDs   []string
	Description       string
	JobStatus         string
}
type createJobResponse struct {
	Response struct {
		Code string `json:"code"`
		Data struct {
			TranslationJobUID string   `json:"translationJobUid"`
			JobName           string   `json:"jobName"`
			JobNumber         string   `json:"jobNumber"`
			TargetLocaleIDs   []string `json:"targetLocaleIds"`
			Description       string   `json:"description"`
			JobStatus         string   `json:"jobStatus"`
		} `json:"data"`
	} `json:"response"`
}

func toCreateJobResponse(r createJobResponse) CreateJobResponse {
	return CreateJobResponse{
		Code:              r.Response.Code,
		TranslationJobUID: r.Response.Data.TranslationJobUID,
		JobName:           r.Response.Data.JobName,
		JobNumber:         r.Response.Data.JobNumber,
		TargetLocaleIDs:   r.Response.Data.TargetLocaleIDs,
		Description:       r.Response.Data.Description,
		JobStatus:         r.Response.Data.JobStatus,
	}
}
