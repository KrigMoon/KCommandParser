package KCommandParser

import (
	"strconv"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/google/uuid"
)

//region Pattern Rules

var commandLexer = lexer.MustSimple([]lexer.SimpleRule{
	{Name: "Command", Pattern: `/[a-zA-Z_][a-zA-Z0-9_]*`},
	{Name: "Address", Pattern: `\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}(:\d{1,5})?`},
	{Name: "String", Pattern: `'[^']*'|"[^"]*"`},
	{Name: "List", Pattern: `[a-zA-Z0-9_-]+(\|[a-zA-Z0-9_-]+)+`},
	{Name: "ID", Pattern: `\$[a-zA-Z_][a-zA-Z0-9_]*`},
	{Name: "Key", Pattern: `-[a-zA-Z_][a-zA-Z0-9_-]*`},
	{Name: "Number", Pattern: `-?\d+(\.\d+)?`},
	{Name: "Ident", Pattern: `[^\s]+`},
	{Name: "Whitespace", Pattern: `\s+`},
})

var programLexer = lexer.MustSimple([]lexer.SimpleRule{
	{Name: "LoopSeq", Pattern: `%!\d{1,3}%>`},
	{Name: "LoopPar", Pattern: `%!\d{1,3}%&`},
	{Name: "Seq", Pattern: `%>`},
	{Name: "Par", Pattern: `%&`},
	{Name: "Body", Pattern: `.+`},
})

//endregion

//region Command, Args and Program

// desc: Local command structure
type command struct {
	Name string `parser:"@Command"`
	Args []*Arg `parser:"@@*"`
}

// desc: Public command structure
type Command struct {
	ID   string
	Name string
	Args []*Arg
}

// desc: Argument structure
type Arg struct {
	Ident  string `parser:"@Ident"`
	Str    string `parser:"| @String"`
	Addr   string `parser:"| @Address"`
	List   string `parser:"| @List"`
	ID     string `parser:"| @ID"`
	Key    string `parser:"| @Key"`
	Number string `parser:"| @Number"`
}

// desc: Local program structure
type program struct {
	Tag  string `parser:"@(LoopSeq | LoopPar | Seq | Par)"`
	Body string `parser:"@Body"`
}

// desc: Program execution mode
type ProgramMode string

const (
	ProgramModeSeq ProgramMode = "seq"
	ProgramModePar ProgramMode = "par"
)

// desc: Public program structure
type Program struct {
	Mode  ProgramMode
	Count int
	Body  []*Command
}

//endregion

//region Help Functions

//region^ Argument value

// desc: Getting the value
func (a *Arg) Value() string {
	if a.Ident != "" {
		return a.Ident
	}
	if a.Str != "" {
		return strings.Trim(a.Str, `"'`)
	}
	if a.Addr != "" {
		return a.Addr
	}
	if a.ID != "" {
		return a.ID[1:]
	}
	if a.Key != "" {
		return a.Key
	}
	if a.Number != "" {
		return a.Number
	}
	return a.List
}

// desc: Splitting a string into a list
func (a *Arg) SplitList() []string {
	if a.List == "" {
		return nil
	}
	return strings.Split(a.List, "|")
}

// desc: Convert To Float64
func (a *Arg) ToFloat64() *float64 {
	if a.Number == "" {
		return nil
	}
	val, err := strconv.ParseFloat(a.Number, 64)
	if err != nil {
		return nil
	}
	return &val
}

//endregion^

//region^ Detect Strategy

// desc: Input detect strategy
type DetectStrategy int

const (
	InputUnknown DetectStrategy = iota
	InputCommand
	InputProgram
)

// desc: Detect input kind
func Detect(input string) DetectStrategy {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return InputUnknown
	}
	if strings.HasPrefix(trimmed, "%>") ||
		strings.HasPrefix(trimmed, "%&") ||
		strings.HasPrefix(trimmed, "%!") {
		return InputProgram
	}
	if strings.Contains(trimmed, "~>") {
		return InputUnknown
	}
	if strings.HasPrefix(trimmed, "/") {
		return InputCommand
	}
	return InputUnknown
}

//endregion^

//endregion

//region Parsing

//region^ Parsing Command

// desc: Parser engine
var parserCommand = participle.MustBuild[command](
	participle.Lexer(commandLexer),
	participle.Elide("Whitespace"),
)

// desc: Parsing from string to command
func ParseCommand(input string) (*Command, error) {
	cmd, err := parserCommand.ParseString("", input)
	if err != nil {
		return nil, err
	}
	id := uuid.New().String()
	return &Command{
		ID:   id,
		Name: cmd.Name,
		Args: cmd.Args,
	}, nil
}

//endregion^

//region^ Parsing Program

// desc: Parser engine for program
var parserProgram = participle.MustBuild[program](
	participle.Lexer(programLexer),
)

// desc: Parsing tag into mode and count
func parseTag(tag string) (ProgramMode, int) {
	mode := ProgramModeSeq
	if strings.HasSuffix(tag, "%&") {
		mode = ProgramModePar
	}
	count := 1
	if strings.HasPrefix(tag, "%!") {
		numStr := strings.TrimSuffix(strings.TrimPrefix(tag, "%!"), "%>")
		numStr = strings.TrimSuffix(numStr, "%&")
		if n, err := strconv.Atoi(numStr); err == nil {
			count = n
		}
	}
	return mode, count
}

// desc: Parsing from string to program
func ParseProgram(input string) (*Program, error) {
	prog, err := parserProgram.ParseString("", input)
	if err != nil {
		return nil, err
	}

	mode, count := parseTag(prog.Tag)

	parts := strings.Split(prog.Body, "~>")
	body := make([]*Command, 0, len(parts))
	for _, part := range parts {
		cmd, err := ParseCommand(part)
		if err != nil {
			return nil, err
		}
		body = append(body, cmd)
	}

	return &Program{
		Mode:  mode,
		Count: count,
		Body:  body,
	}, nil
}

//endregion^

//endregion
