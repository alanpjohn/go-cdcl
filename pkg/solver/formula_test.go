package solver_test

import (
	"testing"

	//handler "github.com/alanpjohn/go-cdcl/pkg/error"
	solver "github.com/alanpjohn/go-cdcl/pkg/solver"
	types "github.com/alanpjohn/go-cdcl/pkg/types"
)

func TestBaseResolve(t *testing.T) {
	d1 := types.Disjunction{-2, -4, -5}
	d2 := types.Disjunction{-1, -3, 5}
	resolved := solver.ResolveBaseClause(d1, d2, -5, 5)

	if !resolved.Contains(-1) {
		t.Errorf("Does not contain -1")
	}
	if !resolved.Contains(-2) {
		t.Errorf("Does not contain -2")
	}
	if !resolved.Contains(-3) {
		t.Errorf("Does not contain -3")
	}
	if !resolved.Contains(-4) {
		t.Errorf("Does not contain -4")
	}

	if resolved.Contains(-5) || resolved.Contains(5) {
		t.Errorf("Should not contain 5 or -5")
	}

}
