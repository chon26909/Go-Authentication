package util

import "strings"

type IError struct {
	Error string `json:"error"`
}

func NewError(err error) IError {
	jerr := IError{"generic error"}
	if err != nil {
		jerr.Error = err.Error()
	}
	return jerr
}

func NormalizeEmail(email string) string {
	return strings.TrimSpace(strings.ToLower(email))
}
