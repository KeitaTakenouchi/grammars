package main

import (
	"fmt"

	"github.com/KeitaTakenouchi/grammars/dsl"
)

func main() {
	doSQL()
	doExp()
}

func doExp() {
	S := dsl.NewSymbol("S")

	exp := dsl.NewSymbol("EXP")

	lparen := dsl.NewSymbol("(")
	rparen := dsl.NewSymbol(")")
	plus := dsl.NewSymbol("+")
	mult := dsl.NewSymbol("*")

	gram := dsl.NewGrammar(S)
	gram.AddRule(S, exp)
	gram.AddRule(exp, lparen, exp, rparen)
	gram.AddRule(exp, exp, plus, exp)
	gram.AddRule(exp, exp, mult, exp)

	fmt.Println(gram)
}

func doSQL() {
	S := dsl.NewSymbol("S")

	T1 := dsl.NewSymbol("T1")
	T2 := dsl.NewSymbol("T2")
	T3 := dsl.NewSymbol("T3")
	T4 := dsl.NewSymbol("T4")

	sel := dsl.NewSymbol("select")
	sort := dsl.NewSymbol("sort")
	fil := dsl.NewSymbol("filter")
	grp := dsl.NewSymbol("groupby")
	join := dsl.NewSymbol("join")
	tbls := dsl.NewSymbol("tbls")

	gram := dsl.NewGrammar(S)
	gram.AddRule(S, sel, T1)
	gram.AddRule(T1, sort, T2)
	gram.AddRule(T2, fil, T3)
	gram.AddRule(T3, grp, T4)
	gram.AddRule(T4, join, tbls)

	fmt.Println(&gram)
}
