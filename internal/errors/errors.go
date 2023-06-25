package errors

import "fmt"

type InvalidIDError struct {
	id interface{}
}

func (e InvalidIDError) Error() string {
	return fmt.Sprintf("Invalid Id : %v", e.id)
}

func NewInvalidIDError(id interface{}) InvalidIDError {
	return InvalidIDError{id}
}
