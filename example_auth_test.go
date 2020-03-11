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

// Smartling SDK v2 Auth Test Example.
//
// Example shows usage of Smartling authentication API
// https://help.smartling.com/v1.0/reference#authentication-1
//
// This example does nothing except the authentication call.
// Useful for testing your user identifier / token.
//
// `UserID` and `TokenSecret` should be specified in the
// example_credentails_test.go before running that test.

package smartling_test

import (
	"fmt"
	"log"

	smartling "github.com/Smartling/api-sdk-go"
)

func ExampleAuth() {
	log.Printf("Initializing smartling client and performing autorization")

	client := smartling.NewClient(UserID, TokenSecret)

	err := client.Authenticate()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Authentication Successfull")

	// Output: Authentication Successfull
}
