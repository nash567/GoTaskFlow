package errors

type AddTaskError struct{}

func (e AddTaskError) Error() string {
	return "error adding task"
}

func NewAddTaskError() AddTaskError {
	return AddTaskError{}
}

type AddTaskStepError struct{}

func (e AddTaskStepError) Error() string {
	return "error adding task step"
}

func NewAddTaskStepError() AddTaskStepError {
	return AddTaskStepError{}
}

type UpdateTaskStepError struct{}

func (e UpdateTaskStepError) Error() string {
	return "error updating task step"
}

func NewUpdateTaskStepError() UpdateTaskStepError {
	return UpdateTaskStepError{}
}

type FailedToInitialiseTransactionError struct{}

func (e FailedToInitialiseTransactionError) Error() string {
	return "failed to initialize transaction"
}

func NewFailedToInitialiseTransactionError() FailedToInitialiseTransactionError {
	return FailedToInitialiseTransactionError{}
}
