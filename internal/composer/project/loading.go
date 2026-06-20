package project

import (
	"encoding/json"

	"github.com/fugalang/fugu/internal/composer/cacher"
	"github.com/fugalang/fugu/internal/library"
)

const (
	DirNameProject = ".fugu/"
	DirConfig      = DirNameProject + "config/"
	DirLibs        = DirNameProject + "cache/libs/"
	PrefixLibrary  = ".flc"

	PrefixFileConfig = ".cgf"
)

func InitProject(name string) *Project {
	path, err := PathOfHome()
	if err != nil {
		return nil
	}

	if CheckProject(DirNameProject) {
		cgf, err := ReadConfig(DirConfig, "")
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
	err = CreateFile(DirConfig, PrefixFileConfig, content)
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
	libs, err := GetLibraries(DirLibs, PrefixLibrary)
	if err != nil {
		return []library.Library{}
	}

	var libraries = []library.Library{}
	for _, libName := range libs {
		pathHome, err := PathOfHome()
		if err != nil {
			return []library.Library{}
		}

		path := pathHome + "/" + DirLibs

		content, err := ReadFile(path, libName, PrefixLibrary)
		libraries = append(libraries, cacher.ParseLibraryCach(content, path+libName+PrefixLibrary))
	}

	return libraries
}
