//go:generate go run generate.go
package action

import (
	"fmt"
	"fugu/pkg/diagnostics"
	"fugu/pkg/parser/action/ast"
	"fugu/pkg/token"
	. "fugu/pkg/token"
	"sort"
	"strings"
)

type ActionType int

const (
	_ ActionType = iota

	Shift  // Сдвиг: поместить текущий токен в стек и перейти в новое состояние
	Reduce // Свёртка: заменить последовательность символов по правилу грамматики
	Accept // Успех: входная цепочка полностью распознана
	Error  // Ошибка: для текущего состояния и токена действие не определено
)

type ActionStruct struct {
	Typ      ActionType
	State    int
	NodeKind ast.NodeKind
	ErrCode  diagnostics.Code
}

type table struct {
	actions   []ActionStruct // действия
	check     []int          // проверка состояния
	base      []int          // смещения для состояний
	gotos     []int          // переходы по нетерминалам
	gotoCheck []int          // проверка для goto
	gotoBase  []int          // смещения для goto
}

// конструкторы
func Sh(state int) ActionStruct {
	return ActionStruct{Typ: Shift, State: state, ErrCode: diagnostics.NoError}
}

// state использовать чтобы понять сколько надо снять элементов
func Red(state int, n ast.NodeKind) ActionStruct {
	return ActionStruct{Typ: Reduce, State: state, NodeKind: n, ErrCode: diagnostics.NoError}
}
func Acc() ActionStruct                   { return ActionStruct{Typ: Accept, State: 0, ErrCode: diagnostics.NoError} }
func Err(e diagnostics.Code) ActionStruct { return ActionStruct{Typ: Error, State: 0, ErrCode: e} }

func BuildActionSlice(src *map[int]map[TokenKind]ActionStruct) *table {
	ms := 0
	for s := range *src {
		if s > ms {
			ms = s
		}
	}
	states := ms + 1

	type couple struct {
		tk  int
		act ActionStruct
	}
	rows := make([][]couple, states)
	for state, row := range *src {
		fullRow := expandMap(row)
		entries := make([]couple, 0, len(fullRow))
		for tk, act := range fullRow {
			entries = append(entries, couple{int(tk), act})
		}
		sort.Slice(entries, func(i, j int) bool { return entries[i].tk < entries[j].tk })
		rows[state] = entries
	}

	var actions []ActionStruct
	var check []int
	base := make([]int, states)
	for i := range base {
		base[i] = -1
	}

	isFree := func(idx int) bool {
		return idx >= len(check) || check[idx] == -1
	}

	order := make([]int, 0, states)
	for i := 0; i < states; i++ {
		if len(rows[i]) > 0 {
			order = append(order, i)
		}
	}
	sort.Slice(order, func(i, j int) bool { return len(rows[order[i]]) > len(rows[order[j]]) })

	for _, state := range order {
		entries := rows[state]
		for b := 0; ; b++ {
			ok := true
			for _, e := range entries {
				if !isFree(b + e.tk) {
					ok = false
					break
				}
			}
			if ok {
				base[state] = b
				maxIdx := b + entries[len(entries)-1].tk
				for len(actions) <= maxIdx {
					actions = append(actions, ActionStruct{Typ: Error, ErrCode: diagnostics.NoError})
					check = append(check, -1)
				}
				for _, e := range entries {
					idx := b + e.tk
					actions[idx] = e.act
					check[idx] = state
				}
				break
			}
		}
	}

	type gpair struct {
		sym int
		to  int
	}

	grows := make([][]gpair, states)
	if GotoSrc != nil {
		for state, row := range *GotoSrc {
			entries := make([]gpair, 0, len(row))
			for k, to := range row {
				entries = append(entries, gpair{int(k), to})
			}
			sort.Slice(entries, func(i, j int) bool { return entries[i].sym < entries[j].sym })
			if state < len(grows) {
				grows[state] = entries
			}
		}
	}

	gotoBase := make([]int, states)
	for i := range gotoBase {
		gotoBase[i] = -1
	}

	gotos := make([]int, 0)
	gotoCheck := make([]int, 0)

	isFreeGoto := func(idx int) bool {
		return idx >= len(gotoCheck) || gotoCheck[idx] == -1
	}

	gorder := make([]int, 0, states)
	for i := 0; i < states; i++ {
		if len(grows[i]) > 0 {
			gorder = append(gorder, i)
		}
	}
	sort.Slice(gorder, func(i, j int) bool { return len(grows[gorder[i]]) > len(grows[gorder[j]]) })

	for _, state := range gorder {
		entries := grows[state]
		for b := 0; ; b++ {
			ok := true
			for _, e := range entries {
				if !isFreeGoto(b + e.sym) {
					ok = false
					break
				}
			}
			if ok {
				gotoBase[state] = b
				maxIdx := b + entries[len(entries)-1].sym
				for len(gotos) <= maxIdx {
					gotos = append(gotos, -1)
					gotoCheck = append(gotoCheck, -1)
				}
				for _, e := range entries {
					idx := b + e.sym
					gotos[idx] = e.to
					gotoCheck[idx] = state
				}
				break
			}
		}
	}

	return &table{
		actions:   actions,
		check:     check,
		base:      base,
		gotos:     gotos,
		gotoCheck: gotoCheck,
		gotoBase:  gotoBase,
	}
}

