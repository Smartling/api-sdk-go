package glossary

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
)

// Get fetches a single glossary by its UID.
// Endpoint: GET /glossary-api/v3/accounts/{accountUid}/glossaries/{glossaryUid}.
func (h httpGlossary) Get(ctx context.Context, accountUID, glossaryUID string) (ReadGlossaryResponse, error) {
	reqURL := path.Join(glossaryBasePath, url.PathEscape(accountUID), "glossaries", url.PathEscape(glossaryUID))

	var row readGlossaryResponseRow
	_, code, err := h.client.GetJSON(ctx, reqURL, nil, &row)
	if err != nil && code == http.StatusNotFound {
		return ReadGlossaryResponse{}, ErrGlossaryNotFound
	}
	if err != nil {
		return ReadGlossaryResponse{}, fmt.Errorf("failed to get glossary: %w", err)
	}

	return ReadGlossaryResponse{
		GlossaryUid: row.GlossaryUid,
		Name:        row.GlossaryName,
		Description: row.Description,
		LocaleIDs:   row.LocaleIDs,
	}, nil
}

// GetByName lists glossaries for the given account, optionally filtered by name.
// Endpoint: POST /glossary-api/v3/accounts/{accountUid}/glossaries/search.
// The Smartling Glossary API exposes listing as a search call: the filter is
// sent in the JSON body rather than as query parameters. An empty name means
// "no name filter" and returns all glossaries.
func (h httpGlossary) GetByName(ctx context.Context, accountUID, name string) ([]ReadGlossaryResponse, error) {
	body := struct {
		Query string `json:"query,omitempty"`
	}{Query: name}

	payload, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal glossary search request: %w", err)
	}

	reqURL := path.Join(glossaryBasePath, url.PathEscape(accountUID), "glossaries", "search")

	var res readGlossaryResponse
	_, code, err := h.client.PostJSON(ctx, reqURL, payload, &res.Response.Data)
	if err != nil && code == http.StatusNotFound {
		return nil, ErrGlossaryNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to list glossaries: %w", err)
	}

	return toReadGlossaryResponses(res), nil
}
