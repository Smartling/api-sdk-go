package job

// JobFile represents a single source file associated with a translation job.
type JobFile struct {
	FileURI   string
	LocaleIDs []string
}

// ListJobFilesResponse is a single page of the job files listing as returned
// by the Jobs API.
type ListJobFilesResponse struct {
	Items      []JobFile
	TotalCount int
}

type listJobFilesResponse struct {
	Response struct {
		Code string `json:"code"`
		Data struct {
			TotalCount int `json:"totalCount"`
			Items      []struct {
				URI       string   `json:"uri"`
				LocaleIDs []string `json:"localeIds"`
			} `json:"items"`
		} `json:"data"`
	} `json:"response"`
}
