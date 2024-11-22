package custom_errors

import "errors"

var (
	ErrNoRows              = errors.New("no rows in result set")
	ErrOffsetOutOfRange    = errors.New("offset exceeds the available number of verses")
	ErrInternalServerError = errors.New("internal server error")
)
