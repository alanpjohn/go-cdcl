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
	m.DecisionLevel = modelList.DecisionLevel
	if modelList.Head == nil {
		modelList.Head = m
		modelList.Tail = m
	} else {
		modelList.Tail.Next = m
		m.Prev = modelList.Tail
		modelList.Tail = m
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

func (modelList *ModelList) PopBack() (ModelElement, error) {
	if modelList.Head == nil {
		return ModelElement{}, handler.Throw("Empty List", nil)
	}
	temp := modelList.Tail
	prev := temp.Prev
	modelList.Tail = prev
	if temp.Decision {
		modelList.DecisionLevel -= 1
	}
	modelList.Size -= 1
	if prev != nil {
		prev.Next = nil
		temp.Prev = nil
	} else {
		modelList.Head = nil
		modelList.Tail = nil
	}

	return *temp, nil
}

func (modelList *ModelList) PopTillLevel(level uint) (ModelElement, error) {

	temp := modelList.Tail
	if temp == nil || temp.DecisionLevel > level {
		return modelList.PopBack()
	} else {
		return ModelElement{}, handler.Throw("No more literals with a higher decision level", nil)
	}
}
