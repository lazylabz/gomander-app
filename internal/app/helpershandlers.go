package app

func (a *App) GetComputedPath(basePath, path string) string {
	return a.pathHelper.GetComputedPath(basePath, path)
}
