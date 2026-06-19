package project

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// File help fn
func CheckProject(dir string) bool {
	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		} else {
			return false
		}

	}

	if info.IsDir() {
		return true
	} else {
		return false
	}
}

func ReadConfig(path, name string) (*Config, error) {
	content, err := ReadFile(path, name)
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
	content, err := ReadFile(path, name)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func CreateFile(folder string, filename string, content []byte) error {
	err := os.MkdirAll(folder, 0755)
	if err != nil {
		return err
	}

	path := folder + "/" + filename

	err = os.WriteFile(path, content, 0644)
	if err != nil {
		return err
	}

	return nil
}

func ReadFile(path, filename string) ([]byte, error) {
	content, err := os.ReadFile(path + "/" + filename)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func PathOfHome() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	rel, err := filepath.Rel(home, cwd)
	if err != nil {
		return "", err
	}

	return rel, nil
}

func GetLibraries(dir, prefix string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var libs []string

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if filepath.Ext(entry.Name()) == prefix {
			name := strings.TrimSuffix(entry.Name(), prefix)
			libs = append(libs, name)
		}
	}

	return libs, nil
}
