package solver

import (
	"fmt"

	logger "github.com/alanpjohn/go-cdcl/pkg/logger"
	types "github.com/alanpjohn/go-cdcl/pkg/types"
)

type BaseClause struct {
	original types.Disjunction
	updated  types.Disjunction
	solved   uint
	learnt   bool
}

func ConstructBaseClause(d types.Disjunction, learnt bool) BaseClause {
	return BaseClause{d, d, 0, learnt}
}

func (c BaseClause) Type() types.ClauseType {
	if c.solved > 0 {
		return types.SOLVED_CLAUSE
	} else if len(c.updated) == 0 {
		return types.EMPTY_CLAUSE
	} else if len(c.updated) == 1 {
		return types.UNIT_CLAUSE
	} else {
		return types.DECISION_CLAUSE
	}
}

func (c BaseClause) Apply(l types.Literal) types.Clause {
	for _, cl := range c.updated {
		if cl == l {
			c.solved++
			return c
		}
	}

	var disjunction types.Disjunction
	for _, cl := range c.updated {
		if cl != -l {
			disjunction = append(disjunction, cl)
		}
	}
	c.updated = disjunction
	return c
}

func (c BaseClause) Undo(l types.Literal) types.Clause {
	for _, cl := range c.original {
		if cl == l {
			c.solved--
			return c
		}
	}
	for _, cl := range c.original {
		if cl == -l {
			c.updated = append(c.updated, cl)
		}
	}
	return c
}

func (c BaseClause) Reset() types.Clause {
	return ConstructBaseClause(c.original, c.learnt)
}

func (c BaseClause) Disjunction() types.Disjunction {
	return c.updated
}

func (c BaseClause) Original() types.Disjunction {
	return c.original
}

func (c BaseClause) IsLearnt() bool {
	return c.learnt
}
func (c BaseClause) IsSolved() bool {
	return c.solved > 0
}

func (c BaseClause) Contains(l types.Literal) bool {
	for i := 0; i < len(c.updated); i++ {
		if c.updated[i] == l {
			return true
		}
	}
	return false
}

func ResolveBaseClause(d1, d2 types.Disjunction, lit types.Literal, atomCount uint) BaseClause {
	var d types.Disjunction

	for _, l := range d1 {
		if l.Atom() != lit.Atom() {
			d = append(d, l)
		}
	}

	for _, l := range d2 {
		if l.Atom() != lit.Atom() {
			d = append(d, l)
		}
	}

	return ConstructBaseClause(d, true)
}

type BaseFormula struct {
	Clauses []types.Clause
}

func (f BaseFormula) NextClause() types.Clause {
	// check
	for _, c := range f.Clauses {
		if c.Type() == types.EMPTY_CLAUSE {
			return c
		}
	}

	for _, c := range f.Clauses {
		if c.Type() == types.UNIT_CLAUSE {
			return c
		}
	}

	for _, c := range f.Clauses {
		if c.Type() == types.DECISION_CLAUSE {
			return c
		}
	}

	return f.Clauses[0]

}

func (f BaseFormula) Assign(l types.Literal) types.Formula {
	for i, c := range f.Clauses {
		f.Clauses[i] = c.Apply(l)
		if !f.Clauses[i].IsSolved() {
			logger.Info(fmt.Sprintf("Updated %v to %v", c.Disjunction(), f.Clauses[i].Disjunction()))
		}
	}
	return f
}

func (f BaseFormula) Unassign(l types.Literal) types.Formula {
	for i, c := range f.Clauses {
		f.Clauses[i] = c.Undo(l)
	}
	return f
}

func (f BaseFormula) Learn(c types.Clause) types.Formula {
	f.Clauses = append(f.Clauses, c)
	return f
}

func (f BaseFormula) Restart() types.Formula {
	for i, c := range f.Clauses {
		f.Clauses[i] = c.Reset()
	}

	return f
}

func (f BaseFormula) Print() string {
	return fmt.Sprintf("%v", f.Clauses)

}
