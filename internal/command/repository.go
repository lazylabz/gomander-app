package command

import (
	"errors"
)

type Repository struct {
	commands map[string]Command
}

func NewCommandRepository(commands map[string]Command) *Repository {
	return &Repository{
		commands: commands,
	}
}

func (r *Repository) AddCommand(newCommand Command) error {
	if _, exists := r.commands[newCommand.Id]; exists {
		return errors.New("Command already exists: " + newCommand.Id)
	}

	r.commands[newCommand.Id] = newCommand

	return nil
}

func (r *Repository) RemoveCommand(id string) error {
	if _, exists := r.commands[id]; !exists {
		return errors.New("Command not found: " + id)
	}

	delete(r.commands, id)

	return nil
}

func (r *Repository) EditCommand(newCommand Command) error {
	if _, exists := r.commands[newCommand.Id]; !exists {
		return errors.New("Command not found: " + newCommand.Id)
	}

	r.commands[newCommand.Id] = newCommand

	return nil
}

func (r *Repository) GetCommands() map[string]Command {
	return r.commands
}

func (r *Repository) GetCommand(id string) (*Command, error) {
	command, exists := r.commands[id]
	if !exists {
		return nil, errors.New("Command not found: " + id)
	}

	return &command, nil
}
