package errorHandler

type AppError struct {
	Code    int
	Message string
}

func (e AppError) Error() string {
	return e.Message
}

func NewAppError(code int, msg string) error {
	return AppError{
		Code:    code,
		Message: msg,
	}
}

var (
	ErrBadRequest   = NewAppError(400, "Bad request")
	ErrUnauthorized = NewAppError(401, "Unauthorized")
	ErrForbidden    = NewAppError(403, "Forbidden")
	ErrConflict     = NewAppError(409, "Conflict")
	ErrInternal     = NewAppError(500, "Internal server error")
)
