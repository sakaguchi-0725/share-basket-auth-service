package apperr

type ApplicationError interface {
	Code() ErrorCode
	Message() string
	Error() string
}

type applicationError struct {
	code    ErrorCode
	message string
	err     error
}

func (a *applicationError) Code() ErrorCode {
	return a.code
}

func (a *applicationError) Error() string {
	return a.err.Error()
}

func (a *applicationError) Message() string {
	return a.message
}

func NewApplicationError(code ErrorCode, msg string, err error) ApplicationError {
	return &applicationError{code, msg, err}
}

func NewInvalidInputError(err error) ApplicationError {
	return NewApplicationError(ErrBadRequest, ErrInvalidData.Error(), err)
}
