package gram

import (
	"fmt"
	"strings"
)

func Parse(src string) (*Grammar, error) {
	g := New()
	var current string

	lines := strings.Split(src, "\n")
	for i, raw := range lines {
		line := cleanLine(raw)
		if line == "" {
			continue
		}
		linen := i + 1

		if strings.HasPrefix(line, "%") {
			if err := parseDirective(g, line); err != nil {
				return nil, fmt.Errorf("строка %d: %w", linen, err)
			}
			continue
		}

		if strings.HasSuffix(line, ":") && !strings.Contains(line, "=>") {
			lhs := strings.TrimSpace(strings.TrimSuffix(line, ":"))
			if lhs == "" {
				return nil, fmt.Errorf("строка %d: пустой заголовок правила", linen)
			}
			current = lhs
			continue
		}

		lhs := current
		body := line
		if strings.HasPrefix(line, "|") {
			body = strings.TrimSpace(strings.TrimPrefix(line, "|"))
			if lhs == "" {
				return nil, fmt.Errorf("строка %d: альтернатива без предыдущего правила", linen)
			}
		} else if current == "" || hasRuleSep(line) {
			var ok bool
			lhs, body, ok = cutRule(line)
			if !ok {
				return nil, fmt.Errorf("строка %d: ждал правило вида A -> B C => Kind", linen)
			}
			current = lhs
		}

		alts := splitAlternatives(body)
		for _, alt := range alts {
			rhs, node, precBy, err := parseAlternative(alt)
			if err != nil {
				return nil, fmt.Errorf("строка %d: %w", linen, err)
			}
			if err := g.AddProduction(lhs, rhs, node, precBy); err != nil {
				return nil, fmt.Errorf("строка %d: %w", linen, err)
			}
		}
	}

	g.AddToken(EOF)
	return g, nil
}

func parseDirective(g *Grammar, line string) error {
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return nil
	}

	switch parts[0] {
	case "%start":
		if len(parts) != 2 {
			return fmt.Errorf("%%start хочет ровно один символ")
		}
		return g.AddStart(parts[1])
	case "%token", "%tokens":
		if len(parts) < 2 {
			return fmt.Errorf("%s хочет хотя бы один токен", parts[0])
		}
		for _, tk := range parts[1:] {
			g.AddToken(tk)
		}
	case "%left":
		g.AddPrecedence(AssocLeft, parts[1:]...)
	case "%right":
		g.AddPrecedence(AssocRight, parts[1:]...)
	case "%nonassoc":
		g.AddPrecedence(AssocNonassoc, parts[1:]...)
	case "%import":
		if len(parts) != 2 {
			return fmt.Errorf("%%import хочет ровно один путь")
		}
		return nil
	default:
		return fmt.Errorf("не знаю директиву %s", parts[0])
	}

	return nil
}

func cleanLine(line string) string {
	if idx := strings.Index(line, "#"); idx >= 0 {
		line = line[:idx]
	}
	if idx := strings.Index(line, "//"); idx >= 0 {
		line = line[:idx]
	}
	line = strings.TrimSpace(line)
	line = strings.TrimSuffix(line, ";")
	return strings.TrimSpace(line)
}

func cutRule(line string) (string, string, bool) {
	for _, sep := range []string{"->", ":"} {
		if before, after, ok := strings.Cut(line, sep); ok {
			lhs := strings.TrimSpace(before)
			body := strings.TrimSpace(after)
			return lhs, body, lhs != "" && body != ""
		}
	}
	return "", "", false
}

func hasRuleSep(line string) bool {
	return strings.Contains(line, "->") || strings.Contains(line, ":")
}

func splitAlternatives(body string) []string {
	raw := strings.Split(body, "|")
	out := make([]string, 0, len(raw))
	for _, alt := range raw {
		alt = strings.TrimSpace(alt)
		if alt != "" {
			out = append(out, alt)
		}
	}
	return out
}

func parseAlternative(alt string) ([]string, string, string, error) {
	var node string
	if before, after, ok := strings.Cut(alt, "=>"); ok {
		alt = strings.TrimSpace(before)
		fields := strings.Fields(strings.TrimSpace(after))
		if len(fields) == 0 {
			return nil, "", "", fmt.Errorf("после => нужен Kind")
		}
		node = fields[0]
		if len(fields) > 1 {
			return nil, "", "", fmt.Errorf("лишнее после Kind: %s", strings.Join(fields[1:], " "))
		}
	}

	fields := strings.Fields(alt)
	rhs := make([]string, 0, len(fields))
	var precBy string
	for i := 0; i < len(fields); i++ {
		switch fields[i] {
		case "%prec":
			if i+1 >= len(fields) {
				return nil, "", "", fmt.Errorf("после %%prec нужен токен")
			}
			precBy = fields[i+1]
			i++
		case "empty", "EMPTY", "_", "ε":
		default:
			rhs = append(rhs, fields[i])
		}
	}

	return rhs, node, precBy, nil
}
