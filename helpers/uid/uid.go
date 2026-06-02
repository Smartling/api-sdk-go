package uid

import (
	"strings"

	"github.com/Smartling/api-sdk-go/helpers/sm_error"
)

type AccountUID string

func (a AccountUID) Validate() error {
	if strings.TrimSpace(string(a)) == "" {
		return smerror.ErrEmptyParam("AccountUID")
	}
	return nil
}

type FileUID string

func (f FileUID) Validate() error {
	if strings.TrimSpace(string(f)) == "" {
		return smerror.ErrEmptyParam("FileUID")
	}
	return nil
}

type MtUID string

func (m MtUID) Validate() error {
	if strings.TrimSpace(string(m)) == "" {
		return smerror.ErrEmptyParam("MtUID")
	}
	return nil
}
