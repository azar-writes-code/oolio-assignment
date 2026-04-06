package apperrors

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Type string

const (
	NotFound     Type = "NOT_FOUND"
	BadRequest   Type = "BAD_REQUEST"
	Internal     Type = "INTERNAL"
	Conflict     Type = "CONFLICT"
	Unauthorized Type = "UNAUTHORIZED"
	Unprocessable Type = "UNPROCESSABLE"
)

type AppError struct {
	Type    Type
	Message string
	Detail  error
}

func (e *AppError) Error() string {
	if e.Detail != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Detail)
	}
	return e.Message
}

func NewNotFound(msg string, args ...any) error {
	return &AppError{Type: NotFound, Message: fmt.Sprintf(msg, args...)}
}

func NewBadRequest(msg string, args ...any) error {
	return &AppError{Type: BadRequest, Message: fmt.Sprintf(msg, args...)}
}

func NewInternal(msg string, err error) error {
	return &AppError{Type: Internal, Message: "An internal server error occurred", Detail: err}
}

func NewConflict(msg string, args ...any) error {
	return &AppError{Type: Conflict, Message: fmt.Sprintf(msg, args...)}
}

func NewUnauthorized(msg string, args ...any) error {
	return &AppError{Type: Unauthorized, Message: fmt.Sprintf(msg, args...)}
}

func NewUnprocessable(msg string, args ...any) error {
	return &AppError{Type: Unprocessable, Message: fmt.Sprintf(msg, args...)}
}

// MapToHTTP translates an AppError to an HTTP status code and response body.
func MapToHTTP(err error) (int, gin.H) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		switch appErr.Type {
		case NotFound:
			return http.StatusNotFound, gin.H{"error": "Not Found", "message": appErr.Message}
		case BadRequest:
			return http.StatusBadRequest, gin.H{"error": "Bad Request", "message": appErr.Message}
		case Conflict:
			return http.StatusConflict, gin.H{"error": "Conflict", "message": appErr.Message}
		case Unauthorized:
			return http.StatusUnauthorized, gin.H{"error": "Unauthorized", "message": appErr.Message}
		case Unprocessable:
			return http.StatusUnprocessableEntity, gin.H{"error": "Unprocessable Entity", "message": appErr.Message}
		case Internal:
			return http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": "An internal error occurred"}
		}
	}

	// Default to 500 for unknown errors
	return http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": "An unexpected error occurred"}
}
