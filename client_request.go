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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	successResponseCode = "success"
	validationErrorCode = "validation_error"
)

// Post performs POST request to the Smartling API. You probably do not want
// to use it.
func (client *Client) Post(
	url string,
	payload []byte,
	result interface{},
	options ...interface{},
) (json.RawMessage, int, error) {
	return client.requestJSON("POST", url, nil, payload, result, options...)
}

// GetJSON performs GET request to the smartling API and tries to decode answer
// as JSON.
func (client *Client) GetJSON(
	url string,
	params url.Values,
	result interface{},
	options ...interface{},
) (json.RawMessage, int, error) {
	return client.requestJSON("GET", url, params, nil, result, options...)
}

// Get performs raw GET request to the Smartling API. You probably do not want
// to use it.
func (client *Client) Get(
	url string,
	params url.Values,
	options ...interface{},
) (io.ReadCloser, int, error) {
	reply, err := client.request("GET", url, params, nil, options...)
	if err != nil {
		return nil, 0, err
	}

	return reply.Body, reply.StatusCode, nil
}

func (client *Client) request(
	method string,
	url string,
	params url.Values,
	payload []byte,
	options ...interface{},
) (*http.Response, error) {
	var (
		authenticate = true
		contentType  = "application/json"
	)

	for _, option := range options {
		switch value := option.(type) {
		case AuthenticationOption:
			authenticate = bool(value)

		case ContentTypeOption:
			contentType = string(value)
		}
	}

	if authenticate {
		err := client.Authenticate()
		if err != nil {
			return nil, fmt.Errorf("unable to authenticate: %w", err)
		}
	}

	token := client.Credentials.AccessToken

	if payload != nil {
		if contentType == "application/json" {
			client.Logger.Debugf(
				"<- %s %s %s [payload %d bytes]\n%s",
				method, url, token, len(payload), payload,
			)
		} else {
			client.Logger.Debugf(
				"<- %s %s %s [payload %d bytes form data]",
				method, url, token, len(payload),
			)
		}
	} else {
		client.Logger.Debugf(
			"<- %s %s %s", method, url, token,
		)
	}

	startTime := time.Now()

	if len(params) > 0 {
		url += "?" + params.Encode()
	}

	request, err := http.NewRequest(
		method,
		client.BaseURL+url,
		bytes.NewBuffer(payload),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create HTTP request: %w", err)
	}

	request.Header.Set("Content-Type", contentType)
	request.Header.Set("User-Agent", client.UserAgent)

	if client.Credentials.AccessToken != nil {
		request.Header.Set("Authorization", "Bearer "+token.Value)
	}

	reply, err := client.HTTP.Do(request)
	if err != nil {
		return nil, fmt.Errorf("unable to perform HTTP request: %w", err)
	}

	requestID := reply.Header.Get("X-Sl-Requestid")
	if requestID == "" {
		requestID = "<none>"
	}

	client.Logger.Debugf(
		"-> %s [took %.2fs] [X-SL-RequestID %s]",
		reply.Status,
		time.Since(startTime).Seconds(),
		requestID,
	)

	return reply, nil
}

func (client *Client) requestJSON(
	method string,
	url string,
	params url.Values,
	payload []byte,
	result interface{},
	options ...interface{},
) (json.RawMessage, int, error) {
	reply, err := client.request(method, url, params, payload, options...)
	if err != nil {
		return nil, 0, err
	}

	defer reply.Body.Close()

	code := reply.StatusCode

	body, err := io.ReadAll(reply.Body)
	if err != nil {
		return nil, code, APIError{
			Cause:   err,
			URL:     url,
			Payload: payload,
		}
	}

	var response struct {
		Response struct {
			Code   string
			Data   json.RawMessage
			Errors []struct {
				Key     string
				Message string
			}
		}
	}

	contentType := reply.Header.Get("Content-Type")

	if strings.HasPrefix(contentType, "application/json") {
		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, code, JSONError{
				Cause:    err,
				Response: body,
			}
		}

		// we don't care about error here, it's only for logging
		message, _ := json.MarshalIndent(response, "", "  ")

		client.Logger.Debugf(
			"=> JSON [status=%s]\n%s",
			response.Response.Code,
			message,
		)
	}

	if code != 200 {
		switch code {
		case 200:
			fallthrough
		case 202:
			// ok

		case 401:
			return nil, code, NotAuthorizedError{}

		case 404:
			return nil, code, NotFoundError{}

		default:
			if strings.ToLower(response.Response.Code) == validationErrorCode {
				return nil, code, ValidationError{
					Errors: response.Response.Errors,
				}
			} else {
				return nil, code, APIError{
					Cause: fmt.Errorf(
						"API call returned unexpected HTTP code: %d", code,
					),
					URL:      url,
					Params:   params,
					Payload:  payload,
					Response: body,
					Headers:  &reply.Header,
				}

			}
		}
	}

	if strings.ToLower(response.Response.Code) != successResponseCode {
		return nil, 0, APIError{
			Cause: fmt.Errorf(
				`unexpected response status (expected "%s"): %#v`,
				successResponseCode,
				response.Response.Code,
			),
			URL:      url,
			Params:   params,
			Payload:  payload,
			Response: body,
			Headers:  &reply.Header,
		}
	}

	if result == nil {
		return response.Response.Data, code, nil
	}

	err = json.Unmarshal(response.Response.Data, result)
	if err != nil {
		return nil, 0, APIError{
			Cause:    fmt.Errorf("unable to decode API response data: %w", err),
			URL:      url,
			Params:   params,
			Payload:  payload,
			Response: body,
			Headers:  &reply.Header,
		}
	}

	return nil, code, nil
}
