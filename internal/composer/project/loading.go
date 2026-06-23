package project

import (
	"encoding/json"

	"github.com/fugalang/fugu/internal/composer/cacher"
	"github.com/fugalang/fugu/internal/diagnostics"
	"github.com/fugalang/fugu/internal/library"
	"github.com/fugalang/fugu/pkg/reader"
)

const (
	DirNameProject = ".fugu/"
	DirConfig      = DirNameProject + "config/"
	DirLibs        = DirNameProject + "cache/libs/"
	PrefixLibrary  = ".flc"

	PrefixFileConfig = ".cgf"
)

func InitProject(a diagnostics.Arena, name string) *Project {
	path, err := reader.PathOfHome()
	if err != nil {
		return nil
	}

	if reader.CheckProject(DirNameProject) {
		cgf, err := ReadConfig(DirConfig, "")
		if err != nil {
			return nil
		}

		// CU1
		return &Project{
			Config:    *cgf,
			Path:      path,
			Libraries: LoadLibraries(a),

			Ad: diagnostics.Arena{},
		}
	}

	content, err := CgfFileContentGen(name)
	if err != nil {
		return nil
	}
	err = reader.CreateFile(DirConfig, PrefixFileConfig, content)
	if err != nil {
		return nil
	}

	// CU1
	return &Project{
		Config: Config{
			NameProject: name,
			IsCache:     true,
		},
		Path:      path,
		Libraries: []library.Library{},

		Ad: diagnostics.Arena{},
	}
}

func CgfFileContentGen(nameProject string) ([]byte, error) {
	cfg := Config{
		NameProject: nameProject,
		IsCache:     true,
	}

	b, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return nil, err
	}

	return b, nil
}

func LoadLibraries(a diagnostics.Arena) []library.Library {
	libs, err := reader.GetLibraries(DirLibs, PrefixLibrary)
	if err != nil {
		return []library.Library{}
	}

	var libraries = []library.Library{}
	for _, libName := range libs {
		pathHome, err := reader.PathOfHome()
		if err != nil {
			return []library.Library{}
		}

		path := pathHome + "/" + DirLibs

		content, err := reader.ReadFile(path, libName, PrefixLibrary)
		libraries = append(libraries, cacher.ParseLibraryCach(a, content, path+libName+PrefixLibrary))
	}

	return libraries
}

func ReadConfig(path, name string) (*Config, error) {
	content, err := reader.ReadFile(path, name, PrefixFileConfig)
	if err != nil {
		return nil, err
	}

	cgf := &Config{}
	err = json.Unmarshal(content, cgf)
	if err != nil {
		return nil, err
	}

	return cgf, nil
}

func ReadLibrary(path, name string) ([]byte, error) {
	content, err := reader.ReadFile(path, name, PrefixLibrary)
	if err != nil {
		return nil, err
	}

	return content, nil
}
