package compiler

import (
	"fmt"
	"strings"

	"github.com/Radical11/Universal-Flow-language/internal/ast"
	"github.com/Radical11/Universal-Flow-language/internal/ir"
	"github.com/Radical11/Universal-Flow-language/internal/parser"
)

const (
	ProfileWorkflow     = "workflow"
	ProfileStory        = "story"
	ProfileStateMachine = "state-machine"
	ProfileSystemMap    = "system-map"
)

var validProfiles = map[string]bool{
	ProfileWorkflow:     true,
	ProfileStory:        true,
	ProfileStateMachine: true,
	ProfileSystemMap:    true,
}

func NormalizeProfile(profile string) (string, error) {
	profile = strings.TrimSpace(strings.ToLower(profile))
	if profile == "" {
		return "", nil
	}
	if !validProfiles[profile] {
		return "", fmt.Errorf("unsupported profile %q", profile)
	}
	return profile, nil
}

func CompileWithProfile(source, profile string) (*ir.Document, error) {
	doc, err := parser.Parse(source)
	if err != nil {
		return nil, err
	}
	profile, err = NormalizeProfile(profile)
	if err != nil {
		return nil, err
	}
	irDoc, err := BuildIR(doc)
	if err != nil {
		return nil, err
	}
	if profile != "" {
		if irDoc.Metadata == nil {
			irDoc.Metadata = map[string]any{}
		}
		irDoc.Metadata["profile"] = profile
	}
	if profile != "" {
		irDoc.Diagnostics = append(irDoc.Diagnostics, ValidateProfile(doc, profile)...)
	}
	return irDoc, nil
}

func ValidateProfile(doc *ast.Document, profile string) []ir.Diagnostic {
	switch profile {
	case ProfileWorkflow:
		return validateWorkflowProfile(doc)
	case ProfileStory:
		return validateStoryProfile(doc)
	case ProfileStateMachine:
		return validateStateMachineProfile(doc)
	case ProfileSystemMap:
		return validateSystemMapProfile(doc)
	default:
		return nil
	}
}

func validateWorkflowProfile(doc *ast.Document) []ir.Diagnostic {
	var diagnostics []ir.Diagnostic
	allowedTypes := stringSet("actor", "artifact", "system", "decision")
	allowedRelations := stringSet("assigned_to", "produces", "depends_on", "approves", "rejects")
	if len(doc.Flows) == 0 {
		diagnostics = append(diagnostics, ir.Diagnostic{Severity: "warning", Message: `workflow profile expects at least one "flow" block`})
	}
	for _, entity := range doc.Entities {
		if entity.Type != "entity" && !allowedTypes[entity.Type] {
			diagnostics = append(diagnostics, profileWarning(entity.Pos, fmt.Sprintf("workflow profile does not recognize entity type %q", entity.Type)))
		}
	}
	for _, rel := range doc.Relations {
		if !allowedRelations[rel.Type] {
			diagnostics = append(diagnostics, profileWarning(rel.Pos, fmt.Sprintf("workflow profile does not recognize relation type %q", rel.Type)))
		}
	}
	return diagnostics
}

func validateStoryProfile(doc *ast.Document) []ir.Diagnostic {
	var diagnostics []ir.Diagnostic
	allowedTypes := stringSet("person", "place", "event", "object", "theme")
	allowedRelations := stringSet("begins_at", "travels_to", "conflicts_with", "reveals", "changes")
	if len(doc.Flows) == 0 && len(doc.Relations) == 0 {
		diagnostics = append(diagnostics, ir.Diagnostic{Severity: "warning", Message: "story profile expects at least one relation or flow"})
	}
	for _, entity := range doc.Entities {
		if entity.Type != "entity" && !allowedTypes[entity.Type] {
			diagnostics = append(diagnostics, profileWarning(entity.Pos, fmt.Sprintf("story profile does not recognize entity type %q", entity.Type)))
		}
	}
	for _, rel := range doc.Relations {
		if !allowedRelations[rel.Type] {
			diagnostics = append(diagnostics, profileWarning(rel.Pos, fmt.Sprintf("story profile does not recognize relation type %q", rel.Type)))
		}
	}
	return diagnostics
}

func validateStateMachineProfile(doc *ast.Document) []ir.Diagnostic {
	var diagnostics []ir.Diagnostic
	if len(doc.States) == 0 {
		diagnostics = append(diagnostics, ir.Diagnostic{Severity: "warning", Message: `state-machine profile expects at least one "state" block`})
		return diagnostics
	}
	initialCount := 0
	for _, state := range doc.States {
		if value, ok := state.Metadata["initial"]; ok && value.Kind == ast.ValueBool && value.Bool {
			initialCount++
		}
		isTerminal := false
		if value, ok := state.Metadata["terminal"]; ok && value.Kind == ast.ValueBool && value.Bool {
			isTerminal = true
		}
		if len(state.Transitions) == 0 && !isTerminal {
			diagnostics = append(diagnostics, profileWarning(state.Pos, fmt.Sprintf("state %q has no outgoing transitions", state.ID)))
		}
	}
	switch {
	case initialCount == 0:
		diagnostics = append(diagnostics, ir.Diagnostic{Severity: "warning", Message: `state-machine profile expects one state with metadata [initial=true]`})
	case initialCount > 1:
		diagnostics = append(diagnostics, ir.Diagnostic{Severity: "warning", Message: "state-machine profile expects only one initial state"})
	}
	return diagnostics
}

func validateSystemMapProfile(doc *ast.Document) []ir.Diagnostic {
	var diagnostics []ir.Diagnostic
	allowedTypes := stringSet("actor", "surface", "service", "store", "queue", "external")
	allowedRelations := stringSet("calls", "reads_from", "writes_to", "publishes", "subscribes", "returns_to", "interacts_with")
	if len(doc.Entities) == 0 || len(doc.Relations) == 0 {
		diagnostics = append(diagnostics, ir.Diagnostic{Severity: "warning", Message: "system-map profile expects both entities and relations"})
	}
	for _, entity := range doc.Entities {
		if entity.Type != "entity" && !allowedTypes[entity.Type] {
			diagnostics = append(diagnostics, profileWarning(entity.Pos, fmt.Sprintf("system-map profile does not recognize entity type %q", entity.Type)))
		}
	}
	for _, rel := range doc.Relations {
		if !allowedRelations[rel.Type] {
			diagnostics = append(diagnostics, profileWarning(rel.Pos, fmt.Sprintf("system-map profile does not recognize relation type %q", rel.Type)))
		}
	}
	return diagnostics
}

func profileWarning(pos ast.SourcePos, message string) ir.Diagnostic {
	return ir.Diagnostic{Severity: "warning", Message: message, Line: pos.Line, Column: pos.Column}
}

func stringSet(values ...string) map[string]bool {
	out := make(map[string]bool, len(values))
	for _, value := range values {
		out[value] = true
	}
	return out
}
