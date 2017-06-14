package smartling

import (
	"net/url"
)

// FileURIRequest represents fileUri query parameter, commonly used in API.
type FileURIRequest struct {
	FileURI string
}

// GetQuery returns URL value representation for file URI.
func (request FileURIRequest) GetQuery() url.Values {
	query := url.Values{}

	query.Set("fileUri", request.FileURI)

	return query
}
