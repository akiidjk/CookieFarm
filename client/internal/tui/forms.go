package tui

import (
	"fmt"
	"strings"

	"github.com/ByteTheCookies/cookieclient/internal/config"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// FormData represents form input data
type FormData struct {
	Fields map[string]string
}

// CreateForm creates input fields for a specific command
func CreateForm(command string) ([]textinput.Model, []string) {
	var inputs []textinput.Model
	var labels []string

	switch command {
	case "config login":
		inputs, labels = createLoginForm()
	case "config update":
		inputs, labels = createConfigUpdateForm()
	case "exploit run":
		inputs, labels = createExploitRunForm()
	case "exploit create", "exploit remove":
		inputs, labels = createExploitNameForm()
	case "exploit stop":
		inputs, labels = createExploitStopForm()
	}

	// Set focus on first input if any exist
	if len(inputs) > 0 {
		inputs[0].Focus()
	}

	return inputs, labels
}

// createLoginForm creates form for login command
func createLoginForm() ([]textinput.Model, []string) {
	var inputs []textinput.Model
	var labels []string

	// Username input
	usernameInput := textinput.New()
	usernameInput.Placeholder = "Username"
	usernameInput.CharLimit = 32
	usernameInput.Width = 30
	usernameInput.SetValue(config.ArgsConfigInstance.Username)
	inputs = append(inputs, usernameInput)
	labels = append(labels, "Username")

	// Password input
	passwordInput := textinput.New()
	passwordInput.Placeholder = "Password"
	passwordInput.EchoMode = textinput.EchoPassword
	passwordInput.CharLimit = 32
	passwordInput.Width = 30
	inputs = append(inputs, passwordInput)
	labels = append(labels, "Password")

	return inputs, labels
}

// createConfigUpdateForm creates form for config update command
func createConfigUpdateForm() ([]textinput.Model, []string) {
	var inputs []textinput.Model
	var labels []string

	// Host input
	hostInput := textinput.New()
	hostInput.Placeholder = "localhost"
	hostInput.CharLimit = 64
	hostInput.Width = 30
	hostInput.SetValue(config.ArgsConfigInstance.Address)
	inputs = append(inputs, hostInput)
	labels = append(labels, "Host")

	// Port input
	portInput := textinput.New()
	portInput.Placeholder = "8080"
	portInput.CharLimit = 5
	portInput.Width = 10
	portInput.SetValue(fmt.Sprintf("%d", config.ArgsConfigInstance.Port))
	inputs = append(inputs, portInput)
	labels = append(labels, "Port")

	// Username input
	usernameInput := textinput.New()
	usernameInput.Placeholder = "cookieguest"
	usernameInput.CharLimit = 32
	usernameInput.Width = 30
	usernameInput.SetValue(config.ArgsConfigInstance.Username)
	inputs = append(inputs, usernameInput)
	labels = append(labels, "Username")

	// HTTPS input
	httpsInput := textinput.New()
	httpsInput.Placeholder = "false"
	httpsInput.CharLimit = 5
	httpsInput.Width = 10
	if config.ArgsConfigInstance.HTTPS {
		httpsInput.SetValue("true")
	} else {
		httpsInput.SetValue("false")
	}
	inputs = append(inputs, httpsInput)
	labels = append(labels, "HTTPS (true/false)")

	return inputs, labels
}

// createExploitRunForm creates form for exploit run command
func createExploitRunForm() ([]textinput.Model, []string) {
	var inputs []textinput.Model
	var labels []string

	// Exploit path input
	exploitPathInput := textinput.New()
	exploitPathInput.Placeholder = "Path to exploit file"
	exploitPathInput.CharLimit = 256
	exploitPathInput.Width = 50
	inputs = append(inputs, exploitPathInput)
	labels = append(labels, "Exploit Path")

	// Service port input
	servicePortInput := textinput.New()
	servicePortInput.Placeholder = "Service port number"
	servicePortInput.CharLimit = 5
	servicePortInput.Width = 10
	inputs = append(inputs, servicePortInput)
	labels = append(labels, "Service Port")

	// Detach mode input
	detachInput := textinput.New()
	detachInput.Placeholder = "false"
	detachInput.CharLimit = 5
	detachInput.Width = 10
	detachInput.SetValue("false")
	inputs = append(inputs, detachInput)
	labels = append(labels, "Detach Mode (true/false)")

	// Tick time input
	tickTimeInput := textinput.New()
	tickTimeInput.Placeholder = "120"
	tickTimeInput.CharLimit = 4
	tickTimeInput.Width = 10
	tickTimeInput.SetValue("120")
	inputs = append(inputs, tickTimeInput)
	labels = append(labels, "Tick Time (seconds)")

	// Thread count input
	threadCountInput := textinput.New()
	threadCountInput.Placeholder = "5"
	threadCountInput.CharLimit = 3
	threadCountInput.Width = 10
	threadCountInput.SetValue("5")
	inputs = append(inputs, threadCountInput)
	labels = append(labels, "Thread Count")

	return inputs, labels
}

// createExploitNameForm creates form for exploit create/remove commands
func createExploitNameForm() ([]textinput.Model, []string) {
	var inputs []textinput.Model
	var labels []string

	// Exploit name input
	nameInput := textinput.New()
	nameInput.Placeholder = "Exploit name"
	nameInput.CharLimit = 64
	nameInput.Width = 40
	inputs = append(inputs, nameInput)
	labels = append(labels, "Exploit Name")

	return inputs, labels
}

// createExploitStopForm creates form for exploit stop command
func createExploitStopForm() ([]textinput.Model, []string) {
	var inputs []textinput.Model
	var labels []string

	// PID input
	pidInput := textinput.New()
	pidInput.Placeholder = "Process ID"
	pidInput.CharLimit = 10
	pidInput.Width = 15
	inputs = append(inputs, pidInput)
	labels = append(labels, "Process ID (PID)")

	return inputs, labels
}

// ValidateForm validates form inputs for a specific command
func ValidateForm(command string, inputs []textinput.Model) error {
	switch command {
	case "config login":
		return validateLoginForm(inputs)
	case "config update":
		return validateConfigUpdateForm(inputs)
	case "exploit run":
		return validateExploitRunForm(inputs)
	case "exploit create", "exploit remove":
		return validateExploitNameForm(inputs)
	case "exploit stop":
		return validateExploitStopForm(inputs)
	}
	return nil
}

// validateLoginForm validates login form
func validateLoginForm(inputs []textinput.Model) error {
	if len(inputs) < 2 {
		return fmt.Errorf("invalid form structure")
	}

	username := strings.TrimSpace(inputs[0].Value())
	password := strings.TrimSpace(inputs[1].Value())

	if username == "" {
		return fmt.Errorf("username is required")
	}
	if password == "" {
		return fmt.Errorf("password is required")
	}

	return nil
}

// validateConfigUpdateForm validates config update form
func validateConfigUpdateForm(inputs []textinput.Model) error {
	if len(inputs) < 4 {
		return fmt.Errorf("invalid form structure")
	}

	host := strings.TrimSpace(inputs[0].Value())
	port := strings.TrimSpace(inputs[1].Value())
	username := strings.TrimSpace(inputs[2].Value())
	httpsStr := strings.TrimSpace(inputs[3].Value())

	// At least one field should be provided
	if host == "" && port == "" && username == "" && httpsStr == "" {
		return fmt.Errorf("at least one field must be provided")
	}

	// Validate HTTPS field
	if httpsStr != "" && strings.ToLower(httpsStr) != "true" && strings.ToLower(httpsStr) != "false" {
		return fmt.Errorf("HTTPS field must be 'true' or 'false'")
	}

	return nil
}

// validateExploitRunForm validates exploit run form
func validateExploitRunForm(inputs []textinput.Model) error {
	if len(inputs) < 5 {
		return fmt.Errorf("invalid form structure")
	}

	exploitPath := strings.TrimSpace(inputs[0].Value())
	servicePort := strings.TrimSpace(inputs[1].Value())
	detachStr := strings.TrimSpace(inputs[2].Value())

	if exploitPath == "" {
		return fmt.Errorf("exploit path is required")
	}
	if servicePort == "" {
		return fmt.Errorf("service port is required")
	}

	// Validate detach field
	if detachStr != "" && strings.ToLower(detachStr) != "true" && strings.ToLower(detachStr) != "false" {
		return fmt.Errorf("detach mode must be 'true' or 'false'")
	}

	return nil
}

// validateExploitNameForm validates exploit name form
func validateExploitNameForm(inputs []textinput.Model) error {
	if len(inputs) < 1 {
		return fmt.Errorf("invalid form structure")
	}

	name := strings.TrimSpace(inputs[0].Value())
	if name == "" {
		return fmt.Errorf("exploit name is required")
	}

	return nil
}

// validateExploitStopForm validates exploit stop form
func validateExploitStopForm(inputs []textinput.Model) error {
	if len(inputs) < 1 {
		return fmt.Errorf("invalid form structure")
	}

	pid := strings.TrimSpace(inputs[0].Value())
	if pid == "" {
		return fmt.Errorf("process ID is required")
	}

	return nil
}

// GetFormData extracts data from form inputs
func GetFormData(inputs []textinput.Model, labels []string) FormData {
	data := FormData{Fields: make(map[string]string)}

	for i, input := range inputs {
		if i < len(labels) {
			data.Fields[labels[i]] = strings.TrimSpace(input.Value())
		}
	}

	return data
}

// NavigateForm handles navigation within forms
func NavigateForm(inputs []textinput.Model, currentFocus int, direction int) int {
	if len(inputs) == 0 {
		return 0
	}

	// Calculate new focus index
	newFocus := currentFocus + direction
	if newFocus < 0 {
		newFocus = len(inputs) - 1
	} else if newFocus >= len(inputs) {
		newFocus = 0
	}

	// Update focus
	for i := range inputs {
		if i == newFocus {
			inputs[i].Focus()
		} else {
			inputs[i].Blur()
		}
	}

	return newFocus
}

// UpdateFormInputs updates form inputs with a tea.Msg
func UpdateFormInputs(inputs []textinput.Model, focusIndex int, msg tea.Msg) ([]textinput.Model, tea.Cmd) {
	var cmd tea.Cmd

	if focusIndex >= 0 && focusIndex < len(inputs) {
		inputs[focusIndex], cmd = inputs[focusIndex].Update(msg)
	}

	return inputs, cmd
}
