# KCommandParser

A package for parsing command and program strings: converts text into `Command` or `Program`. Provides a strategy detection mechanism via `Detect`, as well as a common `Response` format for executors.

## Parsing strategy selection

The `Detect` function determines the input type and returns a `DetectStrategy`: `InputCommand`, `InputProgram`, or `InputUnknown`.

Checks are performed sequentially:

1. Empty string → `InputUnknown`.
2. Prefix `%>`, `%&`, or `%!` → `InputProgram`.
3. Presence of `~>` → `InputUnknown` (a program without a tag is not supported).
4. Prefix `/` → `InputCommand`.
5. Otherwise → `InputUnknown`.

Scenarios:

| Input | Result | Comment |
|---|---|---|
| `/ping` | `InputCommand` | Valid command |
| `%> /ping ~> /status` | `InputProgram` | Program with a tag |
| `%!5%& /ping` | `InputProgram` | Program with a tag and counter |
| `/d ~> /c` | `InputUnknown` | Has `~>`, but no tag — neither command nor program |
| `hello` | `InputUnknown` | Not a command and not a program |
| `` (empty) | `InputUnknown` | Empty input |

## Structures

### Command

A parsing unit obtained from a string like `/name arg1 arg2`. Created by `ParseCommand`, which assigns a unique `ID`.

| Name | Type | Description |
|---|---|---|
| `ID` | `string` | Unique command identifier (UUID) |
| `Name` | `string` | Command name, e.g. `/ping` |
| `Args` | `[]*Arg` | Command arguments |

> A command can be converted to a program using `Convert`. It creates a `Program` with mode `seq`, count `1`, and a body consisting of the single passed command. This is useful when an executor expects a single `Program` type, but a single command was received as input.

### Argument

A single command argument. The specific type is determined by the lexer; `Ident` is used for unrecognized values.

| Name | Type | Description |
|---|---|---|
| `Ident` | `string` | Fallback for unrecognized arguments |
| `Str` | `string` | Quoted string |
| `Addr` | `string` | IP address with optional port |
| `List` | `string` | List separated by `\|` |
| `ID` | `string` | Identifier like `$name` |
| `Key` | `string` | Key like `-name` |
| `Number` | `string` | Number (including fractional and negative) |

The argument value is extracted by `Value`; the list is split by `SplitList`; the number is converted by `ToFloat64`.

### Program

A sequence of commands with a tag. Created by `ParseProgram`; the body is split by `~>`.

| Name | Type | Description |
|---|---|---|
| `Mode` | `ProgramMode` | Execution mode: `seq` or `par` |
| `Count` | `int` | Repeat counter (default 1) |
| `Body` | `[]*Command` | Program commands |

### Response

A common executor response format. Created by `NewResponse`.

| Name | Type | Description |
|---|---|---|
| `CommandID` | `string` | Command identifier |
| `Level` | `string` | Response level |
| `Time` | `time.Time` | Response creation time |
| `ResponseMessage` | `string` | Response text |

## Example

```go
strategy := Detect("/ping") // InputCommand

cmd, _ := ParseCommand("/ping arg 'text' 10.0.0.1 $id -key 42")
fmt.Println(cmd.Name)        // /ping
fmt.Println(cmd.Args[0].Value()) // arg

prog, _ := ParseProgram("%!3%& /ping ~> /status")
fmt.Println(prog.Mode)  // par
fmt.Println(prog.Count) // 3
fmt.Println(len(prog.Body)) // 2

resp := NewResponse(cmd.ID, "info", "ok")
fmt.Println(resp.Level) // info
```

---

This documentation is also available in [Russian](./src/docs/README.ru.md).