package main

import (
	"fmt"

	"github.com/fugalang/fugu/internal/composer/project"
)

func main() {
	pj := project.InitProject("app")
	fmt.Println(pj.Path)
	// TODO: цикл обработки в лексере и парсере всех файлов
}
