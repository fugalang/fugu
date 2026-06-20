package cli

import (
	"os"

	"github.com/fugalang/fugu/internal/composer/project"
	"github.com/fugalang/fugu/pkg/reader"
)

func Run(p *project.Project) (*project.Project, error) {
	if len(os.Args) < 2 {
		return Help(p)
	}
	switch os.Args[2] {
	default:
		content, err := reader.ReadRelativelyFile(os.Args[2])
		if err != nil {
			return p, err
		}
		_ = content
		return p, nil
	}
}
