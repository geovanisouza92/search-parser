package token

//go:generate stringer -type TokenType

type TokenType byte

const (
	Error TokenType = iota
	EOF
)
