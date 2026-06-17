//! DO NOT EDIT
package action

import (
	"fugu/pkg/parser/action/ast"
	"fugu/pkg/reporter"
	. "fugu/pkg/token"
)

var Actions = []ActionStruct{
	Err(reporter.NoError), // 0
	Err(reporter.NoError), // 1
	Err(reporter.NoError), // 2
	Err(reporter.NoError), // 3
	Err(reporter.NoError), // 4
	Red(3, ast.KindAdditiveExpr), // 5
	Red(1, ast.KindLiteral), // 6
	Err(reporter.NoError), // 7
	Sh(1), // 8
	Err(reporter.NoError), // 9
	Red(1, ast.KindLiteral), // 10
	Red(4, ast.KindMultiplicativeExpr), // 11
	Sh(1), // 12
	Err(reporter.NoError), // 13
	Err(reporter.NoError), // 14
	Err(reporter.NoError), // 15
	Err(reporter.NoError), // 16
	Red(5, ast.KindPowerExpr), // 17
	Err(reporter.NoError), // 18
	Err(reporter.NoError), // 19
	Err(reporter.NoError), // 20
	Err(reporter.NoError), // 21
	Err(reporter.NoError), // 22
	Err(reporter.NoError), // 23
	Err(reporter.NoError), // 24
	Err(reporter.NoError), // 25
	Err(reporter.NoError), // 26
	Err(reporter.NoError), // 27
	Err(reporter.NoError), // 28
	Err(reporter.NoError), // 29
	Err(reporter.NoError), // 30
	Err(reporter.NoError), // 31
	Err(reporter.NoError), // 32
	Err(reporter.NoError), // 33
	Err(reporter.NoError), // 34
	Err(reporter.NoError), // 35
	Err(reporter.NoError), // 36
	Err(reporter.NoError), // 37
	Err(reporter.NoError), // 38
	Err(reporter.NoError), // 39
	Err(reporter.NoError), // 40
	Err(reporter.NoError), // 41
	Err(reporter.NoError), // 42
	Err(reporter.NoError), // 43
	Err(reporter.NoError), // 44
	Err(reporter.NoError), // 45
	Err(reporter.NoError), // 46
	Err(reporter.NoError), // 47
	Err(reporter.NoError), // 48
	Err(reporter.NoError), // 49
	Err(reporter.NoError), // 50
	Err(reporter.NoError), // 51
	Err(reporter.NoError), // 52
	Sh(2), // 53
	Red(3, ast.KindAdditiveExpr), // 54
	Red(1, ast.KindLiteral), // 55
	Err(reporter.NoError), // 56
	Sh(2), // 57
	Err(reporter.NoError), // 58
	Err(reporter.NoError), // 59
	Red(4, ast.KindMultiplicativeExpr), // 60
	Err(reporter.NoError), // 61
	Err(reporter.NoError), // 62
	Err(reporter.NoError), // 63
	Err(reporter.NoError), // 64
	Err(reporter.NoError), // 65
	Red(5, ast.KindPowerExpr), // 66
	Err(reporter.NoError), // 67
	Err(reporter.NoError), // 68
	Err(reporter.NoError), // 69
	Err(reporter.NoError), // 70
	Err(reporter.NoError), // 71
	Err(reporter.NoError), // 72
	Err(reporter.NoError), // 73
	Err(reporter.NoError), // 74
	Err(reporter.NoError), // 75
	Red(3, ast.KindAdditiveExpr), // 76
	Red(3, ast.KindAdditiveExpr), // 77
	Sh(4), // 78
	Sh(4), // 79
	Sh(4), // 80
	Sh(5), // 81
	Red(4, ast.KindMultiplicativeExpr), // 82
	Red(4, ast.KindMultiplicativeExpr), // 83
	Sh(4), // 84
	Sh(4), // 85
	Sh(4), // 86
	Sh(5), // 87
	Red(5, ast.KindPowerExpr), // 88
	Red(5, ast.KindPowerExpr), // 89
	Red(5, ast.KindPowerExpr), // 90
	Red(5, ast.KindPowerExpr), // 91
	Red(5, ast.KindPowerExpr), // 92
	Sh(5), // 93
}

var Check = []int{
	-1, // 0
	-1, // 1
	-1, // 2
	-1, // 3
	-1, // 4
	3, // 5
	1, // 6
	-1, // 7
	0, // 8
	-1, // 9
	1, // 10
	4, // 11
	2, // 12
	-1, // 13
	-1, // 14
	-1, // 15
	-1, // 16
	5, // 17
	-1, // 18
	-1, // 19
	-1, // 20
	-1, // 21
	-1, // 22
	-1, // 23
	-1, // 24
	-1, // 25
	-1, // 26
	-1, // 27
	-1, // 28
	-1, // 29
	-1, // 30
	-1, // 31
	-1, // 32
	-1, // 33
	-1, // 34
	-1, // 35
	-1, // 36
	-1, // 37
	-1, // 38
	-1, // 39
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
	0, // 53
	3, // 54
	1, // 55
	-1, // 56
	2, // 57
	-1, // 58
	-1, // 59
	4, // 60
	-1, // 61
	-1, // 62
	-1, // 63
	-1, // 64
	-1, // 65
	5, // 66
	-1, // 67
	-1, // 68
	-1, // 69
	-1, // 70
	-1, // 71
	-1, // 72
	-1, // 73
	-1, // 74
	-1, // 75
	3, // 76
	3, // 77
	3, // 78
	3, // 79
	3, // 80
	3, // 81
	4, // 82
	4, // 83
	4, // 84
	4, // 85
	4, // 86
	4, // 87
	5, // 88
	5, // 89
	5, // 90
	5, // 91
	5, // 92
	5, // 93
}

var Base = []int{
	0, // state 0
	1, // state 1
	4, // state 2
	0, // state 3
	6, // state 4
	12, // state 5
}

func Action(state int, tk TokenKind) ActionStruct {
	if state < 0 || state >= len(Base) {
		return Err(reporter.NoError)
	}
	b := Base[state]
	if b < 0 {
		return Err(reporter.NoError)
	}
	idx := b + int(tk)
	if idx >= 0 && idx < len(Actions) && Check[idx] == state {
		return Actions[idx]
	}
	return Err(reporter.StateDoesNotToken)
}
