// ! DO NOT EDIT
package action

import (
	"fugu/pkg/diagnostics"
	"fugu/pkg/parser/action/ast"
	. "fugu/pkg/token"
)

const ActionTokenCount = 99

var Actions = []ActionStruct{
	Err(diagnostics.NoError),           // 0
	Err(diagnostics.NoError),           // 1
	Err(diagnostics.NoError),           // 2
	Err(diagnostics.NoError),           // 3
	Err(diagnostics.NoError),           // 4
	Red(1, ast.KindLiteral),            // 5
	Sh(1),                              // 6
	Sh(1),                              // 7
	Sh(1),                              // 8
	Red(1, ast.KindLiteral),            // 9
	Sh(1),                              // 10
	Sh(1),                              // 11
	Sh(1),                              // 12
	Sh(1),                              // 13
	Sh(1),                              // 14
	Sh(1),                              // 15
	Sh(1),                              // 16
	Sh(1),                              // 17
	Sh(1),                              // 18
	Sh(1),                              // 19
	Sh(1),                              // 20
	Red(3, ast.KindAdditiveExpr),       // 21
	Sh(1),                              // 22
	Sh(1),                              // 23
	Sh(1),                              // 24
	Sh(1),                              // 25
	Sh(1),                              // 26
	Sh(1),                              // 27
	Sh(1),                              // 28
	Sh(1),                              // 29
	Err(diagnostics.NoError),           // 30
	Err(diagnostics.NoError),           // 31
	Err(diagnostics.NoError),           // 32
	Red(4, ast.KindMultiplicativeExpr), // 33
	Err(diagnostics.NoError),           // 34
	Err(diagnostics.NoError),           // 35
	Err(diagnostics.NoError),           // 36
	Err(diagnostics.NoError),           // 37
	Err(diagnostics.NoError),           // 38
	Red(5, ast.KindPowerExpr),          // 39
	Err(diagnostics.NoError),           // 40
	Err(diagnostics.NoError),           // 41
	Err(diagnostics.NoError),           // 42
	Err(diagnostics.NoError),           // 43
	Err(diagnostics.NoError),           // 44
	Err(diagnostics.NoError),           // 45
	Err(diagnostics.NoError),           // 46
	Err(diagnostics.NoError),           // 47
	Err(diagnostics.NoError),           // 48
	Err(diagnostics.NoError),           // 49
	Err(diagnostics.NoError),           // 50
	Err(diagnostics.NoError),           // 51
	Err(diagnostics.NoError),           // 52
	Sh(2),                              // 53
	Red(1, ast.KindLiteral),            // 54
	Err(diagnostics.NoError),           // 55
	Err(diagnostics.NoError),           // 56
	Err(diagnostics.NoError),           // 57
	Err(diagnostics.NoError),           // 58
	Err(diagnostics.NoError),           // 59
	Err(diagnostics.NoError),           // 60
	Err(diagnostics.NoError),           // 61
	Err(diagnostics.NoError),           // 62
	Err(diagnostics.NoError),           // 63
	Err(diagnostics.NoError),           // 64
	Sh(2),                              // 65
	Err(diagnostics.NoError),           // 66
	Err(diagnostics.NoError),           // 67
	Err(diagnostics.NoError),           // 68
	Err(diagnostics.NoError),           // 69
	Red(3, ast.KindAdditiveExpr),       // 70
	Err(diagnostics.NoError),           // 71
	Err(diagnostics.NoError),           // 72
	Err(diagnostics.NoError),           // 73
	Err(diagnostics.NoError),           // 74
	Err(diagnostics.NoError),           // 75
	Red(1, ast.KindLiteral),            // 76
	Red(1, ast.KindLiteral),            // 77
	Red(1, ast.KindLiteral),            // 78
	Red(1, ast.KindLiteral),            // 79
	Red(1, ast.KindLiteral),            // 80
	Red(1, ast.KindLiteral),            // 81
	Red(4, ast.KindMultiplicativeExpr), // 82
	Err(diagnostics.NoError),           // 83
	Err(diagnostics.NoError),           // 84
	Err(diagnostics.NoError),           // 85
	Err(diagnostics.NoError),           // 86
	Err(diagnostics.NoError),           // 87
	Red(5, ast.KindPowerExpr),          // 88
	Err(diagnostics.NoError),           // 89
	Err(diagnostics.NoError),           // 90
	Err(diagnostics.NoError),           // 91
	Red(3, ast.KindAdditiveExpr),       // 92
	Red(3, ast.KindAdditiveExpr),       // 93
	Sh(4),                              // 94
	Sh(4),                              // 95
	Sh(4),                              // 96
	Sh(5),                              // 97
	Err(diagnostics.NoError),           // 98
	Err(diagnostics.NoError),           // 99
	Err(diagnostics.NoError),           // 100
	Err(diagnostics.NoError),           // 101
	Err(diagnostics.NoError),           // 102
	Err(diagnostics.NoError),           // 103
	Red(4, ast.KindMultiplicativeExpr), // 104
	Red(4, ast.KindMultiplicativeExpr), // 105
	Sh(4),                              // 106
	Sh(4),                              // 107
	Sh(4),                              // 108
	Sh(5),                              // 109
	Red(5, ast.KindPowerExpr),          // 110
	Red(5, ast.KindPowerExpr),          // 111
	Red(5, ast.KindPowerExpr),          // 112
	Red(5, ast.KindPowerExpr),          // 113
	Red(5, ast.KindPowerExpr),          // 114
	Sh(5),                              // 115
}

