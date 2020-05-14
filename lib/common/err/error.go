package err

import "errors"

const (
	//ErrInternalServer type internal server error
	ErrInternalServer = 0
	//ErrBadRequest error type bad request
	ErrBadRequest = 1
)

var (
	ErrCommon = errors.New("A failure occurred")
)
