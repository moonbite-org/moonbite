package formatter

import (
	"fmt"

	parser "github.com/moonbite-org/moonbite/parser/cmd"
)

func FormatStatement(s parser.Statement) string {
	switch s.Kind() {
	case parser.PackageStatementKind:
		return FormatPackageStatement(s.(parser.PackageStatement))
	case parser.UseStatementKind:
		return FormatUseStatement(s.(parser.UseStatement))
	case parser.ReturnStatementKind:
		return FormatReturnStatement(s.(parser.ReturnStatement))
	case parser.DeclarationStatementKind:
		return FormatDeclarationStatement(s.(parser.DeclarationStatement))
	}
	return "unknown statement"
}

func FormatPackageStatement(s parser.PackageStatement) string {
	return fmt.Sprintf("package %s", FormatExpression(s.Name))
}

func FormatUseStatement(s parser.UseStatement) string {
	as := ""

	if s.As != nil {
		as = fmt.Sprintf(" as %s", FormatExpression(*s.As))
	}

	return fmt.Sprintf("use %s%s", FormatExpression(s.Resource), as)
}

func FormatReturnStatement(s parser.ReturnStatement) string {
	return fmt.Sprintf("return%s", empty_expression(s.Value, true, false))
}

func FormatDeclarationStatement(s parser.DeclarationStatement) string {
	kind := ""
	value := ""

	if s.VarKind == parser.VariableKind {
		kind = "var"
	} else {
		kind = "const"
	}

	if s.Value != nil {
		value = fmt.Sprintf(" = %s", FormatExpression(*s.Value))
	}

	return fmt.Sprintf("%s%s%s%s", kind, empty_type(s.Type, true, true), FormatExpression(s.Name), value)
}
