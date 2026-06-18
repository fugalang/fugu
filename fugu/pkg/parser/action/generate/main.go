package main

import (
	"fmt"
	"go/format"
	"os"
	"path/filepath"

	"fugu/pkg/parser/action"
	"fugu/pkg/parser/action/gram"
)

func main() {
	if err := generateAll(); err != nil {
		panic(err)
	}
}

func generateAll() error {
	base := "pkg/parser/action"

	inPath := filepath.Join(base, "ast/ast.gm")
	outDir := base

	raw, err := os.ReadFile(inPath)
	if err != nil {
		return fmt.Errorf("read grammar: %w", err)
	}

	g, err := gram.Parse(string(raw))
	if err != nil {
		return fmt.Errorf("parse grammar: %w", err)
	}

	table, err := gram.Build(g)
	if err != nil {
		return fmt.Errorf("build table: %w", err)
	}

	if err := generateAction(outDir); err != nil {
		return fmt.Errorf("generate action: %w", err)
	}

	if err := generateMap(outDir, table); err != nil {
		return fmt.Errorf("generate map: %w", err)
	}

	return nil
}

func generateAction(outDir string) error {
	content := action.GenerateActionTable(action.ActionSrc)

	return writeFile(
		filepath.Join(outDir, "action.gen.go"),
		content,
	)
}

func generateMap(outDir string, table *gram.Table) error {
	content, err := gram.RenderMaps(table, gram.RenderOptions{
		PackageName: "action",
		Imports: []string{
			`"fugu/pkg/parser/action/ast"`,
			`. "fugu/pkg/token"`,
		},
	})
	if err != nil {
		return err
	}

	return writeFile(
		filepath.Join(outDir, "map.gen.go"),
		content,
	)
}

func writeFile(path string, content string) error {
	formatted, err := format.Source([]byte(content))
	if err != nil {
		return fmt.Errorf("format error %s: %w", path, err)
	}

	return os.WriteFile(path, formatted, 0644)
}
