package token

//go:generate stringer -type TokenType

type TokenType byte

const (
	Error TokenType = iota
	EOF

	// TEXT is a free-form set of characters without whitespace (WS)
	// or . (DOT) within it. The text may represent a variable, string,
	// number, boolean, or alternative literal value and must be handled
	// in a manner consistent with the service's intention.
	Text

	// STRING is a quoted string which may or may not contain a special
	// wildcard `*` character at the beginning or end of the string to
	// indicate a prefix or suffix-based search within a restriction.
	String

	Plus     // +
	Minus    // -
	Wildcard // *
	Not      // NOT
	Eq       // =
	Neq      // !=
	Gt       // >
	Ge       // >=
	Lt       // <
	Le       // <=
	And      // AND
	Or       // OR

	Comma  // ,
	Colon  // :
	LParen // (
	RParen // )
)

type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"NOT": Not,
	"AND": And,
	"OR":  Or,
}

func Lookup(name string) TokenType {
	if t, ok := keywords[name]; ok {
		return t
	}
	return Text
}
