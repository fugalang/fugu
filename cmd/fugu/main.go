package main

import (
	"fmt"
	"os"

	"github.com/fugalang/fugu/internal/cli"
	"github.com/fugalang/fugu/internal/composer/project"
	"github.com/fugalang/fugu/internal/diagnostics"
	"github.com/fugalang/fugu/internal/diagnostics/errors"
)

var proj = &project.Project{}

func main() {
	var er error
	proj, er = cli.HandlerCmd(proj)
	if er != nil {
		da := diagnostics.Arena{}
		err := errors.Errors[4].IU("CLI", []string{""})
		err.Description = []string{
			fmt.Sprintf("ошибка выполнения команды: %s", er.Error()),
		}
		da.Add(err)
		da.Print(os.Stderr)
	}
}
