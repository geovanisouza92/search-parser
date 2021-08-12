package ast

import (
	"bytes"

	"github.com/geovanisouza92/search-parser/token"
)

type Filter struct {
	Exps []Exp
}

func (f Filter) String() string {
	var b bytes.Buffer
	for _, exp := range f.Exps {
		b.WriteString(exp.String())
		b.WriteRune('\n')
	}
	if b.Len() > 0 {
		b.Truncate(b.Len() - 1)
	}
	return b.String()
}

type Exp interface {
	isExp()
	String() string
}

var (
	Empty Exp = empty{}
	_     Exp = Text{}
	_     Exp = String{}
	_     Exp = Infix{}
	_     Exp = Prefix{}
)

type empty struct{}

func (empty) isExp()         {}
func (empty) String() string { return "" }

type Text struct {
	Verbatim token.Token
}

func (Text) isExp()           {}
func (t Text) String() string { return "TEXT " + t.Verbatim.Literal }

type String struct {
	Verbatim token.Token
}

func (String) isExp()           {}
func (s String) String() string { return "STRING " + s.Verbatim.Literal }

type Prefix struct {
	Op  token.Token
	Exp Exp
}

func (Prefix) isExp() {}

func (p Prefix) String() string {
	return p.Op.Literal + " " + p.Exp.String()
}

type Infix struct {
	Op          token.Token
	Left, Right Exp
}

func (Infix) isExp() {}

func (i Infix) String() string {
	return i.Op.Literal + "\n\t" + i.Left.String() + "\n\t" + i.Right.String()
}
