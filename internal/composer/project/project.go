package project

import "github.com/fugalang/fugu/internal/library"

type Project struct {
	Path   string // путь относительный от HOME папки до проекта
	Config Config

	Libraries []library.Library
}

type Config struct {
	NameProject string
	IsCache     bool
}
