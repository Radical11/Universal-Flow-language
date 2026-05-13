package compiler

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExamplesCompile(t *testing.T) {
	files, err := filepath.Glob(filepath.Join("..", "..", "examples", "*.ufl"))
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("no examples found")
	}
	for _, file := range files {
		t.Run(filepath.Base(file), func(t *testing.T) {
			source, err := os.ReadFile(file)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := Compile(string(source)); err != nil {
				t.Fatal(err)
			}
		})
	}
}
