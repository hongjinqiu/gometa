package error

import (
)

type BusinessError struct{
	Message string
	Code int
}

func (e *BusinessError) Error() string {
	return e.Message
}
