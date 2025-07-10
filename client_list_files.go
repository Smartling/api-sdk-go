// Copyright (c) 2020 Smartling
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package smartling

import (
	"fmt"

	"github.com/Smartling/api-sdk-go/helpers/sm_error"
	"github.com/Smartling/api-sdk-go/helpers/sm_file"
)

const (
	endpointFilesList = "/files-api/v2/projects/%s/files/list"
	endpointFileTypes = "/files-api/v2/projects/%s/file-types"
)

// FilesList represents file list reply from Smartling APIa.
type FilesList struct {
	// TotalCount is a total files count.
	TotalCount int

	// Items contains all files matched by request.
	Items []smfile.File
}

// ListFiles returns files list from specified project by specified request.
// Returned result is paginated, so check out TotalCount struct field in the
// reply. API can return only 500 files at once.
func (c *HttpAPIClient) ListFiles(
	projectID string,
	request smfile.FilesListRequest,
) (*FilesList, error) {
	var list FilesList

	_, _, err := c.Client.GetJSON(
		fmt.Sprintf(endpointFilesList, projectID),
		request.GetQuery(),
		&list,
	)
	if err != nil {
		if _, ok := err.(smerror.NotFoundError); ok {
			return nil, err
		}

		return nil, fmt.Errorf("failed to get files list: %w", err)
	}

	return &list, nil
}

// ListAllFiles returns all files by request, even if it requires several API
// calls.
func (c *HttpAPIClient) ListAllFiles(
	projectID string,
	request smfile.FilesListRequest,
) ([]smfile.File, error) {
	var result []smfile.File

	for {
		files, err := c.ListFiles(projectID, request)
		if err != nil {
			return nil, err
		}

		result = append(result, files.Items...)

		if request.Cursor.Limit > 0 {
			request.Cursor.Limit -= len(files.Items)

			if request.Cursor.Limit == 0 {
				return result, nil
			}
		}

		if request.Cursor.Offset+len(files.Items) >= files.TotalCount {
			break
		}

		request.Cursor.Offset += len(files.Items)

		c.Client.Logger.Infof(
			"<= %d/%d files retrieved",
			request.Cursor.Offset,
			files.TotalCount,
		)
	}

	c.Client.Logger.Infof(
		"<= %d/%d files retrieved",
		len(result),
		len(result),
	)

	return result, nil
}
