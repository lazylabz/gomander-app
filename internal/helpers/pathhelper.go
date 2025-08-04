package helpers

import "path/filepath"

type PathHelper struct {
}

func (ph PathHelper) GetComputedPath(basePath, path string) string {
	if path == "" {
		return basePath
	}
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(basePath, path)
}
