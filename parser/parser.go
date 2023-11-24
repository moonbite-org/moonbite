package parser

import (
	"fmt"
	"path"
	"reflect"
)

type parser_s struct {
	input        []byte
	offset       int
	tokens       []Token
	error        Error
	package_done bool
	ast          Ast
	expressions  []Expression
}

func (p *parser_s) current_expression() Expression {
	if len(p.expressions) == 0 {
		return nil
	}
	return p.expressions[len(p.expressions)-1]
}

func (p *parser_s) set_current_expression(expression Expression) {
	if len(p.expressions) == 0 {
		p.expressions = append(p.expressions, expression)
	}

	p.expressions[len(p.expressions)-1] = expression
}

func (p *parser_s) pop_expression() {
	p.expressions = p.expressions[0 : len(p.expressions)-1]
}

func (p *parser_s) push_expression() {
	p.expressions = append(p.expressions, nil)
}

func (p *parser_s) parse_program() {
	defer p.catch()

	if !p.package_done {
		p.must_expect([]token_kind{package_keyword})
		p.backup()
		p.ast.Package = p.parse_package_statement()
		return
	}

	token := p.might_only_expect([]token_kind{use_keyword, type_keyword, trait_keyword, fun_keyword, var_keyword, const_keyword, single_line_comment, multi_line_comment})

	if token == nil {
		return
	}

	switch token.Kind {
	case use_keyword:
		p.backup()
		p.ast.Uses = append([]UseStatement{p.parse_use_statement()}, p.ast.Uses...)
	case fun_keyword:
		p.backup()
		p.ast.Definitions = append([]Definition{p.parse_tl_fun_definition_statement()}, p.ast.Definitions...)
	case type_keyword:
		p.backup()
		p.ast.Definitions = append([]Definition{p.parse_type_definition_statement()}, p.ast.Definitions...)
	case trait_keyword:
		p.backup()
		p.ast.Definitions = append([]Definition{p.parse_trait_definition_statement()}, p.ast.Definitions...)
	case var_keyword, const_keyword:
		p.backup()
		fmt.Println(p.parse_declaration_statement())
	case single_line_comment, multi_line_comment:
		p.backup()
		comment := p.parse_tl_comment()
		p.ast.Definitions = append([]Definition{comment}, p.ast.Definitions...)
		p.ast.Comments = append([]Comment{comment}, p.ast.Comments...)
	}
}

func (p *parser_s) parse_inline_level_statements() []Statement {
	defer p.catch()

	p.skip_whitespace()
	result := []Statement{}

	token := p.might_only_expect([]token_kind{var_keyword, const_keyword, for_keyword, match_keyword, if_keyword, return_keyword, identifier, right_curly_bracks})

	switch token.Kind {
	case return_keyword:
		p.backup()
		result = append(result, p.parse_return_statement())
	case right_curly_bracks:
		p.backup()
		return result
	}

	return result
}

func (p *parser_s) parse_package_statement() PackageStatement {
	defer p.catch()
	location := p.current_token().Location

	p.must_expect([]token_kind{package_keyword})
	p.must_expect([]token_kind{whitespace, new_line})
	p.skip_whitespace()
	ident := p.must_expect([]token_kind{identifier})
	p.must_expect([]token_kind{whitespace, new_line, eof_token_kind})
	p.skip_whitespace()

	statement := PackageStatement{
		Name:     *p.create_ident(ident),
		location: location,
	}
	p.package_done = true
	p.parse_program()
	return statement
}

func (p *parser_s) parse_use_statement() UseStatement {
	defer p.catch()
	location := p.current_token().Location

	p.must_expect([]token_kind{use_keyword})
	p.must_expect([]token_kind{whitespace, new_line})
	p.skip_whitespace()
	ident := p.must_expect([]token_kind{identifier})
	p.must_expect([]token_kind{whitespace, new_line, eof_token_kind})
	p.skip_whitespace()
	as := p.might_expect([]token_kind{as_keyword})

	statement := UseStatement{
		Resource: *p.create_ident(ident),
		location: location,
	}

	if as != nil {
		p.must_expect([]token_kind{whitespace, new_line})
		ident := p.must_expect([]token_kind{identifier})
		p.must_expect([]token_kind{whitespace, new_line, eof_token_kind})
		p.skip_whitespace()
		statement.As = p.create_ident(ident)
	}

	p.parse_program()
	return statement
}

