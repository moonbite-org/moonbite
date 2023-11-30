package parser_test

import (
	"reflect"
	"testing"

	"github.com/moonbite-org/moonbite/common"
	"github.com/moonbite-org/moonbite/parser"
)

func assert_identifier(t *testing.T, given parser.IdentifierExpression, expected string) {
	if given.Value != expected {
		t.Errorf("expected identifier to be %s but found %s", expected, given.Value)
	}
}

func assert_no_error(t *testing.T, err common.Error) {
	if err.Exists {
		t.Errorf("expected no error but got: %s", err)
	}
}

func assert_error(t *testing.T, err common.Error) {
	if !err.Exists {
		t.Errorf("expected error but no error is present")
	}
}

func assert_string(t *testing.T, given, expexted string) {
	if given != expexted {
		t.Errorf("expected string to be %s but got %s", expexted, given)
	}
}

func assert_int(t *testing.T, given, expexted int) {
	if given != expexted {
		t.Errorf("expected int to be %d but got %d", expexted, given)
	}
}

func assert_type(t *testing.T, given, expected any) {
	expected_t := reflect.TypeOf(expected)
	given_t := reflect.TypeOf(given)

	if given_t != expected_t {
		t.Errorf("expected type to be %s but found %s", expected_t, given_t)
	}
}

func TestPackageStatement(t *testing.T) {
	input := []byte("")
	_, err := parser.Parse(input, "test.mb")

	if !err.Exists {
		t.Errorf("expected package keyword error")
	}

	input = []byte("package main")
	ast, err := parser.Parse(input, "test.mb")

	assert_identifier(t, ast.Package.Name, "main")
}

func TestUseStatement(t *testing.T) {
	input := []byte("package main use os use binary as bin")
	ast, err := parser.Parse(input, "test.mb")

	assert_no_error(t, err)

	assert_identifier(t, ast.Uses[0].Resource, "os")
	assert_identifier(t, ast.Uses[1].Resource, "binary")
	assert_identifier(t, *ast.Uses[1].As, "bin")
}

func TestDeclarationStatement(t *testing.T) {
	input := []byte("package main const test = data.count")
	ast, err := parser.Parse(input, "test.mb")

	assert_no_error(t, err)

	definition := ast.Definitions[0]

	assert_type(t, definition, parser.DeclarationStatement{})
}

func TestIdentifierExpression(t *testing.T) {
	input := []byte("package main const test = data")
	ast, err := parser.Parse(input, "test.mb")

	assert_no_error(t, err)

	definition := ast.Definitions[0]

	assert_type(t, definition, parser.DeclarationStatement{})
	assert_type(t, *definition.(parser.DeclarationStatement).Value, parser.IdentifierExpression{})
}

func TestArithmeticExpression(t *testing.T) {
	input := []byte("package main const test = 2 + 3")
	ast, err := parser.Parse(input, "test.mb")

	assert_no_error(t, err)

	definition := ast.Definitions[0]
	assert_type(t, *definition.(parser.DeclarationStatement).Value, parser.ArithmeticExpression{})
	expression := (*definition.(parser.DeclarationStatement).Value).(parser.ArithmeticExpression)

	assert_type(t, expression.LeftHandSide, parser.NumberLiteralExpression{})
	assert_type(t, expression.RightHandSide, parser.NumberLiteralExpression{})
	assert_string(t, expression.Operator, "+")

	input = []byte("package main const test = 2 + 3 * 5")
	ast, err = parser.Parse(input, "test.mb")
	assert_no_error(t, err)

	definition = ast.Definitions[0]
	assert_type(t, *definition.(parser.DeclarationStatement).Value, parser.ArithmeticExpression{})
	expression = (*definition.(parser.DeclarationStatement).Value).(parser.ArithmeticExpression)

	assert_type(t, expression.LeftHandSide, parser.NumberLiteralExpression{})
	assert_type(t, expression.RightHandSide, parser.ArithmeticExpression{})
	assert_string(t, expression.Operator, "+")
	assert_string(t, expression.RightHandSide.(parser.ArithmeticExpression).Operator, "*")

	input = []byte("package main const test = (2 + 3) * 5")
	ast, err = parser.Parse(input, "test.mb")
	assert_no_error(t, err)

	definition = ast.Definitions[0]
	assert_type(t, *definition.(parser.DeclarationStatement).Value, parser.ArithmeticExpression{})
	expression = (*definition.(parser.DeclarationStatement).Value).(parser.ArithmeticExpression)

	assert_type(t, expression.LeftHandSide, parser.GroupExpression{})
	assert_type(t, expression.RightHandSide, parser.NumberLiteralExpression{})
	assert_string(t, expression.Operator, "*")

	group := expression.LeftHandSide.(parser.GroupExpression)

	assert_type(t, group.Expression, parser.ArithmeticExpression{})
	assert_string(t, group.Expression.(parser.ArithmeticExpression).Operator, "+")

	input = []byte("package main const test = (2 + 3 * 5")
	ast, err = parser.Parse(input, "test.mb")
	assert_error(t, err)
}

