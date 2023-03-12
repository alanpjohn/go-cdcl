/*
The solver package contains the BaseCDCLSolver
*/
package solver

import (
	"fmt"

	handler "github.com/alanpjohn/go-cdcl/pkg/error"
	logger "github.com/alanpjohn/go-cdcl/pkg/logger"
	types "github.com/alanpjohn/go-cdcl/pkg/types"
)

/*
The BaseCDCLSolver is a close implementation of SAT Solver v3 from http://poincare.matf.bg.ac.rs/~filip/phd/sat-tutorial.pdf
*/
type BaseCDCLSolver struct {
	Model         ModelList       // Model is stored in a customized LinkedList. Refer `model.go`
	Check         []*ModelElement // Check is used to check is an Atom has been included in the Model
	AtomCount     uint            // No of Atoms
	DecisionCount uint            // No of decisions made
	F             types.Formula   // Formula of clause to solve satisfiability problem

	/*
		Construct wraps Disjunctions into Clauses.

		This is overwritten on initialization to support modularity.
	*/
	Construct func(types.Disjunction, bool) types.Clause // a Clause Construction function which is overwritten on initialization to support modularity
}

// Intializes all the BaseCDCLSolver fields based on SATFile and CLI Flags
func InitializeBaseSolver(satfile types.SATFile, experimental bool) (solver BaseCDCLSolver, err error) {
	var clauses []types.Clause

	if experimental {
		solver.Construct = ConstructMapClause
	} else {
		solver.Construct = ConstructBaseClause
	}

	for _, d := range satfile.Clauses {
		clauses = append(clauses, solver.Construct(d, false))

	}

	/*
		The experimental flag is also meant to replace `BaseFormula` with `PQFormula`
		The `PQFormula` is currently unstable hence the code has been commented out
	*/
	// if experimental {
	// 	solver.F = ConstructHeap(clauses)
	// } else {
	// 	solver.F = BaseFormula{Clauses: clauses}
	// }

	solver.DecisionCount = 0
	solver.AtomCount = satfile.AtomCount
	solver.Check = make([]*ModelElement, satfile.AtomCount+1)

	return solver, nil
}

func (solver BaseCDCLSolver) Solve() (types.Solution, error) {
	var err error

	currentState := types.PROGRESS

	for currentState == types.PROGRESS {

		// Get next clause to process from formula
		currClause := solver.F.NextClause()
		switch currClause.Type() {
		/*
			If clause is an empty clause then we check if we have any decision literals in Model
			It we do, then we perform conflict resolution
			else the problem is unsatisfiable
		*/
		case types.EMPTY_CLAUSE:
			if solver.DecisionCount == 0 {
				return types.UNSATISFIABLE, nil
			} else {
				if err = solver.ResolveConflict(currClause); err != nil {
					return types.UNKNOWN, err
				}
			}
		/*
			If clause is a unit clause, we perform unit propagtion
		*/
		case types.UNIT_CLAUSE:
			if err = solver.UnitPropagate(currClause); err != nil {
				return types.UNKNOWN, handler.Throw("Unit Propagation Failed", err)
			}
		/*
			We dont have any unit clauses or empty clauses, hence we decide on a literal
		*/
		case types.DECISION_CLAUSE:
			if err = solver.Decide(currClause); err != nil {
				return types.UNKNOWN, handler.Throw("Decide Failed", err)
			}
		/*
			If we get a solved clause as our next clause, this means all clauses
			in formula are true in the given model, hence the solution is satisfiable
		*/
		case types.SOLVED_CLAUSE:
			return types.SATISFIABLE, nil
		}
	}
	return currentState, err
}

/*
UnitPropgate takes the only literal and appends that to our model
*/
func (solver *BaseCDCLSolver) UnitPropagate(clause types.Clause) error {
	if solver.Model.Size >= solver.AtomCount {
		return handler.Throw("Model is larger than no. of atoms", nil)
	}

	// Ideally the literal in our unit clause should not be present in the Model
	lit := clause.Disjunction()[0]
	if solver.Check[lit.Atom()] != nil {
		return handler.Throw("Atom Repeated: "+fmt.Sprint(lit), nil)
	}

	logger.Info(fmt.Sprintf("Unit propgating %v", lit))

	modelElem := &ModelElement{
		Reason:   clause,
		Literal:  lit,
		Decision: false,
	}

	solver.Model.Pushback(modelElem)
	solver.Check[lit.Atom()] = modelElem
	solver.F = solver.F.Assign(lit)

	return nil
}

