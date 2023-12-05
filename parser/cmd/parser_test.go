package parser_test

import (
	"reflect"
	"testing"

	parser "github.com/moonbite-org/moonbite/parser/cmd"
	common "github.com/moonbite-org/moonbite/parser/common"
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

func assert_bool(t *testing.T, given, expexted bool) {
	if given != expexted {
		t.Errorf("expected bool to be %t but got %t", expexted, given)
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

func TestToken(t *testing.T) {
	token := parser.Token{
		Location:   common.Location{},
		Literal:    "",
		Raw:        "",
		Offset:     0,
		LineBreaks: 0,
	}

	assert_string(t, token.String(), "End of file()[0:0][0]")

	token = parser.Token{
		Kind:       parser.Whitespace,
		Location:   common.Location{},
		Literal:    "",
		Raw:        "",
		Offset:     0,
		LineBreaks: 0,
	}

	assert_string(t, token.String(), "ws[0:0][0]")

	token = parser.Token{
		Kind:       parser.Keyword,
		Location:   common.Location{},
		Literal:    "as",
		Raw:        "as",
		Offset:     0,
		LineBreaks: 0,
	}

	assert_string(t, token.String(), "<as>[0:0][0]")
}

func TestLexer(t *testing.T) {
	input := []byte("package main const s = \"")
	_, err := parser.Parse(input, "test.mb")

	assert_error(t, err)

	input = []byte("package main const r = '")
	_, err = parser.Parse(input, "test.mb")

	assert_error(t, err)

	input = []byte("package main const r = ''")
	_, err = parser.Parse(input, "test.mb")

	assert_error(t, err)

	input = []byte("package main const r = 'ab'")
	_, err = parser.Parse(input, "test.mb")

	assert_error(t, err)

	input = []byte("package main const r = 'a'")
	_, err = parser.Parse(input, "test.mb")

	assert_no_error(t, err)

	// input = []byte("package main const r = '\\''")
	// _, err = parser.Parse(input, "test.mb")

	// assert_no_error(t, err)

	input = []byte("package main const s = \"test\"")
	_, err = parser.Parse(input, "test.mb")

	assert_no_error(t, err)

	input = []byte("package main const s = \"te\\\\\"st\"")
	_, err = parser.Parse(input, "test.mb")

	assert_no_error(t, err)

	input = []byte(`package main
	const multi = line
	`)
	_, err = parser.Parse(input, "test.mb")

	assert_no_error(t, err)

	input = []byte("package main const multi = `line" + "\n" + "string`")
	_, err = parser.Parse(input, "test.mb")

	assert_no_error(t, err)

	input = []byte("package main const multi = `line")
	_, err = parser.Parse(input, "test.mb")

	assert_error(t, err)

	input = []byte("package main //comment")
	_, err = parser.Parse(input, "test.mb")

	assert_no_error(t, err)

	input = []byte("package main /* comment */")
	_, err = parser.Parse(input, "test.mb")

	assert_no_error(t, err)

	input = []byte("package main /* comment")
	_, err = parser.Parse(input, "test.mb")

	assert_error(t, err)

	input = []byte("package main /* comment" + "\n" + " end */")
	_, err = parser.Parse(input, "test.mb")

	assert_no_error(t, err)

	input = []byte("package main const n = 5")
	_, err = parser.Parse(input, "test.mb")

	assert_no_error(t, err)

	input = []byte("package main const n = 5.5")
	_, err = parser.Parse(input, "test.mb")

	assert_no_error(t, err)

	input = []byte("package main const n = 5.0e8")
	_, err = parser.Parse(input, "test.mb")

	assert_no_error(t, err)

	input = []byte("package main const n = 05")
	_, err = parser.Parse(input, "test.mb")

	assert_error(t, err)

	input = []byte("package main const n = 500")
	_, err = parser.Parse(input, "test.mb")

	assert_no_error(t, err)
}

func TestPackageStatement(t *testing.T) {
	input := []byte("")
	_, err := parser.Parse(input, "test.mb")

	assert_error(t, err)

	input = []byte("package main")
	ast, err := parser.Parse(input, "test.mb")

	assert_no_error(t, err)
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
	input := []byte("package main const test = data")
	ast, err := parser.Parse(input, "test.mb")

	assert_no_error(t, err)

	definition := ast.Definitions[0]

	assert_type(t, definition, parser.DeclarationStatement{})

	input = []byte("package main hidden const test = data")
	ast, err = parser.Parse(input, "test.mb")

	assert_no_error(t, err)

	definition = ast.Definitions[0]

	assert_type(t, definition, parser.DeclarationStatement{})
	assert_bool(t, definition.(parser.DeclarationStatement).Hidden, true)

	input = []byte("package main const test")
	ast, err = parser.Parse(input, "test.mb")

	assert_error(t, err)

	input = []byte("package main var String test")
	ast, err = parser.Parse(input, "test.mb")

	assert_no_error(t, err)

	definition = ast.Definitions[0]
	assert_type(t, definition, parser.DeclarationStatement{})

	if definition.(parser.DeclarationStatement).Value != nil {
		t.Errorf("expected expression to be nil but found %+v", definition.(parser.DeclarationStatement).Value)
	}

	input = []byte("package main var String test = []")
	ast, err = parser.Parse(input, "test.mb")

	assert_no_error(t, err)

	input = []byte("package main var String test = T{}")
	ast, err = parser.Parse(input, "test.mb")

	assert_no_error(t, err)
}

func TestAssignmentStatement(t *testing.T) {
	input := []byte(`package main
	fun main() {
		count = 0
		count += 10
		count -= 1
		count /= 3
		count *= 4
		count %= 10
	}
	`)
	ast, err := parser.Parse(input, "test.mb")

	assert_no_error(t, err)

	definition := ast.Definitions[0]
	body := definition.(*parser.UnboundFunDefinitionStatement).Body

	assert_type(t, body[0], parser.AssignmentStatement{})
	assert_string(t, body[0].(parser.AssignmentStatement).Operator, "=")
	assert_type(t, body[0].(parser.AssignmentStatement).LeftHandSide, parser.IdentifierExpression{})
	assert_string(t, body[0].(parser.AssignmentStatement).LeftHandSide.(parser.IdentifierExpression).Value, "count")
	assert_type(t, body[0].(parser.AssignmentStatement).RightHandSide, parser.NumberLiteralExpression{})
	assert_int(t, body[0].(parser.AssignmentStatement).RightHandSide.(parser.NumberLiteralExpression).Value.Value.(int), 0)

	assert_type(t, body[1], parser.AssignmentStatement{})
	assert_string(t, body[1].(parser.AssignmentStatement).Operator, "+=")
	assert_type(t, body[1].(parser.AssignmentStatement).LeftHandSide, parser.IdentifierExpression{})
	assert_string(t, body[1].(parser.AssignmentStatement).LeftHandSide.(parser.IdentifierExpression).Value, "count")
	assert_type(t, body[1].(parser.AssignmentStatement).RightHandSide, parser.NumberLiteralExpression{})
	assert_int(t, body[1].(parser.AssignmentStatement).RightHandSide.(parser.NumberLiteralExpression).Value.Value.(int), 10)

	assert_type(t, body[2], parser.AssignmentStatement{})
	assert_string(t, body[2].(parser.AssignmentStatement).Operator, "-=")
	assert_type(t, body[2].(parser.AssignmentStatement).LeftHandSide, parser.IdentifierExpression{})
	assert_string(t, body[2].(parser.AssignmentStatement).LeftHandSide.(parser.IdentifierExpression).Value, "count")
	assert_type(t, body[2].(parser.AssignmentStatement).RightHandSide, parser.NumberLiteralExpression{})
	assert_int(t, body[2].(parser.AssignmentStatement).RightHandSide.(parser.NumberLiteralExpression).Value.Value.(int), 1)

	assert_type(t, body[3], parser.AssignmentStatement{})
	assert_string(t, body[3].(parser.AssignmentStatement).Operator, "/=")
	assert_type(t, body[3].(parser.AssignmentStatement).LeftHandSide, parser.IdentifierExpression{})
	assert_string(t, body[3].(parser.AssignmentStatement).LeftHandSide.(parser.IdentifierExpression).Value, "count")
	assert_type(t, body[3].(parser.AssignmentStatement).RightHandSide, parser.NumberLiteralExpression{})
	assert_int(t, body[3].(parser.AssignmentStatement).RightHandSide.(parser.NumberLiteralExpression).Value.Value.(int), 3)

	assert_type(t, body[4], parser.AssignmentStatement{})
	assert_string(t, body[4].(parser.AssignmentStatement).Operator, "*=")
	assert_type(t, body[4].(parser.AssignmentStatement).LeftHandSide, parser.IdentifierExpression{})
	assert_string(t, body[4].(parser.AssignmentStatement).LeftHandSide.(parser.IdentifierExpression).Value, "count")
	assert_type(t, body[4].(parser.AssignmentStatement).RightHandSide, parser.NumberLiteralExpression{})
	assert_int(t, body[4].(parser.AssignmentStatement).RightHandSide.(parser.NumberLiteralExpression).Value.Value.(int), 4)

	assert_type(t, body[5], parser.AssignmentStatement{})
	assert_string(t, body[5].(parser.AssignmentStatement).Operator, "%=")
	assert_type(t, body[5].(parser.AssignmentStatement).LeftHandSide, parser.IdentifierExpression{})
	assert_string(t, body[5].(parser.AssignmentStatement).LeftHandSide.(parser.IdentifierExpression).Value, "count")
	assert_type(t, body[5].(parser.AssignmentStatement).RightHandSide, parser.NumberLiteralExpression{})
	assert_int(t, body[5].(parser.AssignmentStatement).RightHandSide.(parser.NumberLiteralExpression).Value.Value.(int), 10)
}

func TestTypeDefinitionStatement(t *testing.T) {
	input := []byte(`package main
	type String string
	type String implements [Observable<string>] Saturable<string>
	type Data {
		key typ;
		key2 typ2;
	}
	type Data string & bool
	type Data string | bool
	type Data string & (bool | int)
	type Data string("data")
	`)
	ast, err := parser.Parse(input, "test.mb")

	assert_no_error(t, err)

	definition := ast.Definitions[0].(parser.TypeDefinitionStatement)
	assert_type(t, definition, parser.TypeDefinitionStatement{})
	assert_type(t, definition.Name, parser.IdentifierExpression{})
	assert_string(t, definition.Name.Value, "String")
	assert_type(t, definition.Definition, parser.TypeIdentifier{})
	assert_type(t, definition.Definition.(parser.TypeIdentifier).Name, &parser.IdentifierExpression{})
	assert_string(t, definition.Definition.(parser.TypeIdentifier).Name.(*parser.IdentifierExpression).Value, "string")
	assert_int(t, len(definition.Generics), 0)
	assert_int(t, len(definition.Implementations), 0)

	definition = ast.Definitions[1].(parser.TypeDefinitionStatement)
	assert_type(t, definition, parser.TypeDefinitionStatement{})
	assert_type(t, definition.Name, parser.IdentifierExpression{})
	assert_string(t, definition.Name.Value, "String")
	assert_type(t, definition.Definition, parser.TypeIdentifier{})
	assert_type(t, definition.Definition.(parser.TypeIdentifier).Name, &parser.IdentifierExpression{})
	assert_string(t, definition.Definition.(parser.TypeIdentifier).Name.(*parser.IdentifierExpression).Value, "Saturable")
	assert_int(t, len(definition.Generics), 0)
	assert_int(t, len(definition.Implementations), 1)
	assert_type(t, definition.Implementations[0], parser.TypeIdentifier{})
	assert_type(t, definition.Implementations[0].Name, &parser.IdentifierExpression{})
	assert_string(t, definition.Implementations[0].Name.(*parser.IdentifierExpression).Value, "Observable")

	definition = ast.Definitions[2].(parser.TypeDefinitionStatement)
	assert_type(t, definition, parser.TypeDefinitionStatement{})
	assert_type(t, definition.Name, parser.IdentifierExpression{})
	assert_string(t, definition.Name.Value, "Data")
	assert_type(t, definition.Definition, parser.StructLiteral{})
	literal := definition.Definition.(parser.StructLiteral)

	assert_int(t, len(literal), 2)
	assert_type(t, literal[0], parser.ValueTypePair{})
	assert_string(t, literal[0].Key.Value, "key")
	assert_type(t, literal[0].Type, parser.TypeIdentifier{})
	assert_string(t, literal[0].Type.(parser.TypeIdentifier).Name.(*parser.IdentifierExpression).Value, "typ")
	assert_type(t, literal[1], parser.ValueTypePair{})
	assert_string(t, literal[1].Key.Value, "key2")
	assert_type(t, literal[1].Type, parser.TypeIdentifier{})
	assert_string(t, literal[1].Type.(parser.TypeIdentifier).Name.(*parser.IdentifierExpression).Value, "typ2")

	definition = ast.Definitions[3].(parser.TypeDefinitionStatement)
	assert_type(t, definition, parser.TypeDefinitionStatement{})
	assert_type(t, definition.Name, parser.IdentifierExpression{})
	assert_string(t, definition.Name.Value, "Data")
	assert_type(t, definition.Definition, parser.OperatedType{})
	operated := definition.Definition.(parser.OperatedType)

	assert_string(t, operated.Operator, "&")
	assert_type(t, operated.LeftHandSide, parser.TypeIdentifier{})
	assert_string(t, operated.LeftHandSide.(parser.TypeIdentifier).Name.(*parser.IdentifierExpression).Value, "string")
	assert_type(t, operated.RightHandSide, parser.TypeIdentifier{})
	assert_string(t, operated.RightHandSide.(parser.TypeIdentifier).Name.(*parser.IdentifierExpression).Value, "bool")

	definition = ast.Definitions[4].(parser.TypeDefinitionStatement)
	assert_type(t, definition, parser.TypeDefinitionStatement{})
	assert_type(t, definition.Name, parser.IdentifierExpression{})
	assert_string(t, definition.Name.Value, "Data")
	assert_type(t, definition.Definition, parser.OperatedType{})
	operated = definition.Definition.(parser.OperatedType)

	assert_string(t, operated.Operator, "|")
	assert_type(t, operated.LeftHandSide, parser.TypeIdentifier{})
	assert_string(t, operated.LeftHandSide.(parser.TypeIdentifier).Name.(*parser.IdentifierExpression).Value, "string")
	assert_type(t, operated.RightHandSide, parser.TypeIdentifier{})
	assert_string(t, operated.RightHandSide.(parser.TypeIdentifier).Name.(*parser.IdentifierExpression).Value, "bool")

	definition = ast.Definitions[5].(parser.TypeDefinitionStatement)
	assert_type(t, definition, parser.TypeDefinitionStatement{})
	assert_type(t, definition.Name, parser.IdentifierExpression{})
	assert_string(t, definition.Name.Value, "Data")
	assert_type(t, definition.Definition, parser.OperatedType{})
	operated = definition.Definition.(parser.OperatedType)

	assert_string(t, operated.Operator, "&")
	assert_type(t, operated.LeftHandSide, parser.TypeIdentifier{})
	assert_string(t, operated.LeftHandSide.(parser.TypeIdentifier).Name.(*parser.IdentifierExpression).Value, "string")
	assert_type(t, operated.RightHandSide, parser.GroupType{})
	assert_type(t, operated.RightHandSide.(parser.GroupType).Type, parser.OperatedType{})
	operated = operated.RightHandSide.(parser.GroupType).Type.(parser.OperatedType)
	assert_string(t, operated.Operator, "|")

	definition = ast.Definitions[6].(parser.TypeDefinitionStatement)
	assert_type(t, definition, parser.TypeDefinitionStatement{})
	assert_type(t, definition.Name, parser.IdentifierExpression{})
	assert_string(t, definition.Name.Value, "Data")
	assert_type(t, definition.Definition, parser.TypedLiteral{})
	typed_literal := definition.Definition.(parser.TypedLiteral)

	assert_type(t, typed_literal.Type, parser.TypeIdentifier{})
	assert_string(t, typed_literal.Type.Name.(*parser.IdentifierExpression).Value, "string")
	assert_type(t, typed_literal.Literal, parser.StringLiteralExpression{})
	assert_string(t, typed_literal.Literal.(parser.StringLiteralExpression).Value, "data")
}

func TestIfStatement(t *testing.T) {
	input := []byte(`package main
	fun main() {
		if (true) {
			count++
		}

		if (boolean) {
			count++
		}else {
			count--
		}

		if (boolean) {
			count++
		}else if (other) {
			count = 0
		}else {
			count--
		}
	}
	`)
	ast, err := parser.Parse(input, "test.mb")

	assert_no_error(t, err)

	body := ast.Definitions[0].(*parser.UnboundFunDefinitionStatement).Body
	assert_int(t, len(body), 3)

	assert_type(t, body[0], parser.IfStatement{})
	statement := body[0].(parser.IfStatement)
	assert_type(t, statement.MainBlock.Predicate, parser.BoolLiteralExpression{})
	assert_int(t, len(statement.MainBlock.Body), 1)
	assert_int(t, len(statement.ElseIfBlocks), 0)
	assert_int(t, len(statement.ElseBlock), 0)

	statement = body[1].(parser.IfStatement)
	assert_type(t, statement.MainBlock.Predicate, parser.IdentifierExpression{})
	assert_int(t, len(statement.MainBlock.Body), 1)
	assert_int(t, len(statement.ElseIfBlocks), 0)
	assert_int(t, len(statement.ElseBlock), 1)

	statement = body[2].(parser.IfStatement)
	assert_type(t, statement.MainBlock.Predicate, parser.IdentifierExpression{})
	assert_int(t, len(statement.MainBlock.Body), 1)
	assert_int(t, len(statement.ElseIfBlocks), 1)
	assert_type(t, statement.ElseIfBlocks[0], parser.PredicateBlock{})
	assert_int(t, len(statement.ElseIfBlocks[0].Body), 1)
	assert_int(t, len(statement.ElseBlock), 1)
}

func TestTraitDefinitionStatement(t *testing.T) {
	input := []byte(`package main
	trait Greeter {
		fun greet(name String) String
	}

	hidden trait LobbyBoy<T Warning> mimics [Greeter, Repeater] {
		fun welcome(name String) String
	}
	`)
	ast, err := parser.Parse(input, "test.mb")

	assert_no_error(t, err)

	trait := ast.Definitions[0].(parser.TraitDefinitionStatement)
	assert_string(t, trait.Name.Value, "Greeter")
	assert_int(t, len(trait.Definition), 1)
	assert_int(t, len(trait.Generics), 0)
	assert_int(t, len(trait.Mimics), 0)
	assert_string(t, trait.Definition[0].Name.Value, "greet")

	trait = ast.Definitions[1].(parser.TraitDefinitionStatement)
	assert_bool(t, trait.Hidden, true)
	assert_string(t, trait.Name.Value, "LobbyBoy")
	assert_int(t, len(trait.Definition), 1)
	assert_string(t, trait.Definition[0].Name.Value, "welcome")
	assert_int(t, len(trait.Generics), 1)
	assert_int(t, len(trait.Mimics), 2)
	assert_string(t, (*trait.Definition[0].ReturnType).(parser.TypeIdentifier).Name.(*parser.IdentifierExpression).Value, "String")
}

func TestLoopStatement(t *testing.T) {
	input := []byte(`package main
	fun main() {
		for (true) {}

		for (key, value of iterator) {}

		for (, value of iterator) {}

		for (var i = 0; i < 10; i++) {}
	}
	`)
	ast, err := parser.Parse(input, "test.mb")

	assert_no_error(t, err)

	body := ast.Definitions[0].(*parser.UnboundFunDefinitionStatement).Body

	assert_type(t, body[0], parser.LoopStatement{})
	assert_type(t, body[0].(parser.LoopStatement).Predicate, parser.UnipartiteLoopPredicate{})

	assert_type(t, body[1], parser.LoopStatement{})
	assert_type(t, body[1].(parser.LoopStatement).Predicate, parser.BipartiteLoopPredicate{})

	assert_type(t, body[2], parser.LoopStatement{})
	assert_type(t, body[2].(parser.LoopStatement).Predicate, parser.BipartiteLoopPredicate{})
	if body[2].(parser.LoopStatement).Predicate.(parser.BipartiteLoopPredicate).Key != nil {
		t.Errorf("expected key to be nil but found %+v", body[2].(parser.LoopStatement).Predicate.(parser.BipartiteLoopPredicate).Key)
	}
	assert_string(t, body[2].(parser.LoopStatement).Predicate.(parser.BipartiteLoopPredicate).Value.Value, "value")
	assert_type(t, body[2].(parser.LoopStatement).Predicate.(parser.BipartiteLoopPredicate).Iterator, parser.IdentifierExpression{})
	assert_string(t, body[2].(parser.LoopStatement).Predicate.(parser.BipartiteLoopPredicate).Iterator.(parser.IdentifierExpression).Value, "iterator")

	assert_type(t, body[3], parser.LoopStatement{})
	assert_type(t, body[3].(parser.LoopStatement).Predicate, parser.TripartiteLoopPredicate{})
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

	input = []byte("package main const test = 2 - 3")
	ast, err = parser.Parse(input, "test.mb")
	assert_no_error(t, err)

	definition = ast.Definitions[0]
	assert_type(t, *definition.(parser.DeclarationStatement).Value, parser.ArithmeticExpression{})
	expression = (*definition.(parser.DeclarationStatement).Value).(parser.ArithmeticExpression)
	assert_string(t, expression.Operator, "-")

	input = []byte("package main const test = 2 / 3")
	ast, err = parser.Parse(input, "test.mb")
	assert_no_error(t, err)

	definition = ast.Definitions[0]
	assert_type(t, *definition.(parser.DeclarationStatement).Value, parser.ArithmeticExpression{})
	expression = (*definition.(parser.DeclarationStatement).Value).(parser.ArithmeticExpression)
	assert_string(t, expression.Operator, "/")

	input = []byte("package main const test = 2 % 3")
	ast, err = parser.Parse(input, "test.mb")
	assert_no_error(t, err)

	definition = ast.Definitions[0]
	assert_type(t, *definition.(parser.DeclarationStatement).Value, parser.ArithmeticExpression{})
	expression = (*definition.(parser.DeclarationStatement).Value).(parser.ArithmeticExpression)
	assert_string(t, expression.Operator, "%")
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

func TestIndexExpression(t *testing.T) {
	input := []byte("package main const test = list[0]")
	ast, err := parser.Parse(input, "test.mb")

	assert_no_error(t, err)

	definition := ast.Definitions[0]
	assert_type(t, *definition.(parser.DeclarationStatement).Value, parser.IndexExpression{})
	expression := (*definition.(parser.DeclarationStatement).Value).(parser.IndexExpression)

	assert_type(t, expression.Host, parser.IdentifierExpression{})
	assert_type(t, expression.Index, parser.NumberLiteralExpression{})

	input = []byte("package main const test = data.list[count()]")
	ast, err = parser.Parse(input, "test.mb")

	assert_no_error(t, err)

	definition = ast.Definitions[0]
	assert_type(t, *definition.(parser.DeclarationStatement).Value, parser.IndexExpression{})
	expression = (*definition.(parser.DeclarationStatement).Value).(parser.IndexExpression)

	assert_type(t, expression.Host, parser.MemberExpression{})
	assert_type(t, expression.Index, parser.CallExpression{})

	member := expression.Host.(parser.MemberExpression)
	assert_string(t, member.LeftHandSide.(parser.IdentifierExpression).Value, "data")
	assert_string(t, member.RightHandSide.Value, "list")

	index := expression.Index.(parser.CallExpression)
	assert_type(t, index.Callee, parser.IdentifierExpression{})
	assert_string(t, index.Callee.(parser.IdentifierExpression).Value, "count")
}

func TestTypeCastExpression(t *testing.T) {
	input := []byte("package main const test = data.(String)")
	ast, err := parser.Parse(input, "test.mb")

	assert_no_error(t, err)

	definition := ast.Definitions[0]
	assert_type(t, *definition.(parser.DeclarationStatement).Value, parser.TypeCastExpression{})
	expression := (*definition.(parser.DeclarationStatement).Value).(parser.TypeCastExpression)

	assert_type(t, expression.Value, parser.IdentifierExpression{})
	assert_string(t, expression.Value.(parser.IdentifierExpression).Value, "data")
	assert_type(t, expression.Type, parser.TypeIdentifier{})
	assert_type(t, expression.Type.Name, &parser.IdentifierExpression{})
	assert_string(t, expression.Type.Name.(*parser.IdentifierExpression).Value, "String")
}

func TestCaretExpression(t *testing.T) {
	input := []byte("package main const test = ^")
	ast, err := parser.Parse(input, "test.mb")

	assert_no_error(t, err)

	definition := ast.Definitions[0]
	assert_type(t, *definition.(parser.DeclarationStatement).Value, parser.CaretExpression{})

	input = []byte("package main const test = data.^")
	ast, err = parser.Parse(input, "test.mb")

	assert_error(t, err)

	input = []byte("package main const test = ^ instanceof Warning")
	ast, err = parser.Parse(input, "test.mb")

	assert_no_error(t, err)

	definition = ast.Definitions[0]
	assert_type(t, *definition.(parser.DeclarationStatement).Value, parser.InstanceofExpression{})
	assert_type(t, (*definition.(parser.DeclarationStatement).Value).(parser.InstanceofExpression).LeftHandSide, parser.CaretExpression{})
}

func TestInstanceofExpression(t *testing.T) {
	input := []byte("package main const test = data instanceof typ")
	ast, err := parser.Parse(input, "test.mb")

	assert_no_error(t, err)

	definition := ast.Definitions[0]
	assert_type(t, *definition.(parser.DeclarationStatement).Value, parser.InstanceofExpression{})
	expression := (*definition.(parser.DeclarationStatement).Value).(parser.InstanceofExpression)

	assert_type(t, expression.LeftHandSide, parser.IdentifierExpression{})
	assert_type(t, expression.RightHandSide, parser.TypeIdentifier{})

	assert_string(t, expression.LeftHandSide.(parser.IdentifierExpression).Value, "data")
	assert_string(t, expression.RightHandSide.(parser.TypeIdentifier).Name.(*parser.IdentifierExpression).Value, "typ")
}

func TestMatchSelfExpression(t *testing.T) {
	input := []byte("package main const test = .")
	_, err := parser.Parse(input, "test.mb")

	assert_error(t, err)

	input = []byte(`package main
	fun main() {
		match (data) {
			(.) {}
			(.value) {}
			((.).(String)) {}
			base {
				// comment
			}
		}
	}
	`)
	ast, err := parser.Parse(input, "test.mb")

	assert_no_error(t, err)

	body := ast.Definitions[0].(*parser.UnboundFunDefinitionStatement).Body
	match := body[0].(parser.ExpressionStatement).Expression.(parser.MatchExpression)

	assert_type(t, match.Blocks[0].Predicate, parser.MatchSelfExpression{})

	assert_type(t, match.Blocks[1].Predicate, parser.MemberExpression{})
	assert_type(t, match.Blocks[1].Predicate.(parser.MemberExpression).LeftHandSide, parser.MatchSelfExpression{})
	assert_type(t, match.Blocks[1].Predicate.(parser.MemberExpression).RightHandSide, parser.IdentifierExpression{})
	assert_string(t, match.Blocks[1].Predicate.(parser.MemberExpression).RightHandSide.Value, "value")

	assert_type(t, match.Blocks[2].Predicate, parser.TypeCastExpression{})
	assert_type(t, match.Blocks[2].Predicate.(parser.TypeCastExpression).Value.(parser.GroupExpression).Expression, parser.MatchSelfExpression{})

	assert_int(t, len(match.BaseBlock), 0)
}

func TestGroupExpression(t *testing.T) {
	input := []byte("package main const test = (identifier)")
	ast, err := parser.Parse(input, "test.mb")

	assert_no_error(t, err)

	definition := ast.Definitions[0]
	assert_type(t, *definition.(parser.DeclarationStatement).Value, parser.GroupExpression{})
	expression := (*definition.(parser.DeclarationStatement).Value).(parser.GroupExpression)

	assert_type(t, expression.Expression, parser.IdentifierExpression{})
	assert_string(t, expression.Expression.(parser.IdentifierExpression).Value, "identifier")

	input = []byte("package main const test = (identifier")
	ast, err = parser.Parse(input, "test.mb")

	assert_error(t, err)
}

func TestThisExpression(t *testing.T) {
	input := []byte(`package main
	fun for Data test() {
		return this
	}
	`)
	ast, err := parser.Parse(input, "test.mb")

	assert_no_error(t, err)

	definition := ast.Definitions[0]
	assert_type(t, definition, &parser.BoundFunDefinitionStatement{})

	fun := definition.(*parser.BoundFunDefinitionStatement)
	ret := fun.Body[0].(parser.ReturnStatement)

	assert_type(t, *ret.Value, parser.ThisExpression{})

	input = []byte(`package main
	fun test() {
		return this
	}
	`)
	ast, err = parser.Parse(input, "test.mb")

	assert_error(t, err)
}

func TestArithmeticUnaryExpression(t *testing.T) {
	input := []byte(`package main 
	fun main() {
		index++
		index--
	}`)
	ast, err := parser.Parse(input, "test.mb")

	assert_no_error(t, err)

	definition := ast.Definitions[0]
	assert_type(t, definition, &parser.UnboundFunDefinitionStatement{})

	fun := definition.(*parser.UnboundFunDefinitionStatement)

	var expression parser.ArithmeticUnaryExpression = fun.Body[0].(parser.ExpressionStatement).Expression.(parser.ArithmeticUnaryExpression)

	assert_type(t, expression, parser.ArithmeticUnaryExpression{})
	assert_type(t, expression.Operation, parser.IncrementKind)
	assert_bool(t, expression.Pre, false)

	expression = fun.Body[1].(parser.ExpressionStatement).Expression.(parser.ArithmeticUnaryExpression)

	assert_type(t, expression, parser.ArithmeticUnaryExpression{})
	assert_type(t, expression.Operation, parser.DecrementKind)
	assert_bool(t, expression.Pre, false)

	input = []byte(`package main
	fun main() {
		test()++
	}`)
	ast, err = parser.Parse(input, "test.mb")

	assert_error(t, err)

	input = []byte(`package main
	fun main() {
		index++ index index-- index
	}`)
	ast, err = parser.Parse(input, "test.mb")

	assert_no_error(t, err)
}

func TestFunExpression(t *testing.T) {
	input := []byte(`package main 
	const test = fun() {

	}
	`)
	ast, err := parser.Parse(input, "test.mb")

	assert_no_error(t, err)

	definition := ast.Definitions[0]
	assert_type(t, definition, parser.DeclarationStatement{})

	decl := definition.(parser.DeclarationStatement)
	assert_type(t, *decl.Value, parser.AnonymousFunExpression{})

	input = []byte(`package main 
	const test = fun<T, K>() {

	}
	`)
	ast, err = parser.Parse(input, "test.mb")

	assert_no_error(t, err)

	definition = ast.Definitions[0]
	assert_type(t, definition, parser.DeclarationStatement{})

	decl = definition.(parser.DeclarationStatement)
	assert_type(t, *decl.Value, parser.AnonymousFunExpression{})
	assert_int(t, len((*decl.Value).(parser.AnonymousFunExpression).Signature.Generics), 2)
	generics := (*decl.Value).(parser.AnonymousFunExpression).Signature.Generics

	assert_type(t, generics[0].Name, &parser.IdentifierExpression{})
	assert_string(t, generics[0].Name.(*parser.IdentifierExpression).Value, "T")

	assert_type(t, generics[1].Name, &parser.IdentifierExpression{})
	assert_string(t, generics[1].Name.(*parser.IdentifierExpression).Value, "K")

	input = []byte(`package main 
	const test = fun(data Int) String {
		var count = data
		count++
		console.log(count)
	}
	`)
	ast, err = parser.Parse(input, "test.mb")

	assert_no_error(t, err)

	definition = ast.Definitions[0]
	decl = definition.(parser.DeclarationStatement)
	assert_type(t, definition, parser.DeclarationStatement{})

	expression := (*decl.Value).(parser.AnonymousFunExpression)
	assert_int(t, len(expression.Signature.Generics), 0)
	assert_int(t, len(expression.Signature.Parameters), 1)
	assert_int(t, len(expression.Body), 3)

	assert_string(t, expression.Signature.Parameters[0].Name.Value, "data")
	assert_type(t, *expression.Signature.ReturnType, parser.TypeIdentifier{})
	assert_string(t, (*expression.Signature.ReturnType).(parser.TypeIdentifier).Name.(*parser.IdentifierExpression).Value, "String")

	input = []byte(`package main 
	const test = fun(data Int) String {
	`)
	ast, err = parser.Parse(input, "test.mb")

	assert_error(t, err)

	input = []byte(`package main 
	const test = fun(data) String {}
	`)
	ast, err = parser.Parse(input, "test.mb")

	assert_error(t, err)

	input = []byte(`package main 
	const test = fun( String {}
	`)
	ast, err = parser.Parse(input, "test.mb")

	assert_error(t, err)

	input = []byte(`package main 
	const test = fun(): String {}
	`)
	ast, err = parser.Parse(input, "test.mb")

	assert_error(t, err)

	input = []byte(`package main 
	const test = fun()
	`)
	ast, err = parser.Parse(input, "test.mb")

	assert_error(t, err)
}

func TestOrExpression(t *testing.T) {
	input := []byte(`package main 
	const test = read_file() or 0
	const test = read_file() or giveup
	`)
	ast, err := parser.Parse(input, "test.mb")

	assert_no_error(t, err)

	definition := ast.Definitions[0]
	decl := definition.(parser.DeclarationStatement)

	assert_type(t, *decl.Value, parser.OrExpression{})
	expression := (*decl.Value).(parser.OrExpression)

	assert_type(t, expression.LeftHandSide, parser.CallExpression{})
	assert_type(t, expression.RightHandSide, parser.NumberLiteralExpression{})

	definition = ast.Definitions[1]
	decl = definition.(parser.DeclarationStatement)

	assert_type(t, *decl.Value, parser.OrExpression{})
	expression = (*decl.Value).(parser.OrExpression)

	assert_type(t, expression.LeftHandSide, parser.CallExpression{})
	assert_type(t, expression.RightHandSide, parser.GiveupExpression{})

	input = []byte(`package main
	const test = read_file or 0
	`)
	ast, err = parser.Parse(input, "test.mb")

	assert_error(t, err)

	input = []byte(`package main
	const test = 2 + 2 or 0
	`)
	ast, err = parser.Parse(input, "test.mb")

	assert_error(t, err)
}

func TestNotExpression(t *testing.T) {
	input := []byte(`package main 
	const test = !is_admin
	const test = !(2 * 2)
	`)
	ast, err := parser.Parse(input, "test.mb")

	assert_no_error(t, err)

	definition := ast.Definitions[0]
	decl := definition.(parser.DeclarationStatement)

	assert_type(t, *decl.Value, parser.NotExpression{})
	expression := (*decl.Value).(parser.NotExpression)

	assert_type(t, expression.Expression, parser.IdentifierExpression{})
	assert_string(t, expression.Expression.(parser.IdentifierExpression).Value, "is_admin")

	definition = ast.Definitions[1]
	decl = definition.(parser.DeclarationStatement)

	assert_type(t, *decl.Value, parser.NotExpression{})
	expression = (*decl.Value).(parser.NotExpression)

	assert_type(t, expression.Expression, parser.GroupExpression{})
	assert_type(t, expression.Expression.(parser.GroupExpression).Expression, parser.ArithmeticExpression{})
}
