package commonmodel

type Error struct {
	code    string
	message string
}

func (e Error) Error() string {
	return e.message
}

func (e Error) Code() string {
	return e.code
}

// Message allows you to set a custom message for the error.
func (e Error) Message(msg string) Error {
	return Error{
		code:    e.code,
		message: msg,
	}
}

func NewError(code, message string) Error {
	return Error{
		code:    code,
		message: message,
	}
}
