package cli

import (
	"os"

	"github.com/fugalang/fugu/internal/composer/project"
)

func Init(p *project.Project) (*project.Project, error) {
	if len(os.Args) < 3 {
		return Help(p)
	}

	proj := project.InitProject(os.Args[2])
	return proj, nil
}
