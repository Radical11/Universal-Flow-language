package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Radical11/Universal-Flow-language/internal/ast"
	"github.com/Radical11/Universal-Flow-language/internal/lexer"
)

type Parser struct {
	tokens []lexer.Token
	pos    int
}

func Parse(input string) (*ast.Document, error) {
	tokens, err := lexer.Tokenize(input)
	if err != nil {
		return nil, err
	}
	return New(tokens).Parse()
}

func New(tokens []lexer.Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) Parse() (*ast.Document, error) {
	doc := &ast.Document{}
	for !p.check(lexer.TokenEOF) {
		p.skipNewlines()
		if p.check(lexer.TokenEOF) {
			break
		}
		switch {
		case p.match(lexer.TokenHeading):
			tok := p.previous()
			level, title := parseHeading(tok.Value)
			if doc.Title == "" && level == 1 {
				doc.Title = title
			}
			doc.Sections = append(doc.Sections, ast.Section{Level: level, Title: title})
			p.consumeLine()
		case p.matchKeyword("entity"):
			entity, err := p.parseEntity()
			if err != nil {
				return nil, err
			}
			doc.Entities = append(doc.Entities, entity)
		case p.matchKeyword("rel"):
			rel, err := p.parseRelation()
			if err != nil {
				return nil, err
			}
			doc.Relations = append(doc.Relations, rel)
		case p.matchKeyword("flow"):
			flow, err := p.parseFlow()
			if err != nil {
				return nil, err
			}
			doc.Flows = append(doc.Flows, flow)
		case p.matchKeyword("state"):
			state, err := p.parseState()
			if err != nil {
				return nil, err
			}
			doc.States = append(doc.States, state)
		default:
			return nil, p.errorAtCurrent("expected heading, entity, rel, flow, or state")
		}
	}
	return doc, nil
}

func (p *Parser) parseEntity() (ast.Entity, error) {
	id, err := p.consumeIdent("expected entity id")
	if err != nil {
		return ast.Entity{}, err
	}
	entity := ast.Entity{ID: id.Value, Type: "entity", Metadata: ast.Metadata{}, Pos: ast.SourcePos{Line: id.Line, Column: id.Column}}
	for !p.atLineEnd() {
		if p.matchKeyword("as") {
			typ, err := p.consumeIdent("expected entity type after as")
			if err != nil {
				return ast.Entity{}, err
			}
			entity.Type = typ.Value
			continue
		}
		if p.matchKeyword("label") {
			label, err := p.consumeTextValue("expected label value")
			if err != nil {
				return ast.Entity{}, err
			}
			entity.Label = label
			continue
		}
		if p.check(lexer.TokenLBracket) {
			meta, err := p.parseMetadata()
			if err != nil {
				return ast.Entity{}, err
			}
			entity.Metadata = meta
			continue
		}
		return ast.Entity{}, p.errorAtCurrent("unexpected entity field")
	}
	p.consumeLine()
	return entity, nil
}

func (p *Parser) parseRelation() (ast.Relation, error) {
	from, err := p.consumeIdent("expected relation source")
	if err != nil {
		return ast.Relation{}, err
	}
	if _, err := p.consume(lexer.TokenArrow, "expected -> in relation"); err != nil {
		return ast.Relation{}, err
	}
	to, err := p.consumeIdent("expected relation target")
	if err != nil {
		return ast.Relation{}, err
	}
	rel := ast.Relation{From: from.Value, To: to.Value, Type: "relates", Metadata: ast.Metadata{}}
	rel.Pos = ast.SourcePos{Line: from.Line, Column: from.Column}
	if p.matchKeyword("as") {
		typ, err := p.consumeIdent("expected relation type after as")
		if err != nil {
			return ast.Relation{}, err
		}
		rel.Type = typ.Value
	}
	if p.check(lexer.TokenLBracket) {
		meta, err := p.parseMetadata()
		if err != nil {
			return ast.Relation{}, err
		}
		rel.Metadata = meta
	}
	if !p.atLineEnd() {
		return ast.Relation{}, p.errorAtCurrent("unexpected relation field")
	}
	p.consumeLine()
	return rel, nil
}