func (p *parser_s) parse_declaration_statement() DeclarationStatement {
	defer p.catch()

	start := p.must_expect([]token_kind{var_keyword, const_keyword})
	p.must_expect([]token_kind{whitespace, new_line})

	p.advance()
	ws := p.skip_whitespace()
	next := p.must_expect([]token_kind{identifier, assignment})

	var value *Expression
	var typ *TypeLiteral
	var name IdentifierExpression

	p.backup_by(2 + ws)

	switch next.Kind {
	case assignment:
		n := p.must_expect([]token_kind{identifier})
		name = *p.create_ident(n)
		p.skip_whitespace()
	case identifier:
		t := p.parse_type_literal()
		typ = &t
		p.skip_whitespace()

		n := p.must_expect([]token_kind{identifier})
		name = *p.create_ident(n)
		p.skip_whitespace()
	}

	is_assigned := p.might_expect([]token_kind{assignment})

	if is_assigned != nil {
		p.skip_whitespace()
		v := p.parse_expression()

		if v != nil {
			value = &v
		}
	}

	var kind var_kind

	if start.Literal == "var" {
		kind = variable
	} else {
		kind = constant
	}

	return DeclarationStatement{
		VarKind: kind,
		Name:    name,
		Type:    typ,
		Value:   value,
	}
}

func (p *parser_s) parse_tl_comment() Comment {
	p.catch()

	var result Comment
	current := p.current_token()

	switch current.Kind {
	case single_line_comment:
		result = SingleLineCommentStatement{Comment: current.Literal, location: current.Location}
	case multi_line_comment:
		result = MultiLineCommentStatement{Comment: current.Literal, location: current.Location}
	default:
		result = SingleLineCommentStatement{}
	}

	p.advance()
	p.skip_whitespace()

	p.parse_program()
	return result
}

func (p *parser_s) parse_type_literal() TypeLiteral {
	defer p.catch()

	var result TypeLiteral

	switch p.current_token().Kind {
	case left_curly_bracks:
		result = p.parse_struct_literal()
	case left_parens:
		p.advance()
		p.skip_whitespace()
		typ := p.parse_type_literal()
		p.must_expect([]token_kind{right_parens})
		p.skip_whitespace()
		result = GroupType{Type: typ}
	default:
		result = p.parse_type_identifier()
	}

	is_typed_literal := p.might_expect([]token_kind{left_parens})

	if is_typed_literal != nil {
		if reflect.TypeOf(result) == reflect.TypeOf(TypeIdentifier{}) {
			p.skip_whitespace()
			value := p.parse_literal_expression()
			p.skip_whitespace()
			p.must_expect([]token_kind{right_parens})

			result = TypedLiteral{
				Type:    result.(TypeIdentifier),
				Literal: value,
			}
		} else {
			p.throw(fmt.Sprintf(error_messages["i_con"], "use a struct literal to create a type literal"))
		}
	}
	p.skip_whitespace()

	is_operated := p.might_expect([]token_kind{pipe, ampersand})

	if is_operated != nil {
		p.skip_whitespace()
		right_hand := p.parse_type_literal()
		result = OperatedType{
			LeftHandSide:  result,
			RightHandSide: right_hand,
			Operator:      is_operated.Literal,
		}
	}

	return result
}

func (p *parser_s) parse_type_identifier() TypeIdentifier {
	defer p.catch()

	name := p.must_expect([]token_kind{identifier, cardinal_literal})

	result := TypeIdentifier{
		Name:     p.create_ident(name),
		Generics: []TypeLiteral{},
	}

	is_generic := p.might_expect([]token_kind{left_angle_bracks})

	if is_generic != nil {
		p.backup()
		result.Generics = parse_seperated_list(p, p.parse_type_literal, comma, left_angle_bracks, right_angle_bracks, false, false)
	}

	return result
}

func (p *parser_s) parse_value_type_pair() ValueTypePair {
	defer p.catch()

	key := p.must_expect([]token_kind{identifier})
	p.skip_whitespace()
	typ := p.parse_type_literal()

	return ValueTypePair{
		Key:  *p.create_ident(key),
		Type: typ,
	}
}

func (p *parser_s) parse_struct_literal() StructLiteral {
	defer p.catch()

	var pairs StructLiteral = parse_seperated_list(p, p.parse_value_type_pair, semicolon, left_curly_bracks, right_curly_bracks, true, true)

	return pairs
}

func (p *parser_s) parse_constrained_type() ConstrainedType {
	defer p.catch()

	name := p.must_expect([]token_kind{identifier})
	is_spaced := p.might_expect([]token_kind{whitespace, new_line})
	var constraint *TypeLiteral

	if is_spaced != nil {
		p.skip_whitespace()

		is_constrained := p.might_expect([]token_kind{identifier})

		if is_constrained != nil {
			p.backup()
			c := p.parse_type_literal()
			constraint = &c
			p.skip_whitespace()
		}
	}

	return ConstrainedType{
		Name:       p.create_ident(name),
		Constraint: constraint,
	}
}

