package game

import "errors"

var (
	ErrorNotFound  = errors.New("not found")
	ErrorIllegal   = errors.New("illegal action")
	ErrorNotReady  = errors.New("not ready")
	ErrorInvalid   = errors.New("invalid")
	ErrorAmbiguous = errors.New("duplicate")
)
