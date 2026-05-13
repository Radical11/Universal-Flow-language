package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	analysisbackend "github.com/Radical11/Universal-Flow-language/internal/backend/analysis"
	dotbackend "github.com/Radical11/Universal-Flow-language/internal/backend/dot"
	jsonbackend "github.com/Radical11/Universal-Flow-language/internal/backend/json"
	mmdbackend "github.com/Radical11/Universal-Flow-language/internal/backend/mmd"
	runtimebackend "github.com/Radical11/Universal-Flow-language/internal/backend/runtime"
	timelinebackend "github.com/Radical11/Universal-Flow-language/internal/backend/timeline"
	"github.com/Radical11/Universal-Flow-language/internal/compiler"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "ufl:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		return usage()
	}
	switch args[0] {
	case "parse":
		return parseCmd(args[1:])
	case "validate":
		return validateCmd(args[1:])
	case "compile":
		return compileCmd(args[1:])
	case "fmt":
		return fmtCmd(args[1:])
	case "inspect":
		return inspectCmd(args[1:])
	case "help", "-h", "--help":
		return usage()
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func parseCmd(args []string) error {
	fs := flag.NewFlagSet("parse", flag.ContinueOnError)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		return fmt.Errorf("usage: ufl parse input.ufl")
	}
	source, err := os.ReadFile(fs.Arg(0))
	if err != nil {
		return err
	}
	doc, err := compiler.Parse(string(source))
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(doc)
}

func validateCmd(args []string) error {
	profile, normalizedArgs, err := extractProfileOption(args)
	if err != nil {
		return err
	}
	fs := flag.NewFlagSet("validate", flag.ContinueOnError)
	if err := fs.Parse(normalizedArgs); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		return fmt.Errorf("usage: ufl validate input.ufl [--profile workflow]")
	}
	source, err := os.ReadFile(fs.Arg(0))
	if err != nil {
		return err
	}
	doc, err := compiler.CompileWithProfile(string(source), profile)
	if err != nil {
		return err
	}
	if len(doc.Diagnostics) > 0 {
		for _, diagnostic := range doc.Diagnostics {
			if diagnostic.Line > 0 {
				fmt.Fprintf(os.Stderr, "%s:%d:%d: %s\n", diagnostic.Severity, diagnostic.Line, diagnostic.Column, diagnostic.Message)
				continue
			}
			fmt.Fprintf(os.Stderr, "%s: %s\n", diagnostic.Severity, diagnostic.Message)
		}
	}
	fmt.Println("valid")
	return nil
}

func compileCmd(args []string) error {
	target, outPath, profile, normalizedArgs, err := extractCompileOptions(args)
	if err != nil {
		return err
	}
	fs := flag.NewFlagSet("compile", flag.ContinueOnError)
	if err := fs.Parse(normalizedArgs); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		return fmt.Errorf("usage: ufl compile input.ufl --target json --profile workflow --out output.json")
	}
	source, err := os.ReadFile(fs.Arg(0))
	if err != nil {
		return err
	}
	doc, err := compiler.CompileWithProfile(string(source), profile)
	if err != nil {
		return err
	}
	var output bytes.Buffer
	switch target {
	case "json":
		err = jsonbackend.Encode(&output, doc)
	case "mermaid":
		err = mmdbackend.Encode(&output, doc)
	case "dot":
		err = dotbackend.Encode(&output, doc)
	case "timeline":
		err = timelinebackend.Encode(&output, doc)
	case "runtime":
		err = runtimebackend.Encode(&output, doc)
	case "analysis":
		err = analysisbackend.Encode(&output, doc)
	default:
		return fmt.Errorf("unsupported target %q", target)
	}
	if err != nil {
		return err
	}
	return writeOutput(outPath, output.Bytes())
}

func fmtCmd(args []string) error {
	write, normalizedArgs, err := extractFmtOptions(args)
	if err != nil {
		return err
	}
	fs := flag.NewFlagSet("fmt", flag.ContinueOnError)
	if err := fs.Parse(normalizedArgs); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		return fmt.Errorf("usage: ufl fmt input.ufl [--write]")
	}
	sourcePath := fs.Arg(0)
	source, err := os.ReadFile(sourcePath)
	if err != nil {
		return err
	}
	doc, err := compiler.Parse(string(source))
	if err != nil {
		return err
	}
	var output bytes.Buffer
	if err := compiler.Format(&output, doc); err != nil {
		return err
	}
	if write {
		return os.WriteFile(sourcePath, output.Bytes(), 0644)
	}
	_, err = io.Copy(os.Stdout, &output)
	return err
}

