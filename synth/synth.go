package synth

import (
	"fmt"

	"github.com/KeitaTakenouchi/grammars/dsl"
)

type Synthesizer struct {
	grammar   dsl.Grammar
	evaluator dsl.Evaluator
	filler    func(*dsl.Symbol) []interface{}
}

func NewSynthesizer(grammar dsl.Grammar, eval dsl.Evaluator, filler func(*dsl.Symbol) []interface{}) Synthesizer {
	return Synthesizer{
		grammar:   grammar,
		evaluator: eval,
		filler:    filler,
	}
}

func (s *Synthesizer) Execute() {
	worklist := make([]*dsl.ProgramTree, 0)
	start := dsl.NewProgramTree(s.grammar.GetStart())
	worklist = append(worklist, start)

	index, maxIndex := 0, 0
	counterForDebug := 0
	for index <= maxIndex {
		counterForDebug++
		if counterForDebug > 1000 {
			break
		}

		target := worklist[index]
		index++

		nonTerminals := target.NonTerminalLeaves()
		if len(nonTerminals) == 0 {
			for _, completePgm := range s.fillSketch(target) {
				if s.check(completePgm) {
					return
				}
			}
		}

		for i := 0; i < len(nonTerminals); i++ {
			seqs := s.grammar.GetRhs(nonTerminals[i].Symbol)
			for _, seq := range seqs {
				cpy := target.Clone()
				node := cpy.NonTerminalLeaves()[i]
				for _, symbol := range seq {
					pgm := dsl.NewProgramTree(symbol)
					node.AddChildren(pgm)
				}
				worklist = append(worklist, cpy)
				maxIndex++
			}
		}
	}

}

func (s *Synthesizer) fillSketch(pgm *dsl.ProgramTree) []*dsl.ProgramTree {
	valuesList := make([][]interface{}, 0)
	indexesOfHoles := make([]int, 0)
	for i, leaf := range pgm.Leaves() {
		values := s.filler(leaf.Symbol)
		if len(values) > 0 {
			valuesList = append(valuesList, values)
			indexesOfHoles = append(indexesOfHoles, i)
		}
	}

	var ret []*dsl.ProgramTree
	for _, valueComb := range cartesianProduct(valuesList) {
		clone := pgm.Clone()
		for i, indexOfHole := range indexesOfHoles {
			clone.Leaves()[indexOfHole].With(valueComb[i])
		}
		ret = append(ret, clone)
	}
	return ret
}

func (s *Synthesizer) check(pgm *dsl.ProgramTree) bool {
	env := dsl.NewEnv()
	result := s.evaluator.Eval(pgm, env)
	fmt.Println("-----------------------")
	fmt.Println(pgm.FormattedString())
	fmt.Println("result =", result)
	return false
}

func cartesianProduct(lists [][]interface{}) [][]interface{} {
	if len(lists) == 0 {
		return [][]interface{}{[]interface{}{}}
	}
	ret := make([][]interface{}, 0)
	head, tail := lists[0], lists[1:]
	for _, tails := range cartesianProduct(tail) {
		for _, h := range head {
			comb := append([]interface{}{h}, tails...)
			ret = append(ret, comb)
		}
	}
	return ret
}
