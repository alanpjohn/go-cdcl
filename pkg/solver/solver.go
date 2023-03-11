package solver

import (
	"fmt"

	handler "github.com/alanpjohn/go-cdcl/pkg/error"
	logger "github.com/alanpjohn/go-cdcl/pkg/logger"
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
	var clauses []types.Clause
	for _, d := range satfile.Clauses {
		clauses = append(clauses, ConstructBaseClause(d, false))
	}

	solver.F = BaseFormula{Clauses: clauses}

	solver.DecisionCount = 0
	solver.AtomCount = satfile.AtomCount
	solver.Check = make([]*ModelElement, satfile.AtomCount+1)
	solver.Seek = 0

	return solver, nil
}

func (solver BaseCDCLSolver) Solve() (types.Solution, error) {
	var err error

	currentState := types.PROGRESS

	for currentState == types.PROGRESS {
		currClause := solver.F.NextClause()
		logger.Info(fmt.Sprintf("Current clause: %v", currClause.Disjunction()))
		logger.Info(fmt.Sprintf("Originally: %v", currClause.Original()))
		switch currClause.Type() {
		case types.EMPTY_CLAUSE:
			if solver.DecisionCount == 0 {
				return types.UNSATISFIABLE, nil
			} else {
				if err = solver.ResolveConflict(currClause); err != nil {
					return types.UNKNOWN, err
				}
			}
		case types.UNIT_CLAUSE:
			if err = solver.UnitPropagate(currClause); err != nil {
				return types.UNKNOWN, handler.Throw("Unit Propagation Failed", err)
			}
		case types.DECISION_CLAUSE:
			if err = solver.Decide(currClause); err != nil {
				return types.UNKNOWN, handler.Throw("Decide Failed", err)
			}
		case types.SOLVED_CLAUSE:
			return types.SATISFIABLE, nil
		}
	}
	return currentState, err
}

func (solver *BaseCDCLSolver) UnitPropagate(clause types.Clause) error {
	if solver.Model.Size >= solver.AtomCount {
		return handler.Throw("Model is larger than no. of atoms", nil)
	}
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

func (solver *BaseCDCLSolver) ResolveConflict(clause types.Clause) (err error) {

	logger.Info(fmt.Sprintf("Conflict Detected %v", clause.Original()))

	var resolved types.Clause = ConstructBaseClause(clause.Original(), false)
	if resolved, err = solver.AnalyseConflict(resolved); err != nil {
		return err
	}

	solver.F.Learn(resolved)

	if modelElement, err := solver.Model.SearchLastLiteral(resolved); err != nil {
		return err
	} else {
		lastLit := modelElement.Literal
		backJumpLevel := uint(0)

		for _, lit := range resolved.Disjunction() {
			lit = lit.Negate()
			if lit != lastLit {
				decisionLvl := solver.Check[lit.Atom()].DecisionLevel
				if backJumpLevel < decisionLvl {
					backJumpLevel = decisionLvl
				}
			}
		}
		backJumpedLiterals := solver.Model.PopTillLevel(backJumpLevel)

		for m, er := backJumpedLiterals(); er == nil; {

			bLit := m.Literal
			logger.Info(fmt.Sprintf("Popping %v", bLit))

			if m.Decision {
				solver.DecisionCount -= 1
			}

			solver.Check[bLit.Atom()] = nil
			solver.F = solver.F.Unassign(bLit)

			m, er = backJumpedLiterals()
		}

		modelElem := &ModelElement{
			Literal:  lastLit,
			Decision: false,
		}

		logger.Info(fmt.Sprintf("Appending after conflict resolve %v", lastLit))

		solver.Model.Pushback(modelElem)
		solver.Check[lastLit.Atom()] = modelElem
		solver.F = solver.F.Assign(lastLit)

		return nil
	}
}

func (solver *BaseCDCLSolver) AnalyseConflict(clause types.Clause) (types.Clause, error) {
	if modelElement, err := solver.Model.SearchLastLiteral(clause); err != nil {
		return clause, err
	} else {
		lit := modelElement.Literal

		for !solver.UIP(lit.Negate(), clause) {
			reason := modelElement.Reason
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

func (solver *BaseCDCLSolver) UIP(lit types.Literal, clause types.Clause) bool {
	for _, l2 := range clause.Disjunction() {
		litDL := solver.Check[lit.Atom()].DecisionLevel
		l2DL := solver.Check[l2.Atom()].DecisionLevel
		if lit != l2 && l2DL == litDL {
			return false
		}
	}
	return true
}
