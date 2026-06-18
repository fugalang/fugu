package gram

import (
	"fmt"
	"sort"
	"strings"
)

func sortItems(seen map[Item]bool) []Item {
	items := make([]Item, 0, len(seen))
	for it := range seen {
		items = append(items, it)
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Prod != items[j].Prod {
			return items[i].Prod < items[j].Prod
		}
		return items[i].Dot < items[j].Dot
	})
	return items
}

func itemsKey(items []Item) string {
	var b strings.Builder
	for _, it := range items {
		b.WriteString(fmt.Sprintf("%d.%d;", it.Prod, it.Dot))
	}
	return b.String()
}

func sortedSet(set map[string]bool) []string {
	out := make([]string, 0, len(set))
	for s := range set {
		out = append(out, s)
	}
	sort.Strings(out)
	return out
}

func sortedIntKeys[V any](m map[int]V) []int {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys
}

func sortTokens(row map[string]Action, g *Grammar) []string {
	rank := make(map[string]int)
	for i, tk := range g.terminalSeq {
		rank[tk] = i
	}
	keys := make([]string, 0, len(row))
	for tk := range row {
		keys = append(keys, tk)
	}
	sort.Slice(keys, func(i, j int) bool {
		ri, iok := rank[keys[i]]
		rj, jok := rank[keys[j]]
		if iok && jok && ri != rj {
			return ri < rj
		}
		if iok != jok {
			return iok
		}
		return keys[i] < keys[j]
	})
	return keys
}

func sortedGotoKeys(row map[string]int) []string {
	keys := make([]string, 0, len(row))
	for k := range row {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
