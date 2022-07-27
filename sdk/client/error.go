package client

import (
	"fmt"
)

type ErrNotFound struct {
	ErrorMessage string `json:"error_message"`
}

type ErrConflict struct {
	ErrorMessage string `json:"error_message"`
}

type ErrInternalServer struct {
	ErrorMessage string `json:"error_message"`
}

type ErrBadRequest struct {
	ErrorMessage string `json:"error_message"`
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("[Not Found]: %v", e.ErrorMessage)
}

func (e *ErrConflict) Error() string {
	return fmt.Sprintf("[Conflict]: %v", e.ErrorMessage)
}

func (e *ErrInternalServer) Error() string {
	return fmt.Sprintf("[Internal Server]: %v", e.ErrorMessage)
}

func (e *ErrBadRequest) Error() string {
	return fmt.Sprintf("[Bad Request]: %v", e.ErrorMessage)
}
