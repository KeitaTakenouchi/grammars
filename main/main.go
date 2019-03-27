package main

import (
	"fmt"
	"log"

	"github.com/KeitaTakenouchi/grammars/dsl"
	"github.com/KeitaTakenouchi/grammars/synth"
)

func main() {
	doSQL()
	doExp()
}

func doExp() {
	S := dsl.NewSymbol("S")

	exp := dsl.NewSymbol("exp")

	plus := dsl.NewSymbol("add")
	minus := dsl.NewSymbol("minus")
	mult := dsl.NewSymbol("mult")
	cnst := dsl.NewSymbol("const")
	param := dsl.NewSymbol("param")

	gram := dsl.NewGrammar(S)
	gram.AddRule(S, exp)
	gram.AddRule(exp, plus)
	gram.AddRule(exp, minus)
	gram.AddRule(exp, mult)
	gram.AddRule(exp, cnst)
	gram.AddRule(plus, exp, exp)
	gram.AddRule(minus, exp, exp)
	gram.AddRule(mult, exp, exp)

	fmt.Println(gram)

	var eval func(*dsl.ProgramTree, dsl.Env) dsl.EvalResult
	eval = func(node *dsl.ProgramTree, env dsl.Env) dsl.EvalResult {
		var ret dsl.EvalResult
		switch node.Symbol {
		case plus:
			e1, ok := eval(node.Children[0], env).Value()
			if !ok {
				return dsl.NewEvalResult(nil)
			}
			e2, ok := eval(node.Children[1], env).Value()
			if !ok {
				return dsl.NewEvalResult(nil)
			}
			v1 := e1.(int)
			v2 := e2.(int)
			ret = dsl.NewEvalResult(v1 + v2)
		case minus:
			e1, ok := eval(node.Children[0], env).Value()
			if !ok {
				return dsl.NewEvalResult(nil)
			}
			e2, ok := eval(node.Children[1], env).Value()
			if !ok {
				return dsl.NewEvalResult(nil)
			}
			v1 := e1.(int)
			v2 := e2.(int)
			ret = dsl.NewEvalResult(v1 - v2)
		case mult:
			e1, ok := eval(node.Children[0], env).Value()
			if !ok {
				return dsl.NewEvalResult(nil)
			}
			e2, ok := eval(node.Children[1], env).Value()
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
		case param:
			i, err := node.Value()
			if !err {
				log.Fatal("the const doesn't hava the value")
			}
			ret = dsl.NewEvalResult(env.GetArg(i.(int)))
		default:
			// S, exp,
			if len(node.Children) == 0 {
				return dsl.NewEvalResult(nil)
			}
			ret = eval(node.Children[0], env)
		}
		return ret
	}
	evaluator := dsl.NewEvaluator(eval)

	// Create a program tree to be evaluated.
	nodeS := dsl.NewProgramTree(S)
	nodePlus := dsl.NewProgramTree(plus)
	nodeMinus := dsl.NewProgramTree(minus)
	nodeMult := dsl.NewProgramTree(mult)

	nodeC1 := dsl.NewProgramTree(cnst).With(1)
	nodeC2 := dsl.NewProgramTree(cnst).With(2)
	nodeC3 := dsl.NewProgramTree(cnst).With(3)
	nodeC4 := dsl.NewProgramTree(cnst).With(4)

	nodeP0 := dsl.NewProgramTree(param).With(0)
	nodeP1 := dsl.NewProgramTree(param).With(1)
	_, _, _, _, _, _ = nodeC1, nodeC2, nodeC3, nodeC4, nodeP0, nodeP1

	nodeS.AddChildren(nodeMult)
	nodeMult.AddChildren(nodePlus, nodeMinus)
	nodePlus.AddChildren(nodeC1, nodeP0)
	nodeMinus.AddChildren(nodeC4, nodeC2)

	fmt.Println(nodeS.String())
	fmt.Println(nodeS.FormattedString())

	env := dsl.NewEnv(100, 200)
	result := evaluator.Eval(nodeS, env)
	v, _ := result.Value()
	fmt.Printf("RESULT = %v\n", v)

	filler := func(symbol *dsl.Symbol) []interface{} {
		var ret []interface{}
		switch symbol {
		case cnst:
			for i := 0; i <= 2; i++ {
				ret = append(ret, i)
			}
		}
		return ret
	}
	synthesizer := synth.NewSynthesizer(gram, evaluator, filler)
	synthesizer.Execute()
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
