package main

import (
	"os"

	gram "github.com/fugalang/fugu/pkg/gramforge"
	"github.com/fugalang/fugu/pkg/reader"
)

func main() {
	path, _ := reader.PathOfHome()
	path += "/internal/parser/tabular/action/map.go"

	g, err := gram.ParseFile("cmd/generator/action/.gram")
	table, err := gram.Build(g)
	out, err := gram.RenderMaps(table, gram.RenderOptions{
		PackageName: "action",
		Imports: []string{
			`"github.com/fugalang/fugu/internal/ast"`,
			`. "github.com/fugalang/fugu/internal/token"`,
		},
	})

	data := []byte(out)

	err = os.WriteFile(path, data, 0644)
	if err != nil {
		panic(err)
	}
}
