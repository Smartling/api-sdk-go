package smclient

import (
	"fmt"

	smfile "github.com/Smartling/api-sdk-go/helpers/sm_file"
)

const (
	endpointUploadFile = "/files-api/v2/projects/%s/file"
)

// UploadFile uploads file
func (c *Client) UploadFile(
	projectID string,
	request smfile.FileUploadRequest,
) (*smfile.FileUploadResult, error) {
	var result smfile.FileUploadResult

	form, err := request.GetForm()
	if err != nil {
		return nil, fmt.Errorf("failed to create file upload form: %w", err)
	}

	err = form.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close upload file form: %w", err)
	}

	_, _, err = c.PostJSON(
		fmt.Sprintf(endpointUploadFile, projectID),
		form.Bytes(),
		&result,
		ContentTypeOption(form.GetContentType()),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to upload original file: %w",
			err,
		)
	}

	return &result, nil
}
