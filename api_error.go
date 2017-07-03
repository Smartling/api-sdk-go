package smartling

import (
	"fmt"
	"net/url"
)

type APIError struct {
	Cause    error
	Code     string
	URL      string
	Params   url.Values
	Payload  []byte
	Response []byte
}

func (err APIError) Error() string {
	url := err.URL
	if len(err.Params) > 0 {
		url += "?" + err.Params.Encode()
	}

	code := ""
	if err.Code != "" {
		code = fmt.Sprintf(" [code %s]", code)
	}

	payload := ""
	if len(err.Payload) > 0 {
		payload = fmt.Sprintf(
			"\nRequest Body:\n%s",
			string(err.Payload),
		)
	}

	response := ""
	if len(err.Response) > 0 {
		response = fmt.Sprintf(
			"\nResponse Body:\n%s",
			string(err.Response),
		)
	}

	return fmt.Sprintf(
		"%s%s\nURL: %s%s%s",
		err.Cause,
		code,
		url,
		payload,
		response,
	)
}
