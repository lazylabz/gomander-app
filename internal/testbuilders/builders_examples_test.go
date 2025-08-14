package testbuilders

import (
	commanddomain "gomander/internal/command/domain"
	"testing"
)

func TestCommandBuilder(t *testing.T) {
	t.Run("Should create command with default values", func(t *testing.T) {
		cmd := NewCommandBuilder().Build()

		if cmd.Id == "" {
			t.Fatal("Expected command ID to be set")
		}

		if cmd.ProjectId == "" {
			t.Fatal("Expected command ProjectId to be set")
		}

		if cmd.Name != "Default Command" {
			t.Fatalf("Expected command name 'Default Command', got '%s'", cmd.Name)
		}

		if cmd.Command != "echo 'hello'" {
			t.Fatalf("Expected command 'echo \"hello\"', got '%s'", cmd.Command)
		}

		if cmd.WorkingDirectory != "/app" {
			t.Fatalf("Expected working directory '/app', got '%s'", cmd.WorkingDirectory)
		}

		if cmd.Position != 0 {
			t.Fatalf("Expected position 0, got %d", cmd.Position)
		}
	})

	t.Run("Should create simple custom command", func(t *testing.T) {
		cmd := NewCommandBuilder().
			WithName("Test command").
			WithCommand("go test").
			Build()

		if cmd.Name != "Test command" {
			t.Fatalf("Expected command name 'Test command', got '%s'", cmd.Name)
		}

		if cmd.Command != "go test" {
			t.Fatalf("Expected command 'go test', got '%s'", cmd.Command)
		}
	})
}

func TestCommandGroupBuilder(t *testing.T) {
	t.Run("Should create command group with default values", func(t *testing.T) {
		group := NewCommandGroupBuilder().Build()

		if group.Id == "" {
			t.Fatal("Expected command group ID to be set")
		}

		if group.ProjectId == "" {
			t.Fatal("Expected command group ProjectId to be set")
		}

		if group.Name != "Default Command Group" {
			t.Fatalf("Expected command group name 'Default Command Group', got '%s'", group.Name)
		}

		if len(group.Commands) != 0 {
			t.Fatal("Expected command group to have no commands by default")
		}

		if group.Position != 0 {
			t.Fatalf("Expected position 0, got %d", group.Position)
		}
	})

	t.Run("Should create custom command group", func(t *testing.T) {
		cmd1 := NewCommandBuilder().WithName("Command 1").Build()
		cmd2 := NewCommandBuilder().WithName("Command 2").Build()

		group := NewCommandGroupBuilder().
			WithName("Test Group").
			WithCommands([]commanddomain.Command{cmd1, cmd2}).
			WithPosition(1).
			Build()

		if group.Name != "Test Group" {
			t.Fatalf("Expected command group name 'Test Group', got '%s'", group.Name)
		}

		if len(group.Commands) != 2 {
			t.Fatalf("Expected 2 commands in group, got %d", len(group.Commands))
		}

		if group.Commands[0].Name != "Command 1" || group.Commands[1].Name != "Command 2" {
			t.Fatal("Commands in group do not match expected names")
		}

		if group.Position != 1 {
			t.Fatalf("Expected position 1, got %d", group.Position)
		}
	})

	t.Run("Should add command to existing group", func(t *testing.T) {
		cmd1 := NewCommandBuilder().WithName("Command 1").Build()
		group := NewCommandGroupBuilder().
			WithName("Test Group").
			AddCommand(cmd1).
			Build()

		if len(group.Commands) != 1 {
			t.Fatalf("Expected 1 command in group, got %d", len(group.Commands))
		}

		if group.Commands[0].Name != "Command 1" {
			t.Fatal("Added command does not match expected name")
		}
	})

	t.Run("Should ensure project ID consistency when adding command builder to group", func(t *testing.T) {
		cmdb := NewCommandBuilder().
			WithProjectId("command-project-id")
		group := NewCommandGroupBuilder().
			WithProjectId("group-project-id").
			AddCommandBuilder(cmdb).
			Build()

		if len(group.Commands) != 1 {
			t.Fatalf("Expected 1 command in group, got %d", len(group.Commands))
		}
		if group.Commands[0].ProjectId != "group-project-id" {
			t.Fatalf("Expected command ProjectId to match group ProjectId, got '%s'", group.Commands[0].ProjectId)
		}
	})

	t.Run("Should add multiple commands to group", func(t *testing.T) {
		cmd1 := NewCommandBuilder().WithName("Command 1").Build()
		cmd2 := NewCommandBuilder().WithName("Command 2").Build()

		group := NewCommandGroupBuilder().
			WithName("Test Group").AddCommands(cmd1, cmd2).
			Build()

		if len(group.Commands) != 2 {
			t.Fatalf("Expected 2 commands in group, got %d", len(group.Commands))
		}

		if group.Commands[0].Name != "Command 1" || group.Commands[1].Name != "Command 2" {
			t.Fatal("Commands in group do not match expected names")
		}
	})
}
