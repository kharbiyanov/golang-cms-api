package utils

import (
	"github.com/graphql-go/graphql/gqlerrors"
)

type GQError struct {
	Err  error `json:"error"`
	Code int64 `json:"code"`
}

func (e GQError) Error() string {
	return e.Err.Error()
}

func (e GQError) ToFormattedError(err gqlerrors.FormattedError) gqlerrors.FormattedError {
	return gqlerrors.FormattedError{
		Message:   err.Error(),
		Locations: err.Locations,
		Path:      err.Path,
		Extensions: map[string]interface{}{
			"code": 123,
		},
	}
}
