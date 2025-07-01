package mt

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Smartling/api-sdk-go/helpers/sm_client"
)

// TranslationControl defines translation control behaviour
type TranslationControl interface {
	CancelTranslation(accountUID AccountUID, fileUID FileUID, mtUID MtUID) (CancelTranslationResponse, error)
	DetectFileLanguage(accountUID AccountUID, fileUID FileUID) (DetectFileLanguageResponse, error)
	DetectionProgress(accountUID AccountUID, fileUID FileUID, languageDetectionUID string) (DetectionProgressResponse, error)
}

// NewTranslationControl returns new TranslationControl implementation
func NewTranslationControl(client *smclient.Client) TranslationControl {
	return httpTranslationControl{base: newBase(client)}
}

type httpTranslationControl struct {
	base *base
}

// CancelTranslation cancels translation
func (h httpTranslationControl) CancelTranslation(accountUID AccountUID, fileUID FileUID, mtUID MtUID) (CancelTranslationResponse, error) {
	var res CancelTranslationResponse
	startPath := buildCancelTranslationPath(accountUID, fileUID, mtUID)
	path := joinPath(mtBasePath, startPath)
	_, _, err := h.base.client.PostJSON(path, nil, &res, contentTypeApplicationJson)
	if err != nil {
		return CancelTranslationResponse{}, fmt.Errorf("failed to cancel file translation: %w", err)
	}
	return res, nil
}

// DetectFileLanguage detects file language
func (h httpTranslationControl) DetectFileLanguage(accountUID AccountUID, fileUID FileUID) (DetectFileLanguageResponse, error) {
	startPath := buildDetectFileLanguagePath(accountUID, fileUID)
	path := joinPath(mtBasePath, startPath)

	url := h.base.client.BaseURL + path
	request, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return DetectFileLanguageResponse{}, fmt.Errorf("failed to create request: %v", err)
	}

	request.Header.Set("Authorization", "Bearer "+h.base.client.Credentials.AccessToken.Value)

	response, err := h.base.client.HTTP.Do(request)
	if err != nil {
		return DetectFileLanguageResponse{}, fmt.Errorf("failed to detect file language: %w", err)
	}
	if response.StatusCode != http.StatusAccepted {
		return DetectFileLanguageResponse{}, fmt.Errorf("expected 202 status got: %d", response.StatusCode)
	}
	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)

	type detectFileLanguageResponse struct {
		Response struct {
			Code string `json:"code"`
			Data struct {
				LanguageDetectionUID string `json:"languageDetectionUid"`
			} `json:"data"`
		} `json:"response"`
	}
	var res detectFileLanguageResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return DetectFileLanguageResponse{}, fmt.Errorf("failed to unmarshal: %v", err)
	}
	return DetectFileLanguageResponse{
		Code:                 res.Response.Code,
		LanguageDetectionUID: res.Response.Data.LanguageDetectionUID,
	}, nil
}

// DetectionProgress returns info about detection
func (h httpTranslationControl) DetectionProgress(accountUID AccountUID, fileUID FileUID, languageDetectionUID string) (DetectionProgressResponse, error) {
	progressPath := buildDetectionProgressPath(accountUID, fileUID, languageDetectionUID)
	path := joinPath(mtBasePath, progressPath)

	url := h.base.client.BaseURL + path
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return DetectionProgressResponse{}, fmt.Errorf("failed to create request: %v", err)
	}
	request.Header.Set("Authorization", "Bearer "+h.base.client.Credentials.AccessToken.Value)

	response, err := h.base.client.HTTP.Do(request)
	if err != nil {
		return DetectionProgressResponse{}, fmt.Errorf("failed to detect file language: %w", err)
	}
	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)

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
	var res detectionProgressResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return DetectionProgressResponse{}, fmt.Errorf("failed to unmarshal: %v", err)
	}
	detectedSourceLanguages := make([]DetectedSourceLanguage, len(res.Response.Data.DetectedSourceLanguages))
	for i, v := range res.Response.Data.DetectedSourceLanguages {
		detectedSourceLanguages[i] = DetectedSourceLanguage{
			LanguageID:      v.LanguageID,
			DefaultLocaleID: v.DefaultLocaleID,
		}
	}
	var restErr string
	if res.Response.Data.Error != nil {
		restErr = *(res.Response.Data.Error)
	}
	return DetectionProgressResponse{
		Code:                    res.Response.Code,
		State:                   res.Response.Data.State,
		Error:                   restErr,
		DetectedSourceLanguages: detectedSourceLanguages,
	}, nil
}

func buildCancelTranslationPath(accountUID AccountUID, fileUID FileUID, mtUID MtUID) string {
	return "/accounts/" + string(accountUID) + "/files/" + string(fileUID) + "/mt/" + string(mtUID) + "/cancel"
}

func buildDetectFileLanguagePath(accountUID AccountUID, fileUID FileUID) string {
	return "/accounts/" + string(accountUID) + "/files/" + string(fileUID) + "/language-detection"
}

func buildDetectionProgressPath(accountUID AccountUID, fileUID FileUID, languageDetectionUID string) string {
	return "/accounts/" + string(accountUID) + "/files/" + string(fileUID) + "/language-detection/" + languageDetectionUID + "/status"
}
