// Test package to unit test different components of the BaseSolver
package solver_test

import (
	"testing"

	//handler "github.com/alanpjohn/go-cdcl/pkg/error"
	solver "github.com/alanpjohn/go-cdcl/pkg/solver"
	types "github.com/alanpjohn/go-cdcl/pkg/types"
)

func Assign(s *solver.BaseCDCLSolver, lit types.Literal, decide bool, clause types.Clause) {
	modelElem := &solver.ModelElement{
		Literal:  lit,
		Decision: decide,
	}

	if decide {
		s.DecisionCount++
	} else {
		modelElem.Reason = clause
	}

	s.Model.Pushback(modelElem)
	s.Check[lit.Atom()] = modelElem
	s.F = s.F.Assign(lit)
}

func getSampleSolver() solver.BaseCDCLSolver {
	var s solver.BaseCDCLSolver

	var disjunctions []types.Disjunction = []types.Disjunction{
		{-1, 2},
		{-3, 4},
		{-1, -3, 5},
		{-2, -4, -5},
		{-2, 3, 5, -6},
		{-1, 3, -5, -6},
		{1, -6},
		{1, 7},
	}

	var clauses []types.Clause
	for _, d := range disjunctions {
		clauses = append(clauses, solver.ConstructBaseClause(d, false))
	}

	s.F = solver.BaseFormula{Clauses: clauses}

	s.AtomCount = 7
	s.DecisionCount = 0
	s.Check = make([]*solver.ModelElement, 8)

	return s
}

func TestBackJump1(t *testing.T) {

	t.Log("Initializing Sample Solver")
	s := getSampleSolver()

	// [6d 1 2 7d 3d 4 5] is the model we apply
	Assign(&s, 6, true, solver.ConstructBaseClause(types.Disjunction{}, false))
	Assign(&s, 1, false, solver.ConstructBaseClause(types.Disjunction{1, -6}, false))
	Assign(&s, 2, false, solver.ConstructBaseClause(types.Disjunction{-1, 2}, false))
	Assign(&s, 7, true, solver.ConstructBaseClause(types.Disjunction{}, false))
	Assign(&s, 3, true, solver.ConstructBaseClause(types.Disjunction{}, false))
	Assign(&s, 4, false, solver.ConstructBaseClause(types.Disjunction{-3, 4}, false))
	Assign(&s, 5, false, solver.ConstructBaseClause(types.Disjunction{-1, -3, 5}, false))

	t.Log("Recieved Sample Solver")

	flag := s.UIP(3, solver.ConstructBaseClause(types.Disjunction{-1, -2, -3}, false))

	if !flag {
		t.Error("UIP should have been true")
	}

	flag = s.UIP(4, solver.ConstructBaseClause(types.Disjunction{-1, -2, -3, -4}, false))

	if flag {
		t.Error("UIP should have been false")
	}

	clause := solver.ConstructBaseClause(types.Disjunction{-2, -4, -5}, false)
	res, err := s.AnalyseConflict(clause)

	if err != nil {
		t.Error(err)
	}

	t.Log(res.Disjunction())

	if !res.Contains(-1) {
		t.Errorf("Does not contain -1")
	}
	if !res.Contains(-2) {
		t.Errorf("Does not contain -2")
	}
	if !res.Contains(-3) {
		t.Errorf("Does not contain -3")
	}

	err = s.ResolveConflict(clause)
	if err != nil {
		t.Error(err)
	}

	if s.Model.Tail.Literal != types.Literal(-3) {
		t.Error("Wrong conflict resolution")
	}

	// sol, err := s.Solve()

	// if err != nil || sol != types.SATISFIABLE {
	// 	t.Errorf("Soluton found : %v\nError: %v+", sol, err)
	// }

}

func TestBackJump2(t *testing.T) {

	t.Log("Initializing Sample Solver")
	s := getSampleSolver()

	learnt := solver.ConstructBaseClause(types.Disjunction{-1, -2, -3}, true)

	// [6d 1 2 7d 3d 4 5] is the model we apply
	Assign(&s, 6, true, solver.ConstructBaseClause(types.Disjunction{}, false))
	Assign(&s, 1, false, solver.ConstructBaseClause(types.Disjunction{1, -6}, false))
	Assign(&s, 2, false, solver.ConstructBaseClause(types.Disjunction{-1, 2}, false))
	Assign(&s, -3, false, learnt)
	Assign(&s, 5, false, solver.ConstructBaseClause(types.Disjunction{-2, 3, 5, -6}, false))

	s.F.Learn(learnt)

	t.Log("Recieved Sample Solver")

	flag := s.UIP(5, solver.ConstructBaseClause(types.Disjunction{-1, 3, -5, -6}, false))

	if flag {
		t.Error("UIP 5 should have been true")
	}

	flag = s.UIP(-3, solver.ConstructBaseClause(types.Disjunction{-1, -2, 3, -6}, false))

	if flag {
		t.Error("UIP -3 should have been false")
	}

	flag = s.UIP(2, solver.ConstructBaseClause(types.Disjunction{-1, -2, -6}, false))

	if flag {
		t.Error("UIP 2 should have been false")
	}

	flag = s.UIP(1, solver.ConstructBaseClause(types.Disjunction{-1, -6}, false))

	if flag {
		t.Error("UIP 1 should have been false")
	}

	flag = s.UIP(6, solver.ConstructBaseClause(types.Disjunction{-6}, false))

	if !flag {
		t.Error("UIP 6 should have been true")
	}

	clause := solver.ConstructBaseClause(types.Disjunction{-1, 3, -5, -6}, false)
	_, err := s.AnalyseConflict(clause)

	if err != nil {
		t.Error(err)
	}

	err = s.ResolveConflict(clause)
	if err != nil {
		t.Error(err)
	}

	if s.Model.Tail.Literal != types.Literal(-6) {
		t.Error("Wrong conflict resolution")
	}
}
