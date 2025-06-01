package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
)

// Model is the main TUI model
type Model struct {
	activeView     string
	mainMenu       list.Model
	configMenu     list.Model
	exploitMenu    list.Model
	inputs         []textinput.Model
	inputLabels    []string
	focusIndex     int
	inputMode      bool
	err            error
	help           help.Model
	quitting       bool
	width, height  int
	banner         string
	activeCommand  string
	runningCommand bool
	commandOutput  string
	cmdRunner      *CommandRunner
}

// CommandOutput represents the result of a command execution
type CommandOutput struct {
	Output string
	Error  error
}

// menuItem represents an item in the menu
type menuItem struct {
	title       string
	description string
	command     string
}

func (i menuItem) Title() string       { return i.title }
func (i menuItem) Description() string { return i.description }
func (i menuItem) FilterValue() string { return i.title }

// keyMap defines keybindings
type keyMap struct {
	Up        key.Binding
	Down      key.Binding
	Enter     key.Binding
	Quit      key.Binding
	Back      key.Binding
	Tab       key.Binding
	Submit    key.Binding
	NextInput key.Binding
	PrevInput key.Binding
}

// Keys defines the keybindings
var Keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next input"),
	),
	Submit: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "submit"),
	),
	NextInput: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next input"),
	),
	PrevInput: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "previous input"),
	),
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Enter, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Enter},
		{k.Quit, k.Back},
	}
}

// SetError sets an error message for the model
func (m *Model) SetError(err error) {
	m.err = err
}

// ClearError clears any error message
func (m *Model) ClearError() {
	m.err = nil
}

// IsInputMode returns true if the model is in input mode
func (m *Model) IsInputMode() bool {
	return m.inputMode
}

// IsRunningCommand returns true if a command is currently running
func (m *Model) IsRunningCommand() bool {
	return m.runningCommand
}

// SetInputMode sets the input mode state
func (m *Model) SetInputMode(state bool) {
	m.inputMode = state
}

// SetRunningCommand sets the running command state
func (m *Model) SetRunningCommand(state bool) {
	m.runningCommand = state
}

// SetCommandOutput sets the command output
func (m *Model) SetCommandOutput(output string) {
	m.commandOutput = output
}

// GetActiveView returns the current active view
func (m *Model) GetActiveView() string {
	return m.activeView
}

// SetActiveView sets the active view
func (m *Model) SetActiveView(view string) {
	m.activeView = view
}