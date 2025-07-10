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
	"encoding/json"
	"fmt"
	"time"
)

const (
	endpointAuthenticate        = "/auth-api/v2/authenticate"
	endpointAuthenticateRefresh = "/auth-api/v2/authenticate/refresh"
)

// Authenticate checks that access and refresh tokens are valid and refreshes
// them if needed.
func (client *Client) Authenticate() error {
	if client.Credentials.AccessToken.IsSafe() {
		return nil
	}

	if client.Credentials.UserID == "" {
		return fmt.Errorf("user ID in the client is not set")
	}

	if client.Credentials.Secret == "" {
		return fmt.Errorf("token secret in the client is not set")
	}

	client.Credentials.AccessToken = nil

	url := endpointAuthenticateRefresh
	params := map[string]string{
		"refreshToken": client.Credentials.RefreshToken.Value,
	}
	if !client.Credentials.RefreshToken.IsSafe() {
		url = endpointAuthenticate
		params = map[string]string{
			"userIdentifier": client.Credentials.UserID,
			"userSecret":     client.Credentials.Secret,
		}

		client.Credentials.RefreshToken = nil
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

	_, _, err = client.Post(url, payload, &response, WithoutAuthentication)
	if err != nil {
		if _, ok := err.(NotAuthorizedError); ok {
			return err
		}

		return fmt.Errorf("authenticate request failed: %w", err)
	}

	client.Credentials.AccessToken = &Token{
		Value: response.AccessToken,
		ExpirationTime: time.Now().Add(
			response.ExpiresIn * time.Second,
		),
	}

	client.Credentials.RefreshToken = &Token{
		Value: response.RefreshToken,
		ExpirationTime: time.Now().Add(
			time.Second * response.RefreshExpiresIn,
		),
	}

	return nil
}
