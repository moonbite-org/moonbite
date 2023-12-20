package parser

import (
	"fmt"
	"strconv"
	"strings"
)

func (p *parser_s) throw(reason string) {
	var location Location
	current := p.current_token()

	if current.Kind == eof_token_kind {
		split := strings.Split(string(p.input), "\n")
		last_line := len(split)
		last_col := len(split[len(split)-1]) + 1

		location = Location{Line: last_line, Column: last_col}
	} else {
		location = current.Location
	}

	p.error = Error{
		Kind:     SyntaxError,
		Reason:   reason,
		Location: location,
		Exists:   true,
	}
	panic(nil)
}

func (p *parser_s) catch() {
	// if r := recover(); r != nil {
	// 	panic(r)
	// }
	recover()
}

func (p parser_s) current_token() Token {
	if p.offset >= len(p.tokens) {
		return Token{
			Kind: eof_token_kind,
		}
	}

	return p.tokens[p.offset]
}

func (p *parser_s) advance() {
	p.offset++
}

func (p *parser_s) backup() {
	p.offset--
}

func (p *parser_s) backup_by(n int) {
	if p.offset-n >= 0 {
		p.offset -= n
	}
}

func (p *parser_s) skip_whitespace() int {
	count := 0
	for p.current_token().Kind == whitespace || p.current_token().Kind == new_line {
		p.advance()
		count++
	}
	return count
}

func (p *parser_s) skip() int {
	count := 0

	for p.current_token().Kind == single_line_comment || p.current_token().Kind == multi_line_comment || p.current_token().Kind == whitespace || p.current_token().Kind == new_line {
		p.advance()
		count++
	}
	return count
}

func (p *parser_s) must_expect(tokens []token_kind) Token {
	accepted := false

	for _, kind := range tokens {
		if p.current_token().Kind == kind {
			accepted = true
			break
		}
	}

	if accepted {
		token := p.current_token()
		p.advance()
		return token
	}

	expected := []string{}

	for _, kind := range tokens {
		expected = append(expected, token_map[kind])
	}

	if p.current_token().Kind == eof_token_kind {
		p.throw(fmt.Sprintf(ErrorMessages["u_eof"], strings.Join(expected, ", ")))
	} else {
		if len(tokens) == 1 {
			p.throw(fmt.Sprintf(ErrorMessages["u_tok_s"], token_map[tokens[0]], token_map[p.current_token().Kind]))
		} else {
			p.throw(fmt.Sprintf(ErrorMessages["u_tok"], token_map[p.current_token().Kind]))
		}
	}
	return eof_token
}

func (p *parser_s) might_expect(tokens []token_kind) *Token {
	accepted := false

	for _, kind := range tokens {
		if p.current_token().Kind == kind {
			accepted = true
			break
		}
	}

	if accepted {
		token := p.current_token()
		p.advance()
		return &token
	}
	return nil
}

func (p *parser_s) might_only_expect(tokens []token_kind) *Token {
	accepted := false

	for _, kind := range tokens {
		if p.current_token().Kind == kind {
			accepted = true
			break
		}
	}

	if accepted {
		token := p.current_token()
		p.advance()
		return &token
	}

	if p.current_token().Kind != eof_token_kind {
		if len(tokens) == 1 {
			p.throw(fmt.Sprintf(ErrorMessages["u_toks"], token_map[tokens[0]], token_map[p.current_token().Kind]))
		} else {
			p.throw(fmt.Sprintf(ErrorMessages["u_tok"], token_map[p.current_token().Kind]))
		}
	}
	return nil
}

func parse_seperated_list[T any](p *parser_s, parser_func func() T, seperator token_kind, opener, closer token_kind, allow_empty bool, trailing_seperator bool) []T {
	defer p.catch()

	result := []T{}
	p.must_expect([]token_kind{opener})
	p.skip_whitespace()
	p.skip()

	is_closing := p.might_expect([]token_kind{closer})

	if is_closing != nil && allow_empty {
		return result
	}

	done := false

	for !done {
		value := parser_func()
		result = append(result, value)
		p.skip_whitespace()
		p.skip()
		next := p.must_expect([]token_kind{seperator, closer})
		p.skip_whitespace()
		p.skip()

		switch next.Kind {
		case seperator:
			p.skip_whitespace()
			p.skip()
			if p.current_token().Kind == closer && trailing_seperator {
				p.advance()
				done = true
			} else if p.current_token().Kind == closer && !trailing_seperator {
				p.must_expect([]token_kind{})
				done = true
			}
		case closer:
			done = true
		}
	}

	return result
}

func (p *parser_s) create_ident(token Token) *IdentifierExpression {
	return &IdentifierExpression{Value: token.Literal, Kind_: IdentifierExpressionKind, location: token.Location}
}

func create_number_literal(p parser_s, literal string) NumberLiteral {
	defer p.catch()

	v, err := strconv.ParseFloat(literal, 32)

	if err != nil {
		p.throw(fmt.Sprintf(ErrorMessages["i_val"], err.Error()))
	}

	if strings.Contains(literal, ".") {
		return NumberLiteral{
			Type: TypeIdentifier{
				Name: IdentifierExpression{Value: "float32"},
			},
			Value: v,
		}
	}

	return NumberLiteral{
		Type: TypeIdentifier{
			Name:     IdentifierExpression{Value: "int"},
			Generics: map[int]TypeLiteral{},
		},
		Value: int(v),
	}
}

func generate_generics(p *parser_s) []ConstrainedType {
	generics := parse_seperated_list(p, p.parse_constrained_type, comma, left_angle_bracks, right_angle_bracks, false, false)

	for i, generic := range generics {
		generic.Index = i
	}

	return generics
}
