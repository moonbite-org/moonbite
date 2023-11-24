package parser

import "fmt"

type Location struct {
	File   string `json:"file"`
	Offset int    `json:"offset"`
	Line   int    `json:"line"`
	Column int    `json:"column"`
}

type ErrorKind int

const (
	SyntaxError ErrorKind = iota
	TypeMismatchError
)

var error_messages = map[string]string{
	"u_eof":   "Unexpected end of file, file ended but I was expecting %s",
	"u_tok":   "Unexpected token, I was not expecting '%s'",
	"u_tok_s": "Unexpected token, I was expecting '%s' but got '%s'",
	"i_con":   "Illegal construct, I cannot %s",
	"i_val":   "Invalid value, I cannot make sense of this value. %s",
}

var error_kind_map = map[ErrorKind]string{
	SyntaxError:       "Syntax Error",
	TypeMismatchError: "Type Mismatch Error",
}

type Error struct {
	Kind     ErrorKind
	Reason   string
	Location Location
	Exists   bool
}

func (e Error) String() string {
	return fmt.Sprintf("%s: %s at %d:%d in %s", error_kind_map[e.Kind], e.Reason, e.Location.Line, e.Location.Column, e.Location.File)
}
