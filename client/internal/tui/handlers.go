package tui

import (
	"errors"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// CommandHandler manages command execution and handling
type CommandHandler struct {
	cmdRunner *CommandRunner
}

// NewCommandHandler creates a new command handler
func NewCommandHandler() *CommandHandler {
	return &CommandHandler{
		cmdRunner: NewCommandRunner(),
	}
}

// HandleCommand processes a command and returns appropriate tea.Cmd
func (h *CommandHandler) HandleCommand(command string, formData *FormData) tea.Cmd {
	if IsDirectCommand(command) {
		return h.executeDirectCommand(command)
	}

	if RequiresInput(command) && formData != nil {
		return h.executeFormCommand(command, formData)
	}

	return nil
}

// executeDirectCommand executes commands that don't require input
func (h *CommandHandler) executeDirectCommand(command string) tea.Cmd {
	return tea.Sequence(
		tea.Println("Executing command..."),
		func() tea.Msg {
			parts := strings.Split(command, " ")
			if len(parts) != 2 {
				return CommandOutput{
					Output: "",
					Error:  fmt.Errorf("invalid command format: %s", command),
				}
			}

			var output string
			var err error

			switch parts[0] {
			case "config":
				output, err = h.cmdRunner.ExecuteConfigCommand(parts[1])
			case "exploit":
				output, err = h.cmdRunner.ExecuteExploitCommand(parts[1])
			default:
				err = fmt.Errorf("unknown command group: %s", parts[0])
			}

			return CommandOutput{
				Output: output,
				Error:  err,
			}
		},
	)
}

// executeFormCommand executes commands that require form input
func (h *CommandHandler) executeFormCommand(command string, formData *FormData) tea.Cmd {
	return tea.Sequence(
		tea.Println("Processing form data..."),
		func() tea.Msg {
			var output string
			var err error

			switch command {
			case "config login":
				output, err = h.handleLogin(formData)
			case "config update":
				output, err = h.handleConfigUpdate(formData)
			case "exploit run":
				output, err = h.handleExploitRun(formData)
			case "exploit create":
				output, err = h.handleExploitCreate(formData)
			case "exploit remove":
				output, err = h.handleExploitRemove(formData)
			case "exploit stop":
				output, err = h.handleExploitStop(formData)
			default:
				err = fmt.Errorf("unknown form command: %s", command)
			}

			return CommandOutput{
				Output: output,
				Error:  err,
			}
		},
	)
}

// handleLogin processes login command
func (h *CommandHandler) handleLogin(formData *FormData) (string, error) {
	password := formData.Fields["Password"]

	if password == "" {
		return "", errors.New("username and password are required")
	}

	return h.cmdRunner.ExecuteLogin(password)
}

// handleConfigUpdate processes config update command
func (h *CommandHandler) handleConfigUpdate(formData *FormData) (string, error) {
	host := formData.Fields["Host"]
	port := formData.Fields["Port"]
	username := formData.Fields["Username"]
	httpsStr := formData.Fields["HTTPS (true/false)"]

	useHTTPS := strings.ToLower(httpsStr) == "true"

	return h.cmdRunner.ExecuteConfigUpdate(host, port, username, useHTTPS)
}

// handleExploitRun processes exploit run command
func (h *CommandHandler) handleExploitRun(formData *FormData) (string, error) {
	exploitPath := formData.Fields["Exploit Path"]
	servicePort := formData.Fields["Service Port"]
	detachStr := formData.Fields["Detach Mode (true/false)"]
	tickTime := formData.Fields["Tick Time (seconds)"]
	threadCount := formData.Fields["Thread Count"]

	if exploitPath == "" || servicePort == "" {
		return "", errors.New("exploit path and service port are required")
	}

	detach := strings.ToLower(detachStr) == "true"

	return h.cmdRunner.ExecuteExploitRun(exploitPath, servicePort, detach, tickTime, threadCount)
}

// handleExploitCreate processes exploit create command
func (h *CommandHandler) handleExploitCreate(formData *FormData) (string, error) {
	name := formData.Fields["Exploit Name"]

	if name == "" {
		return "", errors.New("exploit name is required")
	}

	return h.cmdRunner.ExecuteExploitCreate(name)
}

// handleExploitRemove processes exploit remove command
func (h *CommandHandler) handleExploitRemove(formData *FormData) (string, error) {
	name := formData.Fields["Exploit Name"]

	if name == "" {
		return "", errors.New("exploit name is required")
	}

	return h.cmdRunner.ExecuteExploitRemove(name)
}

// handleExploitStop processes exploit stop command
func (h *CommandHandler) handleExploitStop(formData *FormData) (string, error) {
	pid := formData.Fields["Process ID (PID)"]

	if pid == "" {
		return "", errors.New("process ID is required")
	}

	return h.cmdRunner.ExecuteExploitStop(pid)
}

// HandleNavigation processes navigation commands
func (*CommandHandler) HandleNavigation(command string, model *Model) (*Model, tea.Cmd) {
	switch command {
	case "quit":
		model.quitting = true
		return model, tea.Quit
	case "back":
		model.SetActiveView("main")
		model.SetInputMode(false)
		model.SetRunningCommand(false)
		model.ClearError()
		return model, nil
	case "config":
		model.SetActiveView("config")
		return model, nil
	case "exploit":
		model.SetActiveView("exploit")
		return model, nil
	}

	return model, nil
}

// ProcessFormSubmission handles form submission
func (h *CommandHandler) ProcessFormSubmission(model *Model) (*Model, tea.Cmd) {
	// Validate form
	if err := ValidateForm(model.activeCommand, model.inputs); err != nil {
		model.SetError(err)
		return model, nil
	}

	// Get form data
	formData := GetFormData(model.inputs, model.inputLabels)

	// Switch to command execution mode
	model.SetInputMode(false)
	model.SetRunningCommand(true)
	model.ClearError()

	// Execute command
	return model, h.HandleCommand(model.activeCommand, &formData)
}

// SetupFormForCommand prepares form inputs for a specific command
func (*CommandHandler) SetupFormForCommand(model *Model, command string) {
	model.activeCommand = command
	model.inputs, model.inputLabels = CreateForm(command)
	model.focusIndex = 0
	model.SetInputMode(true)
	model.ClearError()
}
