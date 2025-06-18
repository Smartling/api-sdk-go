package mt

import (
	"github.com/Smartling/api-sdk-go/helpers/sm_client"
)

type Downloader interface {
	File(url, dest string) error
	Batch(urls []string, destDir string) error
}

func NewDownloader(client *smclient.Client) Downloader {
	return httpDownloader{}
}

type httpDownloader struct {
	base base
}

func (d httpDownloader) File(url, dest string) error {
	//TODO implement me
	panic("implement me")
}

func (d httpDownloader) Batch(urls []string, destDir string) error {
	//TODO implement me
	panic("implement me")
}
