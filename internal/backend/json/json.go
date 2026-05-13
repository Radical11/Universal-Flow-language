package json

import (
	stdjson "encoding/json"
	"io"

	"github.com/Radical11/Universal-Flow-language/internal/ir"
)

func Encode(w io.Writer, doc *ir.Document) error {
	encoder := stdjson.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(doc)
}
