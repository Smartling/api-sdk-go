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

	"github.com/Smartling/api-sdk-go/helpers/sm_error"
	"github.com/Smartling/api-sdk-go/helpers/sm_file"
)

const (
	endpointProjectsList   = "/accounts-api/v2/accounts/%v/projects"
	endpointProjectDetails = "/projects-api/v2/projects/%v"
)

// ProjectsList represents projects list under specified account.
type ProjectsList struct {
	// TotalCount represents total count of projects.
	TotalCount int64

	// Items contains projects list by specified request.
	Items []Project
}

// Project represents detailed project information.
type Project struct {
	// ProjectID is a unique project ID.
	ProjectID string

	// ProjectName is a human-friendly project name.
	ProjectName string

	// AccountUID is undocumented by Smartling API.
	AccountUID string

	// SourceLocaleID represents source locale ID for project.
	SourceLocaleID string

	// SourceLocaleDescription describes project's locale.
	SourceLocaleDescription string

	// Archived will be true if project is archived.
	Archived bool
}

// ProjectDetails extends Project type to contain target locales list.
type ProjectDetails struct {
	Project

	// TargetLocales represents target locales list.
	TargetLocales []Locale
}

// Locale represents locale for translation.
type Locale struct {
	// LocaleID is a unique locale ID.
	LocaleID string

	// Description describes locale.
	Description string

	// Enabled is a flag that represents is locale enabled or not.
	Enabled bool
}

// ListProjects returns projects in specified account matching specified
// request.
func (c *Client) ListProjects(
	accountID string,
	request smfile.ProjectsListRequest,
) (*ProjectsList, error) {
	var list ProjectsList

	_, _, err := c.Client.GetJSON(
		fmt.Sprintf(endpointProjectsList, accountID),
		request.GetQuery(),
		&list,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get projects list: %w", err)
	}

	return &list, nil
}

// GetProjectDetails returns project details for specified project.
func (c *Client) GetProjectDetails(
	projectID string,
) (*ProjectDetails, error) {
	var details ProjectDetails

	_, _, err := c.Client.GetJSON(
		fmt.Sprintf(endpointProjectDetails, projectID),
		nil,
		&details,
	)
	if err != nil {
		if _, ok := err.(smerror.NotFoundError); ok {
			return nil, err
		}

		return nil, fmt.Errorf("failed to get project details: %w", err)
	}

	return &details, nil
}
