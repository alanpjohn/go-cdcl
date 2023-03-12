package solver

import (
	handler "github.com/alanpjohn/go-cdcl/pkg/error"
	types "github.com/alanpjohn/go-cdcl/pkg/types"
)

/*
ModelElement is the Node of the ModelList.
It represents a literal assignment in the Model
*/
type ModelElement struct {
	Literal       types.Literal // The Literal Assigned in the Model
	Reason        types.Clause  // The reason for assignment. Used for conflict resolution and unit propagation
	Decision      bool          // True if this is a decision literal
	DecisionLevel uint          // No of decision literal present in model before this literal counting itself
	Next          *ModelElement // The Next Pointer for ModelElement
	Prev          *ModelElement // The Prev Pointer for ModelElement
}

/*
ModelList is the LinkedList implementation of the Model
Insertion and Deletion is more effecient in Linkedlist so that is why it has been preferred over arrays
*/
type ModelList struct {
	Head          *ModelElement // Head of the List
	Tail          *ModelElement // Tail of the List
	DecisionLevel uint          // Number of decision literals in list
	Size          uint          // Size of list
}

// Pushing Literal assignment to the Model
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

// Finding Last Literal that refutes given Clause. Used for conflict resolution
func (modelList *ModelList) SearchLastLiteral(clause types.Clause) (ModelElement, error) {
	for temp := modelList.Tail; temp != nil; temp = temp.Prev {
		lit := temp.Literal
		if clause.Contains(lit.Negate()) {
			return *temp, nil
		}
	}

	return ModelElement{}, handler.Throw("Literal Not Found", nil)
}

// Pop last element of List
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

// Pop element of list till a given decision level
func (modelList *ModelList) PopTillLevel(level uint) (ModelElement, error) {
	temp := modelList.Tail
	if temp == nil || temp.DecisionLevel > level {
		return modelList.PopBack()
	} else {
		return ModelElement{}, handler.Throw("No more literals with a higher decision level", nil)
	}
}
