package project

import "errors"

type Command struct {
	Id               string `json:"id"`
	Name             string `json:"name"`
	Command          string `json:"project"`
	WorkingDirectory string `json:"workingDirectory"`
}

func (p *Project) AddCommand(newCommand Command) error {
	if _, exists := p.Commands[newCommand.Id]; exists {
		return errors.New("Command already exists: " + newCommand.Id)
	}

	p.Commands[newCommand.Id] = newCommand

	return nil
}

func (p *Project) RemoveCommand(id string) error {
	if _, exists := p.Commands[id]; !exists {
		return errors.New("Command not found: " + id)
	}

	delete(p.Commands, id)

	return nil
}

func (p *Project) EditCommand(newCommand Command) error {
	if _, exists := p.Commands[newCommand.Id]; !exists {
		return errors.New("Command not found: " + newCommand.Id)
	}

	p.Commands[newCommand.Id] = newCommand

	return nil
}

func (p *Project) GetCommands() map[string]Command {
	return p.Commands
}

func (p *Project) GetCommand(id string) (*Command, error) {
	command, exists := p.Commands[id]
	if !exists {
		return nil, errors.New("Command not found: " + id)
	}

	return &command, nil
}
