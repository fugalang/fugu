// ! DO NOT EDIT
package action

import (
	"fugu/pkg/diagnostics"
	"fugu/pkg/parser/action/ast"
	. "fugu/pkg/token"
)

const ActionTokenCount = int(EndToken)

var Actions = []ActionStruct{
	Err(diagnostics.NoError), // 0
	Err(diagnostics.NoError), // 1
	Err(diagnostics.NoError), // 2
	Err(diagnostics.NoError), // 3
	Err(diagnostics.NoError), // 4
	Sh(1),                    // 5
	Acc(),                    // 6
}

var Check = []int{
	-1, // 0
	-1, // 1
	-1, // 2
	-1, // 3
	-1, // 4
	0,  // 5
	1,  // 6
}

var Base = []int{
	0, // state 0
	1, // state 1
}

var Gotos = []int{}

var GotoCheck = []int{}

var GotoBase = []int{
	-1, // state 0
	-1, // state 1
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

func Goto(state int, k ast.NodeKind) int {
	if state < 0 || state >= len(GotoBase) {
		return -1
	}
	b := GotoBase[state]
	if b < 0 {
		return -1
	}
	idx := b + int(k)
	if idx >= 0 && idx < len(Gotos) && GotoCheck[idx] == state {
		return Gotos[idx]
	}
	return -1
}