func TestBinaryExpression(t *testing.T) {
	input := []byte("package main const test = 2 == count")
	ast, err := parser.Parse(input, "test.mb")

	assert_no_error(t, err)

	definition := ast.Definitions[0]
	assert_type(t, *definition.(parser.DeclarationStatement).Value, parser.BinaryExpression{})

	expression := (*definition.(parser.DeclarationStatement).Value).(parser.BinaryExpression)
	assert_type(t, expression.LeftHandSide, parser.NumberLiteralExpression{})
	assert_type(t, expression.RightHandSide, parser.IdentifierExpression{})
	assert_string(t, expression.Operator, "==")

	input = []byte("package main const test = 2 > 5")
	ast, err = parser.Parse(input, "test.mb")

	assert_no_error(t, err)

	definition = ast.Definitions[0]
	assert_type(t, *definition.(parser.DeclarationStatement).Value, parser.BinaryExpression{})

	expression = (*definition.(parser.DeclarationStatement).Value).(parser.BinaryExpression)
	assert_type(t, expression.LeftHandSide, parser.NumberLiteralExpression{})
	assert_type(t, expression.RightHandSide, parser.NumberLiteralExpression{})
	assert_string(t, expression.Operator, ">")

	input = []byte("package main const test = 2 >< 5")
	ast, err = parser.Parse(input, "test.mb")

	assert_error(t, err)

	input = []byte("package main const test = <string>")
	ast, err = parser.Parse(input, "test.mb")

	assert_error(t, err)

	input = []byte("package main const test = 2 == ")
	ast, err = parser.Parse(input, "test.mb")

	assert_error(t, err)

	input = []byte("package main const test = == ")
	ast, err = parser.Parse(input, "test.mb")

	assert_error(t, err)
}

func TestCallExpression(t *testing.T) {
	input := []byte("package main const test = print()")
	ast, err := parser.Parse(input, "test.mb")

	assert_no_error(t, err)

	definition := ast.Definitions[0]
	assert_type(t, *definition.(parser.DeclarationStatement).Value, parser.CallExpression{})
	expression := (*definition.(parser.DeclarationStatement).Value).(parser.CallExpression)

	assert_type(t, expression.Callee, parser.IdentifierExpression{})
	assert_string(t, expression.Callee.(parser.IdentifierExpression).Value, "print")
	assert_int(t, len(expression.Arguments), 0)

	input = []byte("package main const test = console.log(true)")
	ast, err = parser.Parse(input, "test.mb")

	assert_no_error(t, err)

	definition = ast.Definitions[0]
	assert_type(t, *definition.(parser.DeclarationStatement).Value, parser.CallExpression{})
	expression = (*definition.(parser.DeclarationStatement).Value).(parser.CallExpression)

	assert_type(t, expression.Callee, parser.MemberExpression{})
	assert_int(t, len(expression.Arguments), 1)
	assert_type(t, expression.Arguments[0], parser.BoolLiteralExpression{})

	input = []byte("package main const test = console.log(true")
	ast, err = parser.Parse(input, "test.mb")

	assert_error(t, err)

	input = []byte("package main const test = console.log(true,)")
	ast, err = parser.Parse(input, "test.mb")

	assert_error(t, err)

	input = []byte("package main const test = (2 + 2)()")
	ast, err = parser.Parse(input, "test.mb")

	assert_error(t, err)
}

func TestMemberExpression(t *testing.T) {
	input := []byte("package main const test = data.count")
	ast, err := parser.Parse(input, "test.mb")

	assert_no_error(t, err)

	definition := ast.Definitions[0]
	assert_type(t, *definition.(parser.DeclarationStatement).Value, parser.MemberExpression{})
	expression := (*definition.(parser.DeclarationStatement).Value)

	assert_type(t, expression.(parser.MemberExpression).LeftHandSide, parser.IdentifierExpression{})
	assert_type(t, expression.(parser.MemberExpression).RightHandSide, parser.IdentifierExpression{})

	input = []byte("package main const test = console.log()")
	ast, err = parser.Parse(input, "test.mb")

	assert_no_error(t, err)

	definition = ast.Definitions[0]
	assert_type(t, *definition.(parser.DeclarationStatement).Value, parser.CallExpression{})
	expression = (*definition.(parser.DeclarationStatement).Value)

	assert_type(t, expression.(parser.CallExpression).Callee, parser.MemberExpression{})
	member := expression.(parser.CallExpression).Callee.(parser.MemberExpression)
	assert_type(t, member.LeftHandSide, parser.IdentifierExpression{})
	assert_type(t, member.RightHandSide, parser.IdentifierExpression{})
	assert_string(t, member.LeftHandSide.(parser.IdentifierExpression).Value, "console")
	assert_string(t, member.RightHandSide.Value, "log")

	input = []byte("package main const test = console.")
	ast, err = parser.Parse(input, "test.mb")

	assert_error(t, err)

	input = []byte("package main const test = .log")
	ast, err = parser.Parse(input, "test.mb")

	assert_error(t, err)

	input = []byte("package main const test = .log()")
	ast, err = parser.Parse(input, "test.mb")

	assert_error(t, err)

	input = []byte("package main const test = console . log")
	ast, err = parser.Parse(input, "test.mb")

	assert_error(t, err)
}
