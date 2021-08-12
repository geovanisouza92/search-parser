package parser

import (
	"strings"
	"testing"

	"github.com/geovanisouza92/search-parser/lexer"
)

func TestParser(t_ *testing.T) {
	t_.Run("empty search", func(t *testing.T) {
		parse(t, "", "")
	})
	t_.Run("string literal", func(t *testing.T) {
		expected := `STRING "steve jobs"`
		parse(t, `"steve jobs"`, expected)
	})
	t_.Run("or comparison with text", func(t *testing.T) {
		expected := `OR
	TEXT jobs
	TEXT gates`
		parse(t, "jobs OR gates", expected)
	})
	t_.Run("and comparison with text", func(t *testing.T) {
		expected := `AND
	TEXT jobs
	TEXT gates`
		parse(t, "jobs AND gates", expected)
	})
	t_.Run("term exclusion with text", func(t *testing.T) {
		expected := `TEXT jobs
- TEXT apple`
		parse(t, "jobs -apple", expected)
	})
	t_.Run("wildcard operator", func(t *testing.T) {})
	t_.Run("grouping", func(t *testing.T) {})
}

func parse(t *testing.T, input, expected string) {
	t.Helper()
	parser := New(lexer.New(strings.NewReader(input)))
	filter := parser.Parse()
	actual := filter.String()
	if actual != expected {
		t.Errorf("Filter representation does not match, expected\n\n%s\n\n\tgot\n\n%s", expected, actual)
	}
}