/*
Decide selects the first unassigned literal from the selected clause

	[TODO] Random Selection of Decide variable
*/
func (solver *BaseCDCLSolver) Decide(clause types.Clause) error {
	if solver.Model.Size >= solver.AtomCount {
		return handler.Throw("Model is larger than no. of atoms", nil)
	}
	lit := clause.Disjunction()[0]
	if solver.Check[lit.Atom()] != nil {
		return handler.Throw("Atom Repeated: "+fmt.Sprint(lit), nil)
	}

	modelElem := &ModelElement{
		Literal:  lit,
		Decision: true,
	}

	logger.Info(fmt.Sprintf("Deciding %v", lit))

	solver.DecisionCount += 1
	solver.Model.Pushback(modelElem)
	solver.Check[lit.Atom()] = modelElem
	solver.F = solver.F.Assign(lit)

	return nil
}

/*
ResolveConflict works in 3 steps

 1. Calls AnalyseConflict to learn new clause
 2. Find the last literal in Model that makes new clause true when we negate it
 3. Find decision level to which we want to BackJump
*/
func (solver *BaseCDCLSolver) ResolveConflict(clause types.Clause) (err error) {

	logger.Info(fmt.Sprintf("Conflict Detected %v", clause.Original()))

	var resolved types.Clause = solver.Construct(clause.Original(), false)

	if resolved, err = solver.AnalyseConflict(resolved); err != nil {
		return err
	}

	solver.F.Learn(resolved)

	if modelElement, err := solver.Model.SearchLastLiteral(resolved); err != nil {
		return err
	} else {
		lastLit := modelElement.Literal // UIP
		backJumpLevel := uint(0)

		// Searching for Backjump level
		for _, lit := range resolved.Disjunction() {
			lit = lit.Negate()
			if lit != lastLit {
				decisionLvl := solver.Check[lit.Atom()].DecisionLevel
				if backJumpLevel < decisionLvl {
					backJumpLevel = decisionLvl
				}
			}
		}

		// Backjumping to Backjump level
		for m, er := solver.Model.PopTillLevel(backJumpLevel); er == nil; {

			bLit := m.Literal
			logger.Info(fmt.Sprintf("Popping %v", bLit))

			if m.Decision {
				solver.DecisionCount -= 1
			}

			solver.Check[bLit.Atom()] = nil
			solver.F = solver.F.Unassign(bLit)

			m, er = solver.Model.PopTillLevel(backJumpLevel)
		}

		lastLit = lastLit.Negate()

		modelElem := &ModelElement{
			Literal:  lastLit,
			Decision: false,
			Reason:   resolved,
		}

		logger.Info(fmt.Sprintf("Appending after conflict resolve %v", lastLit))

		solver.Model.Pushback(modelElem)
		solver.Check[lastLit.Atom()] = modelElem
		solver.F = solver.F.Assign(lastLit)

		return nil
	}
}

/*
AnalyseConflict finds a resolvent clause by continously resolving the conflict clause with
the reason of last literal of conflict clause to get a new conflict clause till we reach
unique implication point
*/
func (solver *BaseCDCLSolver) AnalyseConflict(clause types.Clause) (types.Clause, error) {
	if modelElement, err := solver.Model.SearchLastLiteral(clause); err != nil {
		return clause, err
	} else {
		lit := modelElement.Literal

		for !solver.UIP(lit, clause) {
			reason := modelElement.Reason
			if reason == nil {
				return clause, handler.Throw("null", nil)
			}
			logger.Info(fmt.Sprintf("Resolving with %v", reason.Original()))
			clause = ResolveBaseClause(reason.Original(), clause.Original(), lit, solver.AtomCount)
			logger.Info(fmt.Sprintf("Resolved %v", clause.Disjunction()))
			modelElement, err = solver.Model.SearchLastLiteral(clause)
			if err != nil {
				return clause, err
			}
			lit = modelElement.Literal
		}
	}

	return clause, nil
}

/*
UIP checks if given literal is the UIP of given clause
*/
func (solver *BaseCDCLSolver) UIP(lit types.Literal, clause types.Clause) bool {
	for _, l2 := range clause.Original() {
		l2 = l2.Negate()
		litDL := solver.Check[lit.Atom()].DecisionLevel
		l2DL := solver.Check[l2.Atom()].DecisionLevel
		if lit != l2 && l2DL == litDL {
			return false
		}
	}
	return true
}
