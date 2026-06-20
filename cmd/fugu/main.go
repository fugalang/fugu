package main

import (
	"fmt"

	"github.com/fugalang/fugu/internal/token"

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
		da := diagnostics.DiagnosticArena{}
		err := errors.Errors[4].Update(token.Token{
			Pos: token.Position{
				FileName: "CLI",
				Line:     0,
			},
		})
		err.Description = []string{
			fmt.Sprintf("ошибка выполнения команды: %s", er.Error()),
		}
		da.AddError(err)
		da.Print()
	}
}
