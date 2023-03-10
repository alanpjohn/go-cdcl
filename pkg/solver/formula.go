package solver

import (
	types "github.com/alanpjohn/go-cdcl/pkg/types"
)

type BaseClause struct {
	original types.Disjunction
	updated  types.Disjunction
	solved   bool
}

func (c *BaseClause) Type() types.ClauseType {
	if c.solved {
		return types.SOLVED_CLAUSE
	} else if len(c.updated) == 0 {
		return types.EMPTY_CLAUSE
	} else if len(c.updated) == 1 {
		return types.UNIT_CLAUSE
	} else {
		return types.DECISION_CLAUSE
	}
}

func (c *BaseClause) Apply(l types.Literal) {
	for _, cl := range c.updated {
		if cl == l {
			c.solved = true
			return
		}
	}

	var disjunction types.Disjunction
	for _, cl := range c.updated {
		if cl != -l {
			disjunction = append(disjunction, cl)
		}
	}
	c.updated = disjunction
}

func (c *BaseClause) Undo(l types.Literal) {

}

func (c *BaseClause) Disjunction() types.Disjunction {
	return c.updated
}

func (c *BaseClause) Original() types.Disjunction {
	return c.original
}

func (c *BaseClause) Contains(l types.Literal) bool {
	for i := 0; i < len(c.updated); i++ {
		if c.updated[i] == l || c.updated[i] == -l {
			return true
		}
	}
	return false
}

type BaseFormula struct {
	Clauses []BaseClause
}

func (f BaseFormula) NextClause() types.Clause {
	// check
	for _, c := range f.Clauses {
		if c.Type() == types.EMPTY_CLAUSE {
			return &c
		}
	}

	for _, c := range f.Clauses {
		if c.Type() == types.UNIT_CLAUSE {
			return &c
		}
	}

	for _, c := range f.Clauses {
		if c.Type() == types.DECISION_CLAUSE {
			return &c
		}
	}

	return &f.Clauses[0]

}

func (f BaseFormula) Apply(l types.Literal) types.Formula {
	for i := range f.Clauses {
		f.Clauses[i].Apply(l)
	}
	return f
}
