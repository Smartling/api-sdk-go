package mt

import (
	"strings"

	"github.com/Smartling/api-sdk-go/helpers/sm_client"
)

const (
	successResponseCode = "success"
)

var (
	mtBasePath                 = "/file-translations-api/v2/"
	contentTypeApplicationJson = "Content-Type:application/json"
)

type base struct {
	client *smclient.Client
}

func newBase(client *smclient.Client) *base {
	return &base{client: client}
}

func joinPath(left, right string) string {
	return strings.TrimRight(left, "/") + "/" + strings.TrimLeft(right, "/")
}
