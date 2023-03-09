package pkg

// Represents an Atom in propositional logic as a number as done in SAT files
type Atom uint

// Represents a Literal in propositional logic as an Atom along with it's sign
type Literal struct {
	Atom Atom
	Sign bool
}

// Negates the Literal
func (L *Literal) Negate() {
	L.Sign = !L.Sign
}

/*
Status denotes the status of formula as implied by Formula.

Status can have the following values EMPTY, UNIT, DECIDE, SOLVED
*/
type Status uint

const (
	EMPTY Status = 0	// An empty clause was found.
	UNIT				// A Unit Clause was found.
	DECIDE 				// Decide operation must be performed.
	SOLVED				// All clauses are solved. This means that the current model satisfies the formula.
)

/*
Clause Represents a Disjuction of literals are represented in Clausal Form in propositional logic
*/
type Clause struct {
	Literals []Literal
	Size     uint
}

/*
Solution denotes the the solution found the solver.

Solution can have the following values SATISFIABLE, UNSATISFIABLE, UNKNOWN
*/
type Solution uint

const (
	SATISFIABLE Solution = iota	// A model exists that satifies the formula
	UNSATISFIABLE				// No model exists that can satisfy the formula
	UNKNOWN						// Solver could not find a model that satifies the formula and could not estabilish unsatisfiability due to some error
)
