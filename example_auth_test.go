/*
	Smartling SDK v2 auth test example

	This sample does nothing except the authentication call.
	Useful for testing your user identifier / token.
	Replace userIdentifier and tokenSecret with your credentials
*/

package smartling_test

import (
	"log"
)

const (
	userIdentifier = "" // put your user identifier here
	tokenSecret    = "" // put your token secret here
)

func main() {

	log.Printf("Initializing smartling client and performing autorization")
	client := NewClient(userIdentifier, tokenSecret)

	err := client.AuthenticationTest()
	if err != nil {
		log.Printf(err.Error())
		return
	}

	log.Printf("Authentication is successful, your credentials are valid!")
}
