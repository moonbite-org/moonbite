package cmd

import (
	"fmt"
	"slices"

	"github.com/moonbite-org/moonbite/common"
	errors "github.com/moonbite-org/moonbite/error"
	parser "github.com/moonbite-org/moonbite/parser/cmd"
)

func (c *PackageCompiler) enter_scope() {
	scoped_table := NewScopedSymbolTable(c.SymbolTable)
	c.SymbolTable = scoped_table
}

func (c *PackageCompiler) leave_scope() {
	c.SymbolTable = c.SymbolTable.Outer
}

func (c *PackageCompiler) resolve_assignee(assignee parser.Expression) (*Symbol, errors.Error) {
	acceptable_assignees := []parser.ExpressionKind{parser.IdentifierExpressionKind, parser.MemberExpressionKind, parser.IndexExpressionKind}

	if !slices.Contains(acceptable_assignees, assignee.Kind()) {
		return nil, errors.CreateCompileError(fmt.Sprintf(errors.ErrorMessages["i_con"], "assign to this value, left hand side includes unassignable path"), assignee.Location())
	}

	switch assignee.Kind() {
	case parser.IdentifierExpressionKind:
		identifier := assignee.(parser.IdentifierExpression)
		symbol := c.SymbolTable.Resolve(identifier.Value)

		if symbol == nil {
			return nil, errors.CreateCompileError(fmt.Sprintf("variable '%s' is not defined", identifier.Value), identifier.Location())
		}

		return symbol, errors.EmptyError
	case parser.MemberExpressionKind:
		return c.resolve_assignee(assignee.(parser.MemberExpression).LeftHandSide)
	case parser.IndexExpressionKind:
		return c.resolve_assignee(assignee.(parser.IndexExpression).Host)
	default:
		return nil, errors.CreateCompileError("could not resolve assignee", assignee.Location())
	}
}

func (c *PackageCompiler) resolve_path(expression parser.Expression) (common.InstructionSet, int, errors.Error) {
	switch expression.Kind() {
	case parser.IdentifierExpressionKind:
		value := expression.(parser.IdentifierExpression).Value
		instructions, err := c.compile_literal_expression(parser.StringLiteralExpression{Value: value})
		if err.Exists {
			return common.InstructionSet{}, 0, err
		}
		return instructions, 1, err
	case parser.IndexExpressionKind:
		result := common.InstructionSet{}
		index, err := c.compile_expression(expression.(parser.IndexExpression).Index, false)
		if err.Exists {
			return common.InstructionSet{}, 0, err
		}
		host, size, err := c.resolve_path(expression.(parser.IndexExpression).Host)
		if err.Exists {
			return common.InstructionSet{}, 0, err
		}
		if err.Exists {
			return common.InstructionSet{}, 0, err
		}
		result = append(result, host...)
		result = append(result, index...)
		return result, size + 1, errors.EmptyError
	case parser.MemberExpressionKind:
		result := common.InstructionSet{}
		rhs, rhs_size, err := c.resolve_path(expression.(parser.MemberExpression).LeftHandSide)
		if err.Exists {
			return common.InstructionSet{}, 0, err
		}
		lhs, lhs_size, err := c.resolve_path(expression.(parser.MemberExpression).RightHandSide)
		if err.Exists {
			return common.InstructionSet{}, 0, err
		}
		if err.Exists {
			return common.InstructionSet{}, 0, err
		}
		result = append(result, rhs...)
		result = append(result, lhs...)
		return result, lhs_size + rhs_size, errors.EmptyError
	default:
		return common.InstructionSet{}, 0, errors.CreateCompileError(fmt.Sprintf(errors.ErrorMessages["i_con"], "assign to this value, left hand side includes unassignable path"), expression.Location())
	}
}

func (c *PackageCompiler) compile_statement(statement parser.Statement) (common.InstructionSet, errors.Error) {
	switch statement.Kind() {
	case parser.ExpressionStatementKind:
		return c.compile_expression(statement.(parser.ExpressionStatement).Expression, true)
	case parser.DeclarationStatementKind:
		return c.compile_declaration_statement(statement.(parser.DeclarationStatement))
	case parser.AssignmentStatementKind:
		return c.compile_assignment_statement(statement.(parser.AssignmentStatement))
	case parser.TypeDefinitionStatementKind:
		return c.compile_type_definition_statement(statement.(parser.TypeDefinitionStatement))
	case parser.UnboundFunDefinitionStatementKind:
		return c.compile_unbound_fun_definition_statement(*statement.(*parser.UnboundFunDefinitionStatement))
	case parser.BoundFunDefinitionStatementKind:
		return c.compile_bound_fun_definition_statement(*statement.(*parser.BoundFunDefinitionStatement))
	case parser.ReturnStatementKind:
		return c.compile_return_statement(statement.(parser.ReturnStatement))
	case parser.BreakStatementKind:
		return c.compile_break_statement(statement.(parser.BreakStatement))
	case parser.ContinueStatementKind:
		return c.compile_continue_statement(statement.(parser.ContinueStatement))
	case parser.YieldStatementKind:
		return c.compile_yield_statement(statement.(parser.YieldStatement))
	case parser.IfStatementKind:
		return c.compile_if_statement(statement.(parser.IfStatement))
	case parser.LoopStatementKind:
		return c.compile_loop_statement(statement.(parser.LoopStatement))
	default:
		result := common.InstructionSet{}
		result = append(result, common.NewInstruction(common.OpNoop))
		return result, errors.EmptyError
	}
}

