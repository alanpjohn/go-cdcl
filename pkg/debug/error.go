package dbg

type SolverError struct {
	Message string
	Err     error
}

func (S SolverError) Error() string {
	return S.Message
}

func ThrowSolverError(msg string, err error) SolverError {
	return SolverError{Message: msg, Err: err}
}

