package cmd

import (
	"fmt"

	"github.com/moonbite-org/moonbite/common"
	errors "github.com/moonbite-org/moonbite/error"
	parser "github.com/moonbite-org/moonbite/parser/cmd"
)

func (c *FileCompiler) compile_statement(statement parser.Statement) ([]common.Instruction, errors.Error) {
	switch statement.Kind() {
	case parser.ExpressionStatementKind:
		return c.compile_expression(statement.(parser.ExpressionStatement).Expression, true)
	case parser.DeclarationStatementKind:
		return c.compile_declaration_statement(statement.(parser.DeclarationStatement))
	case parser.UnboundFunDefinitionStatementKind:
		return c.compile_unbound_fun_definition_statement(*statement.(*parser.UnboundFunDefinitionStatement))
	default:
		result := []common.Instruction{}
		result = append(result, common.NewInstruction(common.OpNoop))
		return result, errors.EmptyError
	}
}

func (c *FileCompiler) compile_declaration_statement(statement parser.DeclarationStatement) ([]common.Instruction, errors.Error) {
	result := []common.Instruction{}

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

	symbol := c.SymbolTable.Define(statement.Name.Value, statement.VarKind, c.current_scope)
	result = append(result, common.NewInstruction(common.OpSet, symbol.Index))
	return result, errors.EmptyError
}

func (c *FileCompiler) compile_unbound_fun_definition_statement(statement parser.UnboundFunDefinitionStatement) ([]common.Instruction, errors.Error) {
	result := []common.Instruction{}

	fun_instructions := common.InstructionSet{}
	c.last_scope = c.current_scope
	c.current_scope = LocalScope
	for _, sub_statement := range statement.Body {
		instructions, err := c.compile_statement(sub_statement)
		if err.Exists {
			return result, err
		}

		fun_instructions = append(fun_instructions, instructions...)
	}
	c.current_scope = c.last_scope

	value := common.FunctionObject{
		Value: fun_instructions,
	}
	index := c.ConstantPool.Add(value)
	result = append(result, common.NewInstruction(common.OpConstant, index))

	symbol := c.SymbolTable.Define(statement.Signature.Name.Value, parser.ConstantKind, c.current_scope)
	result = append(result, common.NewInstruction(common.OpSet, symbol.Index))

	return result, errors.EmptyError
}

func (c *FileCompiler) compile_expression(expression parser.Expression, should_clean bool) ([]common.Instruction, errors.Error) {
	result := []common.Instruction{}
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
	default:
		result = append(result, common.NewInstruction(common.OpNoop))
	}

	if should_clean {
		result = append(result, common.NewInstruction(common.OpPop))
	}

	return result, err
}

func (c *FileCompiler) compile_literal_expression(expression parser.LiteralExpression) ([]common.Instruction, errors.Error) {
	result := []common.Instruction{}
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

func (c *FileCompiler) compile_arithmetic_expression(expression parser.ArithmeticExpression) ([]common.Instruction, errors.Error) {
	result := []common.Instruction{}

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

func (c *FileCompiler) compile_binary_expression(expression parser.BinaryExpression) ([]common.Instruction, errors.Error) {
	result := []common.Instruction{}

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
	case ">":
		result = append(result, common.NewInstruction(common.OpGreaterThan))
	case ">=":
		result = append(result, common.NewInstruction(common.OpGreaterThanOrEqual))
	case "<":
		result = append(result, common.NewInstruction(common.OpGreaterThan))
	case "<=":
		result = append(result, common.NewInstruction(common.OpGreaterThanOrEqual))
	case "==":
		result = append(result, common.NewInstruction(common.OpEqual))
	case "!=":
		result = append(result, common.NewInstruction(common.OpNotEqual))
	}

	return result, errors.EmptyError
}

func (c *FileCompiler) compile_identifier_expression(expression parser.IdentifierExpression) ([]common.Instruction, errors.Error) {
	result := []common.Instruction{}

	symbol := c.SymbolTable.Resolve(expression.Value)

	if symbol == nil {
		return result, errors.CreateCompileError(fmt.Sprintf("cannot find variable '%s'", expression.Value), expression.Location())
	}

	result = append(result, common.NewInstruction(common.OpGet, symbol.Index))

	return result, errors.EmptyError
}
