package action

import (
	"fmt"
	"sort"
	"strings"

	"github.com/fugalang/fugu/internal/ast"
	"github.com/fugalang/fugu/internal/diagnostics/errors"
	"github.com/fugalang/fugu/internal/token"
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
	NodeKind ast.NodeType
	Error    errors.Error
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
	return ActionStruct{Typ: Shift, State: state, Error: errors.Errors[6]}
}

// state использовать чтобы понять сколько надо снять элементов
func Red(state int, n ast.NodeType) ActionStruct {
	return ActionStruct{Typ: Reduce, State: state, NodeKind: n, Error: errors.Errors[6]}
}
func Acc() ActionStruct               { return ActionStruct{Typ: Accept, State: 0, Error: errors.Errors[6]} }
func Err(e errors.Error) ActionStruct { return ActionStruct{Typ: Error, State: 0, Error: e} }

func BuildActionSlice(src *map[int]map[token.Kind]ActionStruct, gotoSrc *map[int]map[ast.NodeType]int) *table {
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
					actions = append(actions, ActionStruct{Typ: Error, Error: errors.Errors[6]})
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
	if gotoSrc != nil {
		for state, row := range *gotoSrc {
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

func expandMap(m map[token.Kind]ActionStruct) map[token.Kind]ActionStruct {
	fullRow := make(map[token.Kind]ActionStruct, len(m))
	keys := make([]int, 0, len(m))

	for tk, act := range m {
		fullRow[tk] = act
		keys = append(keys, int(tk))
	}

	sort.Ints(keys)
	for _, key := range keys {
		tk := token.Kind(key)
		act := m[tk]

		for _, tk := range token.Expand(tk) {
			if _, ok := fullRow[tk]; !ok {
				fullRow[tk] = act
			}
		}
	}

	return fullRow
}

func GenerateActionTable(src *map[int]map[token.Kind]ActionStruct) string {
	table := BuildActionSlice(src, GotoMap)

	var out strings.Builder

	// w := func(s string) {
	//	out.WriteString(s)
	//}
	wn := func(s string) {
		out.WriteString(s)
		out.WriteString("\n")
	}
	wf := func(format string, args ...any) {
		out.WriteString(fmt.Sprintf(format, args...))
	}

	wn("// Code generated by fugu. DO NOT EDIT.")
	wn("package action")
	wn("")
	wn("import (")
	wn("\t\"github.com/fugalang/fugu/internal/ast\"")
	wn("\t\"github.com/fugalang/fugu/internal/diagnostics/errors\"")
	wn("\t. \"github.com/fugalang/fugu/internal/token\"")

	wn(")")
	wn("")

	wn(fmt.Sprintf("var ActionTokenCount int = %d", int(token.EndToken)))
	wn("")

	wn("var Actions = []ActionStruct{")
	for i, act := range table.actions {
		wf("\t%v, // %d\n", renderActionType(act), i)
	}
	wn("}")
	wn("")

	wn("var Check = []int{")
	for i, c := range table.check {
		wf("\t%d, // %d\n", c, i)
	}
	wn("}")
	wn("")

	wn("var Base = []int{")
	for i, b := range table.base {
		wf("\t%d, // state %d\n", b, i)
	}
	wn("}")
	wn("")

	wn("var Gotos = []int{")
	for i, g := range table.gotos {
		wf("\t%d, // %d\n", g, i)
	}
	wn("}")
	wn("")

	wn("var GotoCheck = []int{")
	for i, c := range table.gotoCheck {
		wf("\t%d, // %d\n", c, i)
	}
	wn("}")
	wn("")

	wn("var GotoBase = []int{")
	for i, b := range table.gotoBase {
		wf("\t%d, // state %d\n", b, i)
	}
	wn("}")
	wn("")

	wn("func Action(state int, tk Kind) ActionStruct {")
	wn("\tif state < 0 || state >= len(Base) {")
	wn("\t\treturn Err(errors.Errors[6])")
	wn("\t}")
	wn("\tif tk <= 0 || int(tk) >= ActionTokenCount {")
	wn("\t\treturn Err(errors.Errors[6])")
	wn("\t}")
	wn("\tb := Base[state]")
	wn("\tif b < 0 {")
	wn("\t\treturn Err(errors.Errors[6])")
	wn("\t}")
	wn("\tidx := b + int(tk)")
	wn("\tif idx >= 0 && idx < len(Actions) && Check[idx] == state {")
	wn("\t\treturn Actions[idx]")
	wn("\t}")
	wn("\treturn Err(errors.Errors[6])")
	wn("}")
	wn("")

	wn("func Goto(state int, k ast.NodeType) int {")
	wn("\tif state < 0 || state >= len(GotoBase) {")
	wn("\t\treturn -1")
	wn("\t}")
	wn("\tb := GotoBase[state]")
	wn("\tif b < 0 {")
	wn("\t\treturn -1")
	wn("\t}")
	wn("\tidx := b + int(k)")
	wn("\tif idx >= 0 && idx < len(Gotos) && GotoCheck[idx] == state {")
	wn("\t\treturn Gotos[idx]")
	wn("\t}")
	wn("\treturn -1")
	wn("}")

	return out.String()
}

func renderActionType(a ActionStruct) string {
	switch a.Typ {
	case Shift:
		return fmt.Sprintf("Sh(%d)", a.State)
	case Reduce:
		return fmt.Sprintf("Red(%d, ast.%s)", a.State, a.NodeKind.String())
	case Accept:
		return "Acc()"
	case Error:
		return fmt.Sprintf("Err(errors.Errors[%d])", +a.Error.Code)
	}
	return "Err(errors.Errors[6])"
}
