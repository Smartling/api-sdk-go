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
	"time"
)

const tokenExpirationSafetyDuration = 30 * time.Second

// Token represents authentication token, either access or refresh.
type Token struct {
	// Value is a string representation of token.
	Value string

	// ExpirationTime is a expiration time for token when it becomes invalid.
	ExpirationTime time.Time
}

// IsValid returns true if token still can be used.
func (token *Token) IsValid() bool {
	if token == nil {
		return false
	}

	if token.Value == "" {
		return false
	}

	return time.Now().Before(token.ExpirationTime)
}

// IsSafe returns true if token still can be used and it's expiration time is
// in safe bounds.
func (token *Token) IsSafe() bool {
	if !token.IsValid() {
		return false
	}

	return time.Now().Add(tokenExpirationSafetyDuration).Before(
		token.ExpirationTime,
	)
}

// String returns token representation for logging purposes only.
func (token *Token) String() string {
	if token == nil {
		return "[no token]"
	}

	return fmt.Sprintf(
		"[token=%s...{%d bytes} ttl %.2fs]",
		token.Value[:7],
		len(token.Value),
		time.Until(token.ExpirationTime).Seconds(),
	)
}
