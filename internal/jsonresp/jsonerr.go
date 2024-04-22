package jsonresp

import "fmt"

type JSONError struct {
	Type       string `json:"error"`
	Payload    any    `json:"payload"`
	StatusCode int    `json:"-"`
	Err        error  `json:"-"`
}

func (err *JSONError) Error() string {
	return fmt.Sprintf("JsonErr: %v", err.Err.Error())
}

func (err *JSONError) Unwrap() error {
	return err.Err
}
