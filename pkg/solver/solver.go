package solver

import (
	"fmt"

	handler "github.com/alanpjohn/go-cdcl/pkg/error"
	types "github.com/alanpjohn/go-cdcl/pkg/types"
)

type BaseCDCLSolver struct {
	Seek          types.Literal
	Model         ModelList
	Check         []*ModelElement
	AtomCount     uint
	DecisionCount uint
	F             types.Formula
}

func InitializeBaseSolver(satfile types.SATFile) (solver BaseCDCLSolver, err error) {
	var clauses []BaseClause
	for i := 0; i < len(satfile.Clauses); i++ {
		clauses = append(clauses, BaseClause{satfile.Clauses[i], satfile.Clauses[i], false})
	}

	solver.F = BaseFormula{Clauses: clauses}

	solver.DecisionCount = 0
	solver.AtomCount = satfile.AtomCount
	solver.Check = make([]*ModelElement, satfile.AtomCount)
	solver.Seek = 0

	return solver, nil
}

func (solver BaseCDCLSolver) Solve() (types.Solution, error) {
	currClause := solver.F.NextClause()

	var nextSolver BaseCDCLSolver
	var err error

	switch currClause.Type() {
	case types.EMPTY_CLAUSE:
		if solver.DecisionCount == 0 {
			return types.UNSATISFIABLE, nil
		} else {
			if nextSolver, err = solver.ResolveConflict(currClause); err != nil {
				return types.UNKNOWN, handler.Throw("Conflict Resolution failed", err)
			}
		}
	case types.UNIT_CLAUSE:
		if nextSolver, err = solver.UnitPropagate(currClause); err != nil {
			return types.UNKNOWN, handler.Throw("Unit Propagation Failed", err)
		}
	case types.DECISION_CLAUSE:
		if nextSolver, err = solver.Decide(currClause); err != nil {
			return types.UNKNOWN, handler.Throw("Unit Propagation Failed", err)
		}
	case types.SOLVED_CLAUSE:
		return types.SATISFIABLE, nil
	}

	return nextSolver.Solve()
}

func (solver BaseCDCLSolver) UnitPropagate(clause types.Clause) (BaseCDCLSolver, error) {
	if solver.Model.Size >= solver.AtomCount {
		return solver, handler.Throw("Model is larger than no. of atoms", nil)
	}
	lit := clause.Disjunction()[0]
	if solver.Check[lit.Atom()] != nil {
		return solver, handler.Throw("Atom Repeated: "+fmt.Sprint(lit), nil)
	}

	org := clause.Original()
	var deps []*ModelElement
	for _, v := range org {
		if lit != v && solver.Check[v.Atom()] != nil {
			deps = append(deps, solver.Check[v.Atom()])
		}
	}

	modelElem := &ModelElement{
		Dependents: deps,
		Literal:    lit,
		Decision:   false,
	}

	solver.Model.Pushback(modelElem)
	solver.Check[lit.Atom()] = modelElem
	solver.F.Apply(lit)

	return solver, nil
}

func (solver BaseCDCLSolver) Decide(clause types.Clause) (BaseCDCLSolver, error) {
	if solver.Model.Size >= solver.AtomCount {
		return solver, handler.Throw("Model is larger than no. of atoms", nil)
	}
	lit := clause.Disjunction()[0]
	if solver.Check[lit.Atom()] != nil {
		return solver, handler.Throw("Atom Repeated: "+fmt.Sprint(lit), nil)
	}

	modelElem := &ModelElement{
		Literal:  lit,
		Decision: true,
	}

	solver.Model.Pushback(modelElem)
	solver.Check[lit.Atom()] = modelElem
	solver.F.Apply(lit)

	return solver, nil
}

func (solver BaseCDCLSolver) ResolveConflict(clause types.Clause) (BaseCDCLSolver, error) {
	last := solver.Model.Tail
	if !clause.Contains(types.Literal(last.Literal)) {
		return solver, handler.Throw("Model does not affect equation", nil)
	}

	return solver, handler.Throw("Unimplemented Conflict Resolution", nil)
}
