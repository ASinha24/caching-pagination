package api

import "fmt"

type MartError struct {
	Code        ErrorCode `json:"code,omitempty"`
	Message     string    `json:"message,omitempty"`
	Description string    `json:"description,omitempty"`
}

func (m MartError) Error() string {
	return fmt.Sprintf("code: %d msg: %s", m.Code, m.Description)
}
