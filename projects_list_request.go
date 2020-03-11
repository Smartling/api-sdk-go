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

// ProjectsListRequest is a request used in GetProjectsList method.
type ProjectsListRequest struct {
	// Cursor specifies limit/offset pagination pair.
	Cursor LimitOffsetRequest

	// ProjectNameFilter specifies filter for project name.
	ProjectNameFilter string

	// IncludeArchived specifies should archived items be included or not.
	IncludeArchived bool
}

// GetQuery returns URL-encoded representation of current request.
func (request ProjectsListRequest) GetQuery() url.Values {
	query := request.Cursor.GetQuery()

	if len(request.ProjectNameFilter) > 0 {
		query.Set("projectNameFilter", request.ProjectNameFilter)
	}

	query.Set("includeArchived", fmt.Sprint(request.IncludeArchived))

	return query
}
