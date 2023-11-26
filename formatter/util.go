package formatter

import (
	"fmt"

	"github.com/moonbite-org/moonbite/parser"
)

func empty_expression(e *parser.Expression, pad_left, pad_right bool) string {
	pad_r := ""

	if pad_right {
		pad_r = " "
	}

	if e == nil {
		return pad_r
	}

	pad_l := ""

	if pad_left {
		pad_l = " "
	}

	return fmt.Sprintf("%s%s%s", pad_l, FormatExpression(*e), pad_r)
}

func empty_type(t *parser.TypeLiteral, pad_left, pad_right bool) string {
	pad_r := ""

	if pad_right {
		pad_r = " "
	}

	if t == nil {
		return pad_r
	}

	if t == nil {
		return ""
	}

	pad_l := ""

	if pad_left {
		pad_l = " "
	}

	return fmt.Sprintf("%s%s%s", pad_l, FormatTypeLiteral(*t), pad_r)
}
