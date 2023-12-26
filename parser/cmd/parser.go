package parser

import (
	"fmt"
	"path"
	"reflect"

	errors "github.com/moonbite-org/moonbite/error"
)

// extra allowed keywords inside code blocks
var generator_function_context = []token_kind{return_keyword, yield_keyword}
var function_context = []token_kind{return_keyword}
var loop_context = []token_kind{break_keyword, continue_keyword}
var predicate_body_context = []token_kind{}

type parser_s struct {
	input                 []byte
	offset                int
	tokens                []Token
	error                 errors.Error
	ast                   Ast
	expressions           []Expression
	is_match_context      bool
	is_this_context       bool
	body_context          []token_kind
	previous_body_context []token_kind
}

type TopLevelResult struct {
	Definitions []Definition
	Uses        []UseStatement
	Comments    []Comment
}

func (r *TopLevelResult) merge(result TopLevelResult) {
	r.Definitions = append(r.Definitions, result.Definitions...)
	r.Uses = append(r.Uses, result.Uses...)
	r.Comments = append(r.Comments, result.Comments...)
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

func (p *parser_s) set_context(context []token_kind) {
	p.previous_body_context = p.body_context
	p.body_context = append(p.body_context, context...)
}

func (p *parser_s) reset_context() {
	p.body_context = p.previous_body_context
}

func (p *parser_s) parse_program() {
	defer p.catch()

	p.must_expect([]token_kind{package_keyword})
	p.backup()
	p.ast.Package = p.parse_package_statement()

	result := p.parse_top_level_statements()

	p.ast.Uses = result.Uses
	p.ast.Definitions = result.Definitions
	p.ast.Comments = result.Comments
}

func (p *parser_s) parse_top_level_statements() TopLevelResult {
	defer p.catch()

	token := p.might_only_expect([]token_kind{use_keyword, type_keyword, trait_keyword, fun_keyword, var_keyword, const_keyword, single_line_comment, multi_line_comment, hidden_keyword})

	result := TopLevelResult{
		Definitions: []Definition{},
		Uses:        []UseStatement{},
		Comments:    []Comment{},
	}

	if token == nil {
		return result
	}

	if token.Kind == hidden_keyword {
		p.must_expect([]token_kind{whitespace, new_line})
		ws := p.skip()

		token = p.might_only_expect([]token_kind{type_keyword, trait_keyword, fun_keyword, var_keyword, const_keyword})
		p.backup_by(1 /* current token */ + ws /* other whtitespaces */ + 1 /* required whitespace */)
	}

	switch token.Kind {
	case use_keyword:
		p.backup()
		result.Uses = append(result.Uses, p.parse_use_statement())
		result.merge(p.parse_top_level_statements())
	case fun_keyword:
		p.backup()
		result.Definitions = append(result.Definitions, p.parse_fun_definition_statement())
		result.merge(p.parse_top_level_statements())
	case type_keyword:
		p.backup()
		result.Definitions = append(result.Definitions, p.parse_type_definition_statement())
		result.merge(p.parse_top_level_statements())
	case trait_keyword:
		p.backup()
		result.Definitions = append(result.Definitions, p.parse_trait_definition_statement())
		result.merge(p.parse_top_level_statements())
	case var_keyword, const_keyword:
		p.backup()
		result.Definitions = append(result.Definitions, p.parse_declaration_statement())
		result.merge(p.parse_top_level_statements())
	case single_line_comment, multi_line_comment:
		p.backup()
		comment := p.parse_comment()
		result.Definitions = append(result.Definitions, comment)
		result.Comments = append(result.Comments, comment)
		result.merge(p.parse_top_level_statements())
	}

	return result
}

func (p *parser_s) parse_inline_level_statements() StatementList {
	defer p.catch()

	p.skip()
	result := StatementList{}

	allowed := []token_kind{var_keyword, const_keyword, for_keyword, match_keyword, if_keyword, identifier, right_curly_bracks, single_line_comment, multi_line_comment, hidden_keyword, corout_keyword, gen_keyword}
	allowed = append(allowed, p.body_context...)

	if p.is_this_context {
		allowed = append(allowed, this_keyword)
	}

	token := p.might_only_expect(allowed)

	if token.Kind == hidden_keyword {
		p.must_expect([]token_kind{})
	}

	switch token.Kind {
	case var_keyword, const_keyword:
		p.backup()
		result = append(result, p.parse_declaration_statement())
		result = append(result, p.parse_inline_level_statements()...)
	case return_keyword:
		p.backup()
		result = append(result, p.parse_return_statement())
		result = append(result, p.parse_inline_level_statements()...)
	case if_keyword:
		p.backup()
		result = append(result, p.parse_if_statement())
		result = append(result, p.parse_inline_level_statements()...)
	case yield_keyword:
		p.backup()
		result = append(result, p.parse_yield_statement())
		result = append(result, p.parse_inline_level_statements()...)
	case break_keyword, continue_keyword:
		p.backup()
		result = append(result, p.parse_flow_control_statement())
		result = append(result, p.parse_inline_level_statements()...)
	case single_line_comment, multi_line_comment:
		p.backup()
		result = append(result, p.parse_comment())
		result = append(result, p.parse_inline_level_statements()...)
	case for_keyword:
		p.backup()
		result = append(result, p.parse_loop_statement())
		result = append(result, p.parse_inline_level_statements()...)
	case right_curly_bracks:
		p.backup()
		return result
	default:
		// Could be an expression statement or an assignment statement
		p.backup()
		// Will record the last offset. If an assignment token is found, will go back and parse assignment statement.
		last_offset := p.offset
		expression := p.parse_expression()
		p.skip()

		is_assignment := p.might_expect([]token_kind{assignment, arithmetic_assignment})

		if is_assignment != nil {
			p.backup_by(p.offset - last_offset)

			result = append(result, p.parse_assignment_statement())
			result = append(result, p.parse_inline_level_statements()...)
		} else {
			result = append(result, ExpressionStatement{
				Expression: expression,
				Kind_:      ExpressionStatementKind,
				location:   expression.Location(),
			})
			result = append(result, p.parse_inline_level_statements()...)
		}
	}

	return result
}

func (p *parser_s) parse_package_statement() PackageStatement {
	defer p.catch()
	location := p.current_token().Location

	p.must_expect([]token_kind{package_keyword})
	p.must_expect([]token_kind{whitespace, new_line})
	p.skip()
	ident := p.must_expect([]token_kind{identifier})
	p.must_expect([]token_kind{whitespace, new_line, eof_token_kind})
	p.skip()

	statement := PackageStatement{
		Name:     *p.create_ident(ident),
		Kind_:    PackageStatementKind,
		location: location,
	}
	return statement
}

func (p *parser_s) parse_use_statement() UseStatement {
	defer p.catch()
	location := p.current_token().Location

	p.must_expect([]token_kind{use_keyword})
	p.must_expect([]token_kind{whitespace, new_line})
	p.skip()
	ident := p.must_expect([]token_kind{string_literal})
	p.must_expect([]token_kind{whitespace, new_line, eof_token_kind})
	p.skip()
	as := p.might_expect([]token_kind{as_keyword})

	statement := UseStatement{
		Resource: StringLiteralExpression{
			Value:    ident.Literal,
			location: ident.Location,
		},
		Kind_:    UseStatementKind,
		location: location,
	}

	if as != nil {
		p.must_expect([]token_kind{whitespace, new_line})
		ident := p.must_expect([]token_kind{identifier})
		p.must_expect([]token_kind{whitespace, new_line, eof_token_kind})
		p.skip()
		statement.As = p.create_ident(ident)
	}

	return statement
}

func (p *parser_s) parse_if_statement() IfStatement {
	defer p.catch()

	start := p.must_expect([]token_kind{if_keyword})
	p.skip()

	main_block := p.parse_predicate_block()
	if main_block.Predicate == nil {
		p.must_expect([]token_kind{})
	}
	p.skip()

	else_if_blocks := []PredicateBlock{}
	else_block := StatementList{}

	current := p.current_token()
	for current.Kind == else_keyword {
		p.advance()
		p.must_expect([]token_kind{whitespace, new_line})
		skipped := p.skip()
		is_else_if := p.might_expect([]token_kind{if_keyword})

		if is_else_if != nil {
			p.skip()
			predicate := p.parse_predicate_block()
			if predicate.Predicate == nil {
				break
			}
			else_if_blocks = append(else_if_blocks, predicate)
			current = p.current_token()
		} else {
			p.backup_by(skipped + 2)
			break
		}
	}

	is_else := p.might_expect([]token_kind{else_keyword})

	if is_else != nil {
		p.skip()
		p.set_context(predicate_body_context)
		else_block = p.parse_block()
		p.reset_context()
	}

	return IfStatement{
		MainBlock:    main_block,
		ElseIfBlocks: else_if_blocks,
		ElseBlock:    else_block,
		Kind_:        IfStatementKind,
		location:     start.Location,
	}
}

func (p *parser_s) parse_declaration_statement() DeclarationStatement {
	defer p.catch()

	var start Token
	var kind_n Token
	is_hidden := p.might_expect([]token_kind{hidden_keyword})

	if is_hidden != nil {
		start = *is_hidden
		p.skip()
		kind_n = p.must_expect([]token_kind{var_keyword, const_keyword})
	} else {
		start = p.must_expect([]token_kind{var_keyword, const_keyword})
		kind_n = start
	}

	p.must_expect([]token_kind{whitespace, new_line})
	p.skip()

	p.advance()
	ws := p.skip()
	next := p.must_expect([]token_kind{identifier, assignment})

	var value *Expression
	var typ *TypeLiteral
	var name IdentifierExpression

	p.backup_by(2 + ws)

	switch next.Kind {
	case assignment:
		n := p.must_expect([]token_kind{identifier})
		name = *p.create_ident(n)
		p.skip()
	case identifier:
		t := p.parse_type_literal()
		typ = &t
		p.skip()

		n := p.must_expect([]token_kind{identifier})
		name = *p.create_ident(n)
		p.skip()
	}

	is_assigned := p.might_expect([]token_kind{assignment})

	if is_assigned != nil {
		p.skip()
		v := p.parse_expression()

		if v == nil {
			p.must_expect([]token_kind{})
		}

		value = &v
	}

	p.skip()
	var kind VarKind

	if kind_n.Literal == "var" {
		kind = VariableKind
	} else {
		kind = ConstantKind
	}

	return DeclarationStatement{
		VarKind:  kind,
		Name:     name,
		Type:     typ,
		Value:    value,
		Hidden:   is_hidden != nil,
		Kind_:    DeclarationStatementKind,
		location: start.Location,
	}
}

func (p *parser_s) parse_assignment_statement() AssignmentStatement {
	defer p.catch()

	lhs := p.parse_expression()
	p.skip()
	operator := p.must_expect([]token_kind{assignment, arithmetic_assignment})
	p.skip()
	rhs := p.parse_expression()

	return AssignmentStatement{
		LeftHandSide:  lhs,
		RightHandSide: rhs,
		Operator: OperatorToken{
			Literal:  operator.Literal,
			location: operator.Location,
		},
		Kind_:    AssignmentStatementKind,
		location: lhs.Location(),
	}
}

func (p *parser_s) parse_loop_predicate() LoopPredicate {
	defer p.catch()

	token := p.might_expect([]token_kind{var_keyword, const_keyword, comma, identifier})

	if token == nil {
		return p.parse_unipartite_loop_predicate()
	}

	switch token.Kind {
	case comma, identifier:
		p.backup()
		return p.parse_bipartite_loop_predicate()
	case var_keyword, const_keyword:
		p.backup()
		return p.parse_tripartite_loop_predicate()
	default:
		// there is something wrong with the token, just throw
		p.must_expect([]token_kind{})
	}

	return nil
}

func (p *parser_s) parse_unipartite_loop_predicate() UnipartiteLoopPredicate {
	defer p.catch()

	expression := p.parse_expression()

	if expression == nil {
		p.must_expect([]token_kind{})
	}

	return UnipartiteLoopPredicate{
		Kind_: UnipartiteLoopKind,

		Expression: expression,
	}
}

func (p *parser_s) parse_bipartite_loop_predicate() LoopPredicate {
	defer p.catch()

	token := p.might_expect([]token_kind{identifier, comma})

	result := BipartiteLoopPredicate{
		Kind_: BipartiteLoopKind,
	}

	switch token.Kind {
	case identifier:
		result.Key = p.create_ident(*token)
		is_still_bipirtite := p.might_expect([]token_kind{comma})

		if is_still_bipirtite == nil {
			p.backup()
			return p.parse_unipartite_loop_predicate()
		}
	case comma:
		result.Key = nil
	}

	p.skip()
	token = p.might_expect([]token_kind{identifier, of_keyword})

	switch token.Kind {
	case identifier:
		result.Value = p.create_ident(*token)
		p.skip()
		p.must_expect([]token_kind{of_keyword})
		p.skip()
	case of_keyword:
		p.must_expect([]token_kind{whitespace, new_line})
		p.skip()
	default:
		p.skip()
		result.Value = nil
		p.must_expect([]token_kind{of_keyword})
		p.skip()
	}

	iterator := p.parse_expression()

	if iterator == nil {
		p.must_expect([]token_kind{})
	}

	result.Iterator = iterator

	return result
}

func (p *parser_s) parse_tripartite_loop_predicate() TripartiteLoopPredicate {
	defer p.catch()

	result := TripartiteLoopPredicate{
		Kind_: TripartiteLoopKind,
	}

	is_decl_empty := p.might_expect([]token_kind{semicolon})

	if is_decl_empty == nil {
		decl := p.parse_declaration_statement()
		result.Declaration = &decl
	}

	p.skip()
	p.must_expect([]token_kind{semicolon})
	predicate := p.parse_expression()

	if predicate == nil {
		p.must_expect([]token_kind{})
	}
	result.Predicate = predicate
	p.must_expect([]token_kind{semicolon})
	p.skip()

	procedure := p.parse_expression()

	if procedure != nil {
		result.Procedure = &procedure
	}

	return result
}

func (p *parser_s) parse_loop_statement() LoopStatement {
	defer p.catch()

	start := p.must_expect([]token_kind{for_keyword})
	result := LoopStatement{
		Kind_:    LoopStatementKind,
		location: start.Location,
	}

	p.must_expect([]token_kind{whitespace, new_line})
	p.skip()
	p.must_expect([]token_kind{left_parens})
	p.skip()
	result.Predicate = p.parse_loop_predicate()
	p.must_expect([]token_kind{right_parens})
	p.skip()
	p.set_context(loop_context)
	result.Body = p.parse_block()
	p.reset_context()

	return result
}

func (p *parser_s) parse_comment() Comment {
	p.catch()

	var result Comment
	current := p.current_token()

	switch current.Kind {
	case single_line_comment:
		result = SingleLineCommentStatement{Comment: current.Literal, location: current.Location, Kind_: SingleLineCommentStatementKind}
	case multi_line_comment:
		result = MultiLineCommentStatement{Comment: current.Literal, location: current.Location, Kind_: MultiLineCommentStatementKind}
	}

	p.advance()
	p.skip()

	return result
}

func (p *parser_s) parse_type_literal() TypeLiteral {
	defer p.catch()

	start := p.current_token().Location
	var result TypeLiteral

	switch p.current_token().Kind {
	case left_curly_bracks:
		result = p.parse_struct_literal()
	case fun_keyword:
		result = p.parse_anonymous_fun_signature()
	case left_parens:
		p.advance()
		p.skip()
		typ := p.parse_type_literal()
		p.must_expect([]token_kind{right_parens})
		p.skip()
		result = GroupType{Type: typ, TypeKind_: GroupTypeKind, location: start}
	default:
		result = p.parse_type_identifier()
	}

	is_typed_literal := p.might_expect([]token_kind{left_parens})

	if is_typed_literal != nil {
		if reflect.TypeOf(result) == reflect.TypeOf(TypeIdentifier{}) {
			p.skip()
			value := p.parse_literal_expression()
			p.skip()
			p.must_expect([]token_kind{right_parens})

			result = TypedLiteral{
				TypeKind_: TypedLiteralKind,
				Type:      result.(TypeIdentifier),
				Literal:   value,
				location:  start,
			}
		} else {
			p.throw(fmt.Sprintf(errors.ErrorMessages["i_con"], "use a struct literal to create a type literal"))
		}
	}
	p.skip()

	is_operated := p.might_expect([]token_kind{pipe, ampersand})

	if is_operated != nil {
		p.skip()
		right_hand := p.parse_type_literal()
		result = OperatedType{
			TypeKind_:     OperatedTypeKind,
			LeftHandSide:  result,
			RightHandSide: right_hand,
			Operator: OperatorToken{
				Literal:  is_operated.Literal,
				location: is_operated.Location,
			},
			location: result.Location(),
		}
	}

	return result
}

func (p *parser_s) parse_type_identifier() TypeIdentifier {
	defer p.catch()

	var name Expression
	is_cardinal := p.might_expect([]token_kind{cardinal_literal})

	if is_cardinal != nil {
		name = *p.create_ident(*is_cardinal)
	} else {
		expression := p.parse_type_expression()

		if expression == nil {
			p.must_expect([]token_kind{})
		} else {
			name = expression
		}
	}

	result := TypeIdentifier{
		TypeKind_: TypeIdentifierKind,
		Name:      name,
		Generics:  map[int]TypeLiteral{},
		location:  name.Location(),
	}

	is_generic := p.might_expect([]token_kind{left_angle_bracks})

	if is_generic != nil {
		p.backup()
		generics := parse_seperated_list(p, p.parse_type_literal, comma, left_angle_bracks, right_angle_bracks, false, false)

		for i, generic := range generics {
			result.Generics[i] = generic
		}
	}

	return result
}

func (p *parser_s) parse_simple_type_identifier() TypeIdentifier {
	defer p.catch()

	name := p.must_expect([]token_kind{identifier, cardinal_literal})

	result := TypeIdentifier{
		TypeKind_: TypeIdentifierKind,
		Name:      *p.create_ident(name),
		Generics:  map[int]TypeLiteral{},
		location:  name.Location,
	}

	is_generic := p.might_expect([]token_kind{left_angle_bracks})

	if is_generic != nil {
		p.backup()
		generics := parse_seperated_list(p, p.parse_type_literal, comma, left_angle_bracks, right_angle_bracks, false, false)

		for i, generic := range generics {
			result.Generics[i] = generic
		}
	}

	return result
}

func (p *parser_s) continue_type_expression() Expression {
	defer p.catch()

	exit := func() Expression {
		p.skip()
		result := p.current_expression()
		p.pop_expression()
		return result
	}

	p.skip()

	switch p.current_token().Kind {
	case identifier:
		if p.current_expression() != nil {
			p.backup()
			return p.current_expression()
		}
		ident := p.create_ident(p.current_token())
		p.advance()
		p.set_current_expression(*ident)
		return p.continue_type_expression()
	case dot:
		p.advance()

		if p.current_expression() == nil {
			// If there is no left hand side, it could just be a match self expression
			if p.is_match_context {
				p.set_current_expression(MatchSelfExpression{location: p.current_token().Location, Kind_: MatchSelfExpressionKind})
			} else {
				p.backup()
			}
		}

		if p.is_left_fun(p.current_expression()) {
			p.throw(fmt.Sprintf(errors.ErrorMessages["i_con"], "read a value off of a function"))
		}

		rhs_t := p.must_expect([]token_kind{identifier})
		rhs := p.create_ident(rhs_t)

		p.set_current_expression(MemberExpression{
			LeftHandSide:  p.current_expression(),
			RightHandSide: *rhs,
			Kind_:         MemberExpressionKind,
			location:      p.current_expression().Location(),
		})

		return p.continue_type_expression()
	case new_line, whitespace:
		return exit()
	default:
		return exit()
	}
}

func (p *parser_s) parse_type_expression() Expression {
	defer p.catch()

	if p.current_expression() != nil {
		p.push_expression()
	}

	return p.continue_type_expression()
}

func (p *parser_s) parse_value_type_pair() ValueTypePair {
	defer p.catch()

	var key Token
	start := p.must_expect([]token_kind{identifier, hidden_keyword})
	if start.Kind == hidden_keyword {
		p.skip()
		key = p.must_expect([]token_kind{identifier})
	} else {
		key = start
	}
	p.skip()
	typ := p.parse_type_literal()

	return ValueTypePair{
		Key:      *p.create_ident(key),
		Type:     typ,
		Hidden:   start.Kind == hidden_keyword,
		Location: start.Location,
	}
}

func (p *parser_s) parse_struct_literal() StructLiteral {
	defer p.catch()

	location := p.current_token().Location
	values := parse_seperated_list(p, p.parse_value_type_pair, semicolon, left_curly_bracks, right_curly_bracks, true, true)

	return StructLiteral{
		TypeKind_: StructLiteralKind,
		Values:    values,
		location:  location,
	}
}

func (p *parser_s) parse_constrained_type() ConstrainedType {
	defer p.catch()

	name := p.must_expect([]token_kind{identifier})
	is_spaced := p.might_expect([]token_kind{whitespace, new_line})
	var constraint *TypeLiteral

	if is_spaced != nil {
		p.skip()

		is_constrained := p.might_expect([]token_kind{identifier})

		if is_constrained != nil {
			p.backup()
			c := p.parse_type_literal()
			constraint = &c
			p.skip()
		}
	}

	return ConstrainedType{
		Name:       *p.create_ident(name),
		Constraint: constraint,
		Location:   name.Location,
	}
}

func (p *parser_s) parse_type_definition_statement() TypeDefinitionStatement {
	defer p.catch()

	var start Token
	is_hidden := p.might_expect([]token_kind{hidden_keyword})

	if is_hidden != nil {
		start = *is_hidden
		p.skip()
		p.must_expect([]token_kind{type_keyword})
	} else {
		start = p.must_expect([]token_kind{type_keyword})
	}

	p.must_expect([]token_kind{whitespace, new_line})
	p.skip()
	name := p.must_expect([]token_kind{identifier})

	result := TypeDefinitionStatement{
		Name:            *p.create_ident(name),
		Generics:        map[string]ConstrainedType{},
		Implementations: []TypeIdentifier{},
		Definition:      TypeIdentifier{},
		Hidden:          is_hidden != nil,
		Kind_:           TypeDefinitionStatementKind,
		location:        start.Location,
	}

	is_generic := p.might_expect([]token_kind{left_angle_bracks})

	if is_generic != nil {
		p.backup()
		result.Generics = generate_generics(p)
	}

	p.skip()
	is_implementing := p.might_expect([]token_kind{implements_keyword})

	if is_implementing != nil {
		p.skip()
		result.Implementations = parse_seperated_list(p, p.parse_type_identifier, comma, left_squre_bracks, right_squre_bracks, false, false)
	}

	result.Definition = p.parse_type_literal()
	p.skip()

	return result
}

func (p *parser_s) parse_typed_parameter() TypedParameter {
	defer p.catch()

	start := p.might_expect([]token_kind{variadic_marker})
	name := p.must_expect([]token_kind{identifier})
	p.must_expect([]token_kind{whitespace, new_line})
	p.skip()
	typ := p.parse_type_literal()

	var location errors.Location

	if start != nil {
		location = start.Location
	} else {
		location = name.Location
	}

	return TypedParameter{
		Name:     *p.create_ident(name),
		Type:     typ,
		Variadic: start != nil,
		Location: location,
	}
}

func (p *parser_s) parse_trait_definition_statement() TraitDefinitionStatement {
	defer p.catch()

	var start Token
	is_hidden := p.might_expect([]token_kind{hidden_keyword})

	if is_hidden != nil {
		start = *is_hidden
		p.skip()
		p.must_expect([]token_kind{trait_keyword})
	} else {
		p.skip()
		start = p.must_expect([]token_kind{trait_keyword})
	}

	p.must_expect([]token_kind{whitespace, new_line})
	p.skip()
	name := p.must_expect([]token_kind{identifier})

	result := TraitDefinitionStatement{
		Name:      *p.create_ident(name),
		Generics:  map[string]ConstrainedType{},
		Mimics:    []TypeIdentifier{},
		Hidden:    is_hidden != nil,
		Kind_:     TraitDefinitionStatementKind,
		TypeKind_: TraitTypeKind,
		location:  start.Location,
	}

	is_generic := p.might_expect([]token_kind{left_angle_bracks})

	if is_generic != nil {
		p.backup()
		result.Generics = generate_generics(p)
	}

	p.skip()

	is_mimicked := p.might_expect([]token_kind{mimics_keyword})

	if is_mimicked != nil {
		p.skip()
		result.Mimics = parse_seperated_list(p, p.parse_type_identifier, comma, left_squre_bracks, right_squre_bracks, false, false)
		p.skip()
	}

	p.skip()
	result.Definition = parse_seperated_list(p, p.parse_unbound_fun_signature, semicolon, left_curly_bracks, right_curly_bracks, true, true)
	p.skip()

	return result
}

func (p *parser_s) parse_unbound_fun_signature() UnboundFunctionSignature {
	defer p.catch()

	start := p.must_expect([]token_kind{fun_keyword})
	p.must_expect([]token_kind{whitespace, new_line})
	p.skip()
	name := p.must_expect([]token_kind{identifier})

	generics := map[string]ConstrainedType{}
	is_generic := p.might_expect([]token_kind{left_angle_bracks})

	if is_generic != nil {
		p.backup()
		generics = generate_generics(p)
	}

	p.skip()

	params := parse_seperated_list(p, p.parse_typed_parameter, comma, left_parens, right_parens, true, false)
	p.skip()

	var return_type *TypeLiteral
	p.skip()

	return_type_t := p.might_expect([]token_kind{identifier, cardinal_literal, fun_keyword, left_parens})

	if return_type_t != nil {
		p.backup()
		return_type_p := p.parse_type_literal()
		return_type = &return_type_p
	}

	p.skip()

	return UnboundFunctionSignature{
		TypeKind_:  FunTypeKind,
		Name:       *p.create_ident(name),
		Parameters: params,
		Generics:   generics,
		ReturnType: return_type,
		location:   start.Location,
	}
}

func (p *parser_s) parse_anonymous_fun_signature() AnonymousFunctionSignature {
	defer p.catch()

	start := p.must_expect([]token_kind{fun_keyword})

	generics := map[string]ConstrainedType{}
	is_generic := p.might_expect([]token_kind{left_angle_bracks})

	if is_generic != nil {
		p.backup()
		generics = generate_generics(p)
	}

	p.skip()

	params := parse_seperated_list(p, p.parse_typed_parameter, comma, left_parens, right_parens, true, false)
	p.skip()

	var return_type *TypeLiteral
	p.skip()
	return_type_t := p.might_expect([]token_kind{identifier, cardinal_literal, fun_keyword, left_parens})

	if return_type_t != nil {
		p.backup()
		return_type_p := p.parse_type_literal()
		return_type = &return_type_p
	}

	p.skip()

	return AnonymousFunctionSignature{
		TypeKind_:  FunTypeKind,
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
	p.skip()
	p.must_expect([]token_kind{for_keyword})
	p.must_expect([]token_kind{whitespace, new_line})
	p.skip()
	for_typ := p.parse_simple_type_identifier()
	p.must_expect([]token_kind{whitespace, new_line})
	p.skip()
	name := p.must_expect([]token_kind{identifier})

	generics := map[string]ConstrainedType{}
	is_generic := p.might_expect([]token_kind{left_angle_bracks})

	if is_generic != nil {
		p.backup()
		generics = generate_generics(p)
	}

	p.skip()

	params := parse_seperated_list(p, p.parse_typed_parameter, comma, left_parens, right_parens, true, false)
	p.skip()

	var return_type *TypeLiteral
	p.skip()
	return_type_t := p.might_expect([]token_kind{identifier, cardinal_literal, fun_keyword, left_parens})

	if return_type_t != nil {
		p.backup()
		return_type_p := p.parse_type_literal()
		return_type = &return_type_p
	}

	p.skip()

	return BoundFunctionSignature{
		Name:       *p.create_ident(name),
		For:        for_typ,
		Parameters: params,
		Generics:   generics,
		ReturnType: return_type,
		location:   start.Location,
	}
}

func (p *parser_s) parse_block() StatementList {
	defer p.catch()

	p.must_expect([]token_kind{left_curly_bracks})
	body := p.parse_inline_level_statements()
	p.skip()
	p.must_expect([]token_kind{right_curly_bracks})
	p.skip()

	return body
}

func (p *parser_s) parse_fun_definition_statement() FunDefinitionStatement {
	defer p.catch()

	is_hidden := p.might_expect([]token_kind{hidden_keyword})

	if is_hidden != nil {
		p.skip()
		p.must_expect([]token_kind{fun_keyword})
	} else {
		p.must_expect([]token_kind{fun_keyword})
	}

	p.must_expect([]token_kind{whitespace, new_line})
	spaces := p.skip()

	next := p.must_expect([]token_kind{identifier, for_keyword})
	p.backup_by(spaces + 1 /*identifier of for keyword */ + 1 /* single whitespace that is expected */ + 1 /* the fun keyword */)

	var definition FunDefinitionStatement

	switch next.Kind {
	case identifier:
		signature := p.parse_unbound_fun_signature()

		definition = &UnboundFunDefinitionStatement{
			Signature: signature,
			Body:      []Statement{},
			Hidden:    is_hidden != nil,
			Kind_:     UnboundFunDefinitionStatementKind,
			location:  signature.location,
		}
	case for_keyword:
		signature := p.parse_bound_fun_signature()

		definition = &BoundFunDefinitionStatement{
			Signature: signature,
			Body:      []Statement{},
			Hidden:    is_hidden != nil,
			Kind_:     BoundFunDefinitionStatementKind,
			location:  signature.location,
		}
	default:
		// there is something wrong with the token, just throw
		p.must_expect([]token_kind{})
	}

	p.skip()

	// if this is a bound fun definition, this keyword as an expression is allowed
	if reflect.TypeOf(definition) == reflect.TypeOf(&BoundFunDefinitionStatement{}) {
		p.is_this_context = true
	}
	p.set_context(function_context)
	definition.set_body(p.parse_block())
	p.reset_context()
	p.is_this_context = false

	return definition
}

func (p *parser_s) parse_return_statement() ReturnStatement {
	defer p.catch()

	start := p.must_expect([]token_kind{return_keyword})

	p.must_expect([]token_kind{whitespace, new_line})
	p.skip()
	expression := p.parse_expression()

	p.skip()

	p.parse_inline_level_statements()

	return ReturnStatement{
		Kind_: ReturnStatementKind,

		Value:    &expression,
		location: start.Location,
	}
}

func (p *parser_s) parse_yield_statement() YieldStatement {
	defer p.catch()

	start := p.must_expect([]token_kind{yield_keyword})

	p.must_expect([]token_kind{whitespace, new_line})
	p.skip()
	expression := p.parse_expression()

	p.skip()

	p.parse_inline_level_statements()

	return YieldStatement{
		Kind_: YieldStatementKind,

		Value:    &expression,
		location: start.Location,
	}
}

func (p *parser_s) parse_flow_control_statement() Statement {
	defer p.catch()

	token := p.must_expect([]token_kind{continue_keyword, break_keyword})

	if token.Kind == continue_keyword {
		return ContinueStatement{location: p.current_token().Location, Kind_: ContinueStatementKind}
	} else {
		return BreakStatement{location: p.current_token().Location, Kind_: BreakStatementKind}
	}
}

func (p *parser_s) is_left_fun(expression Expression) bool {
	if expression == nil {
		return false
	}

	switch reflect.TypeOf(expression) {
	case reflect.TypeOf(AnonymousFunExpression{}):
		return true
	case reflect.TypeOf(GroupExpression{}):
		return p.is_left_fun(expression.(GroupExpression).Expression)
	}

	return false
}

func (p *parser_s) is_left_call(expression Expression) bool {
	if expression == nil {
		return false
	}

	switch reflect.TypeOf(expression) {
	case reflect.TypeOf(CallExpression{}):
		return true
	case reflect.TypeOf(GroupExpression{}):
		return p.is_left_call(expression.(GroupExpression).Expression)
	}

	return false
}

func (p *parser_s) is_left_callable(expression Expression) bool {
	if expression == nil {
		return false
	}

	switch reflect.TypeOf(expression) {
	case reflect.TypeOf(IdentifierExpression{}), reflect.TypeOf(MemberExpression{}):
		return true
	case reflect.TypeOf(GroupExpression{}):
		return p.is_left_callable(expression.(GroupExpression).Expression)
	}

	return false
}

func (p *parser_s) is_left_unary_arithmetic(expression Expression) bool {
	if expression == nil {
		return false
	}

	switch reflect.TypeOf(expression) {
	case reflect.TypeOf(ArithmeticUnaryExpression{}):
		return true
	case reflect.TypeOf(GroupExpression{}):
		return p.is_left_unary_arithmetic(expression.(GroupExpression).Expression)
	}

	return false
}

/*
Recursively call the continue expression and save it to a stack.
Next expression will look at its left to determine what to be and update the stack.
When the expression ends, the result will be returned and stack will be popped.
There is a stack because there is a possiblity that you may parse an expression inside another expression.
eg: data.count
identifier 'data' is found, current expression is Identifier(data)
. token is found, will try to parse member expression, advance...
Member expression parser will find 'count' token and look at the last expression.
It will then create a member expression like MemberExpression(data, count),
set this as the current expression and will try to move forward by calling continue_expression.
*/
func (p *parser_s) continue_expression() Expression {
	defer p.catch()

	exit := func() Expression {
		p.skip()
		result := p.current_expression()
		p.pop_expression()
		return result
	}

	switch p.current_token().Kind {
	case rune_literal, string_literal, bool_literal, number_literal:
		p.set_current_expression(p.parse_literal_expression())

		is_ended := p.might_expect([]token_kind{new_line})

		if is_ended != nil {
			p.skip()
			result := p.current_expression()
			p.pop_expression()
			return result
		}
		return p.continue_expression()
	case left_curly_bracks:
		p.set_current_expression(p.parse_literal_expression())

		is_ended := p.might_expect([]token_kind{new_line, whitespace})

		if is_ended != nil {
			return exit()
		}
		return p.continue_expression()
	case identifier:
		if p.current_expression() != nil {
			p.backup()
			return p.current_expression()
		}
		ident := p.create_ident(p.current_token())
		p.advance()
		p.set_current_expression(*ident)
		return p.continue_expression()
	case this_keyword:
		// don't let the keyword as an expression outside the bound fun context
		if !p.is_this_context {
			p.backup_by(2)
			defer p.advance()
			p.must_expect([]token_kind{})
		}
		p.set_current_expression(ThisExpression{location: p.current_token().Location, Kind_: ThisExpressionKind})
		p.advance()
		return p.continue_expression()
	case corout_keyword:
		return p.parse_corout_expression()
	case gen_keyword:
		return p.parse_gen_expression()
	case dot:
		p.advance()

		switch p.current_token().Kind {
		case left_parens:
			return p.parse_type_cast_expression()

		case identifier:
			return p.parse_member_expression()
		default:
			// don't let the keyword as an expression outside the match expression
			if !p.is_match_context {
				p.backup_by(2)
				defer p.advance()
				p.must_expect([]token_kind{})
			}
			p.backup()
			p.set_current_expression(MatchSelfExpression{location: p.current_token().Location, Kind_: MatchSelfExpressionKind})
			p.advance()
			return p.continue_expression()
		}
	case left_parens:
		if p.current_expression() == nil {
			return p.parse_group_expression()
		} else {
			return p.parse_call_expression()
		}
	case binary_operator:
		return p.parse_binary_expression()
	case left_angle_bracks, right_angle_bracks, comparison_operator:
		return p.parse_comparison_expression()
	case plus, minus, star, forward_slash, percent:
		return p.parse_arithmetic_expression()
	case increment, decrement:
		expression := p.parse_arithmetic_unary_expression()
		return expression
	case instanceof_keyword:
		return p.parse_instanceof_expression()
	case match_keyword:
		return p.parse_match_expression()
	case exclamation:
		return p.parse_not_expression()
	case fun_keyword:
		return p.parse_anonymous_fun_expression(function_context)
	case giveup_keyword:
		if p.current_expression() != nil {
			p.must_expect([]token_kind{})
		}
		p.advance()
		p.set_current_expression(GiveupExpression{location: p.current_token().Location, Kind_: GiveupExpressionKind})

		return exit()
	case caret:
		if p.current_expression() != nil {
			p.must_expect([]token_kind{})
		}
		p.advance()
		p.set_current_expression(CaretExpression{location: p.current_token().Location, Kind_: CaretExpressionKind})
		return p.continue_expression()
	case or_keyword:
		return p.parse_or_expression()
	case cardinal_literal:
		defer p.advance()
		p.must_expect([]token_kind{})
		return nil
	case left_squre_bracks:
		if p.current_expression() == nil {
			p.set_current_expression(p.parse_literal_expression())
			is_ended := p.might_expect([]token_kind{new_line, whitespace})

			if is_ended != nil {
				return exit()
			}

			return p.continue_expression()
		} else {
			return p.parse_index_expression()
		}
	case whitespace:
		p.skip()

		if p.is_left_unary_arithmetic(p.current_expression()) {
			result := p.current_expression()
			p.pop_expression()

			return result
		} else {
			return p.continue_expression()
		}
	case new_line:
		return exit()
	default:
		return exit()
	}
}

func (p *parser_s) parse_expression() Expression {
	defer p.catch()

	if p.current_expression() != nil {
		p.push_expression()
	}

	return p.continue_expression()
}

func (p *parser_s) parse_anonymous_fun_expression(context []token_kind) Expression {
	if p.current_expression() != nil {
		p.must_expect([]token_kind{})
	}

	signature := p.parse_anonymous_fun_signature()
	result := AnonymousFunExpression{
		Signature: signature,
		Kind_:     AnonymousFunExpressionKind,
		location:  signature.location,
	}
	p.skip()
	p.set_context(context)
	result.Body = p.parse_block()
	p.reset_context()

	p.set_current_expression(result)

	return p.continue_expression()
}

func (p *parser_s) parse_corout_expression() Expression {
	if p.current_expression() != nil {
		p.must_expect([]token_kind{})
	}

	start := p.must_expect([]token_kind{corout_keyword})
	p.skip()
	fun := p.parse_anonymous_fun_expression(function_context)

	p.set_current_expression(CoroutFunExpression{
		Kind_: CoroutFunExpressionKind,

		Fun:      fun.(AnonymousFunExpression),
		location: start.Location,
	})

	return p.continue_expression()
}

func (p *parser_s) parse_gen_expression() Expression {
	if p.current_expression() != nil {
		p.must_expect([]token_kind{})
	}

	start := p.must_expect([]token_kind{gen_keyword})
	p.skip()
	fun := p.parse_anonymous_fun_expression(generator_function_context)

	p.set_current_expression(GenFunExpression{
		Kind_: GenFunExpressionKind,

		Fun:      fun.(AnonymousFunExpression),
		location: start.Location,
	})

	return p.continue_expression()
}

func (p *parser_s) parse_group_expression() Expression {
	defer p.catch()

	start := p.must_expect([]token_kind{left_parens})
	p.skip()
	expression := p.parse_expression()

	if expression == nil {
		p.must_expect([]token_kind{})
	}

	p.skip()
	p.must_expect([]token_kind{right_parens})

	p.set_current_expression(GroupExpression{Expression: expression, location: start.Location, Kind_: GroupExpressionKind})

	return p.continue_expression()
}

func (p *parser_s) parse_call_expression() Expression {
	defer p.catch()

	if !p.is_left_callable(p.current_expression()) {
		p.throw(errors.ErrorMessages["uc_con"])
	}

	args := parse_seperated_list(p, p.parse_expression, comma, left_parens, right_parens, true, false)

	p.set_current_expression(CallExpression{
		Callee:    p.current_expression(),
		Arguments: args,
		Kind_:     CallExpressionKind,
		location:  p.current_token().Location,
	})

	return p.continue_expression()
}

func (p *parser_s) parse_member_expression() Expression {
	defer p.catch()

	if p.current_expression() == nil {
		// If there is no left hand side, it could just be a match self expression
		if p.is_match_context {
			p.set_current_expression(MatchSelfExpression{location: p.current_token().Location, Kind_: MatchSelfExpressionKind})
		} else {
			p.backup()
		}
	}

	if p.is_left_fun(p.current_expression()) {
		p.throw(fmt.Sprintf(errors.ErrorMessages["i_con"], "read a value off of a function"))
	}

	rhs_t := p.must_expect([]token_kind{identifier})
	rhs := p.create_ident(rhs_t)

	p.set_current_expression(MemberExpression{
		LeftHandSide:  p.current_expression(),
		RightHandSide: *rhs,
		Kind_:         MemberExpressionKind,
		location:      p.current_expression().Location(),
	})

	return p.continue_expression()
}

func (p *parser_s) parse_or_expression() Expression {
	defer p.catch()

	if p.current_expression() == nil {
		// if current expression is non-existent, the next lines will expect an identifier and it will throw
		p.backup()
	}

	if !p.is_left_call(p.current_expression()) {
		// if current expression is not a function call throw
		p.throw(fmt.Sprintf(errors.ErrorMessages["i_con"], "recover from a non-function expression"))
	}

	p.skip()
	p.must_expect([]token_kind{or_keyword})
	p.must_expect([]token_kind{whitespace, new_line})
	p.skip()

	rhs := p.parse_expression()

	p.set_current_expression(OrExpression{
		LeftHandSide:  p.current_expression(),
		RightHandSide: rhs,
		Kind_:         OrExpressionKind,
		location:      p.current_expression().Location(),
	})

	return p.continue_expression()
}

func (p *parser_s) parse_index_expression() Expression {
	defer p.catch()

	if p.current_expression() == nil {
		p.must_expect([]token_kind{})
	}

	if p.is_left_fun(p.current_expression()) {
		p.throw(fmt.Sprintf(errors.ErrorMessages["i_con"], "index a function"))
	}

	p.must_expect([]token_kind{left_squre_bracks})
	index := p.parse_expression()
	if index == nil {
		p.must_expect([]token_kind{})
	}
	p.must_expect([]token_kind{right_squre_bracks})

	p.set_current_expression(IndexExpression{
		Host:     p.current_expression(),
		Index:    index,
		Kind_:    IndexExpressionKind,
		location: p.current_expression().Location(),
	})

	return p.continue_expression()
}

func (p *parser_s) parse_arithmetic_expression() Expression {
	defer p.catch()

	current := p.current_expression()
	if current == nil {
		p.must_expect([]token_kind{})
	}

	operator := p.must_expect([]token_kind{plus, minus, star, forward_slash, percent})
	p.skip()

	rhs := p.parse_expression()

	default_expression := ArithmeticExpression{
		LeftHandSide:  current,
		RightHandSide: rhs,
		Operator: OperatorToken{
			Literal:  operator.Literal,
			location: operator.Location,
		},
		Kind_:    ArithmeticExpressionKind,
		location: current.Location(),
	}

	if operator.Kind == plus || operator.Kind == minus {
		p.set_current_expression(default_expression)
		return p.continue_expression()
	}

	if reflect.TypeOf(rhs) != reflect.TypeOf(ArithmeticExpression{}) {
		p.set_current_expression(default_expression)
		return p.continue_expression()
	}

	rhs_e := rhs.(ArithmeticExpression)

	p.set_current_expression(ArithmeticExpression{
		LeftHandSide: ArithmeticExpression{
			LeftHandSide:  current,
			RightHandSide: rhs_e.LeftHandSide,
			Operator: OperatorToken{
				Literal:  operator.Literal,
				location: operator.Location,
			},
			Kind_:    ArithmeticExpressionKind,
			location: current.Location(),
		},
		RightHandSide: rhs_e.RightHandSide,
		Operator:      rhs_e.Operator,
		Kind_:         ArithmeticExpressionKind,
		location:      current.Location(),
	})

	return p.continue_expression()
}

func (p *parser_s) parse_comparison_expression() Expression {
	defer p.catch()

	start_offset := p.offset
	current := p.current_expression()
	if current == nil {
		p.must_expect([]token_kind{})
	}

	operator := p.must_expect([]token_kind{left_angle_bracks, right_angle_bracks, comparison_operator})
	skipped := p.skip()

	rhs := p.parse_expression()

	if rhs == nil {
		p.backup_by(skipped)
		p.must_expect([]token_kind{})
	}

	/* normally an expression like var0 < var1 > var2 is not acceptable
	but there is a chance that it is a typed instance literal like Type<Generic>{}
	check if right hand side is exaclty a map literal and the operator is >
	if so, turn the map literal to an instance literal and return */

	/* evidently, this has a bug in that it expects no whitespace between type and its generic
	so Type< Generic> or Type<Generic0, Generic1> would be unacceptable */
	if rhs.Kind() == ComparisonExpressionKind {
		right_comparison := rhs.(ComparisonExpression)
		if right_comparison.Operator.Literal == ">" && right_comparison.RightHandSide.Kind() == MapLiteralExpressionKind {
			left_offset := p.offset
			p.offset = start_offset

			for p.current_token().Kind != whitespace {
				p.backup()
			}

			p.skip()
			typ := p.parse_type_literal()

			if typ == nil {
				p.must_expect([]token_kind{})
			}
			if typ.TypeKind() != TypeIdentifierKind {
				p.must_expect([]token_kind{})
			}

			value := right_comparison.RightHandSide.(MapLiteralExpression).Value
			p.offset = left_offset
			p.skip()

			return InstanceLiteralExpression{
				Kind_: InstanceLiteralExpressionKind,

				Type:     typ.(TypeIdentifier),
				Value:    value,
				location: p.current_token().Location,
			}
		}

		current_offset := p.current_token().Location.Offset
		skipped_offset := right_comparison.Operator.location.Offset
		p.backup_by(current_offset - skipped_offset)
		p.must_expect([]token_kind{})
	}

	p.set_current_expression(ComparisonExpression{
		Kind_: ComparisonExpressionKind,

		LeftHandSide:  current,
		RightHandSide: rhs,
		Operator: OperatorToken{
			Literal:  operator.Literal,
			location: operator.Location,
		},
		location: current.Location(),
	})

	return p.continue_expression()
}

func (p *parser_s) parse_binary_expression() Expression {
	defer p.catch()

	current := p.current_expression()
	if current == nil {
		p.must_expect([]token_kind{})
	}

	operator := p.must_expect([]token_kind{binary_operator})
	skipped := p.skip()

	rhs := p.parse_expression()

	if rhs == nil {
		p.backup_by(skipped)
		p.must_expect([]token_kind{})
	}

	p.set_current_expression(BinaryExpression{
		LeftHandSide:  current,
		RightHandSide: rhs,
		Operator: OperatorToken{
			Literal:  operator.Literal,
			location: operator.Location,
		},
		Kind_:    BinaryExpressionKind,
		location: current.Location(),
	})

	return p.continue_expression()
}

func (p *parser_s) parse_not_expression() Expression {
	defer p.catch()

	current := p.current_expression()
	if current != nil {
		p.must_expect([]token_kind{})
	}

	start := p.must_expect([]token_kind{exclamation})
	skipped := p.skip()

	expr := p.parse_expression()

	if expr == nil {
		p.backup_by(skipped)
		p.must_expect([]token_kind{})
	}

	p.set_current_expression(NotExpression{
		Expression: expr,
		Kind_:      NotExpressionKind,
		location:   start.Location,
	})

	return p.continue_expression()
}

func (p *parser_s) parse_instanceof_expression() Expression {
	defer p.catch()

	current := p.current_expression()
	if current == nil {
		p.must_expect([]token_kind{})
	}

	p.must_expect([]token_kind{instanceof_keyword})
	p.skip()

	rhs := p.parse_type_literal()

	p.set_current_expression(InstanceofExpression{
		LeftHandSide:  current,
		RightHandSide: rhs,
		Kind_:         InstanceofExpressionKind,
		location:      current.Location(),
	})

	return p.continue_expression()
}

func (p *parser_s) parse_arithmetic_unary_expression() ArithmeticUnaryExpression {
	defer p.catch()

	var kind ArithmeticUnaryKind
	var expression Expression

	if p.current_expression() == nil {
		p.must_expect([]token_kind{})
	}

	expression = p.current_expression()
	token := p.must_expect([]token_kind{increment, decrement})

	if token.Kind == increment {
		kind = IncrementKind
	} else {
		kind = DecrementKind
	}

	if reflect.TypeOf(expression) == reflect.TypeOf(CallExpression{}) {
		p.backup()
		defer p.advance()
		p.throw(fmt.Sprintf(errors.ErrorMessages["i_con"], "do this operation with a function call"))
	}

	p.skip()

	return ArithmeticUnaryExpression{
		Expression: expression,
		Operation:  kind,
		Pre:        false,
		Kind_:      ArithmeticUnaryExpressionKind,
		location:   expression.Location(),
	}
}

func (p *parser_s) parse_type_cast_expression() Expression {
	defer p.catch()

	if p.current_expression() == nil {
		p.must_expect([]token_kind{})
	}

	p.must_expect([]token_kind{left_parens})
	typ := p.parse_type_identifier()
	p.must_expect([]token_kind{right_parens})

	p.set_current_expression(TypeCastExpression{
		Value:    p.current_expression(),
		Type:     typ,
		Kind_:    TypeCastExpressionKind,
		location: p.continue_expression().Location(),
	})

	return p.continue_expression()
}

func (p *parser_s) parse_key_value_entry() KeyValueEntry {
	key := p.must_expect([]token_kind{identifier})
	p.skip()
	p.must_expect([]token_kind{colon})
	p.skip()
	value := p.parse_expression()

	return KeyValueEntry{
		Key:   p.create_ident(key),
		Value: value,
	}
}

func (p *parser_s) parse_literal_expression() LiteralExpression {
	defer p.catch()

	current := p.current_token()
	var result LiteralExpression

	switch current.Kind {
	case string_literal:
		result = StringLiteralExpression{
			Value:    current.Literal,
			Kind_:    StringLiteralExpressionKind,
			location: current.Location,
		}
		p.advance()
	case rune_literal:
		result = RuneLiteralExpression{
			Value:    rune(current.Literal[0]),
			Kind_:    RuneLiteralExpressionKind,
			location: current.Location,
		}
		p.advance()
	case bool_literal:
		result = BoolLiteralExpression{
			Value:    current.Literal == "true",
			Kind_:    BoolLiteralExpressionKind,
			location: current.Location,
		}
		p.advance()
	case number_literal:
		result = NumberLiteralExpression{
			Value:    create_number_literal(*p, current.Literal),
			Kind_:    NumberLiteralExpressionKind,
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
			Kind_:    ListLiteralExpressionKind,
			location: current.Location,
		}

		p.backup()
		if p.current_token().Kind != whitespace && p.current_token().Kind != new_line {
			p.advance()
		}
	case left_curly_bracks:
		current := p.current_expression()

		if current == nil {
			entries := parse_seperated_list(p, p.parse_key_value_entry, comma, left_curly_bracks, right_curly_bracks, true, true)

			result = MapLiteralExpression{
				Kind_: MapLiteralExpressionKind,

				Value:    entries,
				location: p.current_token().Location,
			}
		} else {
			if reflect.TypeOf(current) != reflect.TypeOf(IdentifierExpression{}) && reflect.TypeOf(current) != reflect.TypeOf(MemberExpression{}) {
				p.must_expect([]token_kind{})
			}

			typ := TypeIdentifier{
				Name:     p.current_expression(),
				Generics: map[int]TypeLiteral{},
			}
			entries := parse_seperated_list(p, p.parse_key_value_entry, comma, left_curly_bracks, right_curly_bracks, true, true)

			result = InstanceLiteralExpression{
				Type:     typ,
				Value:    entries,
				Kind_:    InstanceLiteralExpressionKind,
				location: typ.Name.Location(),
			}
		}

		p.backup()
		if p.current_token().Kind != whitespace && p.current_token().Kind != new_line {
			p.advance()
		}
	default:
		result = StringLiteralExpression{
			Value:    current.Literal,
			Kind_:    StringLiteralExpressionKind,
			location: current.Location,
		}
		p.advance()
	}

	return result
}

func (p *parser_s) parse_predicate_block() PredicateBlock {
	defer p.catch()

	p.must_expect([]token_kind{left_parens})
	predicate := p.parse_expression()
	if predicate == nil {
		return PredicateBlock{}
	}

	p.must_expect([]token_kind{right_parens})

	p.skip()
	p.set_context(predicate_body_context)
	body := p.parse_block()
	p.reset_context()

	return PredicateBlock{
		Predicate: predicate,
		Body:      body,
	}
}

func (p *parser_s) parse_match_expression() Expression {
	defer p.catch()

	start := p.must_expect([]token_kind{match_keyword})
	p.skip()
	p.must_expect([]token_kind{left_parens})
	against := p.parse_expression()

	if against == nil {
		p.must_expect([]token_kind{})
	}

	p.must_expect([]token_kind{right_parens})
	p.skip()
	p.must_expect([]token_kind{left_curly_bracks})
	p.skip()
	// in a match statement, . (Dot) keyword as an expression is allowed
	p.is_match_context = true

	blocks := []PredicateBlock{}
	base_block := StatementList{}

	current := p.current_token()
	for current.Kind != right_curly_bracks && current.Kind != base_keyword {
		p.skip()

		if current.Kind == right_curly_bracks || current.Kind == base_keyword {
			break
		}

		predicate := p.parse_predicate_block()
		if predicate.Predicate == nil {
			break
		}
		blocks = append(blocks, predicate)
		current = p.current_token()
	}

	next := p.might_expect([]token_kind{base_keyword})

	if next != nil {
		p.skip()
		base_block = p.parse_block()
		p.skip()
	}

	p.skip()
	p.must_expect([]token_kind{right_curly_bracks})

	p.set_current_expression(MatchExpression{
		Against:   against,
		Blocks:    blocks,
		BaseBlock: base_block,
		Kind_:     MatchExpressionKind,
		location:  start.Location,
	})
	p.is_match_context = false

	return p.continue_expression()
}

func Parse(input []byte, filepath string) (Ast, errors.Error) {
	filename := path.Base(filepath)

	parser := parser_s{
		input: input,
		ast: Ast{
			FileName: filename,
			FilePath: filepath,
			Uses:     []UseStatement{},
			Comments: []Comment{},
		},
		body_context: []token_kind{},
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
