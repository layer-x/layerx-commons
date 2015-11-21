package lxerrors

import "fmt"

type tpiError struct {
	error
	message string
	err     error
}

func New(message string, err error) *tpiError {
	return &tpiError{
		message: message,
		err:     err,
	}
}

func (e *tpiError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: {%s}", e.message, e.err.Error())
	}
	return fmt.Sprintf("%s", e.message)
}
