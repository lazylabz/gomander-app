package project

import (
	"errors"
	"gomander/internal/helpers/array"
)

type Command struct {
	Id               string `json:"id"`
	Name             string `json:"name"`
	Command          string `json:"command"`
	WorkingDirectory string `json:"workingDirectory"`
}

func (p *Project) AddCommand(newCommand Command) error {
	if _, exists := p.Commands[newCommand.Id]; exists {
		return errors.New("Command already exists: " + newCommand.Id)
	}

	p.Commands[newCommand.Id] = newCommand
	p.OrderedCommandIds = append(p.OrderedCommandIds, newCommand.Id)

	return nil
}

func (p *Project) RemoveCommand(id string) error {
	if _, exists := p.Commands[id]; !exists {
		return errors.New("Command not found: " + id)
	}

	delete(p.Commands, id)
	p.removeFromOrderedIds(id)

	return nil
}

func (p *Project) removeFromOrderedIds(id string) {
	p.OrderedCommandIds = array.Filter(p.OrderedCommandIds, func(orderedId string) bool {
		return orderedId != id
	})
}

func (p *Project) EditCommand(newCommand Command) error {
	if _, exists := p.Commands[newCommand.Id]; !exists {
		return errors.New("Command not found: " + newCommand.Id)
	}

	p.Commands[newCommand.Id] = newCommand

	return nil
}

func (p *Project) GetCommands() []Command {
	ordered := make([]Command, 0, len(p.OrderedCommandIds))
	for _, id := range p.OrderedCommandIds {
		if command, exists := p.Commands[id]; exists {
			ordered = append(ordered, command)
		}
	}
	return ordered
}

func (p *Project) GetCommand(id string) (*Command, error) {
	command, exists := p.Commands[id]
	if !exists {
		return nil, errors.New("Command not found: " + id)
	}

	return &command, nil
}

func (p *Project) ReorderCommands(newOrder []string) error {
	if len(newOrder) != len(p.Commands) {
		return errors.New("new order length doesn't match number of commands")
	}

	// Ensure no duplicates and all commands exist in the new order
	commandSet := make(map[string]bool, len(p.Commands))
	for _, id := range newOrder {
		if commandSet[id] {
			return errors.New("duplicate command ID in new order: " + id)
		}
		if _, exists := p.Commands[id]; !exists {
			return errors.New("command not found: " + id)
		}
		commandSet[id] = true
	}

	p.OrderedCommandIds = make([]string, len(newOrder))
	copy(p.OrderedCommandIds, newOrder)

	return nil
}

func (p *Project) ensureOrderedCommandsConsistency() {
	if p.OrderedCommandIds == nil {
		p.OrderedCommandIds = make([]string, 0, len(p.Commands))
	}

	existingCommandsSet := make(map[string]bool, len(p.Commands))
	for id := range p.Commands {
		existingCommandsSet[id] = true
	}

	// Remove IDs that don't exist in Commands
	validIds := p.OrderedCommandIds[:0]
	for _, id := range p.OrderedCommandIds {
		if existingCommandsSet[id] {
			validIds = append(validIds, id)
		}
	}
	p.OrderedCommandIds = validIds

	// Add missing command IDs to ordered list
	alreadyOrderedSet := make(map[string]bool, len(p.OrderedCommandIds))
	for _, id := range p.OrderedCommandIds {
		alreadyOrderedSet[id] = true
	}
	for id := range p.Commands {
		if !alreadyOrderedSet[id] {
			p.OrderedCommandIds = append(p.OrderedCommandIds, id)
		}
	}
}