func (p *Parser) parseFlow() (ast.Flow, error) {
	id, err := p.consumeIdent("expected flow id")
	if err != nil {
		return ast.Flow{}, err
	}
	flow := ast.Flow{ID: id.Value, Metadata: ast.Metadata{}, Pos: ast.SourcePos{Line: id.Line, Column: id.Column}}
	for !p.atLineEnd() {
		if p.matchKeyword("label") {
			label, err := p.consumeTextValue("expected flow label")
			if err != nil {
				return ast.Flow{}, err
			}
			flow.Label = label
			continue
		}
		if p.check(lexer.TokenLBracket) {
			meta, err := p.parseMetadata()
			if err != nil {
				return ast.Flow{}, err
			}
			flow.Metadata = meta
			continue
		}
		return ast.Flow{}, p.errorAtCurrent("unexpected flow field")
	}
	p.consumeLine()
	if _, err := p.consume(lexer.TokenLBrace, "expected { after flow header"); err != nil {
		return ast.Flow{}, err
	}
	p.consumeLine()
	for !p.check(lexer.TokenRBrace) && !p.check(lexer.TokenEOF) {
		p.skipNewlines()
		if p.check(lexer.TokenRBrace) {
			break
		}
		if !p.matchKeyword("step") {
			return ast.Flow{}, p.errorAtCurrent("expected step inside flow")
		}
		step, err := p.parseStep()
		if err != nil {
			return ast.Flow{}, err
		}
		flow.Steps = append(flow.Steps, step)
	}
	if _, err := p.consume(lexer.TokenRBrace, "expected } after flow"); err != nil {
		return ast.Flow{}, err
	}
	p.consumeLine()
	return flow, nil
}

func (p *Parser) parseStep() (ast.Step, error) {
	id, err := p.consumeIdent("expected step id")
	if err != nil {
		return ast.Step{}, err
	}
	step := ast.Step{ID: id.Value, Metadata: ast.Metadata{}, Pos: ast.SourcePos{Line: id.Line, Column: id.Column}}
	for !p.atLineEnd() {
		if p.matchKeyword("uses") {
			target, err := p.consumeIdent("expected step target after uses")
			if err != nil {
				return ast.Step{}, err
			}
			step.Target = target.Value
			continue
		}
		if p.matchKeyword("label") {
			label, err := p.consumeTextValue("expected step label")
			if err != nil {
				return ast.Step{}, err
			}
			step.Label = label
			continue
		}
		if p.check(lexer.TokenLBracket) {
			meta, err := p.parseMetadata()
			if err != nil {
				return ast.Step{}, err
			}
			step.Metadata = meta
			continue
		}
		return ast.Step{}, p.errorAtCurrent("unexpected step field")
	}
	p.consumeLine()
	return step, nil
}

func (p *Parser) parseState() (ast.State, error) {
	id, err := p.consumeIdent("expected state id")
	if err != nil {
		return ast.State{}, err
	}
	state := ast.State{ID: id.Value, Metadata: ast.Metadata{}, Pos: ast.SourcePos{Line: id.Line, Column: id.Column}}
	for !p.atLineEnd() {
		if p.matchKeyword("label") {
			label, err := p.consumeTextValue("expected state label")
			if err != nil {
				return ast.State{}, err
			}
			state.Label = label
			continue
		}
		if p.check(lexer.TokenLBracket) {
			meta, err := p.parseMetadata()
			if err != nil {
				return ast.State{}, err
			}
			state.Metadata = meta
			continue
		}
		return ast.State{}, p.errorAtCurrent("unexpected state field")
	}
	p.consumeLine()
	if _, err := p.consume(lexer.TokenLBrace, "expected { after state header"); err != nil {
		return ast.State{}, err
	}
	p.consumeLine()
	for !p.check(lexer.TokenRBrace) && !p.check(lexer.TokenEOF) {
		p.skipNewlines()
		if p.check(lexer.TokenRBrace) {
			break
		}
		if !p.matchKeyword("on") {
			return ast.State{}, p.errorAtCurrent("expected on inside state")
		}
		transition, err := p.parseTransition()
		if err != nil {
			return ast.State{}, err
		}
		state.Transitions = append(state.Transitions, transition)
	}
	if _, err := p.consume(lexer.TokenRBrace, "expected } after state"); err != nil {
		return ast.State{}, err
	}
	p.consumeLine()
	return state, nil
}

func (p *Parser) parseTransition() (ast.Transition, error) {
	event, err := p.consumeTextValue("expected transition event")
	if err != nil {
		return ast.Transition{}, err
	}
	if _, err := p.consume(lexer.TokenArrow, "expected -> in transition"); err != nil {
		return ast.Transition{}, err
	}
	to, err := p.consumeIdent("expected transition target")
	if err != nil {
		return ast.Transition{}, err
	}
	transition := ast.Transition{On: event, To: to.Value, Metadata: ast.Metadata{}, Pos: ast.SourcePos{Line: to.Line, Column: to.Column}}
	for !p.atLineEnd() {
		if p.matchKeyword("if") {
			condition, err := p.consumeTextValue("expected transition condition")
			if err != nil {
				return ast.Transition{}, err
			}
			transition.Condition = condition
			continue
		}
		if p.check(lexer.TokenLBracket) {
			meta, err := p.parseMetadata()
			if err != nil {
				return ast.Transition{}, err
			}
			transition.Metadata = meta
			continue
		}
		return ast.Transition{}, p.errorAtCurrent("unexpected transition field")
	}
	p.consumeLine()
	return transition, nil
}

