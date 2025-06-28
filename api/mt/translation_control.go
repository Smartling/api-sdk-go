package mt

import (
	"fmt"

	"github.com/Smartling/api-sdk-go/helpers/sm_client"
	"github.com/Smartling/api-sdk-go/helpers/sm_file"
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
	var res DetectFileLanguageResponse
	startPath := buildDetectFileLanguagePath(accountUID, fileUID)
	path := joinPath(mtBasePath, startPath)
	_, _, err := h.base.client.PostJSON(path, nil, &res, contentTypeApplicationJson)
	if err != nil {
		return DetectFileLanguageResponse{}, fmt.Errorf("failed to detect file language: %w", err)
	}
	return res, nil
}

// DetectionProgress returns info about detection
func (h httpTranslationControl) DetectionProgress(accountUID AccountUID, fileUID FileUID, languageDetectionUID string) (DetectionProgressResponse, error) {
	var res DetectionProgressResponse
	progressPath := buildDetectionProgressPath(accountUID, fileUID, languageDetectionUID)
	path := joinPath(mtBasePath, progressPath)
	_, _, err := h.base.client.GetJSON(
		path,
		smfile.FileURIRequest{FileURI: string(fileUID)}.GetQuery(),
		&res,
	)
	if err != nil {
		return DetectionProgressResponse{}, fmt.Errorf("failed to get detection process: %w", err)
	}
	return res, nil
}

func buildCancelTranslationPath(accountUID AccountUID, fileUID FileUID, mtUID MtUID) string {
	return "/accounts/" + string(accountUID) + "/files/" + string(fileUID) + "/mt/" + string(mtUID) + "/cancel"
}

func buildDetectFileLanguagePath(accountUID AccountUID, fileUID FileUID) string {
	return "/accounts/" + string(accountUID) + "/files/" + string(fileUID) + "/language-detection"
}

func buildDetectionProgressPath(accountUID AccountUID, fileUID FileUID, languageDetectionUID string) string {
	return "/accounts/" + string(accountUID) + "/files/" + string(fileUID) + "/mt/" + languageDetectionUID + "/status"
}
