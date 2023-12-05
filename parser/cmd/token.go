package parser

import (
	"fmt"

	common "github.com/moonbite-org/moonbite/parser/common"
)

type token_kind int

const (
	eof_token_kind token_kind = iota
	Whitespace
	new_line
	whitespace
	any_whitespace

	Punctuation
	left_parens
	right_parens
	left_angle_bracks
	right_angle_bracks
	left_squre_bracks
	right_squre_bracks
	left_curly_bracks
	right_curly_bracks
	dot
	comma
	colon
	semicolon
	variadic_marker // ...

	CommentToken
	multi_line_comment
	single_line_comment

	Operator
	plus
	minus
	star
	forward_slash
	percent
	increment
	decrement
	power
	pipe
	ampersand
	assignment
	arithmetic_assignment
	exclamation
	binary_operator // == != < > <= >= && ||
	caret
	channel
	then

	identifier

	literal
	string_literal
	rune_literal
	number_literal
	bool_literal
	cardinal_literal

	Keyword
	as_keyword
	base_keyword
	break_keyword
	const_keyword
	continue_keyword
	corout_keyword
	else_keyword
	for_keyword
	fun_keyword
	giveup_keyword
	hidden_keyword
	if_keyword
	implements_keyword
	instanceof_keyword
	match_keyword
	map_keyword
	mimics_keyword
	move_keyword
	of_keyword
	or_keyword
	package_keyword
	return_keyword
	this_keyword
	trait_keyword
	type_keyword
	use_keyword
	var_keyword
	yield_keyword
)

var token_map = map[token_kind]string{
	eof_token_kind: "End of file",
	new_line:       "NewLine",
	whitespace:     "Whitespace",
	any_whitespace: "AnyWhitespace",

	left_parens:        "LeftParens",
	right_parens:       "RightParens",
	left_angle_bracks:  "LeftAngleBracks",
	right_angle_bracks: "RightAngleBracks",
	left_squre_bracks:  "LeftSqureBracks",
	right_squre_bracks: "RightSqureBracks",
	left_curly_bracks:  "LeftCurlyBracks",
	right_curly_bracks: "RightCurlyBracks",
	dot:                "Dot",
	comma:              "Comma",
	colon:              "Colon",
	semicolon:          "Semicolon",
	variadic_marker:    "VariadicMarker", // ...

	multi_line_comment:  "MultiLineComment",
	single_line_comment: "SingleLineComment",

	Operator:              "Operator",
	plus:                  "Plus",
	minus:                 "Minus",
	star:                  "Star",
	forward_slash:         "ForwardSlash",
	percent:               "Percent",
	increment:             "Increment",
	decrement:             "Decrement",
	power:                 "Power",
	pipe:                  "Pipe",
	ampersand:             "Ampersand",
	assignment:            "Assignment",
	arithmetic_assignment: "ArithmeticAssignment",
	binary_operator:       "BinaryPredicate", // == != < > <= >=
	caret:                 "Caret",
	channel:               "Channel",
	then:                  "Then",

	identifier: "Identifier",

	literal:          "Literal",
	string_literal:   "StringLiteral",
	rune_literal:     "RuneLiteral",
	number_literal:   "NumberLiteral",
	bool_literal:     "BoolLiteral",
	cardinal_literal: "CardinalLiteral",

	Keyword:            "Keyword",
	as_keyword:         "as",
	base_keyword:       "base",
	break_keyword:      "break",
	const_keyword:      "const",
	continue_keyword:   "continue",
	corout_keyword:     "corout",
	else_keyword:       "else",
	for_keyword:        "for",
	fun_keyword:        "fun",
	giveup_keyword:     "giveup",
	hidden_keyword:     "hidden",
	if_keyword:         "if",
	implements_keyword: "implements",
	instanceof_keyword: "instanceof",
	match_keyword:      "match",
	map_keyword:        "map",
	mimics_keyword:     "mimics",
	move_keyword:       "move",
	of_keyword:         "of",
	or_keyword:         "or",
	package_keyword:    "package",
	return_keyword:     "return",
	this_keyword:       "this",
	trait_keyword:      "trait",
	type_keyword:       "type",
	use_keyword:        "use",
	var_keyword:        "var",
	yield_keyword:      "yield",
}

type Token struct {
	Kind       token_kind      `json:"kind"`
	Location   common.Location `json:"location"`
	Literal    string          `json:"literal"`
	Raw        string          `json:"raw"`
	Offset     int             `json:"offset"`
	LineBreaks int             `json:"line_breaks"`
}

func (t Token) String() string {
	if t.Kind < Keyword {
		if t.Kind >= Whitespace && t.Kind < Punctuation {
			return fmt.Sprintf("ws[%d:%d][%d]", t.Location.Line, t.Location.Column, t.Offset)
		} else {
			return fmt.Sprintf("%s(%s)[%d:%d][%d]", token_map[t.Kind], t.Literal, t.Location.Line, t.Location.Column, t.Offset)
		}
	}

	return fmt.Sprintf("<%s>[%d:%d][%d]", t.Literal, t.Location.Line, t.Location.Column, t.Offset)
}

const eof rune = -1

var eof_token = Token{
	Kind: eof_token_kind,
}
