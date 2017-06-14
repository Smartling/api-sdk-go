package smartling

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	successResponseCode = "success"
)

// Post performs POST request to the Smartling API. You probably do not want
// to use it.
func (client *Client) Post(
	url string,
	body []byte,
	result interface{},
	options ...interface{},
) (json.RawMessage, int, error) {
	return client.requestJSON("POST", url, body, result, options...)
}

// GetJSON performs GET request to the smartling API and tries to decode answer
// as JSON.
func (client *Client) GetJSON(
	url string,
	params url.Values,
	result interface{},
	options ...interface{},
) (json.RawMessage, int, error) {
	if len(params) > 0 {
		url += "?" + params.Encode()
	}

	return client.requestJSON("GET", url, nil, result, options...)
}

// Get performs raw GET request to the Smartling API. You probably do not want
// to use it.
func (client *Client) Get(
	url string,
	params url.Values,
	options ...interface{},
) (io.ReadCloser, int, error) {
	if len(params) > 0 {
		url += "?" + params.Encode()
	}

	return client.request("GET", url, nil, options...)
}

func (client *Client) request(
	method string,
	url string,
	body []byte,
	options ...interface{},
) (io.ReadCloser, int, error) {
	authenticate := true

	for _, option := range options {
		switch value := option.(type) {
		case AuthenticationOption:
			authenticate = bool(value)
		}
	}

	if authenticate {
		err := client.Authenticate()
		if err != nil {
			return nil, 0, fmt.Errorf(
				"unable to authenticate: %s", err,
			)
		}
	}

	request, err := http.NewRequest(
		method,
		client.BaseURL+url,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, 0, fmt.Errorf("unable to create HTTP request: %s", err)
	}

	request.Header.Set("Content-Type", "application/json")

	if client.Credentials.AccessToken != nil {
		request.Header.Set(
			"Authorization",
			"Bearer "+client.Credentials.AccessToken.Value,
		)
	}

	reply, err := client.HTTP.Do(request)
	if err != nil {
		return nil, 0, fmt.Errorf("unable to perform HTTP request: %s", err)
	}

	return reply.Body, reply.StatusCode, nil
}

func (client *Client) requestJSON(
	method string,
	url string,
	body []byte,
	result interface{},
	options ...interface{},
) (json.RawMessage, int, error) {
	reply, code, err := client.request(method, url, body, options...)
	if err != nil {
		return nil, code, err
	}

	defer reply.Close()

	var response struct {
		Response struct {
			Code string
			Data json.RawMessage
		}
	}

	err = json.NewDecoder(reply).Decode(&response)
	if err != nil {
		return nil, 0, fmt.Errorf(
			"unable to decode JSON response: %s", err,
		)
	}

	if strings.ToLower(response.Response.Code) != "success" {
		return nil, 0, fmt.Errorf(
			`unexpected response status (expected "%s"): %#v`,
			successResponseCode,
			response.Response.Code,
		)
	}

	if result == nil {
		return response.Response.Data, code, nil
	}

	if code != 200 {
		return nil, code, fmt.Errorf(
			"API call returned non-200 status code: %s", err,
		)
	}

	err = json.Unmarshal(response.Response.Data, result)
	if err != nil {
		return nil, 0, fmt.Errorf(
			"unable to decode API response data: %s", err,
		)
	}

	return nil, code, nil
}
