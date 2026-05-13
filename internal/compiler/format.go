package compiler

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/Radical11/Universal-Flow-language/internal/ast"
)

func Format(w io.Writer, doc *ast.Document) error {
	if doc.Title != "" {
		if _, err := fmt.Fprintf(w, "# %s\n\n", doc.Title); err != nil {
			return err
		}
	}
	for _, entity := range doc.Entities {
		if _, err := fmt.Fprintf(w, "entity %s", entity.ID); err != nil {
			return err
		}
		if entity.Type != "" && entity.Type != "entity" {
			if _, err := fmt.Fprintf(w, " as %s", entity.Type); err != nil {
				return err
			}
		}
		if entity.Label != "" {
			if _, err := fmt.Fprintf(w, " label %q", entity.Label); err != nil {
				return err
			}
		}
		if err := writeMetadata(w, entity.Metadata); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(w); err != nil {
			return err
		}
	}
	if len(doc.Entities) > 0 && len(doc.Relations) > 0 {
		if _, err := fmt.Fprintln(w); err != nil {
			return err
		}
	}
	for _, rel := range doc.Relations {
		if _, err := fmt.Fprintf(w, "rel %s -> %s", rel.From, rel.To); err != nil {
			return err
		}
		if rel.Type != "" && rel.Type != "relates" {
			if _, err := fmt.Fprintf(w, " as %s", rel.Type); err != nil {
				return err
			}
		}
		if err := writeMetadata(w, rel.Metadata); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(w); err != nil {
			return err
		}
	}
	for _, flow := range doc.Flows {
		if _, err := fmt.Fprintln(w); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(w, "flow %s", flow.ID); err != nil {
			return err
		}
		if flow.Label != "" {
			if _, err := fmt.Fprintf(w, " label %q", flow.Label); err != nil {
				return err
			}
		}
		if err := writeMetadata(w, flow.Metadata); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(w, "\n{"); err != nil {
			return err
		}
		for _, step := range flow.Steps {
			if _, err := fmt.Fprintf(w, "  step %s", step.ID); err != nil {
				return err
			}
			if step.Target != "" {
				if _, err := fmt.Fprintf(w, " uses %s", step.Target); err != nil {
					return err
				}
			}
			if step.Label != "" {
				if _, err := fmt.Fprintf(w, " label %q", step.Label); err != nil {
					return err
				}
			}
			if err := writeMetadata(w, step.Metadata); err != nil {
				return err
			}
			if _, err := fmt.Fprintln(w); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintln(w, "}"); err != nil {
			return err
		}
	}
	for _, state := range doc.States {
		if _, err := fmt.Fprintln(w); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(w, "state %s", state.ID); err != nil {
			return err
		}
		if state.Label != "" {
			if _, err := fmt.Fprintf(w, " label %q", state.Label); err != nil {
				return err
			}
		}
		if err := writeMetadata(w, state.Metadata); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(w, "\n{"); err != nil {
			return err
		}
		for _, transition := range state.Transitions {
			if _, err := fmt.Fprintf(w, "  on %s -> %s", quoteIfNeeded(transition.On), transition.To); err != nil {
				return err
			}
			if transition.Condition != "" {
				if _, err := fmt.Fprintf(w, " if %q", transition.Condition); err != nil {
					return err
				}
			}
			if err := writeMetadata(w, transition.Metadata); err != nil {
				return err
			}
			if _, err := fmt.Fprintln(w); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintln(w, "}"); err != nil {
			return err
		}
	}
	return nil
}

func writeMetadata(w io.Writer, meta ast.Metadata) error {
	if len(meta) == 0 {
		return nil
	}
	keys := make([]string, 0, len(meta))
	for key := range meta {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	pairs := make([]string, 0, len(keys))
	for _, key := range keys {
		pairs = append(pairs, fmt.Sprintf("%s=%s", key, quoteIfNeeded(meta[key])))
	}
	_, err := fmt.Fprintf(w, " [%s]", strings.Join(pairs, ", "))
	return err
}

func quoteIfNeeded(value string) string {
	if value == "" {
		return `""`
	}
	for _, ch := range value {
		if !(ch == '-' || ch == '_' || ch == '.' || ch >= '0' && ch <= '9' || ch >= 'A' && ch <= 'Z' || ch >= 'a' && ch <= 'z') {
			return fmt.Sprintf("%q", value)
		}
	}
	return value
}
