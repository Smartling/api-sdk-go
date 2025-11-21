package job

// GetJobResponse defines get job response
type GetJobResponse struct {
	TranslationJobUID string
	JobName           string
}
type getJobResponse struct {
	Response struct {
		Code string `json:"code"`
		Data struct {
			JobName           string `json:"jobName"`
			TranslationJobUID string `json:"translationJobUid"`
		} `json:"data"`
	} `json:"response"`
}

func toGetJobResponse(r getJobResponse) GetJobResponse {
	return GetJobResponse{
		TranslationJobUID: r.Response.Data.TranslationJobUID,
		JobName:           r.Response.Data.JobName,
	}
}
