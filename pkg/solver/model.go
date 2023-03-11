package solver

import (
	handler "github.com/alanpjohn/go-cdcl/pkg/error"
	types "github.com/alanpjohn/go-cdcl/pkg/types"
)

type ModelElement struct {
	Literal       types.Literal
	Reason        types.Clause
	Decision      bool
	DecisionLevel uint
	Next          *ModelElement
	Prev          *ModelElement
}

type ModelList struct {
	Head          *ModelElement
	Tail          *ModelElement
	DecisionLevel uint
	Size          uint
}

func (modelList *ModelList) Pushback(m *ModelElement) {
	if m.Decision {
		modelList.DecisionLevel += 1
	}

	if modelList.Head == nil {
		modelList.Head = m
		modelList.Tail = m
	} else {
		modelList.Tail.Next = m
		m.Prev = modelList.Tail
		modelList.Tail = m
		m.DecisionLevel = modelList.DecisionLevel
	}
	modelList.Size += 1

}

func (modelList *ModelList) SearchLastLiteral(clause types.Clause) (ModelElement, error) {
	for temp := modelList.Tail; temp != nil; temp = temp.Prev {
		lit := temp.Literal
		if clause.Contains(lit.Negate()) {
			return *temp, nil
		}
	}

	return ModelElement{}, handler.Throw("Literal Not Found", nil)
}

func (modelList *ModelList) PopTillLevel(level uint) func() (ModelElement, error) {
	temp := modelList.Tail

	return func() (ModelElement, error) {
		for temp != nil {
			curr := temp
			temp = temp.Prev
			if curr.DecisionLevel > level {
				modelList.Tail = temp
				temp.Next = nil
				curr.Prev = nil
				modelList.Size -= 1
				return *curr, nil
			}
		}

		return ModelElement{}, handler.Throw("No more elements", nil)
	}
}
