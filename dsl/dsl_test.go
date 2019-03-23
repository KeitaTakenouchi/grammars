package dsl

import (
	"testing"
)

func TestSymbol_String_01(t *testing.T) {
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

func TestSymbol_String_02(t *testing.T) {
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
