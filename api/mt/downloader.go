package mt

import (
	"github.com/Smartling/api-sdk-go/helpers/sm_client"
)

// Downloader defines downloader behaviour
type Downloader interface {
	File(url, dest string) error
	Batch(urls []string, destDir string) error
}

// NewDownloader returns new Downloader implementation
func NewDownloader(client *smclient.Client) Downloader {
	return newHttpDownloader(newBase(client))
}

type httpDownloader struct {
	base *base
}

func newHttpDownloader(base *base) *httpDownloader {
	return &httpDownloader{base: base}
}

// File downloads file
func (d httpDownloader) File(url, dest string) error {
	//TODO implement me
	panic("implement me")
}

// Batch download files
func (d httpDownloader) Batch(urls []string, destDir string) error {
	//TODO implement me
	panic("implement me")
}
