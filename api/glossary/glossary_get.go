package glossary

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/Smartling/api-sdk-go/helpers/uid"
)

// Get fetches a single glossary by its UID.
// Endpoint: GET /glossary-api/v3/accounts/{accountUid}/glossaries/{glossaryUid}.
func (h httpGlossary) Get(ctx context.Context, accountUID uid.AccountUID, glossaryUID string) (GetGlossaryResponse, error) {
	reqURL := path.Join(glossaryBasePath, url.PathEscape(string(accountUID)), "glossaries", url.PathEscape(glossaryUID))

	var response getGlossaryResponse
	_, code, err := h.client.GetJSON(ctx, reqURL, nil, &response.Response.Data)
	if err != nil && code == http.StatusNotFound {
		return GetGlossaryResponse{}, ErrGlossaryNotFound
	}
	if err != nil {
		return GetGlossaryResponse{}, fmt.Errorf("failed to get glossary: %w", err)
	}

	return toGetGlossaryResponse(response.Response.Data), nil
}

// GetByName lists glossaries for the given account, optionally filtered by name.
// Endpoint: POST /glossary-api/v3/accounts/{accountUid}/glossaries/search.
// The Smartling Glossary API exposes listing as a search call: the filter is
// sent in the JSON body rather than as query parameters. An empty name means
// "no name filter" and returns all glossaries.
func (h httpGlossary) GetByName(ctx context.Context, accountUID uid.AccountUID, name string) ([]GetGlossaryResponse, error) {
	body := struct {
		Query string `json:"query,omitempty"`
	}{Query: name}

	payload, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal glossary search request: %w", err)
	}

	reqURL := path.Join(glossaryBasePath, url.PathEscape(string(accountUID)), "glossaries", "search")

	var res getGlossariesResponse
	_, code, err := h.client.PostJSON(ctx, reqURL, payload, &res.Response.Data)
	if err != nil && code == http.StatusNotFound {
		return nil, ErrGlossaryNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to list glossaries: %w", err)
	}

	return toReadGlossariesResponse(res), nil
}
