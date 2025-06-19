package mt

import (
	"github.com/Smartling/api-sdk-go/helpers/sm_client"
)

type base struct {
	client *smclient.Client
}

func newBase(client *smclient.Client) *base {
	return &base{client: client}
}
