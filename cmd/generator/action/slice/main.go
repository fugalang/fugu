package main

import (
	"os"

	"github.com/fugalang/fugu/internal/parser/tabular/action"
	"github.com/fugalang/fugu/pkg/reader"
)

func main() {

	path, _ := reader.PathOfHome()
	path += "/internal/parser/tabular/action/slice.go"

	data := []byte(action.GenerateActionTable(action.ActionMap))

	err := os.WriteFile(path, data, 0644)
	if err != nil {
		panic(err)
	}
}