func (c *PackageCompiler) compile_declaration_statement(statement parser.DeclarationStatement) (common.InstructionSet, errors.Error) {
	result := common.InstructionSet{}

	if statement.Value != nil {
		value, err := c.compile_expression(*statement.Value, false)
		if err.Exists {
			return result, err
		}
		result = append(result, value...)
	} else {
		value_literal := c.Typechecker.GetDefault(*statement.Type)
		index := c.ConstantPool.Add(value_literal)
		result = append(result, common.NewInstruction(common.OpConstant, index))
	}

	symbol, err := c.SymbolTable.Define(statement.Name.Value, statement.VarKind, statement.Hidden)
	if err != nil {
		return result, errors.CreateCompileError(err.Error(), statement.Name.Location())
	}

	if symbol.Scope == GlobalScope {
		result = append(result, common.NewInstruction(common.OpSet, symbol.Index))
	} else {
		result = append(result, common.NewInstruction(common.OpSetLocal, symbol.Index))
	}

	return result, errors.EmptyError
}

func (c *PackageCompiler) compile_assignment_statement(statement parser.AssignmentStatement) (common.InstructionSet, errors.Error) {
	result := common.InstructionSet{}

	symbol, err := c.resolve_assignee(statement.LeftHandSide)

	if err.Exists {
		return result, err
	}

	if symbol.Kind == parser.ConstantKind {
		return result, errors.CreateCompileError(fmt.Sprintf("cannot assign to constant variable '%s'", symbol.Name), statement.Location())
	}

	left, err := c.compile_expression(statement.LeftHandSide, false)
	if err.Exists {
		return result, err
	}

	right, err := c.compile_expression(statement.RightHandSide, false)
	if err.Exists {
		return result, err
	}

	var left_getter common.Op

	if symbol.Scope == GlobalScope {
		left_getter = common.OpGet
	} else {
		left_getter = common.OpGetLocal
	}

	switch statement.LeftHandSide.Kind() {
	case parser.IdentifierExpressionKind:
		switch statement.Operator.Literal {
		case "=":
			result = append(result, right...)
		case "+=":
			result = append(result, common.NewInstruction(left_getter, symbol.Index))
			result = append(result, right...)
			result = append(result, common.NewInstruction(common.OpAdd))
		case "-=":
			result = append(result, common.NewInstruction(left_getter, symbol.Index))
			result = append(result, right...)
			result = append(result, common.NewInstruction(common.OpSub))
		case "*=":
			result = append(result, common.NewInstruction(left_getter, symbol.Index))
			result = append(result, right...)
			result = append(result, common.NewInstruction(common.OpMul))
		case "/=":
			result = append(result, common.NewInstruction(left_getter, symbol.Index))
			result = append(result, right...)
			result = append(result, common.NewInstruction(common.OpDiv))
		case "%=":
			result = append(result, common.NewInstruction(left_getter, symbol.Index))
			result = append(result, right...)
			result = append(result, common.NewInstruction(common.OpMod))
		}

		if symbol.Scope == GlobalScope {
			result = append(result, common.NewInstruction(common.OpAssign, symbol.Index))
		} else {
			result = append(result, common.NewInstruction(common.OpAssignLocal, symbol.Index))
		}
	case parser.MemberExpressionKind, parser.IndexExpressionKind:
		switch statement.Operator.Literal {
		case "=":
			result = append(result, right...)
		case "+=":
			result = append(result, left...)
			result = append(result, right...)
			result = append(result, common.NewInstruction(common.OpAdd))
		case "-=":
			result = append(result, left...)
			result = append(result, right...)
			result = append(result, common.NewInstruction(common.OpSub))
		case "*=":
			result = append(result, left...)
			result = append(result, right...)
			result = append(result, common.NewInstruction(common.OpMul))
		case "/=":
			result = append(result, left...)
			result = append(result, right...)
			result = append(result, common.NewInstruction(common.OpDiv))
		case "%=":
			result = append(result, left...)
			result = append(result, right...)
			result = append(result, common.NewInstruction(common.OpMod))
		}

		path, size, err := c.resolve_path(statement.LeftHandSide)

		if err.Exists {
			return result, err
		}

		size--
		path = path[1:]

		result = append(result, path...)
		result = append(result, common.NewInstruction(common.OpSetItem, symbol.Index, size))
	}

	return result, errors.EmptyError
}

