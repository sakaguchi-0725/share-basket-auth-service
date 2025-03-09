package apperr

import "fmt"

type ErrorCode int

const (
	ErrBadRequest ErrorCode = iota + 1
	ErrUnauthorized
	ErrNotFound
)

func (code ErrorCode) String() string {
	switch code {
	case ErrBadRequest:
		return "BadRequest"
	case ErrUnauthorized:
		return "Unauthorized"
	case ErrNotFound:
		return "NotFound"
	default:
		return fmt.Sprintf("Unknown: %d", code)
	}
}
