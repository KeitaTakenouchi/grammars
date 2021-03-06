package dsl

import (
	"fmt"
	"log"
	"strings"
)

type Symbol struct {
	Id         string
	isTerminal bool
}

func NewSymbol(id string) *Symbol {
	return &Symbol{
		Id:         id,
		isTerminal: true,
	}
}

func (s *Symbol) IsTerminal() bool {
	return s.isTerminal
}

func (s *Symbol) String() string {
	var str string
	if s.isTerminal {
		str += "\"" + s.Id + "\""
	} else {
		str += s.Id
	}
	return str
}

type Grammar struct {
	symbols symbolSet
	rules   rules
	start   *Symbol
}

func NewGrammar(start *Symbol) Grammar {
	start.isTerminal = false
	return Grammar{
		symbols: newSymbols(),
		rules:   newRules(),
		start:   start,
	}
}

func (g *Grammar) GetRhs(leftSymbol *Symbol) [][]*Symbol {
	rhs := g.rules.getRhs(leftSymbol)
	return rhs.getAllSeqs()
}

func (g *Grammar) GetStart() *Symbol {
	return g.start
}

func (g *Grammar) AddRule(left *Symbol, right ...*Symbol) {
	g.symbols.addSymbol(left)
	for _, r := range right {
		g.symbols.addSymbol(r)
	}
	g.rules.addRule(left, newSeqence(right...))
}

func (g Grammar) String() string {
	var str string
	str += "----------------------------------\n"
	str += "START : " + g.start.String() + "\n"
	str += "RULES : \n" + g.rules.String() + "\n"
	str += "SYMBOLS: " + g.symbols.String() + "\n"
	str += "----------------------------------"
	return str
}

type symbolSet struct {
	symbols map[*Symbol]struct{}
}

func newSymbols() symbolSet {
	return symbolSet{
		symbols: make(map[*Symbol]struct{}),
	}
}

func (ss *symbolSet) addSymbol(s *Symbol) {
	ss.symbols[s] = struct{}{}
}

func (ss symbolSet) String() string {
	var strs []string
	for s, _ := range ss.symbols {
		strs = append(strs, s.String())
	}
	return "[" + strings.Join(strs, ", ") + "]"
}

type rules struct {
	ruleMap map[lhs]rhs
}

func newRules() rules {
	return rules{
		ruleMap: make(map[lhs]rhs),
	}
}

func (rs *rules) addRule(lsymbol *Symbol, rsymbols sequence) {
	left := newLhs(lsymbol)
	right, e := rs.ruleMap[left]
	if !e {
		right = newRhs()
	}
	right.addSquence(rsymbols)
	rs.ruleMap[left] = right
}

func (rs *rules) getRhs(lsymbol *Symbol) rhs {
	lhs := newLhs(lsymbol)
	return rs.ruleMap[lhs]
}

func (rs rules) String() string {
	var strs []string
	for left, right := range rs.ruleMap {
		strs = append(strs, left.symbol.String()+" ->"+right.String())
	}
	return " " + strings.Join(strs, "\n ")
}

type lhs struct {
	symbol *Symbol
}

func newLhs(s *Symbol) lhs {
	s.isTerminal = false
	return lhs{
		symbol: s,
	}
}

type rhs struct {
	seqs []sequence
}

func newRhs() rhs {
	return rhs{
		seqs: make([]sequence, 0),
	}
}

func (r *rhs) getAllSeqs() [][]*Symbol {
	ret := make([][]*Symbol, 0)
	for _, seq := range r.seqs {
		ret = append(ret, seq.symbols)
	}
	return ret
}

func (r *rhs) addSquence(s sequence) {
	r.seqs = append(r.seqs, s)
}

func (r rhs) String() string {
	strs := make([]string, len(r.seqs))
	for i, seq := range r.seqs {
		strs[i] = seq.String()
	}
	return " " + strings.Join(strs, "\n\t| ")
}

type sequence struct {
	symbols []*Symbol
}

func newSeqence(symbols ...*Symbol) sequence {
	seq := sequence{}
	for _, s := range symbols {
		seq.addSymbol(s)
	}
	return seq
}

