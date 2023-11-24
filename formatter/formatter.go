package formatter

import (
	"strings"

	"github.com/CanPacis/tokenizer/parser"
)

func NewFormatter(ast parser.Ast) Formatter {
	return Formatter{
		ast:      ast,
		Result:   "",
		new_line: "\n",
		tab:      "  ",
	}
}

type printable interface {
	String() string
}

type Formatter struct {
	ast      parser.Ast
	Result   string
	new_line string
	tab      string
}

func (f *Formatter) add_new_line() {
	f.Result += f.new_line
}

func (f *Formatter) trim() {
	f.Result = strings.TrimSpace(f.Result)
}

func new_line_list[T printable](f *Formatter, items []T, seperation int) {
	seperator := ""

	for i := 0; i < seperation; i++ {
		seperator += f.new_line
	}

	for _, item := range items {
		f.Result += seperator
		f.Result += item.String()
	}
}

func (f *Formatter) Format() {
	f.Result += f.ast.Package.String()
	f.add_new_line()

	new_line_list(f, f.ast.Uses, 1)
	new_line_list(f, f.ast.Definitions, 2)
	f.add_new_line()
	f.trim()
}

func Format(ast parser.Ast) string {
	formatter := NewFormatter(ast)
	formatter.Format()
	return formatter.Result
}
