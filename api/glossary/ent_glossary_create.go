package glossary

import "time"

// CreateGlossaryRequest is the JSON body for a glossary-create call.
type CreateGlossaryRequest struct {
	GlossaryName     string           `json:"glossaryName"`
	Description      string           `json:"description"`
	VerificationMode bool             `json:"verificationMode"`
	LocaleIDs        []string         `json:"localeIds"`
	FallbackLocales  []FallbackLocale `json:"fallbackLocales"`
}

type FallbackLocale struct {
	FallbackLocaleID string   `json:"fallbackLocaleId"`
	LocaleIDs        []string `json:"localeIds"`
}

// CreateGlossaryResponse is the public result of a glossary-create call.
type CreateGlossaryResponse struct {
	Code         int
	GlossaryUID  string
	AccountUID   string
	GlossaryName string
}

// createGlossaryResponse holds the raw JSON envelope returned by the API.
type createGlossaryResponse struct {
	Response struct {
		Code string `json:"code"`
		Data struct {
			GlossaryUID       string           `json:"glossaryUid"`
			AccountUID        string           `json:"accountUid"`
			GlossaryName      string           `json:"glossaryName"`
			Description       string           `json:"description"`
			VerificationMode  bool             `json:"verificationMode"`
			Archived          bool             `json:"archived"`
			CreatedByUserUID  string           `json:"createdByUserUid"`
			ModifiedByUserUID string           `json:"modifiedByUserUid"`
			CreatedDate       time.Time        `json:"createdDate"`
			ModifiedDate      time.Time        `json:"modifiedDate"`
			LocaleIDs         []string         `json:"localeIds"`
			FallbackLocales   []FallbackLocale `json:"fallbackLocales"`
		} `json:"data"`
	} `json:"response"`
}

func toCreateGlossaryResponse(g createGlossaryResponse, code int) CreateGlossaryResponse {
	return CreateGlossaryResponse{
		Code:         code,
		GlossaryUID:  g.Response.Data.GlossaryUID,
		AccountUID:   g.Response.Data.AccountUID,
		GlossaryName: g.Response.Data.GlossaryName,
	}
}
