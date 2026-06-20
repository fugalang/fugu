package main

import (
	"fmt"

	"github.com/fugalang/fugu/internal/cli"
	"github.com/fugalang/fugu/internal/composer/project"
)

func main() {
	pj := project.InitProject("app")
	fmt.Println(pj.Path)

	cli.HandlerCmd(pj)
}