func (sq *sequence) addSymbol(s *Symbol) {
	sq.symbols = append(sq.symbols, s)
}

func (sq sequence) String() string {
	strs := make([]string, len(sq.symbols))
	for i, s := range sq.symbols {
		strs[i] = s.String()
	}
	return strings.Join(strs, " ")
}

type ProgramTree struct {
	Symbol   *Symbol
	Children []*ProgramTree
	value    interface{}
}

func NewProgramTree(s *Symbol) *ProgramTree {
	return &ProgramTree{
		Symbol:   s,
		Children: make([]*ProgramTree, 0),
		value:    nil,
	}
}

func (n *ProgramTree) AddChildren(children ...*ProgramTree) {
	n.Children = append(n.Children, children...)
}

func (n *ProgramTree) With(value interface{}) *ProgramTree {
	n.value = value
	return n
}

func (n *ProgramTree) Value() (interface{}, bool) {
	if n.value == nil {
		return nil, false
	}
	return n.value, true
}

func (n *ProgramTree) Leaves() []*ProgramTree {
	if len(n.Children) == 0 {
		return []*ProgramTree{n}
	}
	ret := make([]*ProgramTree, 0)
	for _, c := range n.Children {
		leaves := c.Leaves()
		ret = append(ret, leaves...)
	}
	return ret
}

func (n *ProgramTree) NonTerminalLeaves() []*ProgramTree {
	ret := make([]*ProgramTree, 0)
	for _, leaf := range n.Leaves() {
		if !leaf.Symbol.IsTerminal() {
			ret = append(ret, leaf)
		}
	}
	return ret
}

func (n *ProgramTree) Clone() *ProgramTree {
	// shallow copy
	thisCpy := NewProgramTree(n.Symbol)
	thisCpy.value = n.value

	// deep copy
	for _, c := range n.Children {
		childCpy := c.Clone()
		thisCpy.AddChildren(childCpy)
	}
	return thisCpy
}

func (n *ProgramTree) String() string {
	var symbolStr string
	if val, ok := n.Value(); ok {
		symbolStr = fmt.Sprintf("%s(%v)", n.Symbol.String(), val)
	} else {
		symbolStr = n.Symbol.String()
	}
	if n.Children == nil || len(n.Children) == 0 {
		return symbolStr
	}
	childStrs := make([]string, len(n.Children))
	for i, c := range n.Children {
		childStrs[i] = c.String()
	}
	return symbolStr + "[" + strings.Join(childStrs, ",") + "]"
}

func (n *ProgramTree) FormattedString() string {
	// function definition
	spaces := func(n int) string {
		var ret string
		for i := 0; i < n-1; i++ {
			ret += "   "
		}
		ret += "---"
		return ret
	}

	// format the string
	var str string
	indent := 0
	for _, r := range n.String() {
		c := string(r)
		switch c {
		case "[":
			indent++
			str += "\n"
			str += spaces(indent)
		case "]":
			indent--
		case ",":
			str += "\n"
			str += spaces(indent)
		default:
			str += c
		}
	}
	return str
}

type Evaluator struct {
	evalFunc func(*ProgramTree, Env) EvalResult
}

func NewEvaluator(eval func(*ProgramTree, Env) EvalResult) Evaluator {
	return Evaluator{
		evalFunc: eval,
	}
}

func (e *Evaluator) Eval(ast *ProgramTree, env Env) EvalResult {
	return e.evalFunc(ast, env)
}

type EvalResult struct {
	value interface{}
}

func NewEvalResult(value interface{}) EvalResult {
	return EvalResult{
		value: value,
	}
}

func (e EvalResult) Value() (interface{}, bool) {
	if e.value != nil {
		return e.value, true
	}
	return nil, false
}

type Env struct {
	args []interface{}
}

func NewEnv(args ...interface{}) Env {
	return Env{
		args: args,
	}
}

func (e *Env) AddArgs(args ...interface{}) {
	e.args = append(e.args, args...)
}

func (e *Env) GetArg(i int) interface{} {
	if i >= len(e.args) {
		log.Fatal("invalid index of args in Environment")
	}
	return e.args[i]
}
