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
	{Name: "Ident", Pattern: `[^\s]+`}, //! Ident is a fallback for unrecognized arguments.
	{Name: "Whitespace", Pattern: `\s+`},
})

//endregion

//region Command, Args

type command struct {
	Name string `parser:"@Command"`
	Args []*Arg `parser:"@@*"`
}

type Command struct {
	ID   string
	Name string
	Args []*Arg
}

type Arg struct {
	Ident  string `parser:"@Ident"` //! Ident is a fallback for unrecognized arguments.
	Str    string `parser:"| @String"`
	Addr   string `parser:"| @Address"`
	List   string `parser:"| @List"`
	ID     string `parser:"| @ID"`
	Key    string `parser:"| @Key"`
	Number string `parser:"| @Number"`
}

//endregion

//region Argument value

func (a *Arg) Value() string {
	if a.List != "" {
		return a.List
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
	return a.Ident
}

func (a *Arg) SplitList() []string {
	if a.List == "" {
		return nil
	}
	return strings.Split(a.List, "|")
}

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

//endregion

//region Parsing

var parserCommand = participle.MustBuild[command](
	participle.Lexer(commandLexer),
	participle.Elide("Whitespace"),
)

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

//endregion
