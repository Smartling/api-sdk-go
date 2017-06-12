package smartling

import (
	"fmt"
)

const (
	endpointFilesList = "/files-api/v2/projects/%v/files/list"
	endpointFileTypes = "/files-api/v2/projects/%v/file-types"
)

// FilesList represents file list reply from Smartling APIa.
type FilesList struct {
	// TotalCount is a total files count.
	TotalCount int

	// Items contains all files matched by request.
	Items []FileStatus
}

// ListFiles returns files list from specified project by specified request.
func (client *Client) ListFiles(
	projectID string,
	request FilesListRequest,
) (*FilesList, error) {
	var list FilesList

	_, _, err := client.Get(
		fmt.Sprintf(endpointFilesList, projectID),
		request.GetQuery(),
		&list,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get files list: %s", err,
		)
	}

	return &list, nil
}

// ListFileTypes returns returns file types list from specified project.
func (client *Client) ListFileTypes(
	projectID string,
) ([]FileType, error) {
	var result struct {
		Items []FileType
	}

	_, _, err := client.Get(
		fmt.Sprintf(endpointFileTypes, projectID),
		nil,
		&result,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get file types: %s", err,
		)
	}

	return result.Items, nil
}
