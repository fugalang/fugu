// Код сгенерирован генератором грамматики; руками не править.
package action

import (
	"fugu/pkg/parser/action/ast"
	. "fugu/pkg/token"
)

var ActionSrc = &map[int]map[TokenKind]ActionStruct{
	0: {
		GLITERAL: Sh(2),
		L_PAREN:  Sh(7),
	},
	1: {
		EOF: Acc(),
	},
	2: {
		R_PAREN:  Red(1, ast.KindLiteral),
		INCREASE: Red(1, ast.KindLiteral),
		DECREASE: Red(1, ast.KindLiteral),
		MULTIPLY: Red(1, ast.KindLiteral),
		DIVIDE:   Red(1, ast.KindLiteral),
		DEGREE:   Red(1, ast.KindLiteral),
		EOF:      Red(1, ast.KindLiteral),
	},
	3: {
		R_PAREN:  Red(1, ast.Expr),
		INCREASE: Sh(9),
		DECREASE: Sh(8),
		EOF:      Red(1, ast.Expr),
	},
	4: {
		R_PAREN:  Red(1, ast.KindMultiplicativeExpr),
		INCREASE: Red(1, ast.KindMultiplicativeExpr),
		DECREASE: Red(1, ast.KindMultiplicativeExpr),
		MULTIPLY: Red(1, ast.KindMultiplicativeExpr),
		DIVIDE:   Red(1, ast.KindMultiplicativeExpr),
		EOF:      Red(1, ast.KindMultiplicativeExpr),
	},
	5: {
		R_PAREN:  Red(1, ast.KindAdditiveExpr),
		INCREASE: Red(1, ast.KindAdditiveExpr),
		DECREASE: Red(1, ast.KindAdditiveExpr),
		MULTIPLY: Sh(11),
		DIVIDE:   Sh(10),
		EOF:      Red(1, ast.KindAdditiveExpr),
	},
	6: {
		R_PAREN:  Red(1, ast.KindDegreeExpr),
		INCREASE: Red(1, ast.KindDegreeExpr),
		DECREASE: Red(1, ast.KindDegreeExpr),
		MULTIPLY: Red(1, ast.KindDegreeExpr),
		DIVIDE:   Red(1, ast.KindDegreeExpr),
		DEGREE:   Sh(12),
		EOF:      Red(1, ast.KindDegreeExpr),
	},
	7: {
		GLITERAL: Sh(2),
		L_PAREN:  Sh(7),
	},
	8: {
		GLITERAL: Sh(2),
		L_PAREN:  Sh(7),
	},
	9: {
		GLITERAL: Sh(2),
		L_PAREN:  Sh(7),
	},
	10: {
		GLITERAL: Sh(2),
		L_PAREN:  Sh(7),
	},
	11: {
		GLITERAL: Sh(2),
		L_PAREN:  Sh(7),
	},
	12: {
		GLITERAL: Sh(2),
		L_PAREN:  Sh(7),
	},
	13: {
		R_PAREN: Sh(19),
	},
	14: {
		R_PAREN:  Red(3, ast.KindAdditiveExpr),
		INCREASE: Red(3, ast.KindAdditiveExpr),
		DECREASE: Red(3, ast.KindAdditiveExpr),
		MULTIPLY: Sh(11),
		DIVIDE:   Sh(10),
		EOF:      Red(3, ast.KindAdditiveExpr),
	},
	15: {
		R_PAREN:  Red(3, ast.KindAdditiveExpr),
		INCREASE: Red(3, ast.KindAdditiveExpr),
		DECREASE: Red(3, ast.KindAdditiveExpr),
		MULTIPLY: Sh(11),
		DIVIDE:   Sh(10),
		EOF:      Red(3, ast.KindAdditiveExpr),
	},
	16: {
		R_PAREN:  Red(3, ast.KindMultiplicativeExpr),
		INCREASE: Red(3, ast.KindMultiplicativeExpr),
		DECREASE: Red(3, ast.KindMultiplicativeExpr),
		MULTIPLY: Red(3, ast.KindMultiplicativeExpr),
		DIVIDE:   Red(3, ast.KindMultiplicativeExpr),
		EOF:      Red(3, ast.KindMultiplicativeExpr),
	},
	17: {
		R_PAREN:  Red(3, ast.KindMultiplicativeExpr),
		INCREASE: Red(3, ast.KindMultiplicativeExpr),
		DECREASE: Red(3, ast.KindMultiplicativeExpr),
		MULTIPLY: Red(3, ast.KindMultiplicativeExpr),
		DIVIDE:   Red(3, ast.KindMultiplicativeExpr),
		EOF:      Red(3, ast.KindMultiplicativeExpr),
	},
	18: {
		R_PAREN:  Red(3, ast.KindDegreeExpr),
		INCREASE: Red(3, ast.KindDegreeExpr),
		DECREASE: Red(3, ast.KindDegreeExpr),
		MULTIPLY: Red(3, ast.KindDegreeExpr),
		DIVIDE:   Red(3, ast.KindDegreeExpr),
		EOF:      Red(3, ast.KindDegreeExpr),
	},
	19: {
		R_PAREN:  Red(3, ast.KindParenExpr),
		INCREASE: Red(3, ast.KindParenExpr),
		DECREASE: Red(3, ast.KindParenExpr),
		MULTIPLY: Red(3, ast.KindParenExpr),
		DIVIDE:   Red(3, ast.KindParenExpr),
		DEGREE:   Red(3, ast.KindParenExpr),
		EOF:      Red(3, ast.KindParenExpr),
	},
}

var GotoSrc = &map[int]map[ast.NodeKind]int{
	0: {
		ast.Expr:                   1,
		ast.KindAdditiveExpr:       3,
		ast.KindDegreeExpr:         4,
		ast.KindMultiplicativeExpr: 5,
		ast.KindPrimaryExpr:        6,
	},
	7: {
		ast.Expr:                   13,
		ast.KindAdditiveExpr:       3,
		ast.KindDegreeExpr:         4,
		ast.KindMultiplicativeExpr: 5,
		ast.KindPrimaryExpr:        6,
	},
	8: {
		ast.KindDegreeExpr:         4,
		ast.KindMultiplicativeExpr: 14,
		ast.KindPrimaryExpr:        6,
	},
	9: {
		ast.KindDegreeExpr:         4,
		ast.KindMultiplicativeExpr: 15,
		ast.KindPrimaryExpr:        6,
	},
	10: {
		ast.KindDegreeExpr:  16,
		ast.KindPrimaryExpr: 6,
	},
	11: {
		ast.KindDegreeExpr:  17,
		ast.KindPrimaryExpr: 6,
	},
	12: {
		ast.KindDegreeExpr:  18,
		ast.KindPrimaryExpr: 6,
	},
}
