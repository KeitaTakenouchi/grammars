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

	nodeS := dsl.NewProgramTree(S)
	nodePlus := dsl.NewProgramTree(plus)
	nodeMult := dsl.NewProgramTree(mult)

	nodeC1 := dsl.NewProgramTree(c1)
	nodeC3 := dsl.NewProgramTree(c3)
	nodeC4 := dsl.NewProgramTree(c4)

	nodeS.AddChildren(nodeMult)
	nodeMult.AddChildren(nodePlus, nodeC3)
	nodePlus.AddChildren(nodeC4, nodeC1)

	fmt.Println(nodeS.String())
	fmt.Println(nodeS.FormattedString())

}

func doSQL() {
	S := dsl.NewSymbol("S")

	sel := dsl.NewSymbol("select")
	sort := dsl.NewSymbol("sort")
	fil := dsl.NewSymbol("filter")
	grp := dsl.NewSymbol("groupby")
	join := dsl.NewSymbol("join")
	tbls := dsl.NewSymbol("tbls")

	cols := dsl.NewSymbol("cols")
	desc := dsl.NewSymbol("desc")
	pred := dsl.NewSymbol("pred")
	grpkey := dsl.NewSymbol("grpkey")
	joinkeys := dsl.NewSymbol("keyPairs")

	gram := dsl.NewGrammar(S)
	gram.AddRule(S, sel)
	gram.AddRule(sel, sort, cols)
	gram.AddRule(sel, fil, cols)
	gram.AddRule(sort, fil, desc)
	gram.AddRule(fil, grp, pred)
	gram.AddRule(grp, join, grpkey)
	gram.AddRule(join, tbls, joinkeys)

	fmt.Println(&gram)
}
