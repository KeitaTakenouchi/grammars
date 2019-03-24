package main

import (
	"fmt"
	"log"

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

	gram := dsl.NewGrammar(S)
	gram.AddRule(S, exp)
	gram.AddRule(exp, plus)
	gram.AddRule(exp, mult)
	gram.AddRule(exp, cnst)
	gram.AddRule(plus, exp, exp)
	gram.AddRule(mult, exp, exp)

	fmt.Println(gram)

	var eval func(node *dsl.ProgramTree) dsl.EvalResult
	eval = func(node *dsl.ProgramTree) dsl.EvalResult {
		var ret dsl.EvalResult
		switch node.Symbol {
		case plus:
			e1, ok := eval(node.Children[0]).Value()
			if !ok {
				return dsl.NewEvalResult(nil)
			}
			e2, ok := eval(node.Children[1]).Value()
			if !ok {
				return dsl.NewEvalResult(nil)
			}
			v1 := e1.(int)
			v2 := e2.(int)
			ret = dsl.NewEvalResult(v1 + v2)
		case mult:
			e1, ok := eval(node.Children[0]).Value()
			if !ok {
				return dsl.NewEvalResult(nil)
			}
			e2, ok := eval(node.Children[1]).Value()
			if !ok {
				return dsl.NewEvalResult(nil)
			}
			v1 := e1.(int)
			v2 := e2.(int)
			ret = dsl.NewEvalResult(v1 * v2)
		case cnst:
			val, err := node.Value()
			if !err {
				log.Fatal("the const doesn't hava the value")
			}
			ret = dsl.NewEvalResult(val)
		default:
			// S, exp,
			if len(node.Children) == 0 {
				return dsl.NewEvalResult(nil)
			}
			ret = eval(node.Children[0])
		}
		return ret
	}
	evaluator := dsl.NewEvaluator(eval)

	// Create a program tree to be evaluated.
	nodeS := dsl.NewProgramTree(S)
	nodePlus := dsl.NewProgramTree(plus)
	nodeMult := dsl.NewProgramTree(mult)

	nodeC1 := dsl.NewProgramTree(cnst).With(1)
	nodeC2 := dsl.NewProgramTree(cnst).With(2)
	nodeC3 := dsl.NewProgramTree(cnst).With(3)
	nodeC4 := dsl.NewProgramTree(cnst).With(4)
	_, _, _, _ = nodeC1, nodeC2, nodeC3, nodeC4

	nodeS.AddChildren(nodeMult)
	nodeMult.AddChildren(nodePlus, nodeC3)
	nodePlus.AddChildren(nodeC1, nodeC4)

	fmt.Println(nodeS.String())
	fmt.Println(nodeS.FormattedString())

	result := evaluator.Eval(nodeS)
	v, _ := result.Value()
	fmt.Printf("RESULT = %v\n", v)
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
