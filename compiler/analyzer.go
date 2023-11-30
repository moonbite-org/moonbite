package compiler

import (
	"fmt"

	"github.com/moonbite-org/moonbite/parser"
)

type Analyzer struct {
	Program parser.Ast
}

func (a Analyzer) Analyze() {
	for _, definition := range a.Program.Definitions {
		a.resolve_statement(definition)
	}
}

func (a Analyzer) resolve_statement(s parser.Statement) {
	switch s.Kind() {
	case parser.DeclarationStatementKind:
		a.resolve_declaration_statement(s.(parser.DeclarationStatement))
	default:
		panic("unknown statement kind")
	}
}

func (a Analyzer) resolve_declaration_statement(s parser.DeclarationStatement) {
	fmt.Println(s.VarKind, a.resolve_identifier_expression(s.Name), s.Type, s.Hidden, *s.Value)
}

func (a Analyzer) resolve_expression(e parser.Expression) {
	switch e.Kind() {
	case parser.IdentifierExpressionKind:
		a.resolve_identifier_expression(e.(parser.IdentifierExpression))
	default:
		panic("unknown statement kind")
	}
}

func (a Analyzer) resolve_identifier_expression(e parser.IdentifierExpression) string {
	return e.Value
}

func (a Analyzer) resolve_literal_expression(e parser.LiteralExpression) parser.LiteralExpression {
	return e
}

func check_type(t1, t2 parser.TypeLiteral, v *parser.Expression) bool {
	return false
}
