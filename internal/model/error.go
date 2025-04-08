package model

import (
	"fmt"
)

type ErrorWithCode struct {
	Code string
	Msg  string
	Err  error // optional wrapped error
}

func (e *ErrorWithCode) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Msg)
}

func (e *ErrorWithCode) Unwrap() error {
	return e.Err
}

var (
	ErrWrongCredentials = &ErrorWithCode{
		Code: "wrong_credentials",
		Msg:  "Wrong credentials",
	}
)
