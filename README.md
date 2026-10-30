# KCommandParser

## Description
KCommandParser is a Go library for parsing commands and programs from a compact string syntax. It turns a string like `/name arg1 arg2 ...` into a `Command` structure with typed arguments, and a string with a tag (`%>`, `%&`, `%!N%>`, `%!N%&`) and commands separated by `~>` — into a `Program` structure with an execution mode and repeat count. Each parsed command gets a UUID, which a `Response` can refer to.

## Entities

### Command

#### Parsing rules
A command must start with `/`, followed by a letter or underscore, then — any number of arguments separated by spaces. Each argument is parsed into the most specific supported type:

| Type | Example | Note |
|---|---|---|
| `Ident` | `timeout=5` | Anything that did not fit other types |
| `String` | `'hello world'` / `"..."` | In quotes |
| `Address` | `192.168.1.10:8080` | IPv4 with an optional port |
| `List` | `a|b|c` | Two or more items separated by `\|` |
| `ID` | `$user` | Prefix `$` |
| `Key` | `-verbose` | Prefix `-` |
| `Number` | `42`, `-3.14` | Integer or float |

To get the argument value, use `Value()`. For lists — `SplitList()`. For numbers — `ToFloat64()`.

> **Warning**
>
> For a list or number argument, `Value()` returns the original string like `a|b|c` or `2.12`, not the separate items or the number. If you need the items/number — use `SplitList()`/`ToFloat64()`.

#### Structure
```go
type Command struct {
	ID   string
	Name string
	Args []*Arg
}
```
- `ID` — UUID, generated automatically during parsing;
- `Name` — the command name;
- `Args` — a slice of parsed arguments.

Created via `ParseCommand(input string) (*Command, error)`.

Argument:

```go
type Arg struct {
	Ident  string
	Str    string
	Addr   string
	List   string
	ID     string
	Key    string
	Number string
}
```

### Program

#### Parsing rules

A program starts with a tag, followed by a body of commands separated by `~>`. Each part of the body is parsed via `ParseCommand`.

| Tag | Mode | Count | Meaning |
|---|---|---|---|
| `%>` | seq | 1 | Sequential execution |
| `%&` | par | 1 | Parallel execution |
| `%!N%>` | seq | N | Sequential, repeat N times |
| `%!N%&` | par | N | Parallel, repeat N times |

`N` — from 1 to 999.

#### Structure
```go
type Program struct {
	Mode  ProgramMode
	Count int
	Body  []*Command
}
```
- `Mode` — `ProgramModeSeq` (`"seq"`) or `ProgramModePar` (`"par"`);
- `Count` — repeat count, default is `1`;
- `Body` — a slice of parsed commands.

Created via `ParseProgram(input string) (*Program, error)`.

### Response
```go
type Response struct {
	CommandID       string
	Level           string
	Time            time.Time
	ResponseMessage string
}
```
- `CommandID` — the `ID` of the command this response belongs to;
- `Level` — a free-form level string;
- `Time` — filled automatically upon creation;
- `ResponseMessage` — the response text.

Created via `NewResponse(commandID, level, responseMessage string) *Response`.

## Choosing a parsing strategy

The `Detect(input)` function classifies the input string and returns one of three values of type `DetectStrategy`:

| Value | When returned |
|---|---|
| `InputProgram` | string starts with `%>`, `%&` or `%!` |
| `InputCommand` | string starts with `/` |
| `InputUnknown` | empty string, a string containing `~>` but no program tag, anything else |

Before checking, the string is trimmed via `strings.TrimSpace`, so leading and trailing spaces do not affect the result.