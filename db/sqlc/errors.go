package db

import (
	"errors"

	"github.com/lib/pq"
)

const (
	ForeignKeyViolation = "23503"
	UniqueViolation     = "23505"
	NoDataFound         = "P0002"
)

var ErrUniqueViolation = &pq.Error{
	Code: UniqueViolation,
}

var ErrRecordNotFound = &pq.Error{
	Code: NoDataFound,
}

func ErrorCode(err error) string {
	var pgErr *pq.Error
	if errors.As(err, &pgErr) {
		return pgErr.Code.Name()
	}
	return ""
}
