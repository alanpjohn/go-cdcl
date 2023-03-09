package pkg

/*
Formula represents the formula of clauses to be solved by the solver.
*/
type Formula interface {
	Apply( l Literal)
	Undo( l Literal)
	Learn( c Clause )
	Propagate() (Status,Clause)
	
}