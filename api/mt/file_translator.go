package mt

type FileTranslator interface {
	Start() error
	Progress() error
}
