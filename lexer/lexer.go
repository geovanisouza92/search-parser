package lexer

import (
	"io"
	"text/scanner"

	"github.com/geovanisouza92/search-parser/token"
)

type Lexer struct {
	s scanner.Scanner
	r rune
}

func New(r io.Reader) *Lexer {
	var s scanner.Scanner
	s.Init(r)
	l := &Lexer{s: s}
	l.readRune()
	return l
}

func (l *Lexer) Next() token.Token {
	var t token.Token

	switch l.r {
	case '+':
		t = l.token(token.Plus)
	case '-':
		t = l.token(token.Minus)
	case '*':
		t = l.token(token.Wildcard)
	case '=':
		t = l.token(token.Eq)
	case '!':
		l.readRune()
		t = l.token(token.Neq)
	case '>':
		t = l.choose('=', token.Ge, token.Gt)
	case '<':
		t = l.choose('=', token.Le, token.Lt)
	case ',':
		t = l.token(token.Comma)
	case ':':
		t = l.token(token.Colon)
	case '(':
		t = l.token(token.LParen)
	case ')':
		t = l.token(token.RParen)
	case scanner.Ident:
		lit := l.s.TokenText()
		t = token.Token{Type: token.Lookup(lit), Literal: lit}
	case scanner.Int, scanner.Float:
		lit := l.s.TokenText()
		t = token.Token{Type: token.Text, Literal: lit}
	case scanner.String:
		lit := l.s.TokenText()
		t = token.Token{Type: token.String, Literal: lit}
	case scanner.EOF:
		t = l.token(token.EOF)
	default:
		lit := l.s.TokenText()
		t = token.Token{Type: token.Error, Literal: lit}
	}

	l.readRune()
	return t
}

func (l *Lexer) readRune() {
	l.r = l.s.Scan()
}

func (l *Lexer) token(ty token.TokenType) token.Token {
	lit := l.s.TokenText()
	return token.Token{Type: ty, Literal: lit}
}

func (l *Lexer) choose(lookAhead rune, left, right token.TokenType) token.Token {
	lit := l.s.TokenText()
	if l.s.Peek() == lookAhead {
		l.readRune()
		lit += l.s.TokenText()
		return token.Token{Type: left, Literal: lit}
	}
	return token.Token{Type: right, Literal: lit}
}
