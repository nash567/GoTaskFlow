package errors

type AddUserError struct{}

func (e AddUserError) Error() string {
	return "error adding user"
}

func NewAddUserError() AddUserError {
	return AddUserError{}
}