func (c *PackageCompiler) compile_unbound_fun_definition_statement(statement parser.UnboundFunDefinitionStatement) (common.InstructionSet, errors.Error) {
	result := common.InstructionSet{}

	fun_instructions := common.InstructionSet{}

	c.enter_scope()

	for _, parameter := range statement.Signature.Parameters {
		c.SymbolTable.Define(parameter.Name.Value, parser.ConstantKind, false)
	}

	for _, sub_statement := range statement.Body {
		instructions, err := c.compile_statement(sub_statement)
		if err.Exists {
			return result, err
		}

		fun_instructions = append(fun_instructions, instructions...)
	}
	c.leave_scope()

	value := common.FunctionObject{
		Value: fun_instructions,
	}
	index := c.ConstantPool.Add(value)
	result = append(result, common.NewInstruction(common.OpConstant, index))

	symbol, err := c.SymbolTable.Define(statement.Signature.Name.Value, parser.ConstantKind, statement.Hidden)
	if err != nil {
		return result, errors.CreateCompileError(err.Error(), statement.Signature.Name.Location())
	}

	if symbol.Scope == GlobalScope {
		result = append(result, common.NewInstruction(common.OpSet, symbol.Index))
	} else {
		result = append(result, common.NewInstruction(common.OpSetLocal, symbol.Index))
	}

	return result, errors.EmptyError
}

func (c *PackageCompiler) compile_bound_fun_definition_statement(statement parser.BoundFunDefinitionStatement) (common.InstructionSet, errors.Error) {
	result := common.InstructionSet{}

	for_ := statement.Signature.For.Name.(parser.IdentifierExpression).Value
	symbol_table := c.TypeSymbolTable[for_]

	if symbol_table == nil {
		return result, errors.CreateCompileError(fmt.Sprintf("type '%s' is not defined", for_), statement.Signature.For.Location())
	}

	fun_instructions := common.InstructionSet{}

	c.enter_scope()

	c.SymbolTable.Define("this", parser.ConstantKind, false)
	for _, parameter := range statement.Signature.Parameters {
		c.SymbolTable.Define(parameter.Name.Value, parser.ConstantKind, false)
	}

	for _, sub_statement := range statement.Body {
		instructions, err := c.compile_statement(sub_statement)
		if err.Exists {
			return result, err
		}

		fun_instructions = append(fun_instructions, instructions...)
	}
	c.leave_scope()

	value := common.FunctionObject{
		Value: fun_instructions,
	}
	index := c.ConstantPool.Add(value)
	result = append(result, common.NewInstruction(common.OpConstant, index))

	name := fmt.Sprintf("%s.%s.%s", c.PackageName, for_, statement.Signature.Name.Value)
	symbol, err := symbol_table.Define(name, parser.ConstantKind, statement.Hidden)
	if err != nil {
		return result, errors.CreateCompileError(err.Error(), statement.Signature.Name.Location())
	}

	result = append(result, common.NewInstruction(common.OpSet, symbol.Index))

	return result, errors.EmptyError
}

func (c *PackageCompiler) compile_type_definition_statement(statement parser.TypeDefinitionStatement) (common.InstructionSet, errors.Error) {
	result := common.InstructionSet{}

	if c.TypeSymbolTable[statement.Name.Value] != nil {
		return result, errors.CreateCompileError(fmt.Sprintf("cannot redeclare type '%s'", statement.Name.Value), statement.Name.Location())
	}

	c.TypeSymbolTable[statement.Name.Value] = NewSymbolTable()

	return result, errors.EmptyError
}

func (c *PackageCompiler) compile_return_statement(statement parser.ReturnStatement) (common.InstructionSet, errors.Error) {
	result := common.InstructionSet{}

	if statement.Value != nil {
		value, err := c.compile_expression(*statement.Value, false)

		if err.Exists {
			return result, err
		}

		result = append(result, value...)
		result = append(result, common.NewInstruction(common.OpReturn))
	} else {
		result = append(result, common.NewInstruction(common.OpReturnEmpty))
	}

	return result, errors.EmptyError
}

func (c *PackageCompiler) compile_continue_statement(statement parser.ContinueStatement) (common.InstructionSet, errors.Error) {
	result := common.InstructionSet{}
	result = append(result, common.NewInstruction(common.OpContinue))

	return result, errors.EmptyError
}

func (c *PackageCompiler) compile_break_statement(statement parser.BreakStatement) (common.InstructionSet, errors.Error) {
	result := common.InstructionSet{}
	result = append(result, common.NewInstruction(common.OpBreak))

	return result, errors.EmptyError
}

func (c *PackageCompiler) compile_yield_statement(statement parser.YieldStatement) (common.InstructionSet, errors.Error) {
	result := common.InstructionSet{}

	if statement.Value != nil {
		value, err := c.compile_expression(*statement.Value, false)

		if err.Exists {
			return result, err
		}

		result = append(result, value...)
		result = append(result, common.NewInstruction(common.OpYield))
	} else {
		return result, errors.CreateCompileError("no empty yield statement is allowed", statement.Location())
	}

	return result, errors.EmptyError
}