func (p *parser_s) parse_type_definition_statement() TypeDefinitionStatement {
	defer p.catch()

	start := p.must_expect([]token_kind{type_keyword})
	p.must_expect([]token_kind{whitespace, new_line})
	p.skip_whitespace()
	name := p.must_expect([]token_kind{identifier})

	result := TypeDefinitionStatement{
		Name:            *p.create_ident(name),
		Generics:        []ConstrainedType{},
		Implementations: []TypeIdentifier{},
		Definition:      TypeIdentifier{},
		location:        start.Location,
	}

	is_generic := p.might_expect([]token_kind{left_angle_bracks})

	if is_generic != nil {
		p.backup()
		result.Generics = parse_seperated_list(p, p.parse_constrained_type, comma, left_angle_bracks, right_angle_bracks, false, false)
	}

	p.skip_whitespace()
	is_implementing := p.might_expect([]token_kind{implements_keyword})

	if is_implementing != nil {
		p.skip_whitespace()
		result.Implementations = parse_seperated_list(p, p.parse_type_identifier, comma, left_squre_bracks, right_squre_bracks, false, false)
	}

	result.Definition = p.parse_type_literal()
	p.skip_whitespace()

	p.parse_program()
	return result
}

func (p *parser_s) parse_typed_parameter() TypedParameter {
	defer p.catch()

	name := p.must_expect([]token_kind{identifier})
	p.must_expect([]token_kind{whitespace, new_line})
	p.skip_whitespace()
	typ := p.parse_type_literal()

	return TypedParameter{
		Name: *p.create_ident(name),
		Type: typ,
	}
}

func (p *parser_s) parse_trait_definition_statement() TraitDefinitionStatement {
	defer p.catch()

	p.must_expect([]token_kind{trait_keyword})
	p.must_expect([]token_kind{whitespace, new_line})
	p.skip_whitespace()
	name := p.must_expect([]token_kind{identifier})

	result := TraitDefinitionStatement{
		Name:     *p.create_ident(name),
		Generics: []ConstrainedType{},
		Mimics:   []TypeIdentifier{},
		location: name.Location,
	}

	is_generic := p.might_expect([]token_kind{left_angle_bracks})

	if is_generic != nil {
		p.backup()
		result.Generics = parse_seperated_list(p, p.parse_constrained_type, comma, left_angle_bracks, right_angle_bracks, false, false)
	}

	p.skip_whitespace()

	is_mimicked := p.might_expect([]token_kind{mimics_keyword})

	if is_mimicked != nil {
		p.skip_whitespace()
		result.Mimics = parse_seperated_list(p, p.parse_type_identifier, comma, left_squre_bracks, right_squre_bracks, false, false)
		p.skip_whitespace()
	}

	result.Definition = parse_seperated_list(p, p.parse_unbound_fun_signature, semicolon, left_curly_bracks, right_curly_bracks, true, true)

	p.skip_whitespace()

	p.parse_program()
	return result
}

func (p *parser_s) parse_unbound_fun_signature() UnboundFunctionSignature {
	defer p.catch()

	start := p.must_expect([]token_kind{fun_keyword})
	p.must_expect([]token_kind{whitespace, new_line})
	p.skip_whitespace()
	name := p.must_expect([]token_kind{identifier})

	generics := []ConstrainedType{}
	is_generic := p.might_expect([]token_kind{left_angle_bracks})

	if is_generic != nil {
		p.backup()
		generics = parse_seperated_list(p, p.parse_constrained_type, comma, left_angle_bracks, right_angle_bracks, false, false)
	}

	p.skip_whitespace()

	params := parse_seperated_list(p, p.parse_typed_parameter, comma, left_parens, right_parens, true, false)
	p.skip_whitespace()

	var return_type *TypeLiteral
	p.skip_whitespace()
	return_type_t := p.might_expect([]token_kind{identifier, cardinal_literal})

	if return_type_t != nil {
		p.backup()
		return_type_p := p.parse_type_literal()
		return_type = &return_type_p
	}

	p.skip_whitespace()

	return UnboundFunctionSignature{
		Name:       *p.create_ident(name),
		Parameters: params,
		Generics:   generics,
		ReturnType: return_type,
		location:   start.Location,
	}
}

