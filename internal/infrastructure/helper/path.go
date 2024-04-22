package helper

import (
	"os"
	"path/filepath"
	"strings"
)

const (
	StaticDirPath    = "/public/statics/"
	TemplatesDirPath = "/html/"
	ResourcesDirPath = "/public/resources/"
	LogsDirPath      = "/var/log/"
)

func TemplatePath(template string, dirs ...string) (string, error) {
	staticFilesDir, err := StaticFilesDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(
		filepath.Join(staticFilesDir, TemplatesDirPath),
		filepath.Join(append(dirs, template)...),
	), nil
}

func StaticFilesDir() (string, error) {
	return path(StaticDirPath)
}

func ResourcesDir() (string, error) {
	return path(ResourcesDirPath)
}

func LogsDir() (string, error) {
	return path(LogsDirPath)
}

// path is a function which builts any path from root dir.
func path(additionalPath string) (string, error) {
	root, err := os.Getwd()
	if err != nil {
		return "", err
	}

	b := strings.Builder{}
	b.WriteString(root)
	b.WriteString(additionalPath)
	fullPath := b.String()

	return fullPath, nil
}
