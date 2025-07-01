package mt

import (
	"encoding/json"
	"fmt"

	"github.com/Smartling/api-sdk-go/helpers/sm_client"
	"github.com/Smartling/api-sdk-go/helpers/sm_file"
)

// FileTranslator defines file behaviour
type FileTranslator interface {
	Start(accountUID AccountUID, fileUID FileUID, p StartParams) (StartResponse, error)
	Progress(accountUID AccountUID, fileUID FileUID, mtUID MtUID) (ProgressResponse, error)
}

// NewFileTranslator returns new FileTranslator implementation
func NewFileTranslator(client *smclient.Client) FileTranslator {
	return httpFileTranslator{base: newBase(client)}
}

type httpFileTranslator struct {
	base *base
}

type StartParams struct {
	SourceLocaleIO  string   `json:"sourceLocaleId"`
	TargetLocaleIDs []string `json:"targetLocaleIds"`
}

// Start starts file translation
func (h httpFileTranslator) Start(accountUID AccountUID, fileUID FileUID, p StartParams) (StartResponse, error) {
	startPath := buildStartPath(accountUID, fileUID)
	path := joinPath(mtBasePath, startPath)

	payload, err := json.Marshal(p)
	if err != nil {
		return StartResponse{}, err
	}

	resp, err := h.base.client.Post(path, payload)
	if err != nil {
		return StartResponse{}, fmt.Errorf("failed to start file translation: %w", err)
	}
	type startResponse struct {
		Response struct {
			Code string `json:"code"`
			Data struct {
				MtUID string `json:"mtUid"`
			} `json:"data"`
		} `json:"response"`
	}
	var res startResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return StartResponse{}, fmt.Errorf("failed to parse response: %w", err)
	}
	return StartResponse{
		Code:  res.Response.Code,
		MtUID: MtUID(res.Response.Data.MtUID),
	}, nil
}

// Progress return progress of file translation
func (h httpFileTranslator) Progress(accountUID AccountUID, fileUID FileUID, mtUID MtUID) (ProgressResponse, error) {
	var res ProgressResponse
	progressPath := buildProgressPath(accountUID, fileUID, mtUID)
	path := joinPath(mtBasePath, progressPath)
	_, _, err := h.base.client.GetJSON(
		path,
		smfile.FileURIRequest{FileURI: string(fileUID)}.GetQuery(),
		&res,
	)
	if err != nil {
		return ProgressResponse{}, fmt.Errorf("failed to get progress file translation: %w", err)
	}
	return res, nil
}

func buildStartPath(accountUID AccountUID, fileUID FileUID) string {
	return "/accounts/" + string(accountUID) + "/files/" + string(fileUID) + "/mt"
}

func buildProgressPath(accountUID AccountUID, fileUID FileUID, mtUID MtUID) string {
	return "/accounts/" + string(accountUID) + "/files/" + string(fileUID) + "/mt/" + string(mtUID) + "/status"
}