func (p *parser_s) parse_bound_fun_signature() BoundFunctionSignature {
	defer p.catch()

	start := p.must_expect([]token_kind{fun_keyword})
	p.must_expect([]token_kind{whitespace, new_line})
	p.skip_whitespace()
	p.must_expect([]token_kind{for_keyword})
	p.must_expect([]token_kind{whitespace, new_line})
	p.skip_whitespace()
	for_typ := p.parse_type_identifier()
	p.must_expect([]token_kind{whitespace, new_line})
	p.skip_whitespace()
	name := p.must_expect([]token_kind{identifier})

	generics := []ConstrainedType{}
	is_generic := p.might_expect([]token_kind{left_angle_bracks})

	if is_generic != nil {
		p.backup()
		generics = parse_seperated_list(p, p.parse_constrained_type, comma, left_angle_bracks, right_angle_bracks, false, false)
	}

	p.skip_whitespace()

	params := parse_seperated_list(p, p.parse_typed_parameter, comma, left_parens, right_parens, true, false)
	p.skip_whitespace()

	var return_type *TypeLiteral
	p.skip_whitespace()
	return_type_t := p.might_expect([]token_kind{identifier, cardinal_literal})

	if return_type_t != nil {
		p.backup()
		return_type_p := p.parse_type_literal()
		return_type = &return_type_p
	}

	p.skip_whitespace()

	return BoundFunctionSignature{
		Name:       *p.create_ident(name),
		For:        for_typ,
		Parameters: params,
		Generics:   generics,
		ReturnType: return_type,
		location:   start.Location,
	}
}

func (p *parser_s) parse_tl_fun_definition_statement() FunDefinitionStatement {
	definition := p.parse_fun_definition_statement()

	p.parse_program()
	return definition
}

func (p *parser_s) parse_fun_definition_statement() FunDefinitionStatement {
	defer p.catch()

	p.must_expect([]token_kind{fun_keyword})
	p.must_expect([]token_kind{whitespace, new_line})
	spaces := p.skip_whitespace()

	next := p.must_expect([]token_kind{identifier, for_keyword})
	p.backup_by(spaces + 1 /*identifier of for keyword */ + 1 /* single whitespace that is expected */ + 1 /* the fun keyword */)

	var definition FunDefinitionStatement

	switch next.Kind {
	case identifier:
		signature := p.parse_unbound_fun_signature()

		definition = &UnboundFunDefinitionStatement{
			Signature: signature,
			Body:      []Statement{},
			location:  signature.location,
		}
	case for_keyword:
		signature := p.parse_bound_fun_signature()

		definition = &BoundFunDefinitionStatement{
			Signature: signature,
			Body:      []Statement{},
			location:  signature.location,
		}
	default:
		p.throw("idk, something with the function")
	}

	p.skip_whitespace()
	p.must_expect([]token_kind{left_curly_bracks})
	definition.set_body(p.parse_inline_level_statements())
	p.must_expect([]token_kind{right_curly_bracks})
	p.skip_whitespace()

	return definition
}

func (p *parser_s) parse_return_statement() ReturnStatement {
	defer p.catch()

	start := p.must_expect([]token_kind{return_keyword})

	p.must_expect([]token_kind{whitespace, new_line})
	p.skip_whitespace()
	expression := p.parse_expression()

	p.skip_whitespace()

	return ReturnStatement{
		Value:    &expression,
		location: start.Location,
	}
}

func (p *parser_s) continue_expression() Expression {
	defer p.catch()

	// token := p.must_expect([]token_kind{identifier, rune_literal, string_literal, bool_literal, number_literal, left_squre_bracks, dot, left_parens, whitespace, new_line})

	switch p.current_token().Kind {
	case rune_literal, string_literal, bool_literal, number_literal, left_squre_bracks:
		p.set_current_expression(p.parse_literal_expression())
		return p.continue_expression()
	case identifier:
		ident := p.create_ident(p.current_token())
		p.advance()
		p.set_current_expression(*ident)
		return p.continue_expression()
	case dot:
		p.advance()

		if p.current_token().Kind == left_parens {
			return p.parse_type_cast_expression()
		} else {
			return p.parse_member_expression()
		}
	case left_parens:
		if p.current_expression() == nil {
			return p.parse_group_expression()
		} else {
			return p.parse_call_expression()
		}
	case plus, minus, star, forward_slash:
		return p.parse_arithmetic_expression()
	case whitespace, new_line:
		p.skip_whitespace()
		return p.continue_expression()
	default:
		result := p.current_expression()
		p.pop_expression()

		return result
	}
}

