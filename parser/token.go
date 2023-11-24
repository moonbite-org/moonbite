package parser

import "fmt"

type token_kind int

const (
	eof_token_kind token_kind = iota
	new_line
	whitespace
	any_whitespace

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

	multi_line_comment
	single_line_comment

	operator
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
	binary_predicate // == != < > <= >=
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

	keyword
	as_keyword
	base_keyword
	break_keyword
	const_keyword
	continue_keyword
	corout_keyword
	else_keyword
	for_keyword
	fun_keyword
	if_keyword
	implements_keyword
	instanceof_keyword
	match_keyword
	mimics_keyword
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

	operator:              "Operator",
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
	binary_predicate:      "BinaryPredicate", // == != < > <= >=
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

	keyword:            "Keyword",
	as_keyword:         "AsKeyword",
	base_keyword:       "BaseKeyword",
	break_keyword:      "BreakKeyword",
	const_keyword:      "ConstKeyword",
	continue_keyword:   "ContinueKeyword",
	corout_keyword:     "CoroutKeyword",
	else_keyword:       "ElseKeyword",
	for_keyword:        "ForKeyword",
	fun_keyword:        "FunKeyword",
	if_keyword:         "IfKeyword",
	implements_keyword: "ImplementsKeyword",
	instanceof_keyword: "InstanceofKeyword",
	match_keyword:      "MatchKeyword",
	mimics_keyword:     "MimicsKeyword",
	of_keyword:         "OfKeyword",
	or_keyword:         "OrKeyword",
	package_keyword:    "PackageKeyword",
	return_keyword:     "ReturnKeyword",
	this_keyword:       "ThisKeyword",
	trait_keyword:      "TraitKeyword",
	type_keyword:       "TypeKeyword",
	use_keyword:        "UseKeyword",
	var_keyword:        "VarKeyword",
	yield_keyword:      "YieldKeyword",
}

type Token struct {
	Kind       token_kind `json:"kind"`
	Location   Location   `json:"location"`
	Literal    string     `json:"literal"`
	Raw        string     `json:"raw"`
	Offset     int        `json:"offset"`
	LineBreaks int        `json:"line_breaks"`
}

func (t Token) String() string {
	if t.Kind < keyword {
		if t.Kind == whitespace {
			return fmt.Sprintf("ws[%d:%d][%d]", t.Location.Line, t.Location.Column, t.Offset)
		} else if t.Kind == new_line {
			return fmt.Sprintf("nl[%d:%d][%d]", t.Location.Line, t.Location.Column, t.Offset)
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
