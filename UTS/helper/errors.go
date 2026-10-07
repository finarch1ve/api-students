package helper

import "net/http"

type AppError struct {
	Status  int
	Message string
	Errors  map[string][]string
}

func (e *AppError) Error() string { return e.Message }

func NewAppError(status int, message string) *AppError {
	return &AppError{Status: status, Message: message}
}

func Unauthorized(msg string) *AppError  { return NewAppError(http.StatusUnauthorized, msg) }
func Forbidden(msg string) *AppError     { return NewAppError(http.StatusForbidden, msg) }
func NotFound(msg string) *AppError      { return NewAppError(http.StatusNotFound, msg) }
func Conflict(msg string) *AppError      { return NewAppError(http.StatusConflict, msg) }
func Unprocessable(msg string) *AppError { return NewAppError(http.StatusUnprocessableEntity, msg) }

func ValidationFailed(errs map[string][]string) *AppError {
	return &AppError{
		Status:  http.StatusUnprocessableEntity,
		Message: "Validasi gagal",
		Errors:  errs,
	}
}