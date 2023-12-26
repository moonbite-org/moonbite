package parser

import (
	"fmt"

	errors "github.com/moonbite-org/moonbite/error"
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
	plus                  // +
	minus                 // -
	star                  // *
	forward_slash         // /
	percent               // %
	increment             // ++
	decrement             // --
	pipe                  // |
	ampersand             // &
	assignment            // =
	arithmetic_assignment // += -= *= /= %=
	exclamation
	power               // **
	binary_operator     // && ||
	comparison_operator // == != < > <= >=
	caret               // ^
	channel             // <-
	then                // ->

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
	gen_keyword
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
	new_line:       "new line",
	whitespace:     "whitespace",
	any_whitespace: "whitespace",

	left_parens:        "(",
	right_parens:       ")",
	left_angle_bracks:  "<",
	right_angle_bracks: ">",
	left_squre_bracks:  "[",
	right_squre_bracks: "]",
	left_curly_bracks:  "{",
	right_curly_bracks: "}",
	dot:                ".",
	comma:              ",",
	colon:              ":",
	semicolon:          ";",
	variadic_marker:    "... (variadic marker)", // ...

	multi_line_comment:  "MultiLineComment",
	single_line_comment: "SingleLineComment",

	Operator:              "Operator",
	plus:                  "+",
	minus:                 "-",
	star:                  "*",
	forward_slash:         "/",
	percent:               "%",
	increment:             "++",
	decrement:             "--",
	power:                 "**",
	pipe:                  "|",
	ampersand:             "&",
	assignment:            "assignment operator",            // =
	arithmetic_assignment: "arithmetic assignment operator", // += -= *= /= %=
	exclamation:           "!",
	binary_operator:       "binary operator",     // && ||
	comparison_operator:   "comparison operator", // == != <= >=
	caret:                 "caret",
	channel:               "Channel (To Be Removed)",
	then:                  "Then (To Be Removed)",

	identifier: "Identifier",

	literal:          "Literal",
	string_literal:   "string literal",
	rune_literal:     "rune literal",
	number_literal:   "number literal",
	bool_literal:     "bool literal",
	cardinal_literal: "cardinal literal",

	Keyword:            "Keyword",
	as_keyword:         "as keyword",
	base_keyword:       "base keyword",
	break_keyword:      "break keyword",
	const_keyword:      "const keyword",
	continue_keyword:   "continue keyword",
	corout_keyword:     "corout keyword",
	else_keyword:       "else keyword",
	for_keyword:        "for keyword",
	fun_keyword:        "fun keyword",
	gen_keyword:        "gen keyword",
	giveup_keyword:     "give up keyword",
	hidden_keyword:     "hidden keyword",
	if_keyword:         "if keyword",
	implements_keyword: "implements keyword",
	instanceof_keyword: "instanceof keyword",
	match_keyword:      "match keyword",
	map_keyword:        "map keyword",
	mimics_keyword:     "mimics keyword",
	move_keyword:       "move keyword",
	of_keyword:         "of keyword",
	or_keyword:         "or keyword",
	package_keyword:    "package keyword",
	return_keyword:     "return keyword",
	this_keyword:       "this keyword",
	trait_keyword:      "trait keyword",
	type_keyword:       "type keyword",
	use_keyword:        "use keyword",
	var_keyword:        "var keyword",
	yield_keyword:      "yield keyword",
}

type Token struct {
	Kind       token_kind      `json:"kind"`
	Location   errors.Location `json:"location"`
	Literal    string          `json:"literal"`
	Raw        string          `json:"raw"`
	Offset     int             `json:"offset"`
	LineBreaks int             `json:"line_breaks"`
}

func (t Token) String() string {
	if t.Kind < Keyword {
		if t.Kind >= Whitespace && t.Kind < Punctuation {
			return fmt.Sprintf("ws[%d:%d][%d]", t.Location.Start.Line, t.Location.Start.Column, t.Offset)
		} else {
			return fmt.Sprintf("%s(%s)[%d:%d][%d]", token_map[t.Kind], t.Literal, t.Location.Start.Line, t.Location.Start.Column, t.Offset)
		}
	}

	return fmt.Sprintf("<%s>[%d:%d][%d]", t.Literal, t.Location.Start.Line, t.Location.Start.Column, t.Offset)
}

const eof rune = -1

var eof_token = Token{
	Kind: eof_token_kind,
}
