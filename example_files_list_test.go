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

// Smartling SDK v2 Files API Example.

// Example shows usage of Smartling files API
// https://help.smartling.com/v1.0/reference#upload
//
// `UserID` and `TokenSecret` should be specified in the
// example_credentails_test.go before running that test.
//
// `AccountID` should be specified in the
//┈example_credentails_test.go┈before┈running┈that┈test.
//
// `ProjectID` should be specified in the
//┈example_credentails_test.go┈before┈running┈that┈test.

package smartling_test

import (
	"fmt"
	"log"
	"time"

	smartling "github.com/Smartling/api-sdk-go"
)

func ExampleListFiles() {
	log.Printf("Initializing smartling client and performing autorization")
	client := smartling.NewClient(UserID, TokenSecret)

	log.Printf("Listing project (%v) files", ProjectID)

	log.Println("Listing constraints:")
	log.Println("\tUriMask: 'master' (file URI must contain 'master' substring)")
	log.Println("\tLastUploadedBefore: one month back")
	log.Println("\tFileTypes: json and Java properties files")

	listRequest := smartling.FilesListRequest{
		Cursor:             smartling.LimitOffsetRequest{Limit: 10},
		URIMask:            "master",
		LastUploadedBefore: smartling.UTC{time.Now().AddDate(0, -1, 0)},
		FileTypes:          []smartling.FileType{smartling.FileTypeJSON, smartling.FileTypeJavaProperties},
	}

	files, err := client.ListFiles(ProjectID, listRequest)
	if err != nil {
		log.Printf("%v", err.Error())
		return
	}
	log.Println("Success!")

	log.Printf("Found %v files", files.TotalCount)

	for _, f := range files.Items {
		log.Printf("%v", f.FileURI)
	}

	fmt.Println("Files List Successfull")

	// Output: Files List Successfull
}
