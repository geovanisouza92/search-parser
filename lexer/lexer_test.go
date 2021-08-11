package lexer

import (
	"strings"
	"testing"

	"github.com/geovanisouza92/search-parser/token"
)

type tl struct {
	ty  token.TokenType
	lit string
}

func TestNextToken(t_ *testing.T) {
	t_.Run("empty search", func(t *testing.T) {
		lex(t, "", []tl{
			{token.EOF, ""},
		})
	})
	t_.Run("string literal", func(t *testing.T) {
		lex(t, `"steve jobs"`, []tl{
			{token.String, `"steve jobs"`},
			{token.EOF, ""},
		})
	})
	t_.Run("or comparison with text", func(t *testing.T) {
		lex(t, "jobs OR gates", []tl{
			{token.Text, "jobs"},
			{token.Or, "OR"},
			{token.Text, "gates"},
			{token.EOF, ""},
		})
	})
	t_.Run("and comparison with text", func(t *testing.T) {
		lex(t, "jobs AND gates", []tl{
			{token.Text, "jobs"},
			{token.And, "AND"},
			{token.Text, "gates"},
			{token.EOF, ""},
		})
	})
	t_.Run("term exclusion with text", func(t *testing.T) {
		lex(t, "jobs -apple", []tl{
			{token.Text, "jobs"},
			{token.Minus, "-"},
			{token.Text, "apple"},
			{token.EOF, ""},
		})
	})
	t_.Run("wildcard operator", func(t *testing.T) {
		lex(t, "steve * apple", []tl{
			{token.Text, "steve"},
			{token.Wildcard, "*"},
			{token.Text, "apple"},
			{token.EOF, ""},
		})
	})
	t_.Run("grouping", func(t *testing.T) {
		lex(t, "(ipad OR iphone) apple", []tl{
			{token.LParen, "("},
			{token.Text, "ipad"},
			{token.Or, "OR"},
			{token.Text, "iphone"},
			{token.RParen, ")"},
			{token.Text, "apple"},
			{token.EOF, ""},
		})
	})
}

func lex(t *testing.T, input string, seq []tl) {
	t.Helper()
	l := New(strings.NewReader(input))

	for _, tok := range seq {
		token := l.Next()
		if token.Type != tok.ty {
			t.Errorf("wrong token type: %v; want %v", token.Type, tok.ty)
		}
		if token.Literal != tok.lit {
			t.Errorf("wrong literal: %v; want %v", token.Literal, tok.lit)
		}
	}
}