func (p *Parser) parseMetadata() (ast.Metadata, error) {
	if _, err := p.consume(lexer.TokenLBracket, "expected ["); err != nil {
		return nil, err
	}
	meta := ast.Metadata{}
	for !p.check(lexer.TokenRBracket) {
		key, err := p.consumeIdent("expected metadata key")
		if err != nil {
			return nil, err
		}
		if _, err := p.consume(lexer.TokenEqual, "expected = after metadata key"); err != nil {
			return nil, err
		}
		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}
		meta[key.Value] = value
		if !p.match(lexer.TokenComma) {
			break
		}
	}
	if _, err := p.consume(lexer.TokenRBracket, "expected ] after metadata"); err != nil {
		return nil, err
	}
	return meta, nil
}

func (p *Parser) parseValue() (ast.Value, error) {
	if p.match(lexer.TokenString) {
		return ast.NewStringValue(p.previous().Value), nil
	}
	if p.match(lexer.TokenNumber) {
		raw := p.previous().Value
		number, ok := ParseNumber(raw)
		if !ok {
			return ast.Value{}, p.errorAtCurrent("invalid numeric metadata value")
		}
		return ast.NewNumberValue(raw, number), nil
	}
	if p.match(lexer.TokenIdentifier) {
		value := p.previous().Value
		switch value {
		case "true":
			return ast.NewBoolValue(true), nil
		case "false":
			return ast.NewBoolValue(false), nil
		default:
			return ast.NewStringValue(value), nil
		}
	}
	if p.match(lexer.TokenLBracket) {
		values := []ast.Value{}
		for !p.check(lexer.TokenRBracket) {
			value, err := p.parseValue()
			if err != nil {
				return ast.Value{}, err
			}
			values = append(values, value)
			if !p.match(lexer.TokenComma) {
				break
			}
		}
		if _, err := p.consume(lexer.TokenRBracket, "expected ] after array value"); err != nil {
			return ast.Value{}, err
		}
		return ast.NewArrayValue(values), nil
	}
	return ast.Value{}, p.errorAtCurrent("expected metadata value")
}

func (p *Parser) consumeTextValue(message string) (string, error) {
	if p.match(lexer.TokenString, lexer.TokenIdentifier, lexer.TokenNumber) {
		return p.previous().Value, nil
	}
	return "", p.errorAtCurrent(message)
}

func (p *Parser) consumeIdent(message string) (lexer.Token, error) {
	return p.consume(lexer.TokenIdentifier, message)
}

func (p *Parser) consume(t lexer.TokenType, message string) (lexer.Token, error) {
	if p.check(t) {
		return p.advance(), nil
	}
	return lexer.Token{}, p.errorAtCurrent(message)
}

func (p *Parser) match(types ...lexer.TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) matchKeyword(keyword string) bool {
	if p.check(lexer.TokenIdentifier) && p.peek().Value == keyword {
		p.advance()
		return true
	}
	return false
}

func (p *Parser) check(t lexer.TokenType) bool {
	return p.peek().Type == t
}

func (p *Parser) advance() lexer.Token {
	if !p.check(lexer.TokenEOF) {
		p.pos++
	}
	return p.previous()
}

func (p *Parser) peek() lexer.Token {
	if p.pos >= len(p.tokens) {
		return p.tokens[len(p.tokens)-1]
	}
	return p.tokens[p.pos]
}

func (p *Parser) previous() lexer.Token {
	return p.tokens[p.pos-1]
}

func (p *Parser) skipNewlines() {
	for p.match(lexer.TokenNewline) {
	}
}

func (p *Parser) consumeLine() {
	for !p.check(lexer.TokenEOF) && !p.check(lexer.TokenNewline) {
		p.advance()
	}
	p.match(lexer.TokenNewline)
}

func (p *Parser) atLineEnd() bool {
	return p.check(lexer.TokenNewline) || p.check(lexer.TokenEOF)
}

func (p *Parser) errorAtCurrent(message string) error {
	tok := p.peek()
	return fmt.Errorf("%s at %d:%d near %s", message, tok.Line, tok.Column, tok.String())
}

func parseHeading(value string) (int, string) {
	parts := strings.SplitN(value, " ", 2)
	level := strings.Count(parts[0], "#")
	if len(parts) == 1 {
		return level, ""
	}
	return level, strings.TrimSpace(parts[1])
}

func ParseNumber(value string) (float64, bool) {
	n, err := strconv.ParseFloat(value, 64)
	return n, err == nil
}
