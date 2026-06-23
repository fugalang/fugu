package cacher

import (
	"os"

	dign "github.com/fugalang/fugu/internal/diagnostics"
	"github.com/fugalang/fugu/internal/diagnostics/errors"
	"github.com/fugalang/fugu/internal/library"
	"github.com/fugalang/fugu/pkg/helper"
)

func ParseLibraryCach(a dign.Arena, content []byte, path string) library.Library {
	cont := helper.ToFields(content)
	if len(cont) < 4 {
		err := errors.Errors[3].IU("LOADER", []string{
			"не удалось загрузить библиотеку. причина ошибки: не коректный формат файла кэша.",
		})
		a.Add(err)
		a.Print(os.Stderr)
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

// TODO ir
func ParseContentLibrary() {
}
