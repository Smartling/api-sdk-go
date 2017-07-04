package smartling

import (
	"encoding/json"
	"fmt"
)

const (
	endpointImport = "/files-api/v2/projects/%s/locales/%s/file/import"
)

// Import imports specified file as translation.
func (client *Client) Import(
	projectID string,
	localeID string,
	request ImportRequest,
) (json.RawMessage, error) {
	var result json.RawMessage

	form, err := request.GetForm()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create import form: %s",
			err,
		)
	}

	err = form.Close()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to close import file form: %s", err,
		)
	}

	_, _, err = client.Post(
		fmt.Sprintf(endpointImport, projectID, localeID),
		form.Bytes(),
		&result,
		ContentTypeOption(form.GetContentType()),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to import: %s",
			err,
		)
	}

	return result, nil
}
