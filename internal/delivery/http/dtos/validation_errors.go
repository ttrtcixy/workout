package dtos

import (
	"encoding/json"
	"log"
)

type ValidationErrors struct {
	FieldErrors []*FieldError `json:"errors,omitempty"`
	ListErrors  []*ListError  `json:"errors-list,omitempty"`
}

type FieldError struct {
	FieldName string `json:"field"`
	Error     string `json:"error"`
}

type ListError struct {
	Index       int           `json:"index"`
	FieldErrors []*FieldError `json:"errors"`
}

func (le *ListError) Error() string {
	const op = "ListError.Error"
	str, err := json.Marshal(le)
	if err != nil {
		log.Printf("op: %s, err: %v", op, err)
		return "server error"
	}

	return string(str)
}

func (le *ListError) AddFieldError(field string, err error) {
	le.FieldErrors = append(le.FieldErrors, &FieldError{
		FieldName: field,
		Error:     err.Error(),
	})
}

func (e *ValidationErrors) AddErrToList(err *ListError) {
	e.ListErrors = append(e.ListErrors, err)
}

func (e *ValidationErrors) AddFieldError(field string, err error) {
	e.FieldErrors = append(e.FieldErrors, &FieldError{
		FieldName: field,
		Error:     err.Error(),
	})
}

func (e *ValidationErrors) Error() string {
	const op = "ValidationErrors.Error"
	str, err := json.Marshal(e)
	if err != nil {
		log.Printf("op: %s, err: %v", op, err)
		return "server error"
	}

	return string(str)
}
