package mt

import (
	"context"
	"fmt"

	"github.com/Smartling/api-sdk-go/helpers/sm_client"
)

// TranslationControl defines translation control behaviour
type TranslationControl interface {
	CancelTranslation(ctx context.Context, accountUID AccountUID, fileUID FileUID, mtUID MtUID) (CancelTranslationResponse, error)
	DetectFileLanguage(ctx context.Context, accountUID AccountUID, fileUID FileUID) (DetectFileLanguageResponse, error)
	DetectionProgress(ctx context.Context, accountUID AccountUID, fileUID FileUID, languageDetectionUID string) (DetectionProgressResponse, error)
}

// NewTranslationControl returns new TranslationControl implementation
func NewTranslationControl(client *smclient.Client) TranslationControl {
	return httpTranslationControl{base: newBase(client)}
}

type httpTranslationControl struct {
	base *base
}

// CancelTranslation cancels translation
func (h httpTranslationControl) CancelTranslation(ctx context.Context, accountUID AccountUID, fileUID FileUID, mtUID MtUID) (CancelTranslationResponse, error) {
	path := joinPath(mtBasePath, buildCancelTranslationPath(accountUID, fileUID, mtUID))
	var response cancelTranslationResponse
	_, code, err := h.base.client.PostJSON(ctx, path, nil, &response.Response.Data)
	if err != nil {
		return CancelTranslationResponse{}, fmt.Errorf("failed to cancel file translation: %w", err)
	}
	response.Response.Code = code
	return toCancelTranslationResponse(response), nil
}

// DetectFileLanguage detects file language
func (h httpTranslationControl) DetectFileLanguage(ctx context.Context, accountUID AccountUID, fileUID FileUID) (DetectFileLanguageResponse, error) {
	path := joinPath(mtBasePath, buildDetectFileLanguagePath(accountUID, fileUID))
	var response detectFileLanguageResponse
	_, code, err := h.base.client.PostJSON(ctx, path, nil, &response.Response.Data)
	if err != nil {
		return DetectFileLanguageResponse{}, fmt.Errorf("failed to detect file language: %w", err)
	}
	response.Response.Code = code
	return toDetectFileLanguageResponse(response), nil
}

// DetectionProgress returns info about detection
func (h httpTranslationControl) DetectionProgress(ctx context.Context, accountUID AccountUID, fileUID FileUID, languageDetectionUID string) (DetectionProgressResponse, error) {
	path := joinPath(mtBasePath, buildDetectionProgressPath(accountUID, fileUID, languageDetectionUID))
	var response detectionProgressResponse
	_, code, err := h.base.client.GetJSON(ctx, path, nil, &response.Response.Data)
	if err != nil {
		return DetectionProgressResponse{}, fmt.Errorf("failed to get detection progress: %w", err)
	}
	response.Response.Code = code
	return toDetectionProgressResponse(response), nil
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
