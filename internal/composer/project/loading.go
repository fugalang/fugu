package project

import (
	"encoding/json"

	"github.com/fugalang/fugu/internal/library"
)

const (
	DirNameProject = ".fugu"
	DirConfig      = DirNameProject + "/config"
	DirLibs        = DirNameProject + "/cache/libs/"
	PrefixLibrary  = ".flc"

	NameFileConfig = "project.cgf"
)

func InitProject(name string) *Project {
	path, err := PathOfHome()
	if err != nil {
		return nil
	}

	if CheckProject(DirNameProject) {
		cgf, err := ReadConfig(DirConfig, NameFileConfig)
		if err != nil {
			return nil
		}

		return &Project{
			Name:      cgf.NameProject,
			Path:      path,
			Libraries: LoadLibraries(),
		}
	}

	content, err := CgfFileContentGen(name)
	if err != nil {
		return nil
	}
	err = CreateFile(DirConfig, NameFileConfig, content)
	if err != nil {
		return nil
	}

	return &Project{
		Name:      name,
		Path:      path,
		Libraries: []library.Library{},
	}
}

func CgfFileContentGen(nameProject string) ([]byte, error) {
	cfg := Config{
		NameProject: nameProject,
	}

	b, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return nil, err
	}

	return b, nil
}

func LoadLibraries() []library.Library {
	libs, err := GetLibraries(DirNameProject, PrefixLibrary)
	if err != nil {
		return []library.Library{}
	}

	var libraries = []library.Library{}
	for _, lib := range libs {
		path, err := PathOfHome()
		if err != nil {
			return []library.Library{}
		}

		libraries = append(libraries, library.Library{
			Name:    lib,
			Path:    path + "/" + DirLibs + lib,
			Version: "TODO", // TODO разобрать файл библиотеки и достать версию
		})
	}

	return libraries
}
