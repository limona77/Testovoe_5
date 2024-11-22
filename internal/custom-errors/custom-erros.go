package custom_errors

import "errors"

var (
	ErrNoRows              = errors.New("no rows in result set")
	ErrOffsetOutOfRange    = errors.New("offset exceeds the available number of verses")
	ErrInternalServerError = errors.New("internal server error")
	ErrNoSongInfo          = errors.New("failed to fetch song information from external API")
	ErrAlreadyExists       = errors.New("song already exists")
)
