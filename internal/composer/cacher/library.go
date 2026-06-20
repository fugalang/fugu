package cacher

import (
	"strings"

	d "github.com/fugalang/fugu/internal/diagnostics"
	"github.com/fugalang/fugu/internal/diagnostics/errors"
	"github.com/fugalang/fugu/internal/library"
)

func ParseLibraryCach(content []byte, path string) library.Library {
	cont := toFields(content)
	if len(cont) < 4 {
		err := errors.Errors[3].IU("LOADER", []string{
			"не удалось загрузить библиотеку. причина ошибки: не коректный формат файла кэша.",
		})
		d.Da.AddError(err)
		d.Da.Print()
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
