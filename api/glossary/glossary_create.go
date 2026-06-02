package glossary

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"path"

	"github.com/Smartling/api-sdk-go/helpers/uid"
)

// Create creates a new glossary under the given account.
// Endpoint: POST /glossary-api/v3/accounts/{accountUid}/glossaries.
func (h httpGlossary) Create(ctx context.Context, accountUID uid.AccountUID, req CreateGlossaryRequest) (CreateGlossaryResponse, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return CreateGlossaryResponse{}, fmt.Errorf("marshal create glossary request: %w", err)
	}

	reqURL := path.Join(glossaryBasePath, url.PathEscape(string(accountUID)), "glossaries")

	var response createGlossaryResponse
	_, code, err := h.client.PostJSON(ctx, reqURL, payload, &response.Response.Data)
	if err != nil {
		return CreateGlossaryResponse{}, fmt.Errorf("failed to create glossary: %w", err)
	}

	return toCreateGlossaryResponse(response, code), nil
}
