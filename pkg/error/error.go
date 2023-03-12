// The dbg package contains utility structs and functions for error handling
package handler

// Default error stuct for all CDCL specific errors
type SolverError struct {
	Message string
	Err     error
}

func (S SolverError) Error() string {
	return S.Message
}

// Function to create SolveError
func Throw(msg string, err error) SolverError {
	return SolverError{Message: msg, Err: err}
}
