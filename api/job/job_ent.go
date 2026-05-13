package job

// GetJobResponse defines get job response
type GetJobResponse struct {
	Code              int
	TranslationJobUID string
	JobName           string
}
type getJobResponse struct {
	Response struct {
		Code int
		Data struct {
			JobName           string `json:"jobName"`
			TranslationJobUID string `json:"translationJobUid"`
		} `json:"data"`
	} `json:"response"`
}

func toGetJobResponse(r getJobResponse) GetJobResponse {
	return GetJobResponse{
		Code:              r.Response.Code,
		TranslationJobUID: r.Response.Data.TranslationJobUID,
		JobName:           r.Response.Data.JobName,
	}
}
