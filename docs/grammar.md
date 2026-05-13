# UFL Grammar

This page describes the MVP grammar in a compact EBNF-like form.

```text
document      = { newline | heading | entity | relation | flow | state } eof ;

heading       = "#" { "#" } text newline ;

entity        = "entity" identifier
                [ "as" identifier ]
                [ "label" value ]
                [ metadata ]
                newline ;

relation      = "rel" identifier "->" identifier
                [ "as" identifier ]
                [ metadata ]
                newline ;

flow          = "flow" identifier
                [ "label" value ]
                [ metadata ]
                newline
                "{"
                { newline | step }
                "}" newline ;

step          = "step" identifier
                [ "uses" identifier ]
                [ "label" value ]
                [ metadata ]
                newline ;

state         = "state" identifier
                [ "label" value ]
                [ metadata ]
                newline
                "{"
                { newline | transition }
                "}" newline ;

transition    = "on" value "->" identifier
                [ "if" value ]
                [ metadata ]
                newline ;

metadata      = "[" metadata_pair { "," metadata_pair } "]" ;
metadata_pair = identifier "=" value ;
value         = identifier | number | string ;
```

## Identifiers

Identifiers may contain letters, numbers, underscores, dashes, and dots. They must start with a letter or underscore.

Examples:

```text
hero
approval.step
state-id
```

## Strings

Strings use double quotes.

```ufl
label "Review request"
```

Supported escapes:

```text
\n
\t
\"
\\
```

## Comments

Line comments begin with `//`.

```ufl
// ignored
entity user as actor
```

## Stability

This grammar describes `ufl.ir.v0`. Future syntax should preserve this core unless a breaking IR version is introduced.

