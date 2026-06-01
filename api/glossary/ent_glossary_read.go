package glossary

// GetGlossaryResponse is the public representation of a single glossary.
type GetGlossaryResponse struct {
	GlossaryUID string
	Name        string
	Description string
	LocaleIDs   []string
}

// getGlossaryResponse is the raw envelope returned by the single-glossary
// read endpoint (GET .../glossaries/{glossaryUid}); data is the glossary object.
type getGlossaryResponse struct {
	Response struct {
		Code string                  `json:"code"`
		Data getGlossaryDataResponse `json:"data"`
	} `json:"response"`
}

// getGlossariesResponse is the raw envelope returned by the list/search endpoint.
type getGlossariesResponse struct {
	Response struct {
		Code string `json:"code"`
		Data struct {
			TotalCount int                       `json:"totalCount"`
			Items      []getGlossaryDataResponse `json:"items"`
		} `json:"data"`
	} `json:"response"`
}

type getGlossaryDataResponse struct {
	GlossaryUid  string   `json:"glossaryUid"`
	GlossaryName string   `json:"glossaryName"`
	Description  string   `json:"description"`
	LocaleIDs    []string `json:"localeIds"`
}

func toGetGlossaryResponse(row getGlossaryDataResponse) GetGlossaryResponse {
	return GetGlossaryResponse{
		GlossaryUID: row.GlossaryUid,
		Name:        row.GlossaryName,
		Description: row.Description,
		LocaleIDs:   row.LocaleIDs,
	}
}

func toReadGlossariesResponse(r getGlossariesResponse) []GetGlossaryResponse {
	res := make([]GetGlossaryResponse, len(r.Response.Data.Items))
	for i, item := range r.Response.Data.Items {
		res[i] = toGetGlossaryResponse(item)
	}
	return res
}
