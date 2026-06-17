package action

import (
	"fugu/pkg/parser/action/ast"
	. "fugu/pkg/token"
)

var ActionSrc = &map[int]map[TokenKind]ActionStruct{
	0: { // стартовое
		GLITERAL: Sh(1),
		L_PAREN:  Sh(2),
	},

	1: { // после литерала
		GARITHMETIC: Red(1, ast.KindLiteral),
		R_PAREN:     Red(1, ast.KindLiteral),
		EOF:         Red(1, ast.KindLiteral),
	},

	2: { // после (
		GLITERAL: Sh(1),
		L_PAREN:  Sh(2),
	},

	3: { // после Выражение + - || eof
		INCREASE:  Red(3, ast.KindAdditiveExpr),
		DECREASE:  Red(3, ast.KindAdditiveExpr),
		MULTIPLY:  Sh(4),
		DIVIDE:    Sh(4),
		REMAINDER: Sh(4),
		DEGREE:    Sh(5),
		R_PAREN:   Red(3, ast.KindAdditiveExpr),
		EOF:       Red(3, ast.KindAdditiveExpr),
	},

	4: { // после * / % все трое на одном уровне
		INCREASE:  Red(4, ast.KindMultiplicativeExpr),
		DECREASE:  Red(4, ast.KindMultiplicativeExpr),
		MULTIPLY:  Sh(4),
		DIVIDE:    Sh(4),
		REMAINDER: Sh(4),
		DEGREE:    Sh(5),
		R_PAREN:   Red(4, ast.KindMultiplicativeExpr),
		EOF:       Red(4, ast.KindMultiplicativeExpr),
	},

	5: {
		INCREASE:  Red(5, ast.KindPowerExpr),
		DECREASE:  Red(5, ast.KindPowerExpr),
		MULTIPLY:  Red(5, ast.KindPowerExpr),
		DIVIDE:    Red(5, ast.KindPowerExpr),
		REMAINDER: Red(5, ast.KindPowerExpr),
		DEGREE:    Sh(5),
		R_PAREN:   Red(5, ast.KindPowerExpr),
		EOF:       Red(5, ast.KindPowerExpr),
	},
}

// ((2 + 4) - 3) + 2 * (1)
