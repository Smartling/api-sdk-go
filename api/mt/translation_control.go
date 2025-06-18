package mt

type TranslationControl interface {
	CancelTranslation() error
	DetectFileLanguage() error
	DetectionProgress() error
}
