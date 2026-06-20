package cli

import (
	"fmt"
	"os"

	"github.com/fugalang/fugu/internal/composer/project"
	"github.com/fugalang/fugu/pkg/color"
	"github.com/fugalang/fugu/pkg/reader"
)

func Run(p *project.Project) (*project.Project, error) {
	if len(os.Args) < 4 {
		return Help(p)
	}
	switch os.Args[2] {
	case "--no-cache":
		os.Args = append(os.Args[:2], os.Args[3:]...)
		if NoCache {
			fmt.Println(color.PastelYellow("[!] Нет необходимости использовать один и тот же флаг дважды"))
		}
		NoCache = true
		return Run(p)
	default:
		content, err := reader.ReadRelativelyFile(os.Args[2])
		if err != nil {
			return p, err
		}
		_ = content
		return p, nil
	}
}
