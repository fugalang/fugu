package cacher

import (
	"strings"

	"github.com/fugalang/fugu/internal/diagnostics"
	"github.com/fugalang/fugu/internal/diagnostics/errors"
	"github.com/fugalang/fugu/internal/library"
	"github.com/fugalang/fugu/internal/token"
)

func ParseLibraryCach(content []byte, path string) library.Library {
	cont := toFields(content)
	if len(cont) < 4 {
		da := diagnostics.DiagnosticArena{}
		err := errors.Errors[3].Update(token.Token{
			Pos: token.Position{
				FileName: "LOADER",
				Line:     0,
			},
		})
		err.Description = []string{
			"не удалось загрузить библиотеку. причина ошибки: не коректный формат файла кэша.",
		}
		da.AddError(err)
		da.Print()
		return library.Library{}
	}

	return library.Library{
		Name:    cont[0],
		Author:  cont[1],
		Version: cont[2],
		Path:    path,
		Content: cont[3:],
	}
}

// TODO
func ParseContentLibrary() {
}

func toFields(content []byte) []string {
	return strings.Fields(string(content))
}
