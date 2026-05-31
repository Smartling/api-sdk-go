package mt

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Smartling/api-sdk-go/helpers/sm_client"
	"github.com/Smartling/api-sdk-go/helpers/sm_file"
	"github.com/Smartling/api-sdk-go/helpers/uid"
)

// FileTranslator defines file behaviour
type FileTranslator interface {
	Start(ctx context.Context, accountUID uid.AccountUID, fileUID uid.FileUID, p StartParams) (StartResponse, error)
	Progress(ctx context.Context, accountUID uid.AccountUID, fileUID uid.FileUID, mtUID uid.MtUID) (ProgressResponse, error)
}

// NewFileTranslator returns new FileTranslator implementation
func NewFileTranslator(client *smclient.Client) FileTranslator {
	return httpFileTranslator{base: newBase(client)}
}

type httpFileTranslator struct {
	base *base
}

type StartParams struct {
	SourceLocaleID  string   `json:"sourceLocaleId"`
	TargetLocaleIDs []string `json:"targetLocaleIds"`
}

// Start starts file translation
func (h httpFileTranslator) Start(ctx context.Context, accountUID uid.AccountUID, fileUID uid.FileUID, p StartParams) (StartResponse, error) {
	path := joinPath(mtBasePath, buildStartPath(accountUID, fileUID))

	payload, err := json.Marshal(p)
	if err != nil {
		return StartResponse{}, fmt.Errorf("failed to marshal start params: %w", err)
	}

	var response startResponse
	_, code, err := h.base.client.PostJSON(ctx, path, payload, &response.Response.Data)
	if err != nil {
		return StartResponse{}, fmt.Errorf("failed to start file translation: %w", err)
	}
	response.Response.Code = code
	return toStartResponse(response), nil
}

// Progress return progress of file translation
func (h httpFileTranslator) Progress(ctx context.Context, accountUID uid.AccountUID, fileUID uid.FileUID, mtUID uid.MtUID) (ProgressResponse, error) {
	path := joinPath(mtBasePath, buildProgressPath(accountUID, fileUID, mtUID))

	var response progressResponse
	_, code, err := h.base.client.GetJSON(
		ctx,
		path,
		smfile.FileURIRequest{FileURI: string(fileUID)}.GetQuery(),
		&response.Response.Data,
	)
	if err != nil {
		return ProgressResponse{}, fmt.Errorf("failed to get progress file translation: %w", err)
	}
	response.Response.Code = code
	return toProgressResponse(response), nil
}

func buildStartPath(accountUID uid.AccountUID, fileUID uid.FileUID) string {
	return "/accounts/" + string(accountUID) + "/files/" + string(fileUID) + "/mt"
}

func buildProgressPath(accountUID uid.AccountUID, fileUID uid.FileUID, mtUID uid.MtUID) string {
	return "/accounts/" + string(accountUID) + "/files/" + string(fileUID) + "/mt/" + string(mtUID) + "/status"
}
