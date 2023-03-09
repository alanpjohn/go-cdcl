package error

type SolverError struct {
	Message string
}

func (S SolverError) Error() string {
	return S.Message
}
