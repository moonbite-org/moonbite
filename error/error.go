package error

import "fmt"

type ErrorKind int

const (
	SyntaxError ErrorKind = iota
	TypeError
	CompileError
)

var ErrorMessages = map[string]string{
	"u_eof":    "Unexpected end of file, file ended but I was expecting %s",
	"u_tok":    "Unexpected token, I was not expecting '%s'",
	"u_tok_s":  "Unexpected token, I was expecting '%s' but got '%s'",
	"u_tok_m":  "Unexpected token, I was not expecting '%s'. %s",
	"uc_con":   "Uncallable construct, I cannot treat this expression as a funcation callee",
	"i_con":    "Illegal construct, I cannot %s",
	"i_val":    "Invalid value, I cannot make sense of this value. %s",
	"w_e_args": "Too many arguments, warn expressions should only have exactly 1 argument",
}

var ErrorKindMap = map[ErrorKind]string{
	SyntaxError:  "Syntax Error",
	TypeError:    "Type Error",
	CompileError: "Compile Error",
}

type Position struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

type Location struct {
	File   string   `json:"file"`
	Offset int      `json:"offset"`
	Start  Position `json:"start"`
	End    Position `json:"end"`
}

type Error struct {
	Kind      ErrorKind `json:"kind"`
	Reason    string    `json:"reason"`
	Location  Location  `json:"location"`
	Exists    bool      `json:"exists"`
	Anonymous bool      `json:"anonymous"`
}

func (e Error) String() string {
	if !e.Exists {
		return ""
	}

	if e.Anonymous {
		return fmt.Sprintf("%s: %s", ErrorKindMap[e.Kind], e.Reason)
	}

	return fmt.Sprintf("%s: %s at %d:%d in %s", ErrorKindMap[e.Kind], e.Reason, e.Location.Start.Line, e.Location.Start.Column, e.Location.File)
}

func CreateAnonError(kind ErrorKind, reason string) Error {
	return Error{
		Kind:   kind,
		Reason: reason,
		Location: Location{
			Start:  Position{},
			End:    Position{},
			Offset: 0,
			File:   "",
		},
		Exists:    true,
		Anonymous: true,
	}
}

func CreateTypeError(reason string, location Location) Error {
	return Error{
		Location: location,
		Kind:     TypeError,
		Exists:   true,
		Reason:   reason,
	}
}

func CreateCompileError(reason string, location Location) Error {
	return Error{
		Location: location,
		Kind:     CompileError,
		Exists:   true,
		Reason:   reason,
	}
}

var EmptyError = Error{
	Kind:   0,
	Reason: "",
	Location: Location{
		Start:  Position{},
		End:    Position{},
		Offset: 0,
		File:   "",
	},
	Exists:    false,
	Anonymous: true,
}
