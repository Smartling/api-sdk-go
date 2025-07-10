package mt

// CancelTranslationResponse defines cancel translation response
type CancelTranslationResponse struct {
	Code string
}

type cancelTranslationResponse struct {
	Response struct {
		Code string `json:"code"`
		Data struct {
		} `json:"data"`
	} `json:"response"`
}

func toCancelTranslationResponse(r cancelTranslationResponse) CancelTranslationResponse {
	return CancelTranslationResponse{
		Code: r.Response.Code,
	}
}

// DetectFileLanguageResponse defines detect file language response
type DetectFileLanguageResponse struct {
	Code                 string
	LanguageDetectionUID string
}

type detectFileLanguageResponse struct {
	Response struct {
		Code string `json:"code"`
		Data struct {
			LanguageDetectionUID string `json:"languageDetectionUid"`
		} `json:"data"`
	} `json:"response"`
}

func toDetectFileLanguageResponse(r detectFileLanguageResponse) DetectFileLanguageResponse {
	return DetectFileLanguageResponse{
		Code:                 r.Response.Code,
		LanguageDetectionUID: r.Response.Data.LanguageDetectionUID,
	}
}

// DetectionProgressResponse defines detection progress response
type DetectionProgressResponse struct {
	Code                    string
	State                   string
	Error                   string
	DetectedSourceLanguages []DetectedSourceLanguage
}

// DetectedSourceLanguage defines detected source language
type DetectedSourceLanguage struct {
	LanguageID      string
	DefaultLocaleID string
}
type detectionProgressResponse struct {
	Response struct {
		Code string `json:"code"`
		Data struct {
			State                   string  `json:"state"`
			Error                   *string `json:"error"`
			DetectedSourceLanguages []struct {
				LanguageID      string `json:"languageId"`
				DefaultLocaleID string `json:"defaultLocaleId"`
			} `json:"detectedSourceLanguages"`
		} `json:"data"`
	} `json:"response"`
}

func toDetectionProgressResponse(r detectionProgressResponse) DetectionProgressResponse {
	var err string
	if r.Response.Data.Error != nil {
		err = *r.Response.Data.Error
	}
	detectedSourceLanguages := make([]DetectedSourceLanguage, len(r.Response.Data.DetectedSourceLanguages))
	for key, val := range r.Response.Data.DetectedSourceLanguages {
		detectedSourceLanguages[key] = DetectedSourceLanguage{
			LanguageID:      val.LanguageID,
			DefaultLocaleID: val.DefaultLocaleID,
		}
	}
	return DetectionProgressResponse{
		Code:                    r.Response.Code,
		State:                   r.Response.Data.State,
		Error:                   err,
		DetectedSourceLanguages: detectedSourceLanguages,
	}
}
