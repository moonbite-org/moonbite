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

	input = []byte("package main var String test = 0")
	ast, err = parser.Parse(input, "test.mb")

	assert_no_error(t, err)
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
