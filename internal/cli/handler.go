package cli

import (
	"os"

	"github.com/fugalang/fugu/internal/composer/project"
)

type Command func(p *project.Project) (*project.Project, error)

var commands = map[string]Command{
	"run":  Run,
	"init": Init,
	"help": Help,
}

func HandlerCmd(p *project.Project) (*project.Project, error) {
	if len(os.Args) < 1 {
		cmd := commands[os.Args[1]]
		return cmd(p)

	}
	cmd := commands[os.Args[1]]
	return cmd(p)
}