func (c *PackageCompiler) compile_if_statement(statement parser.IfStatement) (common.InstructionSet, errors.Error) {
	result := common.InstructionSet{}
	template := common.NewInstruction(common.OpJump, 0, 0)

	main_predicate, err := c.compile_expression(statement.MainBlock.Predicate, false)
	if err.Exists {
		return result, err
	}

	result = append(result, main_predicate...)

	main_block := common.InstructionSet{}
	for _, sub_statement := range statement.MainBlock.Body {
		instructions, err := c.compile_statement(sub_statement)
		if err.Exists {
			return result, err
		}

		main_block = append(main_block, instructions...)
	}

	result = append(result, common.NewInstruction(common.OpJumpIfFalse, main_block.GetSize()+template.GetSize(), 0))
	result = append(result, main_block...)

	for i, block := range statement.ElseIfBlocks {
		else_if_instructions := common.InstructionSet{}

		else_if_predicate, err := c.compile_expression(block.Predicate, false)
		if err.Exists {
			return result, err
		}

		for _, sub_statement := range block.Body {
			instructions, err := c.compile_statement(sub_statement)
			if err.Exists {
				return result, err
			}

			else_if_instructions = append(else_if_instructions, instructions...)
		}

		jump_count := else_if_instructions.GetSize() + else_if_predicate.GetSize() + template.GetSize()

		result = append(result, common.NewInstruction(common.OpJump, jump_count, 0))
		result = append(result, else_if_predicate...)
		if i == len(statement.ElseIfBlocks)-1 {
			result = append(result, common.NewInstruction(common.OpJumpIfFalse, else_if_instructions.GetSize(), 0))
		} else {
			result = append(result, common.NewInstruction(common.OpJumpIfFalse, else_if_instructions.GetSize()+template.GetSize(), 0))
		}
		result = append(result, else_if_instructions...)
	}

	else_block := common.InstructionSet{}
	for _, sub_statement := range statement.ElseBlock {
		instructions, err := c.compile_statement(sub_statement)
		if err.Exists {
			return result, err
		}

		else_block = append(else_block, instructions...)
	}

	result = append(result, common.NewInstruction(common.OpJump, else_block.GetSize(), 0))
	result = append(result, else_block...)

	return result, errors.EmptyError
}

func (c *PackageCompiler) compile_loop_statement(statement parser.LoopStatement) (common.InstructionSet, errors.Error) {
	switch statement.Predicate.LoopKind() {
	case parser.UnipartiteLoopKind:
		return c.compile_unipartite_loop_statement(statement)
	case parser.TripartiteLoopKind:
		return c.compile_tripartite_loop_predicate(statement)
	default:
		return common.InstructionSet{}, errors.CreateCompileError(fmt.Sprintf("unknown loop predicate kind %s", statement.Predicate.LoopKind()), statement.Location())
	}
}

func (c *PackageCompiler) compile_unipartite_loop_statement(statement parser.LoopStatement) (common.InstructionSet, errors.Error) {
	result := common.InstructionSet{}
	template := common.NewInstruction(common.OpJump, 0, 0)

	body := common.InstructionSet{}
	c.enter_scope()
	for _, sub_statement := range statement.Body {
		instructions, err := c.compile_statement(sub_statement)
		if err.Exists {
			return result, err
		}
		body = append(body, instructions...)
	}
	c.leave_scope()

	predicate := common.InstructionSet{}

	instructions, err := c.compile_expression(statement.Predicate.(parser.UnipartiteLoopPredicate).Expression, false)
	if err.Exists {
		return predicate, err
	}

	predicate = append(predicate, instructions...)
	predicate = append(predicate, common.NewInstruction(common.OpJumpIfFalse, body.GetSize()+template.GetSize()))

	result = append(result, predicate...)
	result = append(result, body...)
	result = append(result, common.NewInstruction(common.OpJump, body.GetSize()+predicate.GetSize(), 1))

	return result, errors.EmptyError
}

func (c *PackageCompiler) compile_tripartite_loop_predicate(statement parser.LoopStatement) (common.InstructionSet, errors.Error) {
	result := common.InstructionSet{}
	template := common.NewInstruction(common.OpJump, 0, 0)

	body := common.InstructionSet{}

	c.enter_scope()

	declaration, err := c.compile_statement(*statement.Predicate.(parser.TripartiteLoopPredicate).Declaration)
	if err.Exists {
		return result, err
	}
	result = append(result, declaration...)

	predicate, err := c.compile_expression(statement.Predicate.(parser.TripartiteLoopPredicate).Predicate, false)
	if err.Exists {
		return result, err
	}
	result = append(result, predicate...)

	c.enter_scope()
	for _, sub_statement := range statement.Body {
		instructions, err := c.compile_statement(sub_statement)
		if err.Exists {
			return result, err
		}
		body = append(body, instructions...)
	}
	c.leave_scope()
	c.leave_scope()

	result = append(result, common.NewInstruction(common.OpJumpIfFalse, body.GetSize()+template.GetSize(), 0))
	result = append(result, body...)
	result = append(result, common.NewInstruction(common.OpJump, body.GetSize()+template.GetSize()+predicate.GetSize(), 1))

	return result, errors.EmptyError
}

