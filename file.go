package smartling

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
)

// API endpoints
const (
	fileApiList      = "/files-api/v2/projects/%v/files/list"
	fileApiListTypes = "/files-api/v2/projects/%v/file-types"
)

type FileType string

const (
	Android        FileType = "android"
	Ios            FileType = "ios"
	Gettext        FileType = "gettext"
	Html           FileType = "html"
	JavaProperties FileType = "javaProperties"
	Yaml           FileType = "yaml"
	Xliff          FileType = "xliff"
	Xml            FileType = "xml"
	Json           FileType = "json"
	Docx           FileType = "docx"
	Pptx           FileType = "pptx"
	Xlsx           FileType = "xlsx"
	Idml           FileType = "idml"
	Qt             FileType = "qt"
	Resx           FileType = "resx"
	Plaintext      FileType = "plaintext"
	Csv            FileType = "csv"
	Stringsdict    FileType = "stringsdict"
)

type FileStatus struct {
	FileUri         string
	LastUploaded    *Iso8601Time
	FileType        FileType
	HasInstructions bool
}

// list files api call response data
type FilesList struct {
	TotalCount int64
	Items      []FileStatus
}

type FileListRequest struct {
	UriMask            string
	FileTypes          []FileType
	LastUploadedAfter  Iso8601Time
	LastUploadedBefore Iso8601Time
	Offset             int64
	Limit              int64
}

// list files for a project
func (c *Client) ListFiles(projectId string, listRequest FileListRequest) (list FilesList, err error) {

	header, err := c.auth.AccessHeader(c)
	if err != nil {
		return
	}

	// prepare the url
	urlObject, err := url.Parse(c.baseUrl + fmt.Sprintf(fileApiList, projectId))
	if err != nil {
		return
	}
	urlObject.RawQuery = listRequest.RawQuery()

	log.Printf(urlObject.String())
	bytes, statusCode, err := c.doGetRequest(urlObject.String(), header)
	if err != nil {
		return
	}

	if statusCode != 200 {
		err = fmt.Errorf("ListFiles call returned unexpected status code: %v", statusCode)
		return
	}

	// unmarshal transport header
	apiResponse, err := unmarshalTransportHeader(bytes)
	if err != nil {
		return
	}

	// unmarshal file items array
	err = json.Unmarshal(apiResponse.Response.Data, &list)
	if err != nil {
		// TODO: special error here
		return
	}

	log.Printf("List files - received %v status code", statusCode)

	return
}

// returns all file types currently represented in a project
func (c *Client) ListFileTypes(projectId string) (list []FileType, err error) {

	header, err := c.auth.AccessHeader(c)
	if err != nil {
		return
	}

	// prepare the url
	urlObject, err := url.Parse(c.baseUrl + fmt.Sprintf(fileApiListTypes, projectId))
	if err != nil {
		return
	}

	bytes, statusCode, err := c.doGetRequest(urlObject.String(), header)
	if err != nil {
		return
	}

	if statusCode != 200 {
		err = fmt.Errorf("ListFileTypes call returned unexpected status code: %v", statusCode)
		return
	}

	// unmarshal transport header
	apiResponse, err := unmarshalTransportHeader(bytes)
	if err != nil {
		return
	}

	// declare a transport struct
	type TypeList struct {
		Items []FileType
	}
	typeList := TypeList{}

	// unmarshal file types array
	err = json.Unmarshal(apiResponse.Response.Data, &typeList)
	if err != nil {
		// TODO: special error here
		return
	}

	list = typeList.Items

	log.Printf("List file types - received %v status code", statusCode)

	return
}
