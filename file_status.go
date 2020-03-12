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
)

// FileStatus describes file translation status obtained by GetFileStatus
// method.
type FileStatus struct {
	File

	TotalStringCount int
	TotalWordCount   int
	TotalCount       int

	Items []FileStatusTranslation
}

func (fs FileStatus) GetFileStatusTranslation(locale string) (*FileStatusTranslation, error) {
	for i := range fs.Items {
		if fs.Items[i].LocaleID == locale {
			return &fs.Items[i], nil
		}
	}

	return nil, fmt.Errorf(
		"failed to get file status translation for locale: %s", locale,
	)
}

func (fs FileStatus) AwaitingAuthorizationStringCount() int {
	c := 0

	for i := range fs.Items {
		c += fs.Items[i].AwaitingAuthorizationStringCount(fs.TotalStringCount)
	}

	return c
}