func (c *PackageCompiler) compile_expression(expression parser.Expression, should_clean bool) (common.InstructionSet, errors.Error) {
	result := common.InstructionSet{}
	err := errors.EmptyError

	switch expression.Kind() {
	case parser.NumberLiteralExpressionKind,
		parser.MapLiteralExpressionKind,
		parser.InstanceLiteralExpressionKind,
		parser.BoolLiteralExpressionKind,
		parser.StringLiteralExpressionKind,
		parser.RuneLiteralExpressionKind,
		parser.ListLiteralExpressionKind:
		result, err = c.compile_literal_expression(expression.(parser.LiteralExpression))
	case parser.AnonymousFunExpressionKind:
		result, err = c.compile_fun_expression(expression.(parser.AnonymousFunExpression))
	case parser.GroupExpressionKind:
		result, err = c.compile_expression(expression.(parser.GroupExpression).Expression, false)
	case parser.IdentifierExpressionKind:
		result, err = c.compile_identifier_expression(expression.(parser.IdentifierExpression))
	case parser.ArithmeticExpressionKind:
		result, err = c.compile_arithmetic_expression(expression.(parser.ArithmeticExpression))
	case parser.ArithmeticUnaryExpressionKind:
		result, err = c.compile_arithmetic_unary_expression(expression.(parser.ArithmeticUnaryExpression))
	case parser.NotExpressionKind:
		result, err = c.compile_not_expression(expression.(parser.NotExpression))
	case parser.GiveupExpressionKind:
		result, err = c.compile_giveup_expression(expression.(parser.GiveupExpression))
	case parser.ComparisonExpressionKind:
		result, err = c.compile_comparison_expression(expression.(parser.ComparisonExpression))
	case parser.BinaryExpressionKind:
		result, err = c.compile_binary_expression(expression.(parser.BinaryExpression))
	case parser.IndexExpressionKind:
		result, err = c.compile_index_expression(expression.(parser.IndexExpression))
	case parser.MemberExpressionKind:
		result, err = c.compile_member_expression(expression.(parser.MemberExpression))
	case parser.CallExpressionKind:
		result, err = c.compile_call_expression(expression.(parser.CallExpression))
	case parser.InstanceofExpressionKind:
		result, err = c.compile_instanceof_expression(expression.(parser.InstanceofExpression))
	case parser.TypeCastExpressionKind:
		result, err = c.compile_type_cast_expression(expression.(parser.TypeCastExpression))
	case parser.ThisExpressionKind:
		result, err = c.compile_this_expression(expression.(parser.ThisExpression))
	case parser.MatchSelfExpressionKind:
		result, err = c.compile_match_self_expression(expression.(parser.MatchSelfExpression))
	case parser.MatchExpressionKind:
		result, err = c.compile_match_expression(expression.(parser.MatchExpression))
	default:
		result = append(result, common.NewInstruction(common.OpNoop))
	}

	if should_clean {
		result = append(result, common.NewInstruction(common.OpPop))
	}

	return result, err
}

func (c *PackageCompiler) compile_literal_expression(expression parser.LiteralExpression) (common.InstructionSet, errors.Error) {
	result := common.InstructionSet{}
	err := errors.EmptyError

	switch expression.LiteralKind() {
	case parser.NumberLiteralKind:
		value := expression.(parser.NumberLiteralExpression).Value.Value.(int)
		index := c.ConstantPool.Add(common.Int32Object{
			Value: int32(value),
		})

		result = append(result, common.NewInstruction(common.OpConstant, index))
	case parser.StringLiteralKind:
		value := expression.(parser.StringLiteralExpression).Value
		index := c.ConstantPool.Add(common.StringObject{
			Value: value,
		})

		result = append(result, common.NewInstruction(common.OpConstant, index))
	case parser.RuneLiteralKind:
		value := expression.(parser.RuneLiteralExpression).Value
		index := c.ConstantPool.Add(common.RuneObject{
			Value: value,
		})

		result = append(result, common.NewInstruction(common.OpConstant, index))
	case parser.BoolLiteralKind:
		value := expression.(parser.BoolLiteralExpression).Value

		if value {
			result = append(result, common.NewInstruction(common.OpTrue))
		} else {
			result = append(result, common.NewInstruction(common.OpFalse))
		}
	case parser.ListLiteralKind:
		list := expression.(parser.ListLiteralExpression)

		for _, value := range list.Value {
			instructions, err := c.compile_expression(value.Value, false)
			if err.Exists {
				return result, err
			}
			result = append(result, instructions...)
		}

		result = append(result, common.NewInstruction(common.OpArray, len(list.Value)))
	case parser.MapLiteralKind:
		map_ := expression.(parser.MapLiteralExpression)

		for _, entry := range map_.Value {
			key, err := c.compile_literal_expression(parser.StringLiteralExpression{Value: entry.Key.Value})
			if err.Exists {
				return result, err
			}
			result = append(result, key...)
			value, err := c.compile_expression(entry.Value, false)
			if err.Exists {
				return result, err
			}
			result = append(result, value...)
		}
		result = append(result, common.NewInstruction(common.OpMap, len(map_.Value)))
	case parser.InstanceLiteralKind:
		instance := expression.(parser.InstanceLiteralExpression)

		for _, entry := range instance.Value {
			key, err := c.compile_literal_expression(parser.StringLiteralExpression{Value: entry.Key.Value})
			if err.Exists {
				return result, err
			}
			result = append(result, key...)
			value, err := c.compile_expression(entry.Value, false)
			if err.Exists {
				return result, err
			}
			result = append(result, value...)
		}
		result = append(result, common.NewInstruction(common.OpMap, len(instance.Value)))
	}

	return result, err
}

