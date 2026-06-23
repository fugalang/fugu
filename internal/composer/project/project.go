package project

import (
	"github.com/fugalang/fugu/internal/diagnostics"
	"github.com/fugalang/fugu/internal/library"
)

type KindFile uint8

const (
	_ KindFile = iota
)

type Project struct {
	Path   string // путь относительный от HOME папки до проекта
	Config Config

	Files     []File
	Libraries []library.Library

	Ad diagnostics.Arena
}

// При изменении измените структуры с коминтарием (CU1)
type Config struct {
	NameProject string
	IsCache     bool
}

type File struct {
	Name       string
	Type       KindFile
	RawContent []byte // необработанное содержимое файла
}
