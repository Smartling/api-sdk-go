package mt

import (
	"errors"
	"fmt"
)

// ErrEmptyParam creates a new empty param error
func ErrEmptyParam(name string) error {
	return &errEmptyParam{
		paramName: name,
	}
}

// IsErrEmptyParam checks if error is errEmptyParam
func IsErrEmptyParam(err error) bool {
	if err == nil {
		return false
	}
	_, ok := errors.AsType[*errEmptyParam](err)
	return ok
}

type errEmptyParam struct {
	paramName string
}

// Error returns the string representation of the error
func (e errEmptyParam) Error() string {
	return fmt.Sprintf("parameter `%s` cannot be empty", e.paramName)
}
