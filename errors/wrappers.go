package errors

import "errors"

// As and Is are aliases for errors.As and errors.Is. Gives us access to the errors package functions.
var (
	As = errors.As
	Is = errors.Is
)
