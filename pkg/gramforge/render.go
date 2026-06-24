package gram

import (
	"fmt"
	"strings"
)

type RenderOptions struct {
	PackageName   string
	Imports       []string
	ActionMapName string
	GotoMapName   string
	TokenType     string
	NodeType      string
	NodePrefix    string
}

func RenderMaps(t *Table, opt RenderOptions) (string, error) {
	if t == nil {
		return "", fmt.Errorf("грамматика: RenderMaps получил nil-таблицу (Build вернул ошибку?)")
	}
	opt = fillRenderOptions(opt)

	for _, row := range t.Actions {
		for _, act := range row {
			if act.Kind == ActionReduce && act.Prod.Node == "" {
				return "", fmt.Errorf("грамматика: свёртка без node kind у %s", act.Prod.LHS)
			}
		}
	}

	var out strings.Builder
	out.WriteString("// Код сгенерирован генератором грамматики; руками не править.\n")
	out.WriteString("package " + opt.PackageName + "\n\n")
	if len(opt.Imports) > 0 {
		out.WriteString("import (\n")
		for _, imp := range opt.Imports {
			out.WriteString("\t" + imp + "\n")
		}
		out.WriteString(")\n\n")
	}

	out.WriteString(fmt.Sprintf("var %s = &map[int]map[%s]ActionStruct{\n", opt.ActionMapName, opt.TokenType))
	for _, state := range sortedIntKeys(t.Actions) {
		row := t.Actions[state]
		out.WriteString(fmt.Sprintf("\t%d: {\n", state))
		for _, tk := range sortTokens(row, t.Grammar) {
			out.WriteString(fmt.Sprintf("\t\t%s: %s,\n", tk, renderAction(row[tk], opt)))
		}
		out.WriteString("\t},\n")
	}
	out.WriteString("}\n\n")

	out.WriteString(fmt.Sprintf("var %s = &map[int]map[%s]int{\n", opt.GotoMapName, opt.NodeType))
	for _, state := range sortedIntKeys(t.Gotos) {
		row := t.Gotos[state]
		out.WriteString(fmt.Sprintf("\t%d: {\n", state))
		for _, nt := range sortedGotoKeys(row) {
			out.WriteString(fmt.Sprintf("\t\t%s%s: %d,\n", opt.NodePrefix, nt, row[nt]))
		}
		out.WriteString("\t},\n")
	}
	out.WriteString("}\n")

	return out.String(), nil
}

func fillRenderOptions(opt RenderOptions) RenderOptions {
	if opt.PackageName == "" {
		opt.PackageName = "action"
	}
	if opt.ActionMapName == "" {
		opt.ActionMapName = "ActionMap"
	}
	if opt.GotoMapName == "" {
		opt.GotoMapName = "GotoMap"
	}
	if opt.TokenType == "" {
		opt.TokenType = "Kind"
	}
	if opt.NodeType == "" {
		opt.NodeType = "ast.NodeType"
	}
	if opt.NodePrefix == "" {
		opt.NodePrefix = "ast."
	}
	return opt
}

func renderAction(act Action, opt RenderOptions) string {
	switch act.Kind {
	case ActionShift:
		return fmt.Sprintf("Sh(%d)", act.State)
	case ActionReduce:
		return fmt.Sprintf("Red(%d, %s%s)", len(act.Prod.RHS), opt.NodePrefix, act.Prod.Node)
	case ActionAccept:
		return "Acc()"
	default:
		return "Err(errors.Error[6])"
	}
}
