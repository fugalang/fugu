package gram

import (
	"fmt"
)

func mergeGrammars(dest *Grammar, imported *Grammar) error {
	if dest == nil || imported == nil {
		return nil
	}

	destNonTerms := make(map[string]bool)
	for _, p := range dest.productions {
		destNonTerms[p.LHS] = true
	}

	importedNonTerms := make(map[string]bool)
	for _, p := range imported.productions {
		importedNonTerms[p.LHS] = true
	}

	for nt := range importedNonTerms {
		if destNonTerms[nt] {
			if err := checkNonTerminalCompat(dest, imported, nt); err != nil {
				return fmt.Errorf("конфликт нетерминала %q при импорте: %w", nt, err)
			}
		}
	}

	for _, p := range imported.productions {
		if !containsProduction(dest, p) {
			p.ID = len(dest.productions)
			dest.productions = append(dest.productions, p)
		}
	}

	for tk := range imported.Terminals {
		if tk != "" && tk != EOF {
			dest.AddToken(tk)
		}
	}

	for tk, prec := range imported.precedence {
		if _, exists := dest.precedence[tk]; !exists {
			dest.precedence[tk] = prec
		}
	}

	for _, start := range imported.Starts {
		_ = dest.AddStart(start)
	}

	if len(imported.Starts) == 0 && imported.Start != "" {
		_ = dest.AddStart(imported.Start)
	}

	return nil
}

func checkNonTerminalCompat(dest *Grammar, imported *Grammar, nt string) error {
	destProds := make([]Production, 0)
	for _, p := range dest.productions {
		if p.LHS == nt {
			destProds = append(destProds, p)
		}
	}

	importedProds := make([]Production, 0)
	for _, p := range imported.productions {
		if p.LHS == nt {
			importedProds = append(importedProds, p)
		}
	}

	for _, iProd := range importedProds {
		found := false
		for _, dProd := range destProds {
			if productionsEqual(dProd, iProd) {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("несовместимое определение: imported добавляет новое правило %s -> %s", iProd.LHS, joinStrings(iProd.RHS))
		}
	}

	return nil
}

func containsProduction(g *Grammar, p Production) bool {
	for _, existing := range g.productions {
		if productionsEqual(existing, p) {
			return true
		}
	}
	return false
}

func productionsEqual(a, b Production) bool {
	if a.LHS != b.LHS || a.Node != b.Node || a.PrecBy != b.PrecBy {
		return false
	}
	if len(a.RHS) != len(b.RHS) {
		return false
	}
	for i := range a.RHS {
		if a.RHS[i] != b.RHS[i] {
			return false
		}
	}
	return true
}

func joinStrings(strs []string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += " "
		}
		result += s
	}
	return result
}
