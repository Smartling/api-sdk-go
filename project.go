package smartling

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
)

// API endpoints
const (
	projectApiList    = "/accounts-api/v2/accounts/%v/projects"
	projectApiDetails = "/projects-api/v2/projects/%v"
)

// list project api call response data
type ProjectsList struct {
	TotalCount int64
	Items      []Project
}

type Project struct {
	ProjectId               string
	ProjectName             string
	AccountUid              string
	SourceLocaleId          string
	SourceLocaleDescription string
	Archived                bool
}

type ProjectDetails struct {
	Project
	TargetLocales []Locale
}

type Locale struct {
	LocaleId    string
	Description string
}

type ProjectListRequest struct {
	ProjectNameFilter string
	IncludeArchived   bool
	Limit             int64
	Offset            int64
}

func (c *Client) ListProjects(accountId string, listRequest ProjectListRequest) (list ProjectsList, err error) {

	header, err := c.auth.AccessHeader(c)
	if err != nil {
		return
	}

	// prepare the url
	urlObject, err := url.Parse(c.baseUrl + fmt.Sprintf(projectApiList, accountId))
	if err != nil {
		return
	}
	urlObject.RawQuery = listRequest.RawQuery()

	bytes, statusCode, err := c.doGetRequest(urlObject.String(), header)
	if err != nil {
		return
	}

	if statusCode != 200 {
		err = fmt.Errorf("ListProjects call returned unexpected status code: %v", statusCode)
		return
	}

	// unmarshal transport header
	apiResponse, err := unmarshalTransportHeader(bytes)
	if err != nil {
		return
	}

	// unmarshal projects array
	err = json.Unmarshal(apiResponse.Response.Data, &list)
	if err != nil {
		// TODO: special error here
		return
	}

	log.Printf("List proijects - received %v status code", statusCode)

	return
}

func (c *Client) ProjectDetails(projectId string) (projectDetails ProjectDetails, err error) {

	header, err := c.auth.AccessHeader(c)
	if err != nil {
		return
	}

	bytes, statusCode, err := c.doGetRequest(c.baseUrl+fmt.Sprintf(projectApiDetails, projectId), header)
	if err != nil {
		return
	}

	if statusCode != 200 {
		err = fmt.Errorf("ProjectDetailss call returned unexpected status code: %v", statusCode)
		return
	}

	// unmarshal transport header
	apiResponse, err := unmarshalTransportHeader(bytes)
	if err != nil {
		return
	}

	log.Printf("Project Details - received %v status code", statusCode)

	// unmarshal project details
	err = json.Unmarshal(apiResponse.Response.Data, &projectDetails)
	if err != nil {
		return
	}

	return
}
