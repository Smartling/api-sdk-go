package mt

import (
	"net/url"
	"strings"

	"github.com/Smartling/api-sdk-go/helpers/sm_client"
)

var mtBasePath = "https://api.smartling.com/file-translations-api/v2/accounts/"

type base struct {
	client *smclient.Client
}

func newBase(client *smclient.Client) *base {
	return &base{client: client}
}

func joinURL(base, path string) (string, error) {
	u, err := url.Parse(base)
	if err != nil {
		return "", err
	}

	// Ensure the path is appended properly
	u.Path = strings.TrimRight(u.Path, "/") + "/" + strings.TrimLeft(path, "/")
	return u.String(), nil
}
