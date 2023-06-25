package errors

type AddTaskError struct{}

func (e AddTaskError) Error() string {
	return "error adding task"
}

func NewAddTaskError() AddTaskError {
	return AddTaskError{}
}
