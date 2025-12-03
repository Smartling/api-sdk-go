package job

// GetJobResponse defines get job response
type GetJobResponse struct {
	TranslationJobUID string
	JobName           string
}

// FindFirstJobByName finds the first job by name from the list of jobs
func FindFirstJobByName(jobs []GetJobResponse, name string) GetJobResponse {
	for _, job := range jobs {
		if job.JobName == name {
			return job
		}
	}
	return GetJobResponse{}
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

func toGetJobResponse(r getJobResponse) GetJobResponse {
	return GetJobResponse{
		TranslationJobUID: r.Response.Data.TranslationJobUID,
		JobName:           r.Response.Data.JobName,
	}
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
