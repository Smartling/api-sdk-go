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

package smclient

import (
	"encoding/json"
	"fmt"
	"time"

	smartling "github.com/Smartling/api-sdk-go/helpers/sm_error"
)

const (
	endpointAuthenticate        = "/auth-api/v2/authenticate"
	endpointAuthenticateRefresh = "/auth-api/v2/authenticate/refresh"
)

// Authenticate checks that access and refresh tokens are valid and refreshes
// them if needed.
func (c *Client) Authenticate() error {
	if c.Credentials.AccessToken.IsSafe() {
		return nil
	}

	if c.Credentials.UserID == "" {
		return fmt.Errorf("user ID in the client is not set")
	}

	if c.Credentials.Secret == "" {
		return fmt.Errorf("token secret in the client is not set")
	}

	c.Credentials.AccessToken = nil

	var (
		url    string
		params map[string]string
	)

	if c.Credentials.RefreshToken.IsSafe() {
		url = endpointAuthenticateRefresh
		params = map[string]string{
			"refreshToken": c.Credentials.RefreshToken.Value,
		}
	} else {
		url = endpointAuthenticate
		params = map[string]string{
			"userIdentifier": c.Credentials.UserID,
			"userSecret":     c.Credentials.Secret,
		}

		c.Credentials.RefreshToken = nil
	}

	payload, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf(
			"unable to encode authenticate params: %w", err,
		)
	}

	var response struct {
		AccessToken      string
		ExpiresIn        time.Duration
		RefreshToken     string
		RefreshExpiresIn time.Duration
	}

	_, _, err = c.PostJSON(url, payload, &response, WithoutAuthentication)
	if err != nil {
		if _, ok := err.(smartling.NotAuthorizedError); ok {
			return err
		}

		return fmt.Errorf("authenticate request failed: %w", err)
	}

	c.Credentials.AccessToken = &Token{
		Value: response.AccessToken,
		ExpirationTime: time.Now().Add(
			response.ExpiresIn * time.Second,
		),
	}

	c.Credentials.RefreshToken = &Token{
		Value: response.RefreshToken,
		ExpirationTime: time.Now().Add(
			time.Second * response.RefreshExpiresIn,
		),
	}

	return nil
}
