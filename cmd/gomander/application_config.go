package main

import "gomander/internal/runner"

type applicationConfig struct {
	configFolderName string
	runner           runner.Config
}
