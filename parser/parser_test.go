package parser_test

import (
	"reflect"
	"testing"

	"github.com/moonbite-org/moonbite/parser"
)

func assert_identifier(t *testing.T, given parser.IdentifierExpression, expected string) {
	if given.Value != expected {
		t.Errorf("expected identifier to be %s but found %s", expected, given.Value)
	}
}

func assert_no_error(t *testing.T, err parser.Error) {
	if err.Exists {
		t.Errorf("expected no error but got: %s", err)
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