func (p *parser_s) parse_expression() Expression {
	defer p.catch()

	if p.current_expression() != nil {
		p.push_expression()
	}

	return p.continue_expression()
}

func (p *parser_s) parse_group_expression() Expression {
	defer p.catch()

	p.must_expect([]token_kind{left_parens})
	p.skip_whitespace()
	expression := p.parse_expression()
	p.skip_whitespace()
	p.must_expect([]token_kind{right_parens})

	p.set_current_expression(GroupExpression{Expression: expression})

	return p.continue_expression()
}

func (p *parser_s) parse_call_expression() Expression {
	defer p.catch()

	args := parse_seperated_list(p, p.parse_expression, comma, left_parens, right_parens, true, false)

	p.set_current_expression(CallExpression{
		Callee:    p.current_expression(),
		Arguments: args,
	})

	return p.continue_expression()
}

func (p *parser_s) parse_member_expression() Expression {
	defer p.catch()

	rhs_t := p.must_expect([]token_kind{identifier})
	rhs := p.create_ident(rhs_t)

	p.set_current_expression(MemberExpression{
		LeftHandSide:  p.current_expression(),
		RightHandSide: *rhs,
	})

	return p.continue_expression()
}

func (p *parser_s) parse_arithmetic_expression() Expression {
	defer p.catch()

	operator := p.must_expect([]token_kind{plus, minus, star, forward_slash})
	p.skip_whitespace()

	rhs := p.parse_expression()
	current := p.current_expression()

	switch operator.Kind {
	case star, forward_slash:
		if reflect.TypeOf(rhs) == reflect.TypeOf(ArithmeticExpression{}) {
			rhs := rhs.(ArithmeticExpression)

			if rhs.Operator != "*" && rhs.Operator != "/" {
				p.set_current_expression(ArithmeticExpression{
					LeftHandSide: ArithmeticExpression{
						LeftHandSide:  current,
						RightHandSide: rhs.LeftHandSide,
						Operator:      operator.Literal,
					},
					RightHandSide: rhs.RightHandSide,
					Operator:      rhs.Operator,
				})
			} else {
				p.set_current_expression(ArithmeticExpression{
					LeftHandSide:  current,
					RightHandSide: rhs,
					Operator:      operator.Literal,
				})
			}
		}
	case plus, minus:
		p.set_current_expression(ArithmeticExpression{
			LeftHandSide:  current,
			RightHandSide: rhs,
			Operator:      operator.Literal,
		})
	}

	return p.continue_expression()
}

func (p *parser_s) parse_type_cast_expression() Expression {
	defer p.catch()

	p.must_expect([]token_kind{left_parens})
	typ := p.parse_type_identifier()
	p.must_expect([]token_kind{right_parens})

	p.set_current_expression(TypeCastExpression{
		Value: p.continue_expression(),
		Type:  typ,
	})

	return p.continue_expression()
}

func (p *parser_s) parse_literal_expression() LiteralExpression {
	defer p.catch()

	current := p.current_token()
	var result LiteralExpression

	switch current.Kind {
	case string_literal:
		result = StringLiteralExpression{
			Value:    current.Literal,
			location: current.Location,
		}
		p.advance()
	case rune_literal:
		result = RuneLiteralExpression{
			Value:    rune(current.Literal[0]),
			location: current.Location,
		}
		p.advance()
	case bool_literal:
		result = BoolLiteralExpression{
			Value:    current.Literal == "true",
			location: current.Location,
		}
		p.advance()
	case left_squre_bracks:
		values := parse_seperated_list(p, p.parse_expression, comma, left_squre_bracks, right_squre_bracks, true, false)

		entries := []KeyValueEntry{}

		for i, value := range values {
			entries = append(entries, KeyValueEntry{
				Key:   StringLiteralExpression{Value: fmt.Sprintf("%d", i)},
				Value: value,
			})
		}

		result = ListLiteralExpression{
			Value:    entries,
			location: current.Location,
		}
	default:
		result = StringLiteralExpression{
			Value:    current.Literal,
			location: current.Location,
		}
		p.advance()
	}

	return result
}

func Parse(input []byte, filepath string) (Ast, Error) {
	filename := path.Base(filepath)

	parser := parser_s{
		input: input,
		ast: Ast{
			FileName: filename,
			FilePath: filepath,
			Uses:     []UseStatement{},
			Comments: []Comment{},
		},
	}

	tokens, err := lex(input, filename)
	parser.tokens = tokens

	if err.Exists {
		return Ast{}, err
	}

	parser.parse_program()

	if parser.error.Exists {
		return Ast{}, parser.error
	}

	return parser.ast, parser.error
}
