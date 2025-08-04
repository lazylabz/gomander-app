package helpers

import "path/filepath"

func GetComputedPath(basePath, path string) string {
	if path == "" {
		return basePath
	}
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(basePath, path)
}
