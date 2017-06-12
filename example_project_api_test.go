/*
	Smartling SDK v2 project api sample

	Sample shows usage of smartling project api
	http://docs.smartling.com/pages/API/v2/Projects/
	Replace userIdentifier and tokenSecret with your credentials

*/

package smartling_test

import (
	"log"
)

const (
	userIdentifier = "" // put your user identifier here
	tokenSecret    = "" // put your token secret here
	accountId      = "" // put your account id here
)

func main() {

	log.Printf("Initializing smartling client and performing autorization")
	client := NewClient(userIdentifier, tokenSecret)

	log.Printf("Listing projects for accountId %v:", accountId)

	listRequest := ProjectListRequest{
		ProjectNameFilter: "",
		IncludeArchived:   false,
	}

	projects, err := client.ListProjects(accountId, listRequest)
	if err != nil {
		log.Printf(err.Error())
		return
	}

	log.Printf("Found %v project(s) belonging to user account", projects.TotalCount)

	// now iterate over every project and request it's details
	for _, project := range projects.Items {
		projectDetails, err := client.ProjectDetails(project.ProjectId)
		if err != nil {
			log.Printf(err.Error())
			return
		}
		// print project details struct
		log.Printf("%+v", projectDetails)
	}
}
