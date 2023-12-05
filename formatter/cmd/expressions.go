package formatter

import (
	"fmt"

	parser "github.com/moonbite-org/moonbite/parser/cmd"
)

func FormatExpression(e parser.Expression) string {
	switch e.Kind() {
	case parser.StringLiteralExpressionKind:
		return FormatStringLiteralExpression(e.(parser.StringLiteralExpression))
	case parser.RuneLiteralExpressionKind:
		return FormatRuneLiteralExpression(e.(parser.RuneLiteralExpression))
	case parser.BoolLiteralExpressionKind:
		return FormatBoolLiteralExpression(e.(parser.BoolLiteralExpression))
	case parser.NumberLiteralExpressionKind:
		return FormatNumberLiteralExpression(e.(parser.NumberLiteralExpression))
	case parser.IdentifierExpressionKind:
		return FormatIdentifierLiteralExpression(e.(parser.IdentifierExpression))
	}

	return "unknown expression"
}

func FormatStringLiteralExpression(e parser.StringLiteralExpression) string {
	return fmt.Sprintf("\"%s\"", e.Value)
}

func FormatRuneLiteralExpression(e parser.RuneLiteralExpression) string {
	return fmt.Sprintf("\"%c\"", e.Value)
}

func FormatBoolLiteralExpression(e parser.BoolLiteralExpression) string {
	if e.Value {
		return "true"
	}

	return "false"
}

func FormatNumberLiteralExpression(e parser.NumberLiteralExpression) string {
	return fmt.Sprintf("%+v", e.Value.Value)
}

func FormatIdentifierLiteralExpression(e parser.IdentifierExpression) string {
	return fmt.Sprintf("%+v", e.Value)
}
