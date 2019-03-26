package dsl

import (
	"log"
	"reflect"
	"testing"
)

type PGM = ProgramTree

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
	minus := NewSymbol("minus")
	mult := NewSymbol("mult")
	cnst := NewSymbol("const")
	param := NewSymbol("param")

	var eval func(*ProgramTree, Env) EvalResult
	eval = func(node *ProgramTree, env Env) EvalResult {
		var ret EvalResult
		switch node.Symbol {
		case plus:
			e1, ok := eval(node.Children[0], env).Value()
			if !ok {
				return NewEvalResult(nil)
			}
			e2, ok := eval(node.Children[1], env).Value()
			if !ok {
				return NewEvalResult(nil)
			}
			v1 := e1.(int)
			v2 := e2.(int)
			ret = NewEvalResult(v1 + v2)
		case minus:
			e1, ok := eval(node.Children[0], env).Value()
			if !ok {
				return NewEvalResult(nil)
			}
			e2, ok := eval(node.Children[1], env).Value()
			if !ok {
				return NewEvalResult(nil)
			}
			v1 := e1.(int)
			v2 := e2.(int)
			ret = NewEvalResult(v1 - v2)
		case mult:
			e1, ok := eval(node.Children[0], env).Value()
			if !ok {
				return NewEvalResult(nil)
			}
			e2, ok := eval(node.Children[1], env).Value()
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
		case param:
			i, err := node.Value()
			if !err {
				log.Fatal("the const doesn't hava the value")
			}
			ret = NewEvalResult(env.GetArg(i.(int)))
		default:
			// S, exp,
			if len(node.Children) == 0 {
				return NewEvalResult(nil)
			}
			ret = eval(node.Children[0], env)
		}
		return ret
	}
	evaluator := NewEvaluator(eval)

	type args struct {
		env Env
		pgm *PGM
	}
	tests := []struct {
		name string
		args args
		want EvalResult
	}{
		{
			name: "(1+4)*3",
			want: NewEvalResult(15),
			args: args{
				env: NewEnv(),
				pgm: &PGM{
					Symbol: mult, Children: []*PGM{
						&PGM{
							Symbol: plus, Children: []*PGM{
								&PGM{Symbol: cnst, value: 1},
								&PGM{Symbol: cnst, value: 4},
							},
						},
						&PGM{Symbol: cnst, value: 3},
					},
				},
			},
		},
		{
			name: "10*((1+4)*3)",
			want: NewEvalResult(150),
			args: args{
				env: NewEnv(),
				pgm: &PGM{
					Symbol: mult, Children: []*PGM{
						&PGM{Symbol: cnst, value: 10},
						&PGM{
							Symbol: mult, Children: []*PGM{
								&PGM{
									Symbol: plus, Children: []*PGM{
										&PGM{Symbol: cnst, value: 1},
										&PGM{Symbol: cnst, value: 4},
									},
								},
								&PGM{Symbol: cnst, value: 3},
							},
						},
					},
				},
			},
		},
		{
			name: "10*((1-4)*3)",
			want: NewEvalResult(-90),
			args: args{
				env: NewEnv(),
				pgm: &PGM{
					Symbol: mult, Children: []*PGM{
						&PGM{Symbol: cnst, value: 10},
						&PGM{
							Symbol: mult, Children: []*PGM{
								&PGM{
									Symbol: minus, Children: []*PGM{
										&PGM{Symbol: cnst, value: 1},
										&PGM{Symbol: cnst, value: 4},
									},
								},
								&PGM{Symbol: cnst, value: 3},
							},
						},
					},
				},
			},
		},
		{
			name: "10",
			want: NewEvalResult(10),
			args: args{
				env: NewEnv(),
				pgm: &PGM{Symbol: cnst, value: 10},
			},
		},
		{
			name: "-10",
			want: NewEvalResult(-10),
			args: args{
				env: NewEnv(),
				pgm: &PGM{Symbol: cnst, value: -10},
			},
		},
		{
			name: "S",
			want: NewEvalResult(nil),
			args: args{
				env: NewEnv(),
				pgm: &PGM{Symbol: S},
			},
		},
		{
			name: "exp*4",
			want: NewEvalResult(nil),
			args: args{
				env: NewEnv(),
				pgm: &PGM{
					Symbol: mult, Children: []*PGM{
						&PGM{Symbol: exp},
						&PGM{Symbol: cnst, value: 4},
					},
				},
			},
		},
		{
			name: "param(0) with Env[100]",
			want: NewEvalResult(100),
			args: args{
				env: NewEnv(100),
				pgm: &PGM{Symbol: param, value: 0},
			},
		},
		{
			name: "param(1) with Env[100, 200]",
			want: NewEvalResult(200),
			args: args{
				env: NewEnv(100, 200),
				pgm: &PGM{Symbol: param, value: 1},
			},
		},
		{
			name: "(param(1)+4)*param(3)",
			want: NewEvalResult(-24),
			args: args{
				env: NewEnv(10, 2, 4, -4),
				pgm: &PGM{
					Symbol: mult, Children: []*PGM{
						&PGM{
							Symbol: plus, Children: []*PGM{
								&PGM{Symbol: param, value: 1},
								&PGM{Symbol: cnst, value: 4},
							},
						},
						&PGM{Symbol: param, value: 3},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := evaluator.Eval(tt.args.pgm, tt.args.env); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Evaluator.Eval() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProgramTree_Clone(t *testing.T) {
	plus := NewSymbol("add")
	mult := NewSymbol("mult")
	cnst := NewSymbol("const")

	tests := []struct {
		name string
		org  *ProgramTree
	}{
		{
			name: "(1+4)*3",
			org: &PGM{
				Symbol: mult, Children: []*PGM{
					&PGM{
						Symbol: plus, Children: []*PGM{
							&PGM{Symbol: cnst, value: 1},
							&PGM{Symbol: cnst, value: 4},
						},
					},
					&PGM{Symbol: cnst, value: 3},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.org.Clone(); got.String() != tt.org.String() {
				t.Errorf("ProgramTree.Clone() = %v, want %v", got, tt.org)
			}
		})
	}
}

func TestProgramTree_Leaves(t *testing.T) {
	plus := NewSymbol("add")
	mult := NewSymbol("mult")
	cnst := NewSymbol("const")

	tests := []struct {
		name     string
		target   *ProgramTree
		wantSize int
	}{
		{
			name: "(1+4)*3",
			target: &PGM{
				Symbol: mult, Children: []*PGM{
					&PGM{
						Symbol: plus, Children: []*PGM{
							&PGM{Symbol: cnst, value: 1},
							&PGM{Symbol: cnst, value: 4},
						},
					},
					&PGM{Symbol: cnst, value: 3},
				},
			},
			wantSize: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.target.Leaves(); len(got) != tt.wantSize {
				t.Errorf("ProgramTree.Leaves() = %v, want the size of %v", got, tt.wantSize)
			}
		})
	}
}

func TestProgramTree_NonTerminalLeaves(t *testing.T) {
	exp := NewSymbol("exp")
	exp.isTerminal = false
	plus := NewSymbol("add")
	plus.isTerminal = false
	mult := NewSymbol("mult")
	mult.isTerminal = false
	cnst := NewSymbol("const")
	cnst.isTerminal = true

	tests := []struct {
		name     string
		target   *ProgramTree
		wantSize int
	}{
		{
			name: "(1+4)*3",
			target: &PGM{
				Symbol: mult, Children: []*PGM{
					&PGM{
						Symbol: plus, Children: []*PGM{
							&PGM{Symbol: cnst, value: 1},
							&PGM{Symbol: cnst, value: 4},
						},
					},
					&PGM{Symbol: cnst, value: 3},
				},
			},
			wantSize: 0,
		},
		{
			name: "(exp+exp)*exp",
			target: &PGM{
				Symbol: mult, Children: []*PGM{
					&PGM{
						Symbol: plus, Children: []*PGM{
							&PGM{Symbol: exp},
							&PGM{Symbol: exp},
						},
					},
					&PGM{Symbol: exp},
				},
			},
			wantSize: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.target.NonTerminalLeaves(); len(got) != tt.wantSize {
				t.Errorf("ProgramTree.Leaves() = %v, want the size of %v", got, tt.wantSize)
			}
		})
	}
}
