package lexer

import "fmt"

type TokenType string

const (
	TokenEOF        TokenType = "EOF"
	TokenNewline    TokenType = "NEWLINE"
	TokenHeading    TokenType = "HEADING"
	TokenIdentifier TokenType = "IDENTIFIER"
	TokenString     TokenType = "STRING"
	TokenNumber     TokenType = "NUMBER"
	TokenArrow      TokenType = "ARROW"
	TokenColon      TokenType = "COLON"
	TokenComma      TokenType = "COMMA"
	TokenLBrace     TokenType = "LBRACE"
	TokenRBrace     TokenType = "RBRACE"
	TokenLBracket   TokenType = "LBRACKET"
	TokenRBracket   TokenType = "RBRACKET"
	TokenEqual      TokenType = "EQUAL"
)

type Token struct {
	Type   TokenType
	Value  string
	Line   int
	Column int
}

func (t Token) String() string {
	if t.Value == "" {
		return fmt.Sprintf("%s@%d:%d", t.Type, t.Line, t.Column)
	}
	return fmt.Sprintf("%s(%q)@%d:%d", t.Type, t.Value, t.Line, t.Column)
}
