package io

import (
	formula "github.com/alanpjohn/go-cdcl/pkg"
)

type SATFile struct {
	a int             // No of atoms
	c int             // No of clauses
	f formula.Formula // Formula read from .SAT file
}

func (S *SATFile) AtomCount() int {
	return S.a
}

func (S *SATFile) ClauseCount() int {
	return S.c
}

func (S *SATFile) Formula() formula.Formula {
	return S.f
}
