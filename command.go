package main

import (
	"errors"
)

type Command struct {
	Id               string `json:"id"`
	Name             string `json:"name"`
	Command          string `json:"command"`
	WorkingDirectory string `json:"workingDirectory"`
}

type CommandRepository struct {
	commands map[string]Command
}

func NewCommandRepository(commands map[string]Command) *CommandRepository {
	return &CommandRepository{
		commands: commands,
	}
}

func (r *CommandRepository) AddCommand(newCommand Command) error {
	if _, exists := r.commands[newCommand.Id]; exists {
		return errors.New("Command already exists: " + newCommand.Id)
	}

	r.commands[newCommand.Id] = newCommand

	return nil
}

func (r *CommandRepository) RemoveCommand(id string) error {
	if _, exists := r.commands[id]; !exists {
		return errors.New("Command not found: " + id)
	}

	delete(r.commands, id)

	return nil
}

func (r *CommandRepository) EditCommand(newCommand Command) error {
	if _, exists := r.commands[newCommand.Id]; !exists {
		return errors.New("Command not found: " + newCommand.Id)
	}

	r.commands[newCommand.Id] = newCommand

	return nil
}

func (r *CommandRepository) GetCommands() map[string]Command {
	return r.commands
}

func (r *CommandRepository) Get(id string) (*Command, error) {
	command, exists := r.commands[id]
	if !exists {
		return nil, errors.New("Command not found: " + id)
	}

	return &command, nil
}
