package gram

import "fmt"

type transitionKey struct {
	state int
	sym   string
}

func Build(g *Grammar) (*Table, error) {
	if g.Start == "" {
		return nil, fmt.Errorf("грамматика: нет %%start и нет правил")
	}
	if len(g.productions) == 0 {
		return nil, fmt.Errorf("грамматика: нет правил")
	}

	g.completeTerminals()

	augStart := g.Start + "'"
	for g.isNonTerminal(augStart) || g.Terminals[augStart] {
		augStart += "'"
	}

	prods := make([]Production, 0, len(g.productions)+1)
	prods = append(prods, Production{ID: 0, LHS: augStart, RHS: []string{g.Start}, Node: augStart})
	for _, p := range g.productions {
		p.ID = len(prods)
		prods = append(prods, p)
	}

	follow := buildFollow(g, prods, augStart)
	states, transitions := buildLR0(g, prods)

	table := &Table{
		Grammar: g,
		States:  states,
		Actions: make(map[int]map[string]Action),
		Gotos:   make(map[int]map[string]int),
	}

	for stateID, st := range states {
		for _, it := range st.Items {
			p := prods[it.Prod]
			if it.Dot < len(p.RHS) {
				sym := p.RHS[it.Dot]
				to, ok := transitions[transitionKey{stateID, sym}]
				if !ok {
					continue
				}
				if g.isTerminal(sym) {
					if err := table.setAction(stateID, sym, Action{Kind: ActionShift, State: to}, g); err != nil {
						return nil, err
					}
				} else if sym != augStart {
					putGoto(table.Gotos, stateID, sym, to)
				}
				continue
			}

			if p.LHS == augStart {
				if err := table.setAction(stateID, EOF, Action{Kind: ActionAccept}, g); err != nil {
					return nil, err
				}
				continue
			}

			for _, lookahead := range sortedSet(follow[p.LHS]) {
				if err := table.setAction(stateID, lookahead, Action{Kind: ActionReduce, Prod: p}, g); err != nil {
					return nil, err
				}
			}
		}
	}

	return table, nil
}

func buildFollow(g *Grammar, prods []Production, augStart string) map[string]map[string]bool {
	first, nullable := buildFirst(g, prods)
	follow := make(map[string]map[string]bool)
	for nt := range g.nonTerminals() {
		follow[nt] = make(map[string]bool)
	}
	if follow[g.Start] == nil {
		follow[g.Start] = make(map[string]bool)
	}
	follow[g.Start][EOF] = true

	changed := true
	for changed {
		changed = false
		for _, p := range prods {
			if p.LHS == augStart {
				continue
			}
			for i, sym := range p.RHS {
				if !g.isNonTerminal(sym) {
					continue
				}

				restFirst, restNullable := firstOfSeq(g, p.RHS[i+1:], first, nullable)
				for tk := range restFirst {
					if !follow[sym][tk] {
						follow[sym][tk] = true
						changed = true
					}
				}
				if restNullable {
					for tk := range follow[p.LHS] {
						if !follow[sym][tk] {
							follow[sym][tk] = true
							changed = true
						}
					}
				}
			}
		}
	}

	return follow
}

func buildFirst(g *Grammar, prods []Production) (map[string]map[string]bool, map[string]bool) {
	first := make(map[string]map[string]bool)
	nullable := make(map[string]bool)
	for _, p := range prods {
		if first[p.LHS] == nil {
			first[p.LHS] = make(map[string]bool)
		}
	}

	changed := true
	for changed {
		changed = false
		for _, p := range prods {
			if len(p.RHS) == 0 {
				if !nullable[p.LHS] {
					nullable[p.LHS] = true
					changed = true
				}
				continue
			}

			allNullable := true
			for _, sym := range p.RHS {
				if g.isTerminal(sym) {
					if !first[p.LHS][sym] {
						first[p.LHS][sym] = true
						changed = true
					}
					allNullable = false
					break
				}

				for tk := range first[sym] {
					if !first[p.LHS][tk] {
						first[p.LHS][tk] = true
						changed = true
					}
				}
				if !nullable[sym] {
					allNullable = false
					break
				}
			}
			if allNullable && !nullable[p.LHS] {
				nullable[p.LHS] = true
				changed = true
			}
		}
	}

	return first, nullable
}

func firstOfSeq(g *Grammar, seq []string, first map[string]map[string]bool, nullable map[string]bool) (map[string]bool, bool) {
	out := make(map[string]bool)
	if len(seq) == 0 {
		return out, true
	}

	for _, sym := range seq {
		if g.isTerminal(sym) {
			out[sym] = true
			return out, false
		}
		for tk := range first[sym] {
			out[tk] = true
		}
		if !nullable[sym] {
			return out, false
		}
	}

	return out, true
}

func buildLR0(g *Grammar, prods []Production) ([]State, map[transitionKey]int) {
	byLHS := make(map[string][]int)
	for i, p := range prods {
		byLHS[p.LHS] = append(byLHS[p.LHS], i)
	}

	closure := func(items []Item) []Item {
		seen := make(map[Item]bool)
		queue := append([]Item(nil), items...)
		for len(queue) > 0 {
			it := queue[0]
			queue = queue[1:]
			if seen[it] {
				continue
			}
			seen[it] = true

			p := prods[it.Prod]
			if it.Dot >= len(p.RHS) {
				continue
			}
			sym := p.RHS[it.Dot]
			if !g.isNonTerminal(sym) {
				continue
			}
			for _, prodID := range byLHS[sym] {
				queue = append(queue, Item{Prod: prodID})
			}
		}
		return sortItems(seen)
	}

	gotoItems := func(items []Item, sym string) []Item {
		next := make([]Item, 0)
		for _, it := range items {
			p := prods[it.Prod]
			if it.Dot < len(p.RHS) && p.RHS[it.Dot] == sym {
				next = append(next, Item{Prod: it.Prod, Dot: it.Dot + 1})
			}
		}
		return closure(next)
	}

	states := []State{{Items: closure([]Item{{Prod: 0}})}}
	stateIDs := map[string]int{itemsKey(states[0].Items): 0}
	transitions := make(map[transitionKey]int)

	for i := 0; i < len(states); i++ {
		symbols := make(map[string]bool)
		for _, it := range states[i].Items {
			p := prods[it.Prod]
			if it.Dot < len(p.RHS) {
				symbols[p.RHS[it.Dot]] = true
			}
		}

		for _, sym := range sortedSet(symbols) {
			next := gotoItems(states[i].Items, sym)
			if len(next) == 0 {
				continue
			}
			key := itemsKey(next)
			id, ok := stateIDs[key]
			if !ok {
				id = len(states)
				stateIDs[key] = id
				states = append(states, State{Items: next})
			}
			transitions[transitionKey{i, sym}] = id
		}
	}

	return states, transitions
}

func putGoto(gotos map[int]map[string]int, state int, sym string, to int) {
	if gotos[state] == nil {
		gotos[state] = make(map[string]int)
	}
	gotos[state][sym] = to
}
