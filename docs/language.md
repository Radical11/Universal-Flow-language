# UFL Language Manual

UFL is a small, line-oriented semantic language. It is designed to be readable like structured notes while compiling into a precise IR.

## File Structure

UFL files usually start with Markdown-style headings:

```ufl
# Approval Workflow
## Actors
```

Headings organize the document for humans and tools. A level-one heading becomes the document title.

## Entities

Entities define semantic things.

```ufl
entity request as artifact label "Request"
entity reviewer as actor label "Reviewer" [team=ops]
```

Shape:

```ufl
entity <id> as <type> label "<label>" [key=value]
```

Only the ID is required. If no type is provided, the type is `entity`.

## Relationships

Relationships connect entities.

```ufl
rel request -> reviewer as assigned_to
rel reviewer -> decision as produces [weight=high]
```

Shape:

```ufl
rel <from> -> <to> as <type> [key=value]
```

If no type is provided, the type is `relates`.

## Flows

Flows describe ordered semantic movement.

```ufl
flow approval label "Approval path"
{
  step submit uses request label "Submit request"
  step review uses reviewer label "Review request"
}
```

Each step can optionally point at an entity through `uses`.

## States

States describe state machines and transitions.

```ufl
state idle label "Idle"
{
  on start -> active
}

state active label "Active"
{
  on finish -> complete if "all checks pass"
}
```

Shape:

```ufl
state <id> label "<label>"
{
  on <event> -> <state-id> if "<condition>"
}
```

The condition is optional.

## Metadata

Metadata attaches simple key-value meaning to declarations.

```ufl
entity mountain as place [difficulty=high, region=north]
```

Metadata values can be identifiers, numbers, or strings.

## Comments

Line comments begin with `//`.

```ufl
// This is ignored by the compiler.
entity user as actor
```

## Design Boundary

The base language should stay generic. Do not add keywords like `camera`, `quest`, `screen`, or `database_table` to the core syntax. Those concepts belong in libraries, validators, backends, or conventions built on top of the IR.

