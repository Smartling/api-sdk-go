package mt

import smclient "github.com/Smartling/api-sdk-go/helpers/sm_client"

// TranslationControl defines translation control behaviour
type TranslationControl interface {
	CancelTranslation() (CancelTranslationResponse, error)
	DetectFileLanguage() (DetectFileLanguageResponse, error)
	DetectionProgress() (DetectionProgressResponse, error)
}

// NewTranslationControl returns new TranslationControl implementation
func NewTranslationControl(client *smclient.Client) TranslationControl {
	return httpTranslationControl{base: newBase(client)}
}

type httpTranslationControl struct {
	base *base
}

func (h httpTranslationControl) CancelTranslation() (CancelTranslationResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (h httpTranslationControl) DetectFileLanguage() (DetectFileLanguageResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (h httpTranslationControl) DetectionProgress() (DetectionProgressResponse, error) {
	//TODO implement me
	panic("implement me")
}
