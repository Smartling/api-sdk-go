package smartling

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/url"
)

type FileUploadRequest struct {
	File      []byte
	FileURI   string
	FileType  FileType
	Authorize bool

	LocalesToAuthorize []string

	Smartling struct {
		Namespace   string
		FileCharset string
		Directives  map[string]string
	}
}

// GetQuery returns URL value representation for file URI.
func (request FileUploadRequest) GetQuery() url.Values {
	query := url.Values{}

	query.Set("file", string(request.File))
	query.Set("fileUri", request.FileURI)
	query.Set("fileType", string(request.FileType))

	if request.Authorize {
		query.Set("authorize", "true")
	}

	for _, locale := range request.LocalesToAuthorize {
		query.Add("localeIdsToAuthorize", locale)
	}

	if request.Smartling.Namespace != "" {
		query.Set("smartling.namespace", request.Smartling.Namespace)
	}

	if request.Smartling.FileCharset != "" {
		query.Set("smartling.file_charset", request.Smartling.FileCharset)
	}

	for directive, value := range request.Smartling.Directives {
		query.Set("smartling."+directive, value)
	}

	return query
}

func (request *FileUploadRequest) GetForm() (string, []byte, error) {
	var (
		body   = &bytes.Buffer{}
		writer = multipart.NewWriter(body)
		query  = request.GetQuery()
	)

	for key, _ := range query {
		value := query.Get(key)

		switch key {
		case "file":
			writer, err := writer.CreateFormFile("file", request.FileURI)
			if err != nil {
				return "", nil, err
			}

			_, err = io.WriteString(writer, value)
			if err != nil {
				return "", nil, err
			}

		default:
			err := writer.WriteField(key, value)
			if err != nil {
				return "", nil, err
			}
		}
	}

	err := writer.Close()
	if err != nil {
		return "", nil, err
	}

	return writer.FormDataContentType(), body.Bytes(), nil
}
