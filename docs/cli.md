# UFL CLI Reference

The CLI lives at:

```sh
go run ./cmd/ufl
```

After installing or building a binary, the same commands can be run as `ufl`.

## Parse

Parse source into the AST.

```sh
go run ./cmd/ufl parse examples/story.ufl
```

Use this when debugging syntax-level structure.

## Validate

Validate source and report diagnostics.

```sh
go run ./cmd/ufl validate examples/workflow.ufl
```

Successful output:

```text
valid
```

Warnings are printed to stderr.

## Compile

Compile source to a backend target.

```sh
go run ./cmd/ufl compile examples/network.ufl --target json
```

The only current target is `json`.

## Common Development Commands

Run tests:

```sh
go test ./...
```

Format code:

```sh
gofmt -w cmd internal
```

