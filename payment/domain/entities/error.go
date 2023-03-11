package entities

type Error struct {
	Code    string
	Message string
}

func (e *Error) Error() string {
	return e.Code
}