func (c *PackageCompiler) compile_fun_expression(expression parser.AnonymousFunExpression) (common.InstructionSet, errors.Error) {
	result := common.InstructionSet{}

	fun_instructions := common.InstructionSet{}

	c.enter_scope()

	for _, parameter := range expression.Signature.Parameters {
		c.SymbolTable.Define(parameter.Name.Value, parser.ConstantKind, false)
	}

	for _, sub_statement := range expression.Body {
		instructions, err := c.compile_statement(sub_statement)
		if err.Exists {
			return result, err
		}

		fun_instructions = append(fun_instructions, instructions...)
	}
	c.leave_scope()

	value := common.FunctionObject{
		Value: fun_instructions,
	}
	index := c.ConstantPool.Add(value)
	result = append(result, common.NewInstruction(common.OpConstant, index))

	return result, errors.EmptyError
}

func (c *PackageCompiler) compile_arithmetic_expression(expression parser.ArithmeticExpression) (common.InstructionSet, errors.Error) {
	result := common.InstructionSet{}

	left, err := c.compile_expression(expression.LeftHandSide, false)
	if err.Exists {
		return result, err
	}
	right, err := c.compile_expression(expression.RightHandSide, false)
	if err.Exists {
		return result, err
	}

	result = append(result, left...)
	result = append(result, right...)

	switch expression.Operator.Literal {
	case "+":
		result = append(result, common.NewInstruction(common.OpAdd))
	case "-":
		result = append(result, common.NewInstruction(common.OpSub))
	case "*":
		result = append(result, common.NewInstruction(common.OpMul))
	case "/":
		result = append(result, common.NewInstruction(common.OpDiv))
	case "%":
		result = append(result, common.NewInstruction(common.OpMod))
	}

	return result, errors.EmptyError
}

func (c *PackageCompiler) compile_arithmetic_unary_expression(expression parser.ArithmeticUnaryExpression) (common.InstructionSet, errors.Error) {
	result := common.InstructionSet{}
	var operator string

	if expression.Operation == parser.IncrementKind {
		operator = "+"
	} else {
		operator = "-"
	}

	instructions, err := c.compile_assignment_statement(parser.AssignmentStatement{
		LeftHandSide: expression.Expression,
		RightHandSide: parser.ArithmeticExpression{
			LeftHandSide:  expression.Expression,
			RightHandSide: parser.NumberLiteralExpression{Value: parser.NumberLiteral{Value: 1}},
			Operator: parser.OperatorToken{
				Literal: operator,
			},
		},
		Operator: parser.OperatorToken{
			Literal: "=",
		},
	})

	if err.Exists {
		return result, err
	}

	result = append(result, instructions...)

	instructions, err = c.compile_expression(expression.Expression, false)

	if err.Exists {
		return result, err
	}

	result = append(result, instructions...)

	return result, errors.EmptyError
}

func (c *PackageCompiler) compile_not_expression(expression parser.NotExpression) (common.InstructionSet, errors.Error) {
	result := common.InstructionSet{}

	instructions, err := c.compile_expression(expression.Expression, false)
	if err.Exists {
		return result, err
	}

	result = append(result, instructions...)
	result = append(result, common.NewInstruction(common.OpNegate))

	return result, errors.EmptyError
}

func (c *PackageCompiler) compile_giveup_expression(expression parser.GiveupExpression) (common.InstructionSet, errors.Error) {
	result := common.InstructionSet{}

	result = append(result, common.NewInstruction(common.OpExit, 1))

	return result, errors.EmptyError
}

func (c *PackageCompiler) compile_comparison_expression(expression parser.ComparisonExpression) (common.InstructionSet, errors.Error) {
	result := common.InstructionSet{}

	left, err := c.compile_expression(expression.LeftHandSide, false)
	if err.Exists {
		return result, err
	}
	right, err := c.compile_expression(expression.RightHandSide, false)
	if err.Exists {
		return result, err
	}

	if expression.Operator.Literal == "<" || expression.Operator.Literal == "<=" {
		result = append(result, right...)
		result = append(result, left...)
	} else {
		result = append(result, left...)
		result = append(result, right...)
	}

	switch expression.Operator.Literal {
	case ">", "<":
		result = append(result, common.NewInstruction(common.OpGreaterThan))
	case ">=", "<=":
		result = append(result, common.NewInstruction(common.OpGreaterThanOrEqual))
	case "==":
		result = append(result, common.NewInstruction(common.OpEqual))
	case "!=":
		result = append(result, common.NewInstruction(common.OpNotEqual))
	}

	return result, errors.EmptyError
}

