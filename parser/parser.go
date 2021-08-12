package parser

import (
	"github.com/geovanisouza92/search-parser/ast"
	"github.com/geovanisouza92/search-parser/lexer"
	"github.com/geovanisouza92/search-parser/token"
)

// Priority definitions
const (
	lowest     byte = iota
	logical         // AND OR
	comparator      // <= < >= > != = :
	prefix          // + - NOT
)

var precedences = map[token.TokenType]byte{
	token.And:   logical,
	token.Or:    logical,
	token.Le:    comparator,
	token.Lt:    comparator,
	token.Ge:    comparator,
	token.Gt:    comparator,
	token.Neq:   comparator,
	token.Eq:    comparator,
	token.Colon: comparator,
	token.Plus:  prefix,
	token.Minus: prefix,
	token.Not:   prefix,
}

type prefixParser func() ast.Exp
type infixParser func(ast.Exp) ast.Exp

type Parser struct {
	l    *lexer.Lexer
	c, n token.Token
	pp   map[token.TokenType]prefixParser
	ip   map[token.TokenType]infixParser
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:  l,
		pp: make(map[token.TokenType]prefixParser),
		ip: make(map[token.TokenType]infixParser),
	}
	p.next()
	p.next()

	p.pp[token.Text] = p.parseText
	p.pp[token.String] = p.parseString
	p.pp[token.Plus] = p.parsePrefix
	p.pp[token.Minus] = p.parsePrefix
	p.pp[token.Not] = p.parsePrefix
	p.pp[token.LParen] = p.parseGroup

	p.ip[token.And] = p.parseInfix
	p.ip[token.Or] = p.parseInfix
	p.ip[token.Le] = p.parseInfix
	p.ip[token.Lt] = p.parseInfix
	p.ip[token.Ge] = p.parseInfix
	p.ip[token.Gt] = p.parseInfix
	p.ip[token.Neq] = p.parseInfix
	p.ip[token.Eq] = p.parseInfix
	p.ip[token.Colon] = p.parseInfix

	return p
}

func (p *Parser) Parse() ast.Filter {
	exps := []ast.Exp{}
	for p.c.Type != token.EOF {
		if exp := p.parse(lowest); exp != nil {
			exps = append(exps, exp)
		}
		p.next()
	}
	return ast.Filter{Exps: exps}
}

func (p *Parser) parse(prec byte) ast.Exp {
	prefix := p.pp[p.c.Type]
	if prefix == nil {
		// TODO Error
		return nil
	}

	exp := prefix()
	for prec < p.nextPrec() {
		if infix, ok := p.ip[p.n.Type]; !ok {
			return exp
		} else {
			p.next()
			exp = infix(exp)
		}
	}
	return exp
}

func (p *Parser) parseText() ast.Exp {
	return &ast.Text{Verbatim: p.c}
}

func (p *Parser) parseString() ast.Exp {
	return &ast.String{Verbatim: p.c}
}

func (p *Parser) parsePrefix() ast.Exp {
	e := &ast.Prefix{Op: p.c}
	p.next()
	e.Exp = p.parse(prefix)
	return e
}

func (p *Parser) parseInfix(left ast.Exp) ast.Exp {
	e := ast.Infix{Left: left, Op: p.c}

	prec := p.currPrec()
	p.next()
	e.Right = p.parse(prec)

	return e
}

func (p *Parser) parseGroup() ast.Exp {
	p.next()
	e := p.parse(lowest)
	if !p.assertNext(token.RParen) {
		return nil
	}
	return e
}

func (p *Parser) currPrec() byte {
	if p, ok := precedences[p.c.Type]; ok {
		return p
	}
	return lowest
}

func (p *Parser) nextPrec() byte {
	if p, ok := precedences[p.n.Type]; ok {
		return p
	}
	return lowest
}

func (p *Parser) assertNext(t token.TokenType) bool {
	if p.n.Type == t {
		p.next()
		return true
	}
	// TODO Error
	return false
}

func (p *Parser) next() {
	p.c, p.n = p.n, p.l.Next()
}
