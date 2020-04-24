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
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type APIError struct {
	Cause    error
	Code     string
	URL      string
	Params   url.Values
	Payload  []byte
	Response []byte
	Headers  *http.Header
}

func (err APIError) Is(target error) bool {
	t, ok := target.(APIError)
	if !ok {
		return false
	}
	return (t.Cause == nil || err.Cause.Error() == t.Cause.Error()) &&
		(err.Code == t.Code || t.Code == "") &&
		(err.URL == t.URL || t.URL == "")
}

func (err APIError) Unwrap() error {
	return err.Cause
}

func (err APIError) Error() string {
	url := err.URL
	if len(err.Params) > 0 {
		url += "?" + err.Params.Encode()
	}

	code := ""
	if err.Code != "" {
		code = fmt.Sprintf(" [code %s]", code)
	}

	headers := &bytes.Buffer{}
	_ = err.Headers.Write(headers)

	return fmt.Sprintf(
		"%s%s\nURL: %s%s%s%s",
		err.Cause,
		code,
		url,
		err.section("Request Body", string(err.Payload)),
		err.section("Response Body", string(err.Response)),
		err.section("Response Headers", headers.String()),
	)
}

func (err *APIError) section(header string, contents string) string {
	if contents == "" {
		return ""
	}

	return fmt.Sprintf(
		"\n\n%s:\n%s",
		header,
		regexp.MustCompile(`(?m)^`).ReplaceAllLiteralString(
			strings.TrimSuffix(contents, "\n"),
			"  ",
		),
	)
}
