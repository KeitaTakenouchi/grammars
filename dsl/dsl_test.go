package dsl

import (
	"log"
	"reflect"
	"testing"
)

func TestSymbol_01(t *testing.T) {
	S := NewSymbol("S")
	exp := NewSymbol("EXP")
	plus := NewSymbol("+")
	mult := NewSymbol("*")
	gram := NewGrammar(S)
	gram.AddRule(S, exp)
	gram.AddRule(exp, exp, plus, exp)
	gram.AddRule(exp, exp, mult, exp)

	t.Run("test the symbols", func(t *testing.T) {
		got := len(gram.symbols.symbols)
		want := 4
		if got != want {
			t.Errorf("the number of sumbols is incorrect. got=%d, want=%d", got, want)
		}
	})

	t.Run("test the isTerminal", func(t *testing.T) {
		gotTreminalCnt := 0
		gotNonTerminalCnt := 0
		for s, _ := range gram.symbols.symbols {
			if s.isTerminal {
				gotTreminalCnt++
			} else {
				gotNonTerminalCnt++
			}
		}
		wantTreminalCnt := 2
		if gotTreminalCnt != wantTreminalCnt {
			t.Errorf("the number of isTerminal is incorrect. got=%d, want=%d", gotTreminalCnt, wantTreminalCnt)
		}
		wantNonTreminalCnt := 2
		if gotNonTerminalCnt != wantNonTreminalCnt {
			t.Errorf("the number of isTerminal is incorrect. got=%d, want=%d", gotNonTerminalCnt, wantNonTreminalCnt)
		}
	})
}

func TestSymbol_02(t *testing.T) {
	S := NewSymbol("S")
	exp := NewSymbol("EXP")
	lparen := NewSymbol("(")
	rparen := NewSymbol(")")
	plus := NewSymbol("+")
	mult := NewSymbol("*")
	gram := NewGrammar(S)
	gram.AddRule(S, exp)
	gram.AddRule(exp, lparen, exp, rparen)
	gram.AddRule(exp, exp, plus, exp)
	gram.AddRule(exp, exp, plus, exp)
	gram.AddRule(exp, exp, plus, exp)
	gram.AddRule(exp, exp, mult, exp)
	gram.AddRule(exp, exp, mult, exp)
	gram.AddRule(exp, exp, mult, exp)
	gram.AddRule(exp, exp, mult, exp)
	gram.AddRule(exp, exp, mult, exp)

	t.Run("test the symbols", func(t *testing.T) {
		got := len(gram.symbols.symbols)
		want := 6
		if got != want {
			t.Errorf("the number of sumbols is incorrect. got=%d, want=%d", got, want)
		}
	})

	t.Run("test the isTerminal", func(t *testing.T) {
		gotTreminalCnt := 0
		gotNonTerminalCnt := 0
		for s, _ := range gram.symbols.symbols {
			if s.isTerminal {
				gotTreminalCnt++
			} else {
				gotNonTerminalCnt++
			}
		}
		wantTreminalCnt := 4
		if gotTreminalCnt != wantTreminalCnt {
			t.Errorf("the number of isTerminal is incorrect. got=%d, want=%d", gotTreminalCnt, wantTreminalCnt)
		}
		wantNonTreminalCnt := 2
		if gotNonTerminalCnt != wantNonTreminalCnt {
			t.Errorf("the number of isTerminal is incorrect. got=%d, want=%d", gotNonTerminalCnt, wantNonTreminalCnt)
		}
	})
}

func TestEvaluator_Eval_Exp(t *testing.T) {
	S := NewSymbol("S")
	exp := NewSymbol("exp")
	plus := NewSymbol("add")
	mult := NewSymbol("mult")
	cnst := NewSymbol("const")
	_, _, _, _, _ = S, exp, plus, mult, cnst

	var eval func(node *ProgramTree) EvalResult
	eval = func(node *ProgramTree) EvalResult {
		var ret EvalResult
		switch node.Symbol {
		case plus:
			e1, ok := eval(node.Children[0]).Value()
			if !ok {
				return NewEvalResult(nil)
			}
			e2, ok := eval(node.Children[1]).Value()
			if !ok {
				return NewEvalResult(nil)
			}
			v1 := e1.(int)
			v2 := e2.(int)
			ret = NewEvalResult(v1 + v2)
		case mult:
			e1, ok := eval(node.Children[0]).Value()
			if !ok {
				return NewEvalResult(nil)
			}
			e2, ok := eval(node.Children[1]).Value()
			if !ok {
				return NewEvalResult(nil)
			}
			v1 := e1.(int)
			v2 := e2.(int)
			ret = NewEvalResult(v1 * v2)
		case cnst:
			val, err := node.Value()
			if !err {
				log.Fatal("the const doesn't hava the value")
			}
			ret = NewEvalResult(val)
		default:
			// S, exp,
			if len(node.Children) == 0 {
				return NewEvalResult(nil)
			}
			ret = eval(node.Children[0])
		}
		return ret
	}
	evaluator := NewEvaluator(eval)

	type N = ProgramTree
	tests := []struct {
		name string
		arg  *ProgramTree
		want EvalResult
	}{
		{
			name: "(1+4)*3",
			want: NewEvalResult(15),
			arg: &N{
				Symbol: mult, Children: []*N{
					&N{
						Symbol: plus, Children: []*N{
							&N{Symbol: cnst, value: 1},
							&N{Symbol: cnst, value: 4},
						},
					},
					&N{Symbol: cnst, value: 3},
				},
			},
		},
		{
			name: "10*((1+4)*3)",
			want: NewEvalResult(150),
			arg: &N{
				Symbol: mult, Children: []*N{
					&N{Symbol: cnst, value: 10},
					&N{
						Symbol: mult, Children: []*N{
							&N{
								Symbol: plus, Children: []*N{
									&N{Symbol: cnst, value: 1},
									&N{Symbol: cnst, value: 4},
								},
							},
							&N{Symbol: cnst, value: 3},
						},
					},
				},
			},
		},
		{
			name: "10",
			want: NewEvalResult(10),
			arg:  &N{Symbol: cnst, value: 10},
		},
		{
			name: "-10",
			want: NewEvalResult(-10),
			arg:  &N{Symbol: cnst, value: -10},
		},
		{
			name: "S",
			want: NewEvalResult(nil),
			arg:  &N{Symbol: S},
		},
		{
			name: "exp*4",
			want: NewEvalResult(nil),
			arg: &N{
				Symbol: mult, Children: []*N{
					&N{Symbol: exp},
					&N{Symbol: cnst, value: 4},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := evaluator.Eval(tt.arg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Evaluator.Eval() = %v, want %v", got, tt.want)
			}
		})
	}
}
