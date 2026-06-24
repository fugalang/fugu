package gram

import (
	"fmt"
	"strings"
)

func (t *Table) setAction(state int, tk string, next Action, g *Grammar) error {
	if t.Actions[state] == nil {
		t.Actions[state] = make(map[string]Action)
	}
	prev, exists := t.Actions[state][tk]
	if !exists || sameAction(prev, next) {
		t.Actions[state][tk] = next
		return nil
	}

	resolved, ok := resolveConflict(prev, next, tk, g)
	if !ok {
		return fmt.Errorf("грамматика: конфликт state=%d token=%s: %s / %s", state, tk, describeAction(prev), describeAction(next))
	}
	if resolved.Kind != 0 {
		t.Actions[state][tk] = resolved
		return nil
	}
	delete(t.Actions[state], tk)
	return nil
}

func sameAction(a, b Action) bool {
	return a.Kind == b.Kind && a.State == b.State && a.Prod.ID == b.Prod.ID
}

func resolveConflict(a, b Action, tk string, g *Grammar) (Action, bool) {
	if a.Kind == ActionReduce && b.Kind == ActionShift {
		return resolveShiftReduce(b, a, tk, g)
	}
	if a.Kind == ActionShift && b.Kind == ActionReduce {
		return resolveShiftReduce(a, b, tk, g)
	}
	return Action{}, false
}

func resolveShiftReduce(shift Action, reduce Action, tk string, g *Grammar) (Action, bool) {
	tokenPrec, okToken := g.precedence[tk]
	prodPrec, okProd := productionPrec(reduce.Prod, g)
	if !okToken || !okProd {
		return Action{}, false
	}
	if tokenPrec.Level > prodPrec.Level {
		return shift, true
	}
	if tokenPrec.Level < prodPrec.Level {
		return reduce, true
	}
	switch tokenPrec.Assoc {
	case AssocLeft:
		return reduce, true
	case AssocRight:
		return shift, true
	case AssocNonassoc:
		return Action{}, true
	default:
		return Action{}, false
	}
}

func productionPrec(p Production, g *Grammar) (Prec, bool) {
	if p.PrecBy != "" {
		prec, ok := g.precedence[p.PrecBy]
		return prec, ok
	}
	for i := len(p.RHS) - 1; i >= 0; i-- {
		if prec, ok := g.precedence[p.RHS[i]]; ok {
			return prec, true
		}
	}
	return Prec{}, false
}

func describeAction(a Action) string {
	switch a.Kind {
	case ActionShift:
		return fmt.Sprintf("сдвиг %d", a.State)
	case ActionReduce:
		return fmt.Sprintf("свёртка %s -> %s", a.Prod.LHS, strings.Join(a.Prod.RHS, " "))
	case ActionAccept:
		return "успех"
	default:
		return "неизвестное действие"
	}
}
