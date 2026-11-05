package KCommandParser

import (
	"strconv"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

//region Pattern Rules

var programLexer = lexer.MustSimple([]lexer.SimpleRule{
	{Name: "LoopSeq", Pattern: `%!\d{1,3}%>`},
	{Name: "LoopPar", Pattern: `%!\d{1,3}%&`},
	{Name: "Seq", Pattern: `%>`},
	{Name: "Par", Pattern: `%&`},
	{Name: "Body", Pattern: `.+`},
})

//endregion

//region Program

type program struct {
	Tag  string `parser:"@(LoopSeq | LoopPar | Seq | Par)"`
	Body string `parser:"@Body"`
}

type ProgramMode string

const (
	ProgramModeSeq ProgramMode = "seq"
	ProgramModePar ProgramMode = "par"
)

type Program struct {
	Mode  ProgramMode
	Count int
	Body  []*Command
}

//endregion

//region Detect Strategy

type DetectStrategy int

const (
	InputUnknown DetectStrategy = iota
	InputCommand
	InputProgram
)

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

//endregion

//region Parsing

var parserProgram = participle.MustBuild[program](
	participle.Lexer(programLexer),
)

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

//endregion

//region Convert Command to Program

func Convert(command *Command) *Program {
	return &Program{
		Mode:  ProgramModeSeq,
		Count: 1,
		Body:  []*Command{command},
	}
}

//endregion
