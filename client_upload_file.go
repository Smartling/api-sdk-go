package smartling

import (
	"fmt"

	hierr "github.com/reconquest/hierr-go"
)

const (
	endpointUploadFile = "/files-api/v2/projects/%s/file"
)

type FileUploadResult struct {
	Overwritten bool
	StringCount int
	WordCount   int
}

// DownloadFile downloads original file from project.
func (client *Client) UploadFile(
	projectID string,
	request FileUploadRequest,
) (*FileUploadResult, error) {
	var result FileUploadResult

	contentType, body, err := request.GetForm()
	if err != nil {
		return nil, hierr.Errorf(
			err,
			"failed to create file upload form",
		)
	}

	_, _, err = client.Post(
		fmt.Sprintf(endpointUploadFile, projectID),
		body,
		&result,
		ContentTypeOption(contentType),
	)
	if err != nil {
		return nil, hierr.Errorf(
			err,
			"failed to download original file",
		)
	}

	return &result, nil
}
