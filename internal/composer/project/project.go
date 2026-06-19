package project

import "github.com/fugalang/fugu/internal/library"

type Project struct {
	Name string // название проекта
	Path string // путь относительный от HOME папки до проекта

	Libraries []library.Library
}

type Config struct {
	NameProject string
}