var Check = []int{
	-1, // 0
	-1, // 1
	-1, // 2
	-1, // 3
	-1, // 4
	1,  // 5
	0,  // 6
	0,  // 7
	0,  // 8
	1,  // 9
	0,  // 10
	0,  // 11
	0,  // 12
	0,  // 13
	0,  // 14
	0,  // 15
	0,  // 16
	0,  // 17
	2,  // 18
	2,  // 19
	2,  // 20
	3,  // 21
	2,  // 22
	2,  // 23
	2,  // 24
	2,  // 25
	2,  // 26
	2,  // 27
	2,  // 28
	2,  // 29
	-1, // 30
	-1, // 31
	-1, // 32
	4,  // 33
	-1, // 34
	-1, // 35
	-1, // 36
	-1, // 37
	-1, // 38
	5,  // 39
	-1, // 40
	-1, // 41
	-1, // 42
	-1, // 43
	-1, // 44
	-1, // 45
	-1, // 46
	-1, // 47
	-1, // 48
	-1, // 49
	-1, // 50
	-1, // 51
	-1, // 52
	0,  // 53
	1,  // 54
	-1, // 55
	-1, // 56
	-1, // 57
	-1, // 58
	-1, // 59
	-1, // 60
	-1, // 61
	-1, // 62
	-1, // 63
	-1, // 64
	2,  // 65
	-1, // 66
	-1, // 67
	-1, // 68
	-1, // 69
	3,  // 70
	-1, // 71
	-1, // 72
	-1, // 73
	-1, // 74
	-1, // 75
	1,  // 76
	1,  // 77
	1,  // 78
	1,  // 79
	1,  // 80
	1,  // 81
	4,  // 82
	-1, // 83
	-1, // 84
	-1, // 85
	-1, // 86
	-1, // 87
	5,  // 88
	-1, // 89
	-1, // 90
	-1, // 91
	3,  // 92
	3,  // 93
	3,  // 94
	3,  // 95
	3,  // 96
	3,  // 97
	-1, // 98
	-1, // 99
	-1, // 100
	-1, // 101
	-1, // 102
	-1, // 103
	4,  // 104
	4,  // 105
	4,  // 106
	4,  // 107
	4,  // 108
	4,  // 109
	5,  // 110
	5,  // 111
	5,  // 112
	5,  // 113
	5,  // 114
	5,  // 115
}

var Base = []int{
	0,  // state 0
	0,  // state 1
	12, // state 2
	16, // state 3
	28, // state 4
	34, // state 5
}

func Action(state int, tk TokenKind) ActionStruct {
	if state < 0 || state >= len(Base) {
		return Err(diagnostics.NoError)
	}
	if tk <= 0 || int(tk) >= ActionTokenCount {
		return Err(diagnostics.StateDoesNotToken)
	}
	b := Base[state]
	if b < 0 {
		return Err(diagnostics.NoError)
	}
	idx := b + int(tk)
	if idx >= 0 && idx < len(Actions) && Check[idx] == state {
		return Actions[idx]
	}
	return Err(diagnostics.StateDoesNotToken)
}