func expandMap(m map[TokenKind]ActionStruct) map[TokenKind]ActionStruct {
	fullRow := make(map[TokenKind]ActionStruct, len(m))
	keys := make([]int, 0, len(m))

	for tk, act := range m {
		fullRow[tk] = act
		keys = append(keys, int(tk))
	}

	sort.Ints(keys)
	for _, key := range keys {
		tk := TokenKind(key)
		act := m[tk]

		for _, tk := range token.Expand(tk) {
			if _, ok := fullRow[tk]; !ok {
				fullRow[tk] = act
			}
		}
	}

	return fullRow
}

func GenerateActionTable(src *map[int]map[TokenKind]ActionStruct) string {
	table := BuildActionSlice(src)

	var out strings.Builder
	out.WriteString("//! DO NOT EDIT\n")
	out.WriteString("package action\n\n")
	out.WriteString("import (\n")
	out.WriteString("\t\"fugu/pkg/parser/action/ast\"\n")
	out.WriteString("\t\"fugu/pkg/diagnostics\"\n")
	out.WriteString("\t. \"fugu/pkg/token\"\n")
	out.WriteString(")\n\n")
	out.WriteString(fmt.Sprintf("const ActionTokenCount = int(EndToken)\n\n"))
	out.WriteString("var Actions = []ActionStruct{\n")
	for i, act := range table.actions {
		out.WriteString(fmt.Sprintf("\t%v, // %d\n", renderActionType(act), i))
	}
	out.WriteString("}\n\n")
	out.WriteString("var Check = []int{\n")
	for i, c := range table.check {
		out.WriteString(fmt.Sprintf("\t%d, // %d\n", c, i))
	}
	out.WriteString("}\n\n")
	out.WriteString("var Base = []int{\n")
	for i, b := range table.base {
		out.WriteString(fmt.Sprintf("\t%d, // state %d\n", b, i))
	}
	out.WriteString("}\n\n")
	out.WriteString("var Gotos = []int{\n")
	for i, g := range table.gotos {
		out.WriteString(fmt.Sprintf("\t%d, // %d\n", g, i))
	}
	out.WriteString("}\n\n")
	out.WriteString("var GotoCheck = []int{\n")
	for i, c := range table.gotoCheck {
		out.WriteString(fmt.Sprintf("\t%d, // %d\n", c, i))
	}
	out.WriteString("}\n\n")
	out.WriteString("var GotoBase = []int{\n")
	for i, b := range table.gotoBase {
		out.WriteString(fmt.Sprintf("\t%d, // state %d\n", b, i))
	}
	out.WriteString("}\n\n")
	out.WriteString("func Action(state int, tk TokenKind) ActionStruct {\n")
	out.WriteString("\tif state < 0 || state >= len(Base) {\n")
	out.WriteString("\t\treturn Err(diagnostics.NoError)\n")
	out.WriteString("\t}\n")
	out.WriteString("\tif tk <= 0 || int(tk) >= ActionTokenCount {\n")
	out.WriteString("\t\treturn Err(diagnostics.StateDoesNotToken)\n")
	out.WriteString("\t}\n")
	out.WriteString("\tb := Base[state]\n")
	out.WriteString("\tif b < 0 {\n")
	out.WriteString("\t\treturn Err(diagnostics.NoError)\n")
	out.WriteString("\t}\n")
	out.WriteString("\tidx := b + int(tk)\n")
	out.WriteString("\tif idx >= 0 && idx < len(Actions) && Check[idx] == state {\n")
	out.WriteString("\t\treturn Actions[idx]\n")
	out.WriteString("\t}\n")
	out.WriteString("\treturn Err(diagnostics.StateDoesNotToken)\n")
	out.WriteString("}\n")
	out.WriteString("\n")
	out.WriteString("func Goto(state int, k ast.NodeKind) int {\n")
	out.WriteString("\tif state < 0 || state >= len(GotoBase) {\n")
	out.WriteString("\t\treturn -1\n")
	out.WriteString("\t}\n")
	out.WriteString("\tb := GotoBase[state]\n")
	out.WriteString("\tif b < 0 {\n")
	out.WriteString("\t\treturn -1\n")
	out.WriteString("\t}\n")
	out.WriteString("\tidx := b + int(k)\n")
	out.WriteString("\tif idx >= 0 && idx < len(Gotos) && GotoCheck[idx] == state {\n")
	out.WriteString("\t\treturn Gotos[idx]\n")
	out.WriteString("\t}\n")
	out.WriteString("\treturn -1\n")
	out.WriteString("}\n")
	return out.String()
}

func renderActionType(a ActionStruct) string {
	switch a.Typ {
	case Shift:
		return fmt.Sprintf("Sh(%d)", a.State)
	case Reduce:
		return fmt.Sprintf("Red(%d, ast.%s)", a.State, a.NodeKind)
	case Accept:
		return "Acc()"
	case Error:
		return fmt.Sprintf("Err(%s)", "diagnostics."+a.ErrCode.Code())
	}
	return "Err(diagnostics.NoError)"
}
