package command

import (
	"fmt"
)

// Command represents a CLI command
type Command struct {
	Name        string
	Description string
	Execute     func(args []string) error
}

// Manager handles command registration and execution
type Manager struct {
	commands map[string]*Command
}

// NewManager creates a new command manager
func NewManager() *Manager {
	return &Manager{
		commands: make(map[string]*Command),
	}
}

// Register adds a new command to the manager
func (m *Manager) Register(cmd *Command) error {
	if _, exists := m.commands[cmd.Name]; exists {
		return fmt.Errorf("command %s already registered", cmd.Name)
	}

	m.commands[cmd.Name] = cmd
	return nil
}

// Execute runs a command with the given arguments
func (m *Manager) Execute(name string, args []string) error {
	cmd, exists := m.commands[name]
	if !exists {
		return fmt.Errorf("command %s not found", name)
	}

	return cmd.Execute(args)
}

// List returns a list of all registered commands
func (m *Manager) List() []*Command {
	commands := make([]*Command, 0, len(m.commands))
	for _, cmd := range m.commands {
		commands = append(commands, cmd)
	}
	return commands
}

// Help prints help information for a command
func (m *Manager) Help(name string) error {
	if name == "" {
		fmt.Println("Available commands:")
		for _, cmd := range m.List() {
			fmt.Printf("  %s: %s\n", cmd.Name, cmd.Description)
		}
		return nil
	}

	cmd, exists := m.commands[name]
	if !exists {
		return fmt.Errorf("command %s not found", name)
	}

	fmt.Printf("Command: %s\n", cmd.Name)
	fmt.Printf("Description: %s\n", cmd.Description)
	return nil
} 