package mt

// StartResponse defines start translation response
type StartResponse struct {
	Code  string
	MtUID string
}

// StartResponse defines start translation response
type startResponse struct {
	Response struct {
		Code string `json:"code"`
		Data struct {
			MtUID string `json:"mtUid"`
		} `json:"data"`
	} `json:"response"`
}

func toStartResponse(r startResponse) StartResponse {
	return StartResponse{
		Code:  r.Response.Code,
		MtUID: r.Response.Data.MtUID,
	}
}

// ProgressResponse defines progress translation response
type ProgressResponse struct {
	Code                  string
	State                 string
	RequestedStringCount  int
	Error                 string
	LocaleProcessStatuses []LocaleProcessStatusResponse
}

// LocaleProcessStatusResponse defines locale process status response
type LocaleProcessStatusResponse struct {
	LocaleID             string
	State                string
	ProcessedStringCount int
	Error                ErrorResponse
}

// ErrorResponse defines error response
type ErrorResponse struct {
	Key     string
	Message string
	ErrorID string
}

// IsSet checks if ErrorResponse is set
func (e ErrorResponse) IsSet() bool {
	return e.Key != ""
}

type progressResponse struct {
	Response struct {
		Code string `json:"code"`
		Data struct {
			State                 string `json:"state"`
			RequestedStringCount  int    `json:"requestedStringCount"`
			Error                 string `json:"error;omitempty"`
			LocaleProcessStatuses []struct {
				LocaleID             string `json:"localeId"`
				State                string `json:"state"`
				ProcessedStringCount int    `json:"processedStringCount"`
				Error                *struct {
					Key     string `json:"key"`
					Message string `json:"message"`
					Details struct {
						ErrorID string `json:"errorId"`
					} `json:"details"`
				} `json:"error"`
			} `json:"localeProcessStatuses"`
		} `json:"data"`
	} `json:"response"`
}

func toProgressResponse(r progressResponse) ProgressResponse {
	localeProcessStatuses := make([]LocaleProcessStatusResponse, len(r.Response.Data.LocaleProcessStatuses))
	for key, val := range r.Response.Data.LocaleProcessStatuses {
		localeProcessStatuses[key] = LocaleProcessStatusResponse{
			LocaleID:             val.LocaleID,
			State:                val.State,
			ProcessedStringCount: val.ProcessedStringCount,
			Error: ErrorResponse{
				Key:     val.Error.Key,
				Message: val.Error.Message,
				ErrorID: val.Error.Details.ErrorID,
			},
		}
	}
	return ProgressResponse{
		Code:                  r.Response.Code,
		State:                 r.Response.Data.State,
		RequestedStringCount:  r.Response.Data.RequestedStringCount,
		Error:                 r.Response.Data.Error,
		LocaleProcessStatuses: localeProcessStatuses,
	}
}
