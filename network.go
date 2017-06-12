// methods that wrap around http communication

package smartling

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type SmartlingApiResponse struct {
	Response SmartlingResult
}

type SmartlingResult struct {
	Code string
	Data json.RawMessage
}

// parses smartling header wrapper
func unmarshalTransportHeader(bytes []byte) (apiResponse SmartlingApiResponse, err error) {

	// unmarshal transport header
	err = json.Unmarshal(bytes, &apiResponse)
	return
}

// helper function to actually do a post request
func (c *Client) doPostRequest(apiCall string, authHeader string, data []byte) (response []byte, statusCode int, err error) {

	request, err := http.NewRequest("POST", apiCall, bytes.NewBuffer(data))
	if err != nil {
		return
	}
	request.Header.Set("Content-Type", "application/json")
	if len(authHeader) > 0 {
		request.Header.Set("Authorization", authHeader)
	}

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	response, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	statusCode = resp.StatusCode
	return
}

// helper function to actually do a get request
func (c *Client) doGetRequest(apiCall string, authHeader string) (response []byte, statusCode int, err error) {

	request, err := http.NewRequest("GET", apiCall, nil)
	if err != nil {
		return
	}
	request.Header.Set("Content-Type", "application/json")
	if len(authHeader) > 0 {
		request.Header.Set("Authorization", authHeader)
	}

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	response, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	statusCode = resp.StatusCode
	return
}
