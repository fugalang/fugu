package cacher

import (
	"strings"

	"github.com/fugalang/fugu/internal/library"
)

func ParseLibraryCach(content []byte, path string) library.Library {
	cont := toFields(content)
	if len(cont) < 4 {
		// TODO ошибка
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

func ParseContentLibrary() {
}

func toFields(content []byte) []string {
	return strings.Fields(string(content))
}
