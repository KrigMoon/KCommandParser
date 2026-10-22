package KCommandParser

import (
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/google/uuid"
)

var commandLexer = lexer.MustSimple([]lexer.SimpleRule{
	{Name: "Command", Pattern: `/[a-zA-Z_][a-zA-Z0-9_]*`},
	{Name: "Address", Pattern: `\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}:\d{1,5}`},
	{Name: "String", Pattern: `'[^']*'|"[^"]*"`},
	{Name: "List", Pattern: `[a-zA-Z0-9_-]+(\|[a-zA-Z0-9_-]+)+`},
	{Name: "Ident", Pattern: `[^\s]+`},
	{Name: "Whitespace", Pattern: `\s+`},
})

type сommand struct {
	Name string `parser:"@Command"`
	Args []*arg `parser:"@@*"`
}

type Command struct {
	ID   string
	Name string
	Args []*arg
}

type arg struct {
	Ident string `parser:"@Ident"`
	Str   string `parser:"| @String"`
	Addr  string `parser:"| @Address"`
	List  string `parser:"| @List"`
}

func (a *arg) Value() string {
	if a.Ident != "" {
		return a.Ident
	}
	if a.Str != "" {
		return strings.Trim(a.Str, `"'`)
	}
	if a.Addr != "" {
		return a.Addr
	}
	return a.List
}

func (a *arg) SplitList() []string {
	if a.List == "" {
		return nil
	}
	return strings.Split(a.List, "|")
}

var parser = participle.MustBuild[сommand](
	participle.Lexer(commandLexer),
	participle.Elide("Whitespace"),
)

func parseCommand(input string) (*сommand, error) {
	return parser.ParseString("", input)
}

func ParseCommand(input string) (*Command, error) {
	cmd, err := parseCommand(input)
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
