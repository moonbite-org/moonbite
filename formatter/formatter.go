package formatter

import (
	"strings"

	"github.com/moonbite-org/moonbite/parser"
)

const space = " "
const new_line = "\n"

type FormatConfig struct {
	Tabsize int
}
type Formatter struct {
	Ast    parser.Ast
	Config FormatConfig
	Result string
}

func (f Formatter) tab() string {
	result := ""

	for i := 0; i < f.Config.Tabsize; i++ {
		result += space
	}

	return result
}

func (f *Formatter) Format() {
	f.Result += FormatPackageStatement(f.Ast.Package)
	f.Result += new_line
	f.Result += new_line

	for _, use := range f.Ast.Uses {
		f.Result += FormatUseStatement(use)
		f.Result += new_line
	}
	f.Result += new_line

	for _, definition := range f.Ast.Definitions {
		f.Result += FormatStatement(definition.(parser.Statement))
		f.Result += new_line
	}

	f.Result = strings.TrimSpace(f.Result)
}

func NewFormatter(ast parser.Ast) Formatter {
	return Formatter{
		Ast: ast,
		Config: FormatConfig{
			Tabsize: 2,
		},
	}
}

func Format(ast parser.Ast) string {
	formatter := NewFormatter(ast)
	formatter.Format()
	return formatter.Result
}
