package cli

import (
	"fmt"
	"os"

	"github.com/fugalang/fugu/internal/composer/project"
	"github.com/fugalang/fugu/pkg/reader"
)

func Run(p *project.Project) error {
	if len(os.Args) < 2 {
		// TODO help
		return nil
	}
	switch os.Args[2] {
	default:
		content, err := reader.ReadRelativelyFile(os.Args[2])
		if err != nil {
			return err
		}
		_ = content
		fmt.Printf(string(content))
		return nil
	}
}