func (c *PackageCompiler) compile_binary_expression(expression parser.BinaryExpression) (common.InstructionSet, errors.Error) {
	result := common.InstructionSet{}

	left, err := c.compile_expression(expression.LeftHandSide, false)
	if err.Exists {
		return result, err
	}
	right, err := c.compile_expression(expression.RightHandSide, false)
	if err.Exists {
		return result, err
	}

	result = append(result, left...)
	result = append(result, right...)

	switch expression.Operator.Literal {
	case "&&":
		result = append(result, common.NewInstruction(common.OpAnd))
	case "||":
		result = append(result, common.NewInstruction(common.OpOr))
	}

	return result, errors.EmptyError
}

func (c *PackageCompiler) compile_identifier_expression(expression parser.IdentifierExpression) (common.InstructionSet, errors.Error) {
	result := common.InstructionSet{}

	symbol := c.SymbolTable.Resolve(expression.Value)

	if symbol == nil {
		return result, errors.CreateCompileError(fmt.Sprintf("variable '%s' is not defined", expression.Value), expression.Location())
	}

	if symbol.Scope == GlobalScope {
		result = append(result, common.NewInstruction(common.OpGet, symbol.Index))
	} else {
		result = append(result, common.NewInstruction(common.OpGetLocal, symbol.Index))
	}

	return result, errors.EmptyError
}

func (c *PackageCompiler) compile_index_expression(expression parser.IndexExpression) (common.InstructionSet, errors.Error) {
	result := common.InstructionSet{}

	host, err := c.compile_expression(expression.Host, false)
	if err.Exists {
		return result, err
	}

	index, err := c.compile_expression(expression.Index, false)
	if err.Exists {
		return result, err
	}

	result = append(result, host...)
	result = append(result, index...)
	result = append(result, common.NewInstruction(common.OpIndex))

	return result, errors.EmptyError
}

func (c *PackageCompiler) compile_member_expression(expression parser.MemberExpression) (common.InstructionSet, errors.Error) {
	result := common.InstructionSet{}

	left, err := c.compile_expression(expression.LeftHandSide, false)
	if err.Exists {
		return result, err
	}

	value := expression.RightHandSide.Value
	index := c.ConstantPool.Add(common.StringObject{
		Value: value,
	})

	result = append(result, left...)
	result = append(result, common.NewInstruction(common.OpConstant, index))
	result = append(result, common.NewInstruction(common.OpIndex))

	return result, errors.EmptyError
}

func (c *PackageCompiler) compile_call_expression(expression parser.CallExpression) (common.InstructionSet, errors.Error) {
	result := common.InstructionSet{}

	callee, err := c.compile_expression(expression.Callee, false)
	if err.Exists {
		return result, err
	}

	for _, argument := range expression.Arguments {
		instructions, err := c.compile_expression(argument, false)
		if err.Exists {
			return result, err
		}
		result = append(result, instructions...)
	}

	result = append(result, callee...)
	result = append(result, common.NewInstruction(common.OpCall, len(expression.Arguments)))

	return result, errors.EmptyError
}

func (c *PackageCompiler) compile_instanceof_expression(expression parser.InstanceofExpression) (common.InstructionSet, errors.Error) {
	return c.compile_literal_expression(parser.BoolLiteralExpression{
		Value: true,
	})
}

func (c *PackageCompiler) compile_type_cast_expression(expression parser.TypeCastExpression) (common.InstructionSet, errors.Error) {
	result := common.InstructionSet{}

	instructions, err := c.compile_expression(expression.Value, false)
	if err.Exists {
		return result, err
	}
	result = append(result, instructions...)

	// add type somehow
	result = append(result, common.NewInstruction(common.OpCast))

	return result, errors.EmptyError
}

func (c *PackageCompiler) compile_this_expression(expression parser.ThisExpression) (common.InstructionSet, errors.Error) {
	result := common.InstructionSet{}
	result = append(result, common.NewInstruction(common.OpGetLocal, 0))

	return result, errors.EmptyError
}

func (c *PackageCompiler) compile_match_self_expression(expression parser.MatchSelfExpression) (common.InstructionSet, errors.Error) {
	if c.current_match_target == nil {
		return common.InstructionSet{}, errors.CreateTypeError("match self expressions are not allowed outside of match expressions", expression.Location())
	}

	return c.current_match_target, errors.EmptyError
}

type match_block struct {
	predicate common.InstructionSet
	body      common.InstructionSet
}

