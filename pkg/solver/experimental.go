package solver

/*
The experimental file contains alternate implementations of the Clause interface and Formula interface
To experiment with new ideas
*/

import (
	"container/heap"
	"fmt"

	types "github.com/alanpjohn/go-cdcl/pkg/types"
)

type LiteralState int

const (
	UNASSIGNED LiteralState = iota
	ASSIGNED
	REFUTED
)

type MapClause struct {
	original map[types.Atom]int
	litMap   map[types.Atom]LiteralState
	solved   uint
	learnt   bool
}

func ConstructMapClause(d types.Disjunction, learnt bool) types.Clause {

	original := make(map[types.Atom]int)
	litMap := make(map[types.Atom]LiteralState)

	for _, lit := range d {
		var sign int = int(lit / types.Literal(lit.Atom()))
		if val, ok := original[lit.Atom()]; ok && sign != val {
			return MapClause{
				make(map[types.Atom]int),
				make(map[types.Atom]LiteralState),
				2,
				false,
			}
		}
		original[lit.Atom()] = sign
		litMap[lit.Atom()] = UNASSIGNED
	}

	return MapClause{
		original,
		litMap,
		0,
		false,
	}
}

func (c MapClause) Type() types.ClauseType {
	if c.solved > 0 {
		return types.SOLVED_CLAUSE
	}

	un := 0
	for _, val := range c.litMap {
		if val == UNASSIGNED {
			un += 1
		}
	}

	if un == 0 {
		return types.EMPTY_CLAUSE
	}
	if un == 1 {
		return types.UNIT_CLAUSE
	}

	return types.DECISION_CLAUSE
}

func (c MapClause) Apply(lit types.Literal) types.Clause {
	if val, ok := c.original[lit.Atom()]; ok {
		var sign int = int(lit / types.Literal(lit.Atom()))
		if sign == val {
			c.solved += 1
			c.litMap[lit.Atom()] = ASSIGNED
		} else {
			c.litMap[lit.Atom()] = REFUTED
		}
	}
	return c
}

func (c MapClause) Undo(lit types.Literal) types.Clause {
	if val, ok := c.litMap[lit.Atom()]; ok && val != UNASSIGNED {
		var sign int = int(lit / types.Literal(lit.Atom()))
		if sign == c.original[lit.Atom()] {
			c.solved -= 1
			c.litMap[lit.Atom()] = UNASSIGNED
		} else {
			c.litMap[lit.Atom()] = UNASSIGNED
		}
	}
	return c
}

func (c MapClause) Reset() types.Clause {
	c.original = make(map[types.Atom]int)
	c.litMap = make(map[types.Atom]LiteralState)
	c.solved = 0
	return c
}

func (c MapClause) Disjunction() types.Disjunction {
	var d types.Disjunction

	for k, v := range c.original {
		newlit := types.Literal(int(k) * v)
		if c.litMap[k] != REFUTED {
			d = append(d, newlit)
		}
	}

	return d
}

func (c MapClause) Original() types.Disjunction {
	var d types.Disjunction

	for k, v := range c.original {
		newlit := types.Literal(int(k) * v)
		d = append(d, newlit)
	}

	return d
}

func (c MapClause) IsLearnt() bool {
	return c.learnt
}
func (c MapClause) IsSolved() bool {
	return c.solved > 0
}

func (c MapClause) Contains(lit types.Literal) bool {
	if val, ok := c.original[lit.Atom()]; ok {
		var sign int = int(lit / types.Literal(lit.Atom()))
		if val == sign {
			return true
		}
	}
	return false
}

func ResolveMapClause(d1, d2 types.Disjunction, lit types.Literal, atomCount uint) types.Clause {
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

	return ConstructMapClause(d, true)
}

type PQFormula struct {
	Clauses PriorityQueue
}

type PriorityQueue []types.Clause

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].Type() < pq[j].Type()
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x any) {
	item := x.(types.Clause)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	*pq = old[0 : n-1]
	return item
}

func ConstructHeap(clauses []types.Clause) PQFormula {
	var f PQFormula
	f.Clauses = clauses
	for _, clause := range clauses {
		if !clause.IsSolved() {
			f.Clauses = append(f.Clauses, clause)
		}
	}
	heap.Init(&f.Clauses)
	return f
}

func (f PQFormula) NextClause() types.Clause {
	x := heap.Pop(&f.Clauses)
	clause := x.(types.Clause)
	return clause
}

func (f PQFormula) Assign(l types.Literal) types.Formula {
	for i, c := range f.Clauses {
		f.Clauses[i] = c.Apply(l)
	}
	heap.Init(&f.Clauses)
	return f
}

func (f PQFormula) Unassign(l types.Literal) types.Formula {
	for i, c := range f.Clauses {
		f.Clauses[i] = c.Undo(l)
	}
	heap.Init(&f.Clauses)
	return f
}

func (f PQFormula) Learn(c types.Clause) types.Formula {
	heap.Push(&f.Clauses, c)
	return f
}

func (f PQFormula) Restart() types.Formula {
	for i, c := range f.Clauses {
		f.Clauses[i] = c.Reset()
	}
	heap.Init(&f.Clauses)
	return f
}

func (f PQFormula) Print() string {
	return fmt.Sprintf("%v", f.Clauses)

}
