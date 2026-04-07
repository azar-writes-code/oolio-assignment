package unit

import (
	"errors"
	"net/http"
	"testing"

	"github.com/azar-writes-code/oolio-products-backend/pkg/server/rest/utils/apperrors"
)

func TestAppError_Error(t *testing.T) {
	err := apperrors.NewInternal("main error", errors.New("detail error"))
	// Note: NewInternal formats message as "An internal server error occurred"
	expected := "An internal server error occurred: detail error"
	if err.Error() != expected {
		t.Errorf("Expected %q, got %q", expected, err.Error())
	}

	err2 := apperrors.NewNotFound("product %d not found", 123)
	expected2 := "product 123 not found"
	if err2.Error() != expected2 {
		t.Errorf("Expected %q, got %q", expected2, err2.Error())
	}
}

func TestMapToHTTP(t *testing.T) {
	tests := []struct {
		err          error
		expectedCode int
	}{
		{apperrors.NewNotFound("not found"), http.StatusNotFound},
		{apperrors.NewBadRequest("bad request"), http.StatusBadRequest},
		{apperrors.NewConflict("conflict"), http.StatusConflict},
		{apperrors.NewUnauthorized("unauthorized"), http.StatusUnauthorized},
		{apperrors.NewUnprocessable("unprocessable"), http.StatusUnprocessableEntity},
		{apperrors.NewInternal("internal", errors.New("some error")), http.StatusInternalServerError},
		{errors.New("generic error"), http.StatusInternalServerError},
	}

	for _, tt := range tests {
		code, resp := apperrors.MapToHTTP(tt.err)
		if code != tt.expectedCode {
			t.Errorf("For error %v, expected code %d, got %d", tt.err, tt.expectedCode, code)
		}
		if _, ok := resp["error"]; !ok {
			t.Errorf("Expected 'error' field in response for error %v", tt.err)
		}
	}
}