func (c *PackageCompiler) compile_match_block(block parser.PredicateBlock, base_exists bool) (match_block, errors.Error) {
	// To reliably measure the sizes of the blocks later, dummy jump instructions are added
	result := match_block{
		predicate: common.InstructionSet{},
		body:      common.InstructionSet{},
	}

	predicate, err := c.compile_expression(block.Predicate, false)
	if err.Exists {
		return result, err
	}

	predicate = append(predicate, common.NewInstruction(common.OpNegate))
	result.predicate = append(result.predicate, predicate...)

	for _, sub_statement := range block.Body {
		instructions, err := c.compile_statement(sub_statement)
		if err.Exists {
			return result, err
		}
		result.body = append(result.body, instructions...)
	}

	result.predicate = append(result.predicate, common.NewInstruction(common.OpJumpIfFalse, 0, 0))

	if base_exists {
		// if predicate is false this part will run at the end of the predicate
		false_ := c.ConstantPool.Add(common.Uint32Object{Value: 0})
		result.predicate = append(result.predicate, common.NewInstruction(common.OpConstant, false_))

		// if predicate is true this part will run at the end of the body
		true_ := c.ConstantPool.Add(common.Uint32Object{Value: 1})
		result.body = append(result.body, common.NewInstruction(common.OpConstant, true_))
	}

	result.body = append(result.body, common.NewInstruction(common.OpJump, 0, 1))

	return result, errors.EmptyError
}

func (c *PackageCompiler) compile_match_expression(expression parser.MatchExpression) (common.InstructionSet, errors.Error) {
	/* this works by putting the predicates one after the other
	and the bodies after them with order. When a predicate is executed
	it will jump to its body and in the end it will push either a 0 or 1
	depending on if its condition was true or not. In between the predicates
	and the bodies there is a block that will check if all the added values
	add up to 0, basically meaning no condition was met. If so it will jump to
	the base block, if not it will jump out. Because all the bodies jump back
	to their predicates, this in between block will be the last thing that is run.
	*/
	result := common.InstructionSet{}
	blocks := []match_block{}
	base_exists := len(expression.BaseBlock) != 0

	against, err := c.compile_expression(expression.Against, false)
	if err.Exists {
		return common.InstructionSet{}, err
	}

	previous_match_target := c.current_match_target
	c.current_match_target = against

	for _, block := range expression.Blocks {
		match_block, err := c.compile_match_block(block, base_exists)
		if err.Exists {
			return result, err
		}
		blocks = append(blocks, match_block)
	}

	base_predicate := common.InstructionSet{}
	if base_exists {
		/* reduce the accumulated values by adding them up
		 basically if there were 4 predicates, 4 values would
		would accumulate here, running the add instruction n - 1
		times will reduce them into a single value */
		for i := 0; i < len(blocks)-1; i++ {
			base_predicate = append(base_predicate, common.NewInstruction(common.OpAdd))
		}

		// if the reduced value is 0, all predicates failed so jump to the base block otherwise jump out
		false_ := c.ConstantPool.Add(common.Uint32Object{Value: 0})
		base_predicate = append(base_predicate, common.NewInstruction(common.OpConstant, false_))
		base_predicate = append(base_predicate, common.NewInstruction(common.OpJumpIfFalse, 0, 0))
		base_predicate = append(base_predicate, common.NewInstruction(common.OpJump, 0, 0))
	}

	c.current_match_target = previous_match_target

	base_jump_size := 0
	// figure out the actual jump counts for conditions
	for i, block := range blocks {
		base_jump_size += block.body.GetSize()
		predicate_jump_size := 0

		/* a predicate's jump size is the size of the predicates
		after it + bodies before it + base predicate's size */
		for _, consecutive_block := range blocks[i+1:] {
			predicate_jump_size += consecutive_block.predicate.GetSize()
		}
		predicate_jump_size += base_predicate.GetSize()
		for _, previous_block := range blocks[0:i] {
			predicate_jump_size += previous_block.body.GetSize()
		}

		jump_template := common.NewInstruction(common.OpJump, 0, 0)
		constant_template := common.NewInstruction(common.OpConstant, 0)
		/* the body's jump size is the predicate's jump size +
		the body's own size. Predicate includes a jump at the
		end so we will remove that. A jump instruction size is 19 */
		body_jump_size := predicate_jump_size + block.body.GetSize() - jump_template.GetSize()
		/* jump the last constant that pushes 0 onto the stack as well.
		This is a constant instruction that has the size 5 */
		predicate_jump_size += constant_template.GetSize()

		// replace the jumps
		block.predicate[len(block.predicate)-2] = common.NewInstruction(common.OpJumpIfFalse, predicate_jump_size, 0)
		block.body[len(block.body)-1] = common.NewInstruction(common.OpJump, body_jump_size, 1)

		result = append(result, block.predicate...)
	}

	if base_exists {
		// figure out the actual jump counts for exiting
		base_predicate[len(base_predicate)-2] = common.NewInstruction(common.OpJumpIfFalse, base_jump_size+base_predicate.GetSize(), 0)
		base_predicate[len(base_predicate)-1] = common.NewInstruction(common.OpJump, base_jump_size, 0)

		result = append(result, base_predicate...)
	}

	for _, block := range blocks {
		result = append(result, block.body...)
	}

	return result, err
}
