package cli

import (
	"strings"

	"github.com/fugalang/fugu/internal/composer/project"
)

func Help(p *project.Project) (*project.Project, error) {
	var sb strings.Builder
	sb.WriteString("Команды:\n")
	sb.WriteString("fugu init name_project         - создания проекта\n")
	sb.WriteString("fugu run path_src              - запуск файла.\n	Флаги: --no-cache собирает всё с нуля\n")
	sb.WriteString("fugu build path_src path_save  - сборка целевого файла.\n	Флаги: --no-cache собирает всё с нуля\n")
	return p, nil
}
