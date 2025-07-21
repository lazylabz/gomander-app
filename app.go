package main

import (
	"context"
)

// App struct
type App struct {
	ctx                context.Context
	logger             *Logger
	eventEmitter       *EventEmitter
	commandRunner      *CommandRunner
	commandsRepository *CommandRepository
	config             *Config
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	logger := NewLogger(ctx)
	eventEmitter := NewEventEmitter(ctx)
	commandRunner := NewCommandRunner(logger, eventEmitter)

	a.logger = logger
	a.eventEmitter = eventEmitter
	a.commandRunner = commandRunner

	a.logger.info("Loading configuration...")

	a.config = loadConfigOrPanic()

	a.logger.info("Configuration loaded successfully")

	a.commandsRepository = NewCommandRepository(a.config.Commands)
}

func (a *App) GetCommands() map[string]Command {
	return a.commandsRepository.commands
}

func (a *App) AddCommand(newCommand Command) {
	err := a.commandsRepository.AddCommand(newCommand)
	if err != nil {
		a.logger.error(err.Error())
		a.eventEmitter.emitEvent(ErrorNotification, err.Error())
		return
	}

	a.logger.info("Command added: " + newCommand.Id)
	a.eventEmitter.emitEvent(SuccessNotification, "Command added")

	// Update the commands map in the frontend
	a.eventEmitter.emitEvent(GetCommands, nil)

	a.saveConfig()
}

func (a *App) RemoveCommand(id string) {
	err := a.commandsRepository.RemoveCommand(id)
	if err != nil {
		a.logger.error(err.Error())
		a.eventEmitter.emitEvent(ErrorNotification, err.Error())
		return
	}

	a.logger.info("Command removed: " + id)
	a.eventEmitter.emitEvent(SuccessNotification, "Command removed")

	// Update the commands map in the frontend
	a.eventEmitter.emitEvent(GetCommands, nil)

	a.saveConfig()
}

func (a *App) EditCommand(newCommand Command) {
	err := a.commandsRepository.EditCommand(newCommand)
	if err != nil {
		a.logger.error(err.Error())
		a.eventEmitter.emitEvent(ErrorNotification, err.Error())
		return
	}

	a.logger.info("Command edited: " + newCommand.Id)
	a.eventEmitter.emitEvent(SuccessNotification, "Command edited")

	// Update the commands map in the frontend
	a.eventEmitter.emitEvent(GetCommands, nil)

	a.saveConfig()
}

func (a *App) RunCommand(id string) map[string]Command {
	command, err := a.commandsRepository.Get(id)

	if err != nil {
		a.logger.error(err.Error())
		a.eventEmitter.emitEvent(ErrorNotification, err.Error())
		return nil
	}

	err = a.commandRunner.RunCommand(*command, a.config.ExtraPaths)

	if err != nil {
		a.logger.error(err.Error())
		a.eventEmitter.emitEvent(ErrorNotification, "Failed to run command: "+id+" - "+err.Error())
		return nil
	}

	a.logger.info("Command executed: " + id)

	return nil
}

func (a *App) StopCommand(id string) {
	_, err := a.commandsRepository.Get(id)

	if err != nil {
		a.logger.error(err.Error())
		a.eventEmitter.emitEvent(ErrorNotification, err.Error())
		return
	}

	err = a.commandRunner.StopRunningCommand(id)

	if err != nil {
		a.logger.error(err.Error())
		a.eventEmitter.emitEvent(ErrorNotification, "Failed to stop command gracefully: "+id+" - "+err.Error())
		return
	}

	a.logger.info("Command stopped: " + id)

	a.eventEmitter.emitEvent(ProcessFinished, id)
}

func (a *App) SaveUserConfig(userConfig UserConfig) {
	a.config.ExtraPaths = userConfig.ExtraPaths

	saveConfigOrPanic(a.config)

	a.logger.info("Extra paths saved successfully")
	a.eventEmitter.emitEvent(SuccessNotification, "Extra paths saved successfully")
	a.eventEmitter.emitEvent(GetUserConfig, nil)
}

func (a *App) GetUserConfig() UserConfig {
	return UserConfig{
		ExtraPaths: a.config.ExtraPaths,
	}
}

func (a *App) saveConfig() {
	saveConfigOrPanic(&Config{
		Commands: a.commandsRepository.commands,
	})
}
