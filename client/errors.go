package client

import (
	"errors"
	"fmt"
)

// ErrInvalidID is for when an id is provided but it is invalid
var ErrInvalidID = errors.New("invalid id provided; not a valid bson")

// ErrorResponse error response
type ErrorResponse struct { //nolint:errname // Tech debt, rename errors at a later date
	Code      string `json:"code"`
	ErrorCode string `json:"error"`
	RequestID string `json:"requestID"`
}

func (e *ErrorResponse) Error() string {
	errorCode := e.Code

	if e.ErrorCode != "" {
		errorCode = e.ErrorCode
	}

	return fmt.Sprintf("[Error] %s -- RequestId: %s", errorCode, e.RequestID)
}

type ErrNotFound struct{} //nolint:errname // Tech debt, rename errors at a later date

func (e *ErrNotFound) Error() string {
	return "No results found/returned"
}
