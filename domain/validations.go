package domain

import (
	"fmt"
	"strings"
)

type ErrValidation struct {
	Field   string
	Message string
}

func (e *ErrValidation) Code() string {
	return "validation"
}
func (e *ErrValidation) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

func Required(field string, v string) error {
	if strings.TrimSpace(v) != "" {
		return nil
	}

	return &ErrValidation{
		Field:   field,
		Message: "must not be empty",
	}
}
