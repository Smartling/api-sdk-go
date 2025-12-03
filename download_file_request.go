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

	"github.com/Smartling/api-sdk-go/helpers/sm_file"
)

// RetrievalType describes type of file download.
// https://help.smartling.com/v1.0/reference#get_projects-projectid-locales-localeid-file
type RetrievalType string

const (
	// RetrieveDefault specifies that Smartling will decide what type of
	// translation will be returned.
	RetrieveDefault RetrievalType = ""

	// RetrievePending specifies that Smartling returns any translations
	// (including non-published translations)
	RetrievePending = "pending"

	// RetrievePublished specifies that Smartling returns only
	// published/pre-published translations.
	RetrievePublished = "published"

	// RetrievePseudo specifies that Smartling returns a modified version of
	// the original text with certain characters transformed and the text
	// expanded.
	RetrievePseudo = "pseudo"

	// RetrieveChromeInstrumented specifies that Smartling returns a modified
	// version of the original file with strings wrapped in a specific set of
	// Unicode symbols that can later be recognized and matched by the Chrome
	// Context Capture Extension
	RetrieveChromeInstrumented = "contextMatchingInstrumented"
)

// FileDownloadRequest represents optional parameters for file download
// operation.
type FileDownloadRequest struct {
	smfile.FileURIRequest

	Type            RetrievalType
	IncludeOriginal bool
}

// GetQuery returns URL values representation of download file query params.
func (request FileDownloadRequest) GetQuery() url.Values {
	query := request.FileURIRequest.GetQuery()

	query.Set("fileUri", request.FileURI)

	if request.Type != RetrieveDefault {
		query.Set("retrievalType", fmt.Sprint(request.Type))
	}

	if request.IncludeOriginal {
		query.Set("includeOriginalStrings", "true")
	}

	return query
}
