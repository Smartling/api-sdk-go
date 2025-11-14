package mt

import "strings"

type AccountUID string

func (a AccountUID) Validate() error {
	if strings.TrimSpace(string(a)) == "" {
		return ErrEmptyParam("AccountUID")
	}
	return nil
}

type FileUID string

func (f FileUID) Validate() error {
	if strings.TrimSpace(string(f)) == "" {
		return ErrEmptyParam("FileUID")
	}
	return nil
}

type MtUID string

func (m MtUID) Validate() error {
	if strings.TrimSpace(string(m)) == "" {
		return ErrEmptyParam("MtUID")
	}
	return nil
}
