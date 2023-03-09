package pkg

// Solver interface solves a given Formula
type Solver interface {
	solve(f Formula) Solution
}