func inspectCmd(args []string) error {
	profile, normalizedArgs, err := extractProfileOption(args)
	if err != nil {
		return err
	}
	fs := flag.NewFlagSet("inspect", flag.ContinueOnError)
	if err := fs.Parse(normalizedArgs); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		return fmt.Errorf("usage: ufl inspect input.ufl [--profile workflow]")
	}
	source, err := os.ReadFile(fs.Arg(0))
	if err != nil {
		return err
	}
	doc, err := compiler.CompileWithProfile(string(source), profile)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "title: %s\n", doc.Title)
	fmt.Fprintf(os.Stdout, "nodes: %d\n", len(doc.Nodes))
	fmt.Fprintf(os.Stdout, "edges: %d\n", len(doc.Edges))
	fmt.Fprintf(os.Stdout, "flows: %d\n", len(doc.Flows))
	fmt.Fprintf(os.Stdout, "states: %d\n", len(doc.States))
	if profile != "" {
		fmt.Fprintf(os.Stdout, "profile: %s\n", profile)
	}
	fmt.Fprintf(os.Stdout, "diagnostics: %d\n", len(doc.Diagnostics))
	return nil
}

func usage() error {
	fmt.Fprintln(os.Stderr, "usage: ufl <parse|validate|compile|fmt|inspect> [options] input.ufl")
	return nil
}

func extractProfileOption(args []string) (string, []string, error) {
	profile := ""
	normalized := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--profile" {
			if i+1 >= len(args) {
				return "", nil, fmt.Errorf("--profile requires a value")
			}
			profile = args[i+1]
			i++
			continue
		}
		if len(arg) > len("--profile=") && arg[:len("--profile=")] == "--profile=" {
			profile = arg[len("--profile="):]
			continue
		}
		normalized = append(normalized, arg)
	}
	profile, err := compiler.NormalizeProfile(profile)
	if err != nil {
		return "", nil, err
	}
	return profile, normalized, nil
}

func extractCompileOptions(args []string) (string, string, string, []string, error) {
	target := "json"
	outPath := ""
	profile := ""
	normalized := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--target" {
			if i+1 >= len(args) {
				return "", "", "", nil, fmt.Errorf("--target requires a value")
			}
			target = args[i+1]
			i++
			continue
		}
		if len(arg) > len("--target=") && arg[:len("--target=")] == "--target=" {
			target = arg[len("--target="):]
			continue
		}
		if arg == "--out" || arg == "-o" {
			if i+1 >= len(args) {
				return "", "", "", nil, fmt.Errorf("%s requires a value", arg)
			}
			outPath = args[i+1]
			i++
			continue
		}
		if len(arg) > len("--out=") && arg[:len("--out=")] == "--out=" {
			outPath = arg[len("--out="):]
			continue
		}
		if arg == "--profile" {
			if i+1 >= len(args) {
				return "", "", "", nil, fmt.Errorf("--profile requires a value")
			}
			profile = args[i+1]
			i++
			continue
		}
		if len(arg) > len("--profile=") && arg[:len("--profile=")] == "--profile=" {
			profile = arg[len("--profile="):]
			continue
		}
		normalized = append(normalized, arg)
	}
	profile, err := compiler.NormalizeProfile(profile)
	if err != nil {
		return "", "", "", nil, err
	}
	return target, outPath, profile, normalized, nil
}

func extractFmtOptions(args []string) (bool, []string, error) {
	write := false
	normalized := make([]string, 0, len(args))
	for _, arg := range args {
		if arg == "--write" || arg == "-w" {
			write = true
			continue
		}
		normalized = append(normalized, arg)
	}
	return write, normalized, nil
}

func writeOutput(path string, content []byte) error {
	if path == "" {
		_, err := os.Stdout.Write(content)
		return err
	}
	return os.WriteFile(path, content, 0644)
}
