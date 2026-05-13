package lexer

import (
	"fmt"
	"strings"
	"unicode"
)

type Lexer struct {
	input  []rune
	pos    int
	line   int
	column int
}

func New(input string) *Lexer {
	return &Lexer{input: []rune(input), line: 1, column: 1}
}

func Tokenize(input string) ([]Token, error) {
	return New(input).Tokens()
}

func (l *Lexer) Tokens() ([]Token, error) {
	var tokens []Token
	for {
		tok, err := l.Next()
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, tok)
		if tok.Type == TokenEOF {
			return tokens, nil
		}
	}
}

func (l *Lexer) Next() (Token, error) {
	for l.hasNext() {
		ch := l.peek()
		if ch == ' ' || ch == '\t' || ch == '\r' {
			l.advance()
			continue
		}
		if ch == '/' && l.peekN(1) == '/' {
			l.skipLineComment()
			continue
		}
		if ch == '#' {
			return l.heading()
		}
		break
	}

	if !l.hasNext() {
		return Token{Type: TokenEOF, Line: l.line, Column: l.column}, nil
	}

	line, col := l.line, l.column
	ch := l.peek()
	switch {
	case ch == '\n':
		l.advance()
		return Token{Type: TokenNewline, Line: line, Column: col}, nil
	case ch == '-' && l.peekN(1) == '>':
		l.advance()
		l.advance()
		return Token{Type: TokenArrow, Value: "->", Line: line, Column: col}, nil
	case ch == ':':
		l.advance()
		return Token{Type: TokenColon, Value: ":", Line: line, Column: col}, nil
	case ch == ',':
		l.advance()
		return Token{Type: TokenComma, Value: ",", Line: line, Column: col}, nil
	case ch == '{':
		l.advance()
		return Token{Type: TokenLBrace, Value: "{", Line: line, Column: col}, nil
	case ch == '}':
		l.advance()
		return Token{Type: TokenRBrace, Value: "}", Line: line, Column: col}, nil
	case ch == '[':
		l.advance()
		return Token{Type: TokenLBracket, Value: "[", Line: line, Column: col}, nil
	case ch == ']':
		l.advance()
		return Token{Type: TokenRBracket, Value: "]", Line: line, Column: col}, nil
	case ch == '=':
		l.advance()
		return Token{Type: TokenEqual, Value: "=", Line: line, Column: col}, nil
	case ch == '"':
		return l.string()
	case unicode.IsDigit(ch):
		return l.number()
	case isIdentStart(ch):
		return l.identifier()
	default:
		return Token{}, fmt.Errorf("unexpected character %q at %d:%d", ch, line, col)
	}
}

func (l *Lexer) heading() (Token, error) {
	line, col := l.line, l.column
	depth := 0
	for l.hasNext() && l.peek() == '#' {
		depth++
		l.advance()
	}
	for l.hasNext() && (l.peek() == ' ' || l.peek() == '\t') {
		l.advance()
	}
	var b strings.Builder
	for l.hasNext() && l.peek() != '\n' {
		b.WriteRune(l.advance())
	}
	value := strings.TrimSpace(b.String())
	if value == "" {
		return Token{}, fmt.Errorf("empty heading at %d:%d", line, col)
	}
	return Token{Type: TokenHeading, Value: strings.Repeat("#", depth) + " " + value, Line: line, Column: col}, nil
}

func (l *Lexer) string() (Token, error) {
	line, col := l.line, l.column
	l.advance()
	var b strings.Builder
	for l.hasNext() {
		ch := l.advance()
		if ch == '"' {
			return Token{Type: TokenString, Value: b.String(), Line: line, Column: col}, nil
		}
		if ch == '\\' && l.hasNext() {
			next := l.advance()
			switch next {
			case 'n':
				b.WriteRune('\n')
			case 't':
				b.WriteRune('\t')
			case '"', '\\':
				b.WriteRune(next)
			default:
				return Token{}, fmt.Errorf("unsupported escape \\%c at %d:%d", next, l.line, l.column)
			}
			continue
		}
		b.WriteRune(ch)
	}
	return Token{}, fmt.Errorf("unterminated string at %d:%d", line, col)
}

func (l *Lexer) number() (Token, error) {
	line, col := l.line, l.column
	var b strings.Builder
	for l.hasNext() && (unicode.IsDigit(l.peek()) || l.peek() == '.') {
		b.WriteRune(l.advance())
	}
	return Token{Type: TokenNumber, Value: b.String(), Line: line, Column: col}, nil
}

func (l *Lexer) identifier() (Token, error) {
	line, col := l.line, l.column
	var b strings.Builder
	for l.hasNext() && isIdentPart(l.peek()) {
		b.WriteRune(l.advance())
	}
	return Token{Type: TokenIdentifier, Value: b.String(), Line: line, Column: col}, nil
}

func (l *Lexer) skipLineComment() {
	for l.hasNext() && l.peek() != '\n' {
		l.advance()
	}
}

func (l *Lexer) hasNext() bool {
	return l.pos < len(l.input)
}

func (l *Lexer) peek() rune {
	return l.input[l.pos]
}

func (l *Lexer) peekN(n int) rune {
	if l.pos+n >= len(l.input) {
		return 0
	}
	return l.input[l.pos+n]
}

func (l *Lexer) advance() rune {
	ch := l.input[l.pos]
	l.pos++
	if ch == '\n' {
		l.line++
		l.column = 1
	} else {
		l.column++
	}
	return ch
}

func isIdentStart(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_'
}

func isIdentPart(ch rune) bool {
	return isIdentStart(ch) || unicode.IsDigit(ch) || ch == '-' || ch == '.'
}
