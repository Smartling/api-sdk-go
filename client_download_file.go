package smartling

import (
	"fmt"
	"io"
)

const (
	endpointDownloadFile = "/files-api/v2/projects/%s/locales/%s/file"
)

// DownloadFile downloads specified translated file for specified locale.
// Check FileDownloadRequest for more options.
func (client *Client) DownloadFile(
	projectID string,
	localeID string,
	request FileDownloadRequest,
) (io.Reader, error) {
	reader, _, err := client.Get(
		fmt.Sprintf(endpointDownloadFile, projectID, localeID),
		request.GetQuery(),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to download file: %s", err,
		)
	}

	return reader, nil
}
