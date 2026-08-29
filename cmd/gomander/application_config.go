package main

import "gomander/internal/runner"

type applicationConfig struct {
	configFolderName string
	runner           runner.Config
}

var appConfig = applicationConfig{
	configFolderName: "gomander",
	runner: runner.Config{
		ConPTYEnvironments: []runner.HostEnvironment{
			runner.HostEnvironmentWindows10,
		},
	},
}
