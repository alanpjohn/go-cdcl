package pkg

type Model interface{
	UnitPropagate(l Literal)
	Decide(l Literal)
	Backjump() (Status, []Literal) 		// Doubtful
	LemmaLearn() Clause					// Doubtful
	IsComplete() bool
	IsEmpty() bool
	Size() uint
}