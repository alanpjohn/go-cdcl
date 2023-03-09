package base

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
Status denotes the status of formula as implied by Formula.

Status can have the following values EMPTY, UNIT, DECIDE, SOLVED
*/
type Status uint

const (
	EMPTY  Status = 0 // An empty clause was found.
	UNIT              // A Unit Clause was found.
	DECIDE            // Decide operation must be performed.
	SOLVED            // All clauses are solved. This means that the current model satisfies the formula.
)

/*
Represents Disjunction of literals
*/
type Disjunction []Literal

/*
Clause Represents a clause in CNF in propositional logic
*/
type Clause struct {
	Literals Disjunction
	Size     uint
}

/*
Formula represents the formula of clauses to be solved by the solver.
*/
type Formula interface {
	Apply(l Literal)
	Undo(l Literal)
	Learn(c Clause)
	Propagate() (Status, Clause)
}

type Model interface {
	UnitPropagate(l Literal)
	Decide(l Literal)
	Backjump() (Status, []Literal) // Doubtful
	LemmaLearn() Clause            // Doubtful
	IsComplete() bool
	IsEmpty() bool
	Size() uint
}

type SATFile struct {
	AtomCount   uint          // No of atoms
	ClauseCount uint          // No of clauses
	Clauses     []Disjunction // Formula read from .SAT file
}

// Solver interface solves a given Formula
type Solver interface {
	solve(satFile SATFile) (Solution, error)
}

/*
Solution denotes the the solution found the solver.

Solution can have the following values SATISFIABLE, UNSATISFIABLE, UNKNOWN
*/
type Solution uint

const (
	SATISFIABLE   Solution = iota // A model exists that satifies the formula
	UNSATISFIABLE                 // No model exists that can satisfy the formula
	UNKNOWN                       // Solver could not find a model that satifies the formula and could not estabilish unsatisfiability due to some error
)
