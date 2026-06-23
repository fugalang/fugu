package project

import (
	"github.com/fugalang/fugu/internal/diagnostics"
	"github.com/fugalang/fugu/internal/library"
)

type Project struct {
	Path   string // путь относительный от HOME папки до проекта
	Config Config

	Libraries []library.Library

	Ad diagnostics.Arena
}

// При изменении измените структуры с коминтарием (CU1)
type Config struct {
	NameProject string
	IsCache     bool
}
