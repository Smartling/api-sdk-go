// Simple smartling sdk usage example

package smartling_test

import (
	"log"
	"time"

	smartling "github.com/Smartling/api-sdk-go"
)

func ExampleBasic() {
	const (
		userIdentifier = ""
		TokenSecret    = ""
		accountId      = ""
		projectId      = ""
	)

	log.Printf("Initializing smartling client and performing autorization")

	client := smartling.NewClient(userIdentifier, TokenSecret)

	log.Printf("Listing projects:")

	listRequest := smartling.ProjectListRequest{
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

	projectDetails, err := client.ProjectDetails(projectId)
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
