package types

// Represents an Atom in propositional logic as an unsigned int as done in SAT files.
type Atom uint

// Represents a Literal in propositional logic as an signed int.
type Literal int

// Negates the Literal
func (L Literal) Negate() Literal {
	return -L
}

// Returns the Atom of the Literal.
func (L Literal) Atom() Atom {
	if L >= 0 {
		return Atom(uint(L))
	} else {
		return Atom(uint(-L))
	}
}

/*
Disjunction Represents a clause in CNF in propositional logic.
*/
type Disjunction []Literal

type SATFile struct {
	AtomCount   uint          // No of atoms
	ClauseCount uint          // No of clauses
	Clauses     []Disjunction // Formula read from .SAT file
}

/*
ClauseType is an enum defining the types of clauses that the solver
must cater for.
*/
type ClauseType uint

const (
	EMPTY_CLAUSE    ClauseType = iota // Represents an empty clause
	UNIT_CLAUSE                       // Represents a unit clause
	DECISION_CLAUSE                   // Any clause that is neither empty or unit
	SOLVED_CLAUSE                     // A clause that is true in the given model
)

/*
Clause is a wrapper interface around Disjunction that gives more functionality
*/
type Clause interface {
	Type() ClauseType         // Get type of Clause
	Apply(l Literal) Clause   // Assert Literal to clause
	Undo(l Literal) Clause    // Undo the assertion made by literal (used for backjumping)
	Reset() Clause            // Reset model
	Contains(l Literal) bool  // Returns true if literal is present in Clause
	IsSolved() bool           // Returns true if clause is true in given model
	IsLearnt() bool           // Returns true if clause is a learnt clause
	Disjunction() Disjunction // The state of Clause as the solver progresses
	Original() Disjunction    // The original Dijunction stored as a backup
}

/*
Formula represents the formula of clauses to be solved by the solver.
*/
type Formula interface {
	/*
		 Selects the next clause in the formula for the solver to consider the formula has a preference order in which it returns clauses

			- Clauses which are false in the given model, i.e empty clauses are selected first
			- Unit clauses are selected second
			- Decision clauses are selected third
			- Solved clauses are left for last.
	*/
	NextClause() Clause
	Assign(l Literal) Formula   // Asserts the Literal to all clauses in Formula
	Unassign(l Literal) Formula // Undos the assertion by literal for all clauses in Formula
	Learn(c Clause) Formula     // Learn new clause after conflict
	Restart() Formula           // Reset all clauses in Formula but remembering the learnt clauses
	Print() string              // Prints all clauses for debugging purpose
}

// Solver interface solves a given Formula
type Solver interface {
	Solve() (Solution, error) // Solve
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
