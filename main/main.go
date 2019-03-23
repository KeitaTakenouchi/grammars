package main

import (
	"fmt"

	"github.com/KeitaTakenouchi/grammars/grammar"
)

func main() {
	doSQL()
}

func doSQL() {
	S := grammar.NewSymbol("S")

	T1 := grammar.NewSymbol("T1")
	T2 := grammar.NewSymbol("T2")
	T3 := grammar.NewSymbol("T3")
	T4 := grammar.NewSymbol("T4")

	sel := grammar.NewSymbol("select")
	sort := grammar.NewSymbol("sort")
	fil := grammar.NewSymbol("filter")
	grp := grammar.NewSymbol("groupby")
	join := grammar.NewSymbol("join")
	tbls := grammar.NewSymbol("tbls")

	gram := grammar.NewGrammar(S)
	gram.AddRule(S, sel, T1)
	gram.AddRule(T1, sort, T2)
	gram.AddRule(T2, fil, T3)
	gram.AddRule(T3, grp, T4)
	gram.AddRule(T4, join, tbls)

	fmt.Println(&gram)
}
