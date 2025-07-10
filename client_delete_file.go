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

	"github.com/Smartling/api-sdk-go/helpers/sm_client"
	"github.com/Smartling/api-sdk-go/helpers/sm_file"
)

const (
	endpointFileDelete = "/files-api/v2/projects/%s/file/delete"
)

// DeleteFile removes specified files from project.
func (c *HttpAPIClient) DeleteFile(
	projectID string,
	uri string,
) error {
	request := smfile.FileURIRequest{
		FileURI: uri,
	}

	form, err := request.GetForm()
	if err != nil {
		return fmt.Errorf("failed to create file delete form: %w", err)
	}

	err = form.Close()
	if err != nil {
		return fmt.Errorf("failed to close file delete form: %w", err)
	}

	_, _, err = c.Client.Post(
		fmt.Sprintf(endpointFileDelete, projectID),
		form.Bytes(),
		nil,
		smclient.ContentTypeOption(form.GetContentType()),
	)
	if err != nil {
		return fmt.Errorf("failed to remove file: %w", err)
	}

	return nil
}
