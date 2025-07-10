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

// Simple smartling sdk usage example

package smartling_test

import (
	"log"
	"time"

	sdk "github.com/Smartling/api-sdk-go"
	"github.com/Smartling/api-sdk-go/helpers/sm_file"
)

func ExampleBasic() {
	const (
		userID      = ""
		TokenSecret = ""
		accountId   = ""
		projectId   = ""
	)

	log.Printf("Initializing sdk client and performing autorization")

	client := sdk.NewClient(userID, TokenSecret)

	log.Printf("Listing projects:")

	listRequest := smfile.ProjectsListRequest{
		Cursor:            smfile.LimitOffsetRequest{Limit: 10},
		ProjectNameFilter: "VCS",
		IncludeArchived:   false,
	}

	projects, err := client.ListProjects(accountId, listRequest)
	if err != nil {
		log.Printf("%v", err.Error())
		return
	}
	log.Printf("Success")

	log.Printf("Projects belonging to user account:")
	log.Printf("%+v", projects)

	projectDetails, err := client.GetProjectDetails(projectId)
	if err != nil {
		log.Printf("%v", err.Error())
		return
	}
	log.Printf("Success")
	log.Printf("Projects details are")
	log.Printf("%+v", projectDetails)

	for {
		// sleep 6 minutes to issue reauth call
		time.Sleep(time.Minute * 6)
		_, err = client.ListProjects(accountId, listRequest)
		if err != nil {
			log.Printf("%v", err.Error())
			return
		}
		log.Printf("Success")
	}
}
