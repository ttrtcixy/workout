package apperrors

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrServer UserErr = "server error"
)

func Wrap(op string, err error) error {
	if err == nil {
		return nil
	}
	var userErr UserError
	if errors.As(err, &userErr) {
		return err
	}
	return fmt.Errorf("%s: %w", op, err)
}

// UserErr used to create user errors
type UserErr string

func (e UserErr) Error() string {
	return string(e)
}

func (e UserErr) UserError() {}

// UserError used to check custom errors created using UserErr
type UserError interface {
	error
	UserError()
}

type ValidationError struct {
	Field string
	Err   error
}

func (e *ValidationError) Error() string {
	return e.Err.Error()
}
func (e *ValidationError) UserError() {}

type ValidationErrors []ValidationError

func (e *ValidationErrors) UserError() {}
func (e *ValidationErrors) Add(field, err string) {
	*e = append(*e, ValidationError{
		Field: field,
		Err:   errors.New(err),
	})

}
func (e *ValidationErrors) Error() string {
	var err strings.Builder

	for _, v := range *e {
		err.WriteString(v.Field)
		err.WriteString(": ")
		err.WriteString(v.Err.Error() + "; ")
	}
	return strings.TrimSpace(err.String())
}
