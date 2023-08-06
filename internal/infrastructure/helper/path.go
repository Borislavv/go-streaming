package helper

import (
	"os"
	"path/filepath"
	"strings"
)

const (
	TemplatesDir     = "/html/"
	ResourcesDirPath = "/internal/infrastructure/static/"
)

func TemplatePath(template string, dirs ...string) (string, error) {
	resDir, err := ResourcesDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(
		filepath.Join(resDir, TemplatesDir),
		filepath.Join(append(dirs, template)...),
	), nil
}

func ResourcesDir() (string, error) {
	execPath, err := os.Getwd()
	if err != nil {
		return "", err
	}

	b := strings.Builder{}
	b.WriteString(execPath)
	b.WriteString(ResourcesDirPath)
	path := b.String()

	return path, nil
}
