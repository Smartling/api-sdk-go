package job

// GetJobResponse defines get job response
type GetJobResponse struct {
	Code              int
	TranslationJobUID string
	JobName           string
	TargetLocaleIDs   []string
}

// FindFirstJobByName finds the first job by name from the list of jobs
func FindFirstJobByName(jobs []GetJobResponse, name string) (GetJobResponse, bool) {
	for _, job := range jobs {
		if job.JobName == name {
			return job, true
		}
	}
	return GetJobResponse{}, false
}

type getJobResponse struct {
	Response struct {
		Code int
		Data struct {
			JobName           string   `json:"jobName"`
			TranslationJobUID string   `json:"translationJobUid"`
			TargetLocaleIDs   []string `json:"targetLocaleIds"`
		} `json:"data"`
	} `json:"response"`
}

func toGetJobResponse(r getJobResponse) GetJobResponse {
	return GetJobResponse{
		Code:              r.Response.Code,
		TranslationJobUID: r.Response.Data.TranslationJobUID,
		JobName:           r.Response.Data.JobName,
		TargetLocaleIDs:   r.Response.Data.TargetLocaleIDs,
	}
}

type getJobsResponse struct {
	Response struct {
		Code string `json:"code"`
		Data struct {
			Items []struct {
				JobName           string `json:"jobName"`
				TranslationJobUID string `json:"translationJobUid"`
			} `json:"items"`
		} `json:"data"`
	} `json:"response"`
}

func toGetJobsResponse(r getJobsResponse) []GetJobResponse {
	res := make([]GetJobResponse, len(r.Response.Data.Items))
	for i, job := range r.Response.Data.Items {
		res[i] = GetJobResponse{
			TranslationJobUID: job.TranslationJobUID,
			JobName:           job.JobName,
		}
	}
	return res
}
