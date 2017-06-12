package smartling

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	return client.request("POST", url, body, result, options...)
}

// Get performs GET request to the Smartling API. You probably do not want
// to use it.
func (client *Client) Get(
	url string,
	params url.Values,
	result interface{},
	options ...interface{},
) (json.RawMessage, int, error) {
	if len(params) > 0 {
		url += "?" + params.Encode()
	}

	return client.request("GET", url, nil, result, options...)
}

func (client *Client) request(
	method string,
	url string,
	body []byte,
	result interface{},
	options ...interface{},
) (json.RawMessage, int, error) {
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

	defer reply.Body.Close()

	data, err := ioutil.ReadAll(reply.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("unable to read HTTP response: %s", err)
	}

	var response struct {
		Response struct {
			Code string
			Data json.RawMessage
		}
	}

	err = json.Unmarshal(data, &response)
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
		return response.Response.Data, reply.StatusCode, nil
	}

	if reply.StatusCode != 200 {
		return nil, reply.StatusCode, fmt.Errorf(
			"API call returned non-200 status code: %s", err,
		)
	}

	err = json.Unmarshal(response.Response.Data, result)
	if err != nil {
		return nil, 0, fmt.Errorf(
			"unable to decode API response data: %s", err,
		)
	}

	return nil, reply.StatusCode, nil
}
