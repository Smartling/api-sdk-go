package glossary

// ReadGlossaryResponse is the public representation of a single glossary
// returned by the list/search endpoint.
type ReadGlossaryResponse struct {
	GlossaryUid string
	Name        string
	Description string
	LocaleIDs   []string
}

// readGlossaryResponse is the raw envelope returned by the list/search endpoint.
type readGlossaryResponse struct {
	Response struct {
		Code string `json:"code"`
		Data struct {
			TotalCount int                       `json:"totalCount"`
			Items      []readGlossaryResponseRow `json:"items"`
		} `json:"data"`
	} `json:"response"`
}

type readGlossaryResponseRow struct {
	GlossaryUid  string   `json:"glossaryUid"`
	GlossaryName string   `json:"glossaryName"`
	Description  string   `json:"description"`
	LocaleIDs    []string `json:"localeIds"`
}

func toReadGlossaryResponses(r readGlossaryResponse) []ReadGlossaryResponse {
	res := make([]ReadGlossaryResponse, len(r.Response.Data.Items))
	for i, item := range r.Response.Data.Items {
		res[i] = ReadGlossaryResponse{
			GlossaryUid: item.GlossaryUid,
			Name:        item.GlossaryName,
			Description: item.Description,
			LocaleIDs:   item.LocaleIDs,
		}
	}
	return res
}
