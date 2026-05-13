package lexer

import "testing"

func TestTokenizeCoreSyntax(t *testing.T) {
	tokens, err := Tokenize("# Journey\nentity hero as person label \"Hero\" [mood=curious]\nrel hero -> gate as approaches\n")
	if err != nil {
		t.Fatal(err)
	}
	want := []TokenType{TokenHeading, TokenNewline, TokenIdentifier, TokenIdentifier, TokenIdentifier, TokenIdentifier, TokenIdentifier, TokenString, TokenLBracket, TokenIdentifier, TokenEqual, TokenIdentifier, TokenRBracket, TokenNewline, TokenIdentifier, TokenIdentifier, TokenArrow, TokenIdentifier, TokenIdentifier, TokenIdentifier, TokenNewline, TokenEOF}
	if len(tokens) != len(want) {
		t.Fatalf("got %d tokens, want %d: %#v", len(tokens), len(want), tokens)
	}
	for i, typ := range want {
		if tokens[i].Type != typ {
			t.Fatalf("token %d got %s, want %s", i, tokens[i].Type, typ)
		}
	}
}

func TestTokenizeUnterminatedString(t *testing.T) {
	if _, err := Tokenize("entity a label \"open"); err == nil {
		t.Fatal("expected error")
	}
}
