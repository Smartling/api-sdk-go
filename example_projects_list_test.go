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

// Smartling SDK v2 Project API Example
//
// This example lists projects and project details from specified account.
// Useful for testing your user identifier / token.
//
// `UserID` and `TokenSecret` should be specified in the
// example_credentails_test.go before running that test.
//
// `AccountID` should be specified in the
//┈example_credentails_test.go┈before┈running┈that┈test.

package smartling_test

import (
	"encoding/json"
	"fmt"
	"log"

	sdk "github.com/Smartling/api-sdk-go"
	"github.com/Smartling/api-sdk-go/helpers/sm_file"
)

func ExampleProjects_List() {
	log.Printf("Initializing smartling client and performing autorization")

	client := sdk.NewClient(UserID, TokenSecret)

	log.Printf("Listing projects for account ID %v:", AccountID)

	listRequest := smfile.ProjectsListRequest{
		Cursor:            smfile.LimitOffsetRequest{Limit: 10},
		ProjectNameFilter: "",
		IncludeArchived:   false,
	}

	projects, err := client.ListProjects(AccountID, listRequest)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Printf(
		"Found %v project(s) belonging to user account",
		projects.TotalCount,
	)

	for _, project := range projects.Items {
		projectDetails, err := client.GetProjectDetails(project.ProjectID)
		if err != nil {
			log.Fatal(err)
			return
		}

		data, err := json.MarshalIndent(projectDetails, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		log.Print(string(data))
	}

	fmt.Println("Projects List Successfull")

	// Output: Projects List Successfull
}
