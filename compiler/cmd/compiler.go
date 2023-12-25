package cmd

import (
	"fmt"

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

func (c *PackageCompiler) compile_statement(statement parser.Statement) (common.InstructionSet, errors.Error) {
	switch statement.Kind() {
	case parser.ExpressionStatementKind:
		return c.compile_expression(statement.(parser.ExpressionStatement).Expression, true)
	case parser.DeclarationStatementKind:
		return c.compile_declaration_statement(statement.(parser.DeclarationStatement))
	case parser.UnboundFunDefinitionStatementKind:
		return c.compile_unbound_fun_definition_statement(*statement.(*parser.UnboundFunDefinitionStatement))
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

	symbol := c.SymbolTable.Define(statement.Name.Value, statement.VarKind)

	if symbol.Scope == GlobalScope {
		result = append(result, common.NewInstruction(common.OpSet, symbol.Index))
	} else {
		result = append(result, common.NewInstruction(common.OpSetLocal, symbol.Index))
	}

	return result, errors.EmptyError
}

func (c *PackageCompiler) compile_unbound_fun_definition_statement(statement parser.UnboundFunDefinitionStatement) (common.InstructionSet, errors.Error) {
	result := common.InstructionSet{}

	fun_instructions := common.InstructionSet{}

	c.enter_scope()
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

	symbol := c.SymbolTable.Define(statement.Signature.Name.Value, parser.ConstantKind)

	if symbol.Scope == GlobalScope {
		result = append(result, common.NewInstruction(common.OpSet, symbol.Index))
	} else {
		result = append(result, common.NewInstruction(common.OpSetLocal, symbol.Index))
	}

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
	result := common.InstructionSet{}
	template := common.NewInstruction(common.OpJump, 0, 0)

	var predicate common.InstructionSet
	var err errors.Error

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

	switch statement.Predicate.LoopKind() {
	case parser.UnipartiteLoopKind:
		predicate, err = c.compile_uinpartite_loop_predicate(statement.Predicate.(parser.UnipartiteLoopPredicate), body.GetSize()+template.GetSize())
		if err.Exists {
			return result, err
		}
	}

	result = append(result, predicate...)
	result = append(result, body...)
	result = append(result, common.NewInstruction(common.OpJump, body.GetSize()+predicate.GetSize(), 1))

	return result, errors.EmptyError
}

func (c *PackageCompiler) compile_uinpartite_loop_predicate(predicate parser.UnipartiteLoopPredicate, size int) (common.InstructionSet, errors.Error) {
	result := common.InstructionSet{}

	instructions, err := c.compile_expression(predicate.Expression, false)
	if err.Exists {
		return result, err
	}

	result = append(result, instructions...)
	result = append(result, common.NewInstruction(common.OpJumpIfFalse, size, 0))

	return result, errors.EmptyError
}

func (c *PackageCompiler) compile_expression(expression parser.Expression, should_clean bool) (common.InstructionSet, errors.Error) {
	result := common.InstructionSet{}
	err := errors.EmptyError

	switch expression.Kind() {
	case parser.NumberLiteralExpressionKind,
		parser.BoolLiteralExpressionKind,
		parser.StringLiteralExpressionKind,
		parser.ListLiteralExpressionKind:
		result, err = c.compile_literal_expression(expression.(parser.LiteralExpression))
	case parser.GroupExpressionKind:
		result, err = c.compile_expression(expression.(parser.GroupExpression).Expression, false)
	case parser.IdentifierExpressionKind:
		result, err = c.compile_identifier_expression(expression.(parser.IdentifierExpression))
	case parser.ArithmeticExpressionKind:
		result, err = c.compile_arithmetic_expression(expression.(parser.ArithmeticExpression))
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
	case parser.BoolLiteralKind:
		value := expression.(parser.BoolLiteralExpression).Value

		if value {
			result = append(result, common.NewInstruction(common.OpTrue))
		} else {
			result = append(result, common.NewInstruction(common.OpFalse))
		}
	}

	return result, err
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

	switch expression.Operator {
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

	if expression.Operator == "<" || expression.Operator == "<=" {
		result = append(result, right...)
		result = append(result, left...)
	} else {
		result = append(result, left...)
		result = append(result, right...)
	}

	switch expression.Operator {
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

	result = append(result, callee...)
	result = append(result, common.NewInstruction(common.OpCall))

	return result, errors.EmptyError
}

func (c *PackageCompiler) compile_instanceof_expression(expression parser.InstanceofExpression) (common.InstructionSet, errors.Error) {
	return c.compile_literal_expression(parser.BoolLiteralExpression{
		Value: true,
	})
}
