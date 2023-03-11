package types

// Represents an Atom in propositional logic as an unsigned int as done in SAT files
type Atom uint

// Represents a Literal in propositional logic as an signed int
type Literal int

// Negates the Literal
func (L Literal) Negate() Literal {
	return -L
}

func (L Literal) Atom() Atom {
	if L >= 0 {
		return Atom(uint(L))
	} else {
		return Atom(uint(-L))
	}
}

/*
Disjunction Represents a clause in CNF in propositional logic
*/
type Disjunction []Literal

type SATFile struct {
	AtomCount   uint          // No of atoms
	ClauseCount uint          // No of clauses
	Clauses     []Disjunction // Formula read from .SAT file
}

type ClauseType uint

const (
	EMPTY_CLAUSE ClauseType = iota
	UNIT_CLAUSE
	DECISION_CLAUSE
	SOLVED_CLAUSE
)

type Clause interface {
	Type() ClauseType
	Apply(l Literal) Clause
	Undo(l Literal) Clause
	Reset() Clause
	Contains(l Literal) bool
	IsSolved() bool
	IsLearnt() bool
	Disjunction() Disjunction
	Original() Disjunction
}

/*
Formula represents the formula of clauses to be solved by the solver.
*/
type Formula interface {
	NextClause() Clause
	Assign(l Literal) Formula
	Unassign(l Literal) Formula
	Learn(c Clause) Formula
	Restart() Formula
}

// Solver interface solves a given Formula
type Solver interface {
	Solve() (Solution, error)
}

/*
Solution denotes the the solution found the solver.

Solution can have the following values SATISFIABLE, UNSATISFIABLE, UNKNOWN
*/
type Solution uint

const (
	SATISFIABLE   Solution = iota // A model exists that satifies the formula
	UNSATISFIABLE                 // No model exists that can satisfy the formula
	PROGRESS                      // The Solver is currently working on the solution
	UNKNOWN                       // Solver could not find a model that satifies the formula and could not estabilish unsatisfiability due to some error
)

func (s Solution) String() string {
	switch s {
	case 0:
		return "SATISFIABLE"
	case 1:
		return "UNSATISFIABLE"
	}
	return "UNKNOWN"
}
