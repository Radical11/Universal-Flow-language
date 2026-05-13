package ast

import (
	"fmt"
	"strconv"
	"strings"
)

type Document struct {
	Title     string
	Sections  []Section
	Entities  []Entity
	Relations []Relation
	Flows     []Flow
	States    []State
}

type Section struct {
	Level int
	Title string
}

type SourcePos struct {
	Line   int
	Column int
}

type ValueKind string

const (
	ValueString ValueKind = "string"
	ValueNumber ValueKind = "number"
	ValueBool   ValueKind = "bool"
	ValueArray  ValueKind = "array"
)

type Value struct {
	Kind   ValueKind
	Raw    string
	String string
	Number float64
	Bool   bool
	Array  []Value
}

func NewStringValue(value string) Value {
	return Value{Kind: ValueString, String: value}
}

func NewNumberValue(raw string, value float64) Value {
	return Value{Kind: ValueNumber, Raw: raw, Number: value}
}

func NewBoolValue(value bool) Value {
	raw := "false"
	if value {
		raw = "true"
	}
	return Value{Kind: ValueBool, Raw: raw, Bool: value}
}

func NewArrayValue(values []Value) Value {
	return Value{Kind: ValueArray, Array: values}
}

func (v Value) Any() any {
	switch v.Kind {
	case ValueString:
		return v.String
	case ValueNumber:
		return v.Number
	case ValueBool:
		return v.Bool
	case ValueArray:
		items := make([]any, 0, len(v.Array))
		for _, item := range v.Array {
			items = append(items, item.Any())
		}
		return items
	default:
		return nil
	}
}

func (v Value) Format() string {
	switch v.Kind {
	case ValueString:
		return fmt.Sprintf("%q", v.String)
	case ValueNumber:
		if v.Raw != "" {
			return v.Raw
		}
		return strconv.FormatFloat(v.Number, 'f', -1, 64)
	case ValueBool:
		if v.Raw != "" {
			return v.Raw
		}
		if v.Bool {
			return "true"
		}
		return "false"
	case ValueArray:
		parts := make([]string, 0, len(v.Array))
		for _, item := range v.Array {
			parts = append(parts, item.Format())
		}
		return "[" + strings.Join(parts, ", ") + "]"
	default:
		return `""`
	}
}

type Metadata map[string]Value

type Entity struct {
	ID       string
	Type     string
	Label    string
	Metadata Metadata
	Pos      SourcePos
}

type Relation struct {
	From     string
	To       string
	Type     string
	Metadata Metadata
	Pos      SourcePos
}

type Flow struct {
	ID       string
	Label    string
	Steps    []Step
	Metadata Metadata
	Pos      SourcePos
}

type Step struct {
	ID       string
	Target   string
	Label    string
	Metadata Metadata
	Pos      SourcePos
}

type State struct {
	ID          string
	Label       string
	Transitions []Transition
	Metadata    Metadata
	Pos         SourcePos
}

type Transition struct {
	To        string
	On        string
	Condition string
	Metadata  Metadata
	Pos       SourcePos
}
