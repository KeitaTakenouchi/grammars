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

	exp := dsl.NewSymbol("exp")

	plus := dsl.NewSymbol("add")
	mult := dsl.NewSymbol("mult")
	cnst := dsl.NewSymbol("const")

	c1 := dsl.NewSymbol("1")
	c2 := dsl.NewSymbol("2")
	c3 := dsl.NewSymbol("3")
	c4 := dsl.NewSymbol("4")

	gram := dsl.NewGrammar(S)
	gram.AddRule(S, exp)
	gram.AddRule(exp, plus)
	gram.AddRule(exp, mult)
	gram.AddRule(exp, cnst)
	gram.AddRule(plus, exp, exp)
	gram.AddRule(mult, exp, exp)
	gram.AddRule(cnst, c1)
	gram.AddRule(cnst, c2)
	gram.AddRule(cnst, c3)
	gram.AddRule(cnst, c4)

	fmt.Println(gram)

	nodeS := dsl.NewAstNode(S)
	nodeE1 := dsl.NewAstNode(exp)
	nodeE2 := dsl.NewAstNode(exp)

	nodePlus := dsl.NewAstNode(plus)
	nodeMult := dsl.NewAstNode(mult)

	nodeC1 := dsl.NewAstNode(c1)
	nodeC3 := dsl.NewAstNode(c3)
	nodeC4 := dsl.NewAstNode(c4)

	nodeS.AddChildren(nodeE1)
	nodeE1.AddChildren(nodeMult, nodeE2, nodeC3)
	nodeE2.AddChildren(nodePlus, nodeC4, nodeC1)

	fmt.Println(nodeS.String())
	fmt.Println(nodeS.FormattedString())

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
