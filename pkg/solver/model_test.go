package solver_test

import (
	"testing"

	//handler "github.com/alanpjohn/go-cdcl/pkg/error"
	solver "github.com/alanpjohn/go-cdcl/pkg/solver"
	types "github.com/alanpjohn/go-cdcl/pkg/types"
)

func TestSearchLastLiteral(t *testing.T) {
	var modelList solver.ModelList

	for i := 1; i < 20; i++ {
		modelList.Pushback(&solver.ModelElement{
			Literal:  types.Literal(i),
			Decision: i%5 == 3,
		})
	}

	d1 := types.Disjunction{-2, -4, 5}
	c1 := solver.ConstructBaseClause(d1, false)
	m1, err1 := modelList.SearchLastLiteral(c1)

	if err1 != nil {
		t.Errorf("%v\n", err1)
	}

	if m1.Literal != types.Literal(4) {
		t.Errorf("Wrong last element %v of %v", m1.Literal, d1)
	}

	d2 := types.Disjunction{-1, 30, 8}
	c2 := solver.ConstructBaseClause(d2, false)
	m2, err2 := modelList.SearchLastLiteral(c2)

	if err2 != nil {
		t.Errorf("%v\n", err2)
	}

	if m2.Literal != types.Literal(1) {
		t.Errorf("Wrong last element %v of %v", m2.Literal, d2)
	}

	d3 := types.Disjunction{1, 2, 3}
	c3 := solver.ConstructBaseClause(d3, false)
	_, err3 := modelList.SearchLastLiteral(c3)

	if err3 == nil {
		t.Errorf("Should have thrown not found error")
	}
}

func TestPop(t *testing.T) {
	var modelList solver.ModelList

	for i := 1; i < 20; i++ {
		modelList.Pushback(&solver.ModelElement{
			Literal:  types.Literal(i),
			Decision: i%5 == 3,
		})
	}

	backJump := modelList.PopTillLevel(2)

	var (
		m   solver.ModelElement
		err error
	)

	lit := types.Literal(19)
	for m, err = backJump(); err == nil; {
		t.Logf("Popped literal %v at Decision Level %v\n", m.Literal, m.DecisionLevel)
		if lit != m.Literal {
			t.Errorf("Incorrect Item Popped %v with decision level %v\n", m.Literal, m.DecisionLevel)
		}
		lit -= 1
		m, err = backJump()
	}

	if lit != types.Literal(12) || modelList.Tail.Literal != types.Literal(12) {
		t.Errorf("Incorrect tail %v\n", modelList.Tail.Literal)
	}
}
