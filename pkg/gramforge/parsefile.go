package gram

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ParseFile(path string) (*Grammar, error) {
	visited := make(map[string]bool)
	return parseFileRecursive(path, visited)
}

func parseFileRecursive(path string, visited map[string]bool) (*Grammar, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("не могу получить абсолютный путь %q: %w", path, err)
	}

	if visited[absPath] {
		return nil, fmt.Errorf("циклический импорт: %q", absPath)
	}
	visited[absPath] = true

	raw, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("не могу прочитать файл %q: %w", absPath, err)
	}

	g, err := parseWithImports(string(raw), filepath.Dir(absPath), visited)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга %q: %w", absPath, err)
	}

	return g, nil
}

func parseWithImports(src string, baseDir string, visited map[string]bool) (*Grammar, error) {
	g := New()
	var current string

	lines := strings.Split(src, "\n")
	var importPaths []string

	for i, raw := range lines {
		line := cleanLine(raw)
		if line == "" {
			continue
		}
		linen := i + 1

		if strings.HasPrefix(line, "%") {
			if strings.HasPrefix(line, "%import") {
				parts := strings.Fields(line)
				if len(parts) != 2 {
					return nil, fmt.Errorf("строка %d: %%import хочет ровно один путь", linen)
				}
				importPath := parts[1]
				importPath = strings.Trim(importPath, `"`)
				importPaths = append(importPaths, importPath)
				continue
			}

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

	for _, importPath := range importPaths {
		fullPath := filepath.Join(baseDir, importPath)

		importedGram, err := parseFileRecursive(fullPath, visited)
		if err != nil {
			return nil, fmt.Errorf("ошибка при импорте %q: %w", importPath, err)
		}

		if err := mergeGrammars(g, importedGram); err != nil {
			return nil, fmt.Errorf("ошибка при слиянии импорта %q: %w", importPath, err)
		}
	}

	return g, nil
}
