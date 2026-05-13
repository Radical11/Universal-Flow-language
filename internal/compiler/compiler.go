package compiler

import (
	"fmt"

	"github.com/Radical11/Universal-Flow-language/internal/ast"
	"github.com/Radical11/Universal-Flow-language/internal/ir"
	"github.com/Radical11/Universal-Flow-language/internal/parser"
)

const IRVersion = "ufl.ir.v0"

func Parse(source string) (*ast.Document, error) {
	return parser.Parse(source)
}

func Compile(source string) (*ir.Document, error) {
	doc, err := parser.Parse(source)
	if err != nil {
		return nil, err
	}
	return BuildIR(doc)
}

func BuildIR(doc *ast.Document) (*ir.Document, error) {
	result := &ir.Document{
		Version:  IRVersion,
		Title:    doc.Title,
		Metadata: map[string]any{},
		Nodes:    []ir.Node{},
		Edges:    []ir.Edge{},
		Flows:    []ir.Flow{},
		States:   []ir.State{},
	}
	seen := map[string]bool{}
	for _, entity := range doc.Entities {
		if seen[entity.ID] {
			return nil, fmt.Errorf("duplicate entity id %q", entity.ID)
		}
		seen[entity.ID] = true
		result.Nodes = append(result.Nodes, ir.Node{
			ID:       entity.ID,
			Type:     entity.Type,
			Label:    entity.Label,
			Metadata: copyMap(entity.Metadata),
		})
	}
	for i, rel := range doc.Relations {
		if !seen[rel.From] {
			result.Diagnostics = append(result.Diagnostics, ir.Diagnostic{Severity: "warning", Message: fmt.Sprintf("relation source %q is not declared", rel.From), Line: rel.Pos.Line, Column: rel.Pos.Column})
		}
		if !seen[rel.To] {
			result.Diagnostics = append(result.Diagnostics, ir.Diagnostic{Severity: "warning", Message: fmt.Sprintf("relation target %q is not declared", rel.To), Line: rel.Pos.Line, Column: rel.Pos.Column})
		}
		result.Edges = append(result.Edges, ir.Edge{
			ID:       fmt.Sprintf("edge:%d", i+1),
			From:     rel.From,
			To:       rel.To,
			Type:     rel.Type,
			Metadata: copyMap(rel.Metadata),
		})
	}
	for _, flow := range doc.Flows {
		steps := make([]ir.Step, 0, len(flow.Steps))
		for _, step := range flow.Steps {
			if step.Target != "" && !seen[step.Target] {
				result.Diagnostics = append(result.Diagnostics, ir.Diagnostic{Severity: "warning", Message: fmt.Sprintf("flow %q step %q uses undeclared target %q", flow.ID, step.ID, step.Target), Line: step.Pos.Line, Column: step.Pos.Column})
			}
			steps = append(steps, ir.Step{ID: step.ID, Target: step.Target, Label: step.Label, Metadata: copyMap(step.Metadata)})
		}
		result.Flows = append(result.Flows, ir.Flow{ID: flow.ID, Label: flow.Label, Steps: steps, Metadata: copyMap(flow.Metadata)})
	}
	stateIDs := map[string]bool{}
	for _, state := range doc.States {
		stateIDs[state.ID] = true
	}
	for _, state := range doc.States {
		transitions := make([]ir.Transition, 0, len(state.Transitions))
		for _, transition := range state.Transitions {
			if !stateIDs[transition.To] {
				result.Diagnostics = append(result.Diagnostics, ir.Diagnostic{Severity: "warning", Message: fmt.Sprintf("state %q transitions to undeclared state %q", state.ID, transition.To), Line: transition.Pos.Line, Column: transition.Pos.Column})
			}
			transitions = append(transitions, ir.Transition{To: transition.To, On: transition.On, Condition: transition.Condition, Metadata: copyMap(transition.Metadata)})
		}
		result.States = append(result.States, ir.State{ID: state.ID, Label: state.Label, Transitions: transitions, Metadata: copyMap(state.Metadata)})
	}
	return result, nil
}

func copyMap(in ast.Metadata) map[string]any {
	if len(in) == 0 {
		return nil
	}
	out := make(map[string]any, len(in))
	for k, v := range in {
		out[k] = v.Any()
	}
	return out
}
