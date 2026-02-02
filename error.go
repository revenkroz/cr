package cr

import "fmt"

type Violation struct {
	Field string `json:"field"`
	Error string `json:"message"`
}

type Error struct {
	Code       int         `json:"code"`
	Message    string      `json:"message"`
	Violations []Violation `json:"violations,omitempty"`
}

func (e Error) Error() string {
	return fmt.Sprintf("code: %d, message: %s, violations: %v", e.Code, e.Message, e.Violations)
}

func NewError(message string, code int) Error {
	return Error{
		Code:    code,
		Message: message,
	}
}

func NewParseError() Error {
	return NewError("parse error", 422)
}

func NewValidationError(violations []Violation) Error {
	return Error{
		Code:       400,
		Message:    "validation error",
		Violations: violations,
	}
}

func NewNotFoundError() Error {
	return NewError("method not found", 404)
}

func NewOtherError(message string) Error {
	return NewError(message, -1)
}
