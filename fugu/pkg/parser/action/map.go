package action

import (
	"fugu/pkg/parser/action/ast"
	. "fugu/pkg/token"
)

var ActionSrc = &map[int]map[TokenKind]ActionStruct{
	// state 0 (start)
	0: {
		GLITERAL: Sh(1),
		L_PAREN:  Sh(2),
	},

	// state 1: Factor (literal)
	1: {
		INCREASE: Red(1, ast.KindLiteral),
		DECREASE: Red(1, ast.KindLiteral),

		MULTIPLY:  Red(1, ast.KindLiteral),
		DIVIDE:    Red(1, ast.KindLiteral),
		REMAINDER: Red(1, ast.KindLiteral),
		DEGREE:    Red(1, ast.KindLiteral),

		R_PAREN: Red(1, ast.KindLiteral),
		EOF:     Red(1, ast.KindLiteral),
	},

	// state 2: "(" Expr
	2: {
		GLITERAL: Sh(1),
		L_PAREN:  Sh(2),
	},

	// state 3: Expr (after + -)
	3: {
		INCREASE: Red(3, ast.KindAdditiveExpr),
		DECREASE: Red(3, ast.KindAdditiveExpr),

		MULTIPLY:  Sh(4),
		DIVIDE:    Sh(4),
		REMAINDER: Sh(4),

		DEGREE: Sh(5),

		R_PAREN: Red(3, ast.KindAdditiveExpr),
		EOF:     Red(3, ast.KindAdditiveExpr),
	},

	// state 4: Term (* / %)
	4: {
		INCREASE: Red(3, ast.KindMultiplicativeExpr),
		DECREASE: Red(3, ast.KindMultiplicativeExpr),

		MULTIPLY:  Sh(4),
		DIVIDE:    Sh(4),
		REMAINDER: Sh(4),

		DEGREE: Sh(5),

		R_PAREN: Red(3, ast.KindMultiplicativeExpr),
		EOF:     Red(3, ast.KindMultiplicativeExpr),
	},

	// state 5: Power (^)
	5: {
		INCREASE: Red(3, ast.KindPowerExpr),
		DECREASE: Red(3, ast.KindPowerExpr),

		MULTIPLY:  Red(3, ast.KindPowerExpr),
		DIVIDE:    Red(3, ast.KindPowerExpr),
		REMAINDER: Red(3, ast.KindPowerExpr),

		DEGREE: Sh(5),

		R_PAREN: Red(3, ast.KindPowerExpr),
		EOF:     Red(3, ast.KindPowerExpr),
	},
}

var GotoSrc = &map[int]map[ast.NodeKind]int{

	0: {
		ast.KindLiteral: 1,
	},

	1: {
		ast.KindAdditiveExpr: 3,
	},

	2: {
		ast.KindLiteral:      1,
		ast.KindAdditiveExpr: 3,
	},

	3: {
		ast.KindMultiplicativeExpr: 4,
	},

	4: {
		ast.KindPowerExpr: 5,
	},
}
