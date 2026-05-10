package sqlcraft

import "errors"

var (
	ErrEmptyValues      = errors.New("empty values in query")
	ErrEmptyColumns     = errors.New("empty columns in query")
	ErrMissMatchValues  = errors.New("miss match values for given columns")
	ErrInvalidOperator  = errors.New("invalid dafi operator")
	ErrInvalidFieldName = errors.New("invalid field name")
)
