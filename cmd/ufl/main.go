package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	jsonbackend "github.com/Radical11/Universal-Flow-language/internal/backend/json"
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
	fs := flag.NewFlagSet("validate", flag.ContinueOnError)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		return fmt.Errorf("usage: ufl validate input.ufl")
	}
	source, err := os.ReadFile(fs.Arg(0))
	if err != nil {
		return err
	}
	doc, err := compiler.Compile(string(source))
	if err != nil {
		return err
	}
	if len(doc.Diagnostics) > 0 {
		for _, diagnostic := range doc.Diagnostics {
			fmt.Fprintf(os.Stderr, "%s: %s\n", diagnostic.Severity, diagnostic.Message)
		}
	}
	fmt.Println("valid")
	return nil
}

func compileCmd(args []string) error {
	fs := flag.NewFlagSet("compile", flag.ContinueOnError)
	target := fs.String("target", "json", "compile target")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		return fmt.Errorf("usage: ufl compile input.ufl --target json")
	}
	if *target != "json" {
		return fmt.Errorf("unsupported target %q", *target)
	}
	source, err := os.ReadFile(fs.Arg(0))
	if err != nil {
		return err
	}
	doc, err := compiler.Compile(string(source))
	if err != nil {
		return err
	}
	return jsonbackend.Encode(os.Stdout, doc)
}

func usage() error {
	fmt.Fprintln(os.Stderr, "usage: ufl <parse|validate|compile> [options] input.ufl")
	return nil
}
