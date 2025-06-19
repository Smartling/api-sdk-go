package mt

import (
	"fmt"
	"io"

	"github.com/Smartling/api-sdk-go/helpers/sm_client"
	"github.com/Smartling/api-sdk-go/helpers/sm_file"
)

// Downloader defines downloader behaviour
type Downloader interface {
	File(accountUID AccountUID, fileUID FileUID,
		mtUID MtUID, localeID string) (io.Reader, error)
	Batch(accountUID AccountUID, fileUID FileUID, mtUID MtUID) (io.Reader, error)
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
func (d httpDownloader) File(accountUID AccountUID, fileUID FileUID,
	mtUID MtUID, localeID string) (io.Reader, error) {
	path := buildFilePath(accountUID, fileUID, mtUID, localeID)
	url, err := joinURL(mtBasePath, path)
	if err != nil {
		return nil, err
	}
	reader, _, err := d.base.client.Get(
		url,
		smfile.FileURIRequest{FileURI: string(fileUID)}.GetQuery(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to download original file: %w", err)
	}

	return reader, nil
}

// Batch download files
func (d httpDownloader) Batch(accountUID AccountUID, fileUID FileUID, mtUID MtUID) (io.Reader, error) {
	path := buildBatchPath(accountUID, fileUID, mtUID)
	url, err := joinURL(mtBasePath, path)
	if err != nil {
		return nil, err
	}
	reader, _, err := d.base.client.Get(
		url,
		smfile.FileURIRequest{FileURI: string(fileUID)}.GetQuery(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to download original file: %w", err)
	}

	return reader, nil
}

func buildFilePath(accountUID AccountUID, fileUID FileUID, mtUID MtUID, localeID string) string {
	return string(accountUID) + "/files/" + string(fileUID) + "/mt/" + string(mtUID) + "/locales/" + localeID + "/file"
}

func buildBatchPath(accountUID AccountUID, fileUID FileUID, mtUID MtUID) string {
	return string(accountUID) + "/files/" + string(fileUID) + "/mt/" + string(mtUID) + "/locales/all/file/zip"
}
