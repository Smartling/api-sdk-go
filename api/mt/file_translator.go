package mt

import smclient "github.com/Smartling/api-sdk-go/helpers/sm_client"

// FileTranslator defines file behaviour
type FileTranslator interface {
	Start() (StartResponse, error)
	Progress() (ProgressResponse, error)
}

// NewFileTranslator returns new FileTranslator implementation
func NewFileTranslator(client *smclient.Client) FileTranslator {
	return httpFileTranslator{base: newBase(client)}
}

type httpFileTranslator struct {
	base *base
}

// Start starts file translation
func (h httpFileTranslator) Start() (StartResponse, error) {
	//TODO implement me
	panic("implement me")
}

// Progress return progress of file translation
func (h httpFileTranslator) Progress() (ProgressResponse, error) {
	//TODO implement me
	panic("implement me")
}
