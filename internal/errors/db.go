package apperrors

import (
	"github.com/ttrtcixy/workout/internal/core/repository/query"
)

type DBError struct {
	Err   error
	query *query.Query
}

func (e DBError) Error() string {
	return ""
}
