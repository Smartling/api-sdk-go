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
	"net/url"
)

// FilesListRequest represents request used to filter files returned by
// list files API call.
type FilesListRequest struct {
	// Cursor is a limit/offset pair, used to paginate reply.
	Cursor LimitOffsetRequest

	// URIMask instructs API to return only files with a URI containing the
	// given substring. Case is ignored.
	URIMask string

	// FileTypes instructs API to return only specified file types.
	FileTypes []FileType

	// LastUploadedAfter instructs API to return files uploaded after specified
	// date.
	LastUploadedAfter UTC

	// LastUploadedBefore instructs API to return files uploaded after
	// specified date.
	LastUploadedBefore UTC
}

// GetQuery returns URL values representation of files list request.
func (request *FilesListRequest) GetQuery() url.Values {
	query := request.Cursor.GetQuery()
	if len(request.URIMask) > 0 {
		query.Set("uriMask", request.URIMask)
	}

	for _, fileType := range request.FileTypes {
		query.Add("fileTypes[]", fmt.Sprint(fileType))
	}

	if !request.LastUploadedAfter.IsZero() {
		query.Set("lastUploadedAfter", request.LastUploadedAfter.String())
	}

	if !request.LastUploadedBefore.IsZero() {
		query.Set("lastUploadedBefore", request.LastUploadedBefore.String())
	}

	return query
}
