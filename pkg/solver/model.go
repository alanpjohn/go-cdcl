package solver

import types "github.com/alanpjohn/go-cdcl/pkg/types"

type ModelElement struct {
	Literal      types.Literal
	Dependents   []*ModelElement
	Decision     bool
	LastDecision *ModelElement
	Next         *ModelElement
	Prev         *ModelElement
}

type ModelList struct {
	Head *ModelElement
	Tail *ModelElement
	Size uint
}

func (modelList *ModelList) Pushback(m *ModelElement) {
	if modelList.Head == nil {
		modelList.Head = m
		modelList.Tail = m
	} else {
		modelList.Tail.Next = m
		m.Prev = modelList.Tail
		if modelList.Tail.Decision {
			m.LastDecision = modelList.Tail
		} else {
			m.LastDecision = modelList.Tail.LastDecision
		}
		modelList.Tail = m
	}
}
