package cli

import (
	"fmt"

	"github.com/fugalang/fugu/internal/composer/project"
)

func Help(p *project.Project) (*project.Project, error) {
	fmt.Println(``)
	return p, nil
}
