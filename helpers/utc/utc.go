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

package utc

import (
	"encoding/json"
	"time"
)

const utcFormat = "2006-01-02T15:04:05Z"

// UTC represents time in UTC format (zero timezone).
type UTC struct {
	time.Time
}

// MarshalJSON returns JSON representation of UTC.
func (utc UTC) MarshalJSON() ([]byte, error) {
	return json.Marshal(utc.String())
}

// UnmarshalJSON parses JSON representation of UTC.
func (utc *UTC) UnmarshalJSON(data []byte) error {
	var formatted string

	err := json.Unmarshal(data, &formatted)
	if err != nil {
		return err
	}

	location, err := time.LoadLocation("UTC")
	if err != nil {
		return err
	}

	parsed, err := time.ParseInLocation(utcFormat, formatted, location)
	if err != nil {
		return err
	}

	*utc = UTC{parsed}

	return nil
}

// String returns string representation of UTC.
func (utc UTC) String() string {
	return utc.Format(utcFormat)
}
