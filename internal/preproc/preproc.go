package preproc

import (
	"strings"

	"github.com/fugalang/fugu/pkg/helper"
)

type Arena struct {
	Defines map[string]string
}

func New() *Arena {
	return &Arena{
		Defines: map[string]string{},
	}
}

func (a *Arena) Define(content *[]byte) {
	lines := strings.Split(helper.BytesToString(*content), "\n")

	var out []string

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "#define") {
			rest := strings.TrimSpace(strings.TrimPrefix(line, "#define"))

			parts := strings.SplitN(rest, "=", 2)

			if len(parts) != 2 {
				continue
			}

			name := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			a.Defines[name] = value

			continue
		}

		parts := strings.Fields(line)

		for i, p := range parts {
			if value, ok := a.Defines[p]; ok {
				parts[i] = value
			}
		}

		out = append(out, strings.Join(parts, " "))
	}

	*content = []byte(strings.Join(out, "\n"))
}
