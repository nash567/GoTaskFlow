package model

type NoFieldsToUpdateError struct{}

func (e NoFieldsToUpdateError) Error() string {
	return "No fields to update."
}

func NewNoFieldsToUpdateError() NoFieldsToUpdateError {
	return NoFieldsToUpdateError{}
}
