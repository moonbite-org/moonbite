package parser

import (
	"fmt"
	"slices"
	"strings"
	"unicode"

	"github.com/moonbite-org/moonbite/common"
)

type lexer struct {
	input    []rune
	offset   int
	tokens   []Token
	location common.Location
	error    common.Error
}

var keywords = map[string]token_kind{
	"as":         as_keyword,
	"base":       base_keyword,
	"break":      break_keyword,
	"const":      const_keyword,
	"continue":   continue_keyword,
	"corout":     corout_keyword,
	"else":       else_keyword,
	"for":        for_keyword,
	"fun":        fun_keyword,
	"giveup":     giveup_keyword,
	"hidden":     hidden_keyword,
	"if":         if_keyword,
	"implements": implements_keyword,
	"instanceof": instanceof_keyword,
	"match":      match_keyword,
	"map":        map_keyword,
	"mimics":     mimics_keyword,
	"of":         of_keyword,
	"or":         or_keyword,
	"package":    package_keyword,
	"return":     return_keyword,
	"this":       this_keyword,
	"trait":      trait_keyword,
	"type":       type_keyword,
	"use":        use_keyword,
	"var":        var_keyword,
	"yield":      yield_keyword,
}

var bool_literals = []string{"true", "false"}
var cardinal_literals = []string{"string", "bool", "rune", "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "float32", "float64"}

func (l *lexer) throw(reason string) {
	l.error = common.Error{
		Kind:     common.SyntaxError,
		Reason:   reason,
		Location: l.location,
		Exists:   true,
	}
}

func (l lexer) next_rune() rune {
	if l.offset+1 >= len(l.input) {
		return eof
	}

	return l.input[l.offset+1]
}

func (l lexer) next_runes(amount int) []rune {
	result := []rune{}

	for i := 0; i < amount; i++ {
		result = append(result, l.next_rune())
		l.advance()
	}

	l.backup_by(amount)

	return result
}

func (l *lexer) peek(at int) rune {
	runes := l.next_runes(at)
	return runes[len(runes)-1]
}

func (l lexer) current_rune() rune {
	if l.offset >= len(l.input) {
		return eof
	}

	return l.input[l.offset]
}

func (l *lexer) advance() {
	l.offset++
}

func (l *lexer) advance_by(n int) {
	if l.offset+n <= len(l.input) {
		l.offset += n
	}
}

func (l *lexer) backup() {
	l.offset--
}

func (l *lexer) backup_by(n int) {
	if l.offset-n >= 0 {
		l.offset -= n
	}
}

func (l lexer) create_token(kind token_kind, length int) Token {
	raw := l.input[l.offset : l.offset+length]
	literal := raw

	line_breaks := 0
	for _, r := range literal {
		if r == '\n' || r == '\r' {
			line_breaks++
		}
	}

	if kind == string_literal {
		literal = []rune(strings.ReplaceAll(string(literal), "\\", ""))
	}

	return Token{
		Kind:       kind,
		Location:   l.location,
		Literal:    string(literal),
		Raw:        string(raw),
		Offset:     l.offset,
		LineBreaks: line_breaks,
	}
}

func (l *lexer) register_token(token Token) {
	if token.LineBreaks != 0 {
		l.location.Column = 1
	} else {
		l.location.Column += len(token.Raw)
	}

	l.location.Line += token.LineBreaks
	l.location.Offset += len(token.Raw)

	l.tokens = append(l.tokens, token)
	l.advance_by(len(token.Raw))
}

func lex(input []byte, filename string) ([]Token, common.Error) {
	lexer := lexer{
		input: []rune(string(input)),
		location: common.Location{
			Line:   1,
			Column: 1,
			Offset: 0,
			File:   filename,
		},
	}

	control_chars := []rune{'(', ')', '<', '>', '[', ']', '{', '}', '.', ',', ':', ';'}
	operator_chars := []rune{'+', '-', '*', '/', '%', '=', '^', '&', '|', '!'}

	current := lexer.current_rune()
	for lexer.current_rune() != eof {
		if lexer.error.Exists {
			return []Token{}, lexer.error
		}

		switch {
		case unicode.IsSpace(current):
			lexer.lex_whitespace()
		case slices.Contains(control_chars, current):
			lexer.lex_control_chars()
		case slices.Contains(operator_chars, current):
			lexer.lex_operator_chars()
		case unicode.IsDigit(current):
			lexer.lex_number_literal()
		case current == '"':
			lexer.lex_string_literal()
		case current == '`':
			lexer.lex_multi_line_string_literal()
		case current == '\'':
			lexer.lex_rune_literal()
		default:
			lexer.lex_alpha_numeric()
		}
		current = lexer.current_rune()
	}

	return lexer.tokens, common.Error{}
}

func (l *lexer) lex_whitespace() {
	length := 0

	for unicode.IsSpace(l.current_rune()) {
		if l.current_rune() == '\n' || l.current_rune() == '\r' {
			length++
			l.advance()
			l.backup_by(length)

			token := l.create_token(new_line, length)
			l.register_token(token)
			return
		}
		length++
		l.advance()
	}

	l.backup_by(length)

	token := l.create_token(whitespace, length)
	l.register_token(token)
}

func (l *lexer) lex_string_literal() {
	length := 0
	for l.next_rune() != '"' && l.next_rune() != eof && l.next_rune() != '\n' && l.next_rune() != '\r' {
		if l.current_rune() == '\\' {
			length++
			l.advance()
		}

		length++
		l.advance()
	}

	if l.next_rune() == eof || l.next_rune() == '\n' || l.next_rune() == '\r' {
		l.throw(fmt.Sprintf(common.ErrorMessages["u_eof"], "a '\"' (double quote) to close the string literal"))
	} else {
		l.backup_by(length - 1)
		l.register_token(l.create_token(string_literal, length))
		l.advance()
	}
}

func (l *lexer) lex_multi_line_string_literal() {
	length := 0
	for l.next_rune() != '`' && l.next_rune() != eof {
		if l.current_rune() == '\\' {
			length++
			l.advance()
		}

		length++
		l.advance()
	}

	if l.next_rune() == eof {
		l.throw(fmt.Sprintf(common.ErrorMessages["u_eof"], "a '`' (back quote) to close the multiline string literal"))
	} else {
		l.backup_by(length - 1)
		l.register_token(l.create_token(string_literal, length))
		l.advance()
	}
}

func (l *lexer) lex_rune_literal() {
	length := 0
	for l.next_rune() != '\'' && l.next_rune() != eof && l.next_rune() != '\n' && l.next_rune() != '\r' {
		if l.current_rune() == '\\' {
			length++
			l.advance()
		}

		length++
		l.advance()
	}

	if l.next_rune() == eof || l.next_rune() == '\n' || l.next_rune() == '\r' {
		l.throw(fmt.Sprintf(common.ErrorMessages["u_eof"], "a \"'\" (single quote) to close the rune literal"))
	} else {
		if length == 1 {
			l.backup_by(length - 1)
			l.register_token(l.create_token(rune_literal, length))
			l.advance()
		} else {
			l.throw(fmt.Sprintf(common.ErrorMessages["i_val"], "Rune literals must exactly be 1 character"))
		}
	}
}

func (l *lexer) lex_number_literal() {
	length := 1

	if l.current_rune() == '0' {
		if unicode.IsDigit(l.next_rune()) {
			l.throw("malformed number")
			return
		} else {
			l.register_token(l.create_token(number_literal, 1))
			return
		}
	}

	for unicode.IsDigit(l.next_rune()) {
		length++
		l.advance()
	}

	if l.next_rune() == '.' && unicode.IsDigit(l.peek(2)) {
		l.advance()
		length++

		for unicode.IsDigit(l.next_rune()) {
			length++
			l.advance()
		}

		if l.next_rune() == 'e' || l.next_rune() == 'E' {
			l.advance()
			length++

			for unicode.IsDigit(l.next_rune()) {
				length++
				l.advance()
			}
		}
	}

	l.backup_by(length - 1)
	l.register_token(l.create_token(number_literal, length))
}

func (l *lexer) lex_single_line_comment() {
	length := 0

	l.backup()

	if string(l.next_runes(2)) != "//" {
		l.advance()
		return
	}

	l.advance_by(2)

	current := l.current_rune()
	for current != '\n' && current != '\r' && current != eof {
		length++
		l.advance()
		current = l.current_rune()
	}

	l.backup_by(length + 1)
	l.register_token(l.create_token(single_line_comment, length+1))
}

func (l *lexer) lex_multiL_line_comment() {
	length := 0

	for string(l.next_runes(2)) != "*/" && l.next_rune() != eof {
		length++
		l.advance()
	}

	if l.next_rune() == eof {
		l.throw("unexpected eof unclosed multiline comment")
	} else {
		l.backup_by(length)
		l.register_token(l.create_token(multi_line_comment, length+3))
	}
}

func (l *lexer) lex_operator_chars() {
	var token Token

	switch l.current_rune() {
	case '=':
		if l.next_rune() == '=' {
			token = l.create_token(binary_operator, 2)
		} else {
			token = l.create_token(assignment, 1)
		}
	case '+':
		next := l.next_rune()
		if next == '=' {
			token = l.create_token(arithmetic_assignment, 2)
		} else if next == '+' {
			token = l.create_token(increment, 2)
		} else {
			token = l.create_token(plus, 1)
		}
	case '-':
		next := l.next_rune()
		if next == '=' {
			token = l.create_token(arithmetic_assignment, 2)
		} else if next == '-' {
			token = l.create_token(decrement, 2)
		} else if next == '>' {
			token = l.create_token(then, 2)
		} else if unicode.IsDigit(next) {
			l.lex_number_literal()
			return
		} else {
			token = l.create_token(minus, 1)
		}
	case '*':
		next := l.next_rune()
		if next == '=' {
			token = l.create_token(arithmetic_assignment, 2)
		} else if next == '*' {
			token = l.create_token(power, 2)
		} else {
			token = l.create_token(star, 1)
		}
	case '/':
		next := l.next_rune()
		if next == '=' {
			token = l.create_token(arithmetic_assignment, 2)
		} else if next == '/' {
			l.lex_single_line_comment()
			return
		} else if next == '*' {
			l.lex_multiL_line_comment()
			return
		} else {
			token = l.create_token(forward_slash, 1)
		}
	case '%':
		token = l.create_token(percent, 1)
	case '!':
		token = l.create_token(exclamation, 1)
	case '&':
		next := l.next_rune()

		if next == '&' {
			token = l.create_token(binary_operator, 2)
		} else {
			token = l.create_token(ampersand, 1)
		}
	case '|':
		next := l.next_rune()

		if next == '|' {
			token = l.create_token(binary_operator, 2)
		} else {
			token = l.create_token(pipe, 1)
		}
	case '^':
		token = l.create_token(caret, 1)
	}

	l.register_token(token)
}

func (l *lexer) lex_control_chars() {
	var token Token

	switch l.current_rune() {
	case '(':
		token = l.create_token(left_parens, 1)
	case ')':
		token = l.create_token(right_parens, 1)
	case '[':
		token = l.create_token(left_squre_bracks, 1)
	case ']':
		token = l.create_token(right_squre_bracks, 1)
	case '{':
		token = l.create_token(left_curly_bracks, 1)
	case '}':
		token = l.create_token(right_curly_bracks, 1)
	case ',':
		token = l.create_token(comma, 1)
	case ':':
		token = l.create_token(colon, 1)
	case ';':
		token = l.create_token(semicolon, 1)
	case '<':
		if l.next_rune() == '=' {
			token = l.create_token(binary_operator, 2)
		} else if l.next_rune() == '-' {
			token = l.create_token(channel, 2)
		} else {
			token = l.create_token(left_angle_bracks, 1)
		}
	case '>':
		if l.next_rune() == '=' {
			token = l.create_token(binary_operator, 2)
		} else {
			token = l.create_token(right_angle_bracks, 1)
		}
	case '.':
		if string(l.next_runes(2)) == ".." {
			token = l.create_token(variadic_marker, 3)
		} else {
			token = l.create_token(dot, 1)
		}
	}

	l.register_token(token)
}

func (l *lexer) lex_alpha_numeric() {
	length := 1

	for unicode.IsDigit(l.next_rune()) || unicode.IsLetter(l.next_rune()) || l.next_rune() == '_' {
		length++
		l.advance()
	}

	l.backup_by(length - 1)

	literal := string(l.input[l.offset : l.offset+length])

	if keywords[literal] != 0 {
		l.register_token(l.create_token(keywords[literal], length))
	} else if slices.Contains(bool_literals, literal) {
		l.register_token(l.create_token(bool_literal, length))
	} else if slices.Contains(cardinal_literals, literal) {
		l.register_token(l.create_token(cardinal_literal, length))
	} else {
		l.register_token(l.create_token(identifier, length))
	}
}
