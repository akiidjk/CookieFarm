package tui

import (
	"strings"
	"time"

	"github.com/ByteTheCookies/cookieclient/internal/config"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// New creates a new TUI model
func New(banner string) Model {
	// Initialize menus
	mainMenu, configMenu, exploitMenu := InitializeMenus()

	return Model{
		activeView:  "main",
		mainMenu:    mainMenu,
		configMenu:  configMenu,
		exploitMenu: exploitMenu,
		help:        help.New(),
		banner:      banner,
		cmdRunner:   NewCommandRunner(),
	}
}

type tickMsg struct{}

// Init initializes the TUI
func (Model) Init() tea.Cmd {
	return nil
}

// Update handles TUI state updates
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		UpdateMenuSize(&m.mainMenu, &m.configMenu, &m.exploitMenu, msg.Width, msg.Height)
		return m, nil

	case tea.KeyMsg:
		// Handle input mode first
		if m.IsInputMode() {
			return m.handleInputMode(msg)
		}

		// Handle special keys
		newM, cmdH, handled := m.handleKeyPress(msg)
		if handled {
			return newM, cmdH
		}

		// If not handled, let menus handle the key (like arrow keys)
		switch m.GetActiveView() {
		case "main":
			m.mainMenu, cmdH = m.mainMenu.Update(msg)
		case "config":
			m.configMenu, cmdH = m.configMenu.Update(msg)
		case "exploit":
			m.exploitMenu, cmdH = m.exploitMenu.Update(msg)
		}
		return m, cmdH

	case CommandOutput:
		m.SetCommandOutput(msg.Output)
		m.SetRunningCommand(true)
		if msg.Error != nil {
			m.SetError(msg.Error)
		}
		return m, nil

	case ExploitOutput:
		// Append to existing output for exploit commands
		if m.IsRunningCommand() && strings.HasPrefix(m.activeCommand, "exploit") {
			if msg.Content != "" {
				m.AppendCommandOutput(msg.Content)
			}
			if msg.Error != nil {
				m.SetError(msg.Error)
			}
			// Continue streaming output - this is critical for continuous updates
			return m, tea.Batch(
				tea.Tick(100*time.Millisecond, func(time.Time) tea.Msg {
					return tickMsg{}
				}),
				m.cmdRunner.GetExploitOutputCmd(),
			)
		}
		return m, nil
	}

	// Handle other message types and menu updates for non-KeyMsg events
	switch m.GetActiveView() {
	case "main":
		m.mainMenu, cmd = m.mainMenu.Update(msg)
	case "config":
		m.configMenu, cmd = m.configMenu.Update(msg)
	case "exploit":
		m.exploitMenu, cmd = m.exploitMenu.Update(msg)
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

// handleKeyPress handles keyboard input
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd, bool) {
	handler := NewCommandHandler()

	switch {
	case key.Matches(msg, Keys.Quit):
		m.quitting = true
		return m, tea.Quit, true

	case key.Matches(msg, Keys.Back):
		newM, cmd := m.handleBackAction()
		return newM, cmd, true

	case key.Matches(msg, Keys.Enter):
		newM, cmd := m.handleEnterAction(handler)
		return newM, cmd, true
	}

	// Return false to indicate the key wasn't handled
	return m, nil, false
}

// handleBackAction handles the back/escape key
func (m Model) handleBackAction() (tea.Model, tea.Cmd) {
	if m.IsInputMode() {
		m.SetInputMode(false)
		m.ClearError()
		return m, nil
	}

	if m.IsRunningCommand() {
		m.SetRunningCommand(false)
		m.SetCommandOutput("")
		m.ClearError()
		return m, nil
	}

	if m.GetActiveView() != "main" {
		m.SetActiveView("main")
		return m, nil
	}

	return m, nil
}

// handleEnterAction handles the enter key
func (m Model) handleEnterAction(handler *CommandHandler) (tea.Model, tea.Cmd) {
	if m.IsInputMode() {
		newModel, cmd := handler.ProcessFormSubmission(&m)

		// If we're running an exploit command, also start the streaming output
		if strings.HasPrefix(newModel.activeCommand, "exploit run") {
			return newModel, tea.Batch(cmd, newModel.GetExploitStreamCmd())
		}

		return newModel, cmd
	}

	selectedItem, ok := m.getSelectedMenuItem()
	if !ok {
		return m, nil
	}

	newModel, cmd := m.processMenuSelection(selectedItem, handler)

	// If we're running an exploit command, also start the streaming output
	if mTyped, ok := newModel.(Model); ok && mTyped.IsRunningCommand() && strings.HasPrefix(mTyped.activeCommand, "exploit run") {
		return newModel, tea.Batch(cmd, mTyped.GetExploitStreamCmd())
	}

	return newModel, cmd
}

// getSelectedMenuItem gets the currently selected menu item
func (m Model) getSelectedMenuItem() (menuItem, bool) {
	switch m.GetActiveView() {
	case "main":
		return GetSelectedItem(m.mainMenu)
	case "config":
		return GetSelectedItem(m.configMenu)
	case "exploit":
		return GetSelectedItem(m.exploitMenu)
	default:
		return menuItem{}, false
	}
}

// processMenuSelection processes the selected menu item
func (m Model) processMenuSelection(selectedItem menuItem, handler *CommandHandler) (tea.Model, tea.Cmd) {
	command := selectedItem.command

	// Handle navigation commands
	if IsNavigationCommand(command) {
		return handler.HandleNavigation(command, &m)
	}

	// Handle direct commands (no input required)
	if IsDirectCommand(command) {
		m.activeCommand = command
		m.SetRunningCommand(true)
		return m, handler.HandleCommand(command, nil)
	}

	// Handle commands that require input
	if RequiresInput(command) {
		handler.SetupFormForCommand(&m, command)
		return m, nil
	}

	return m, nil
}

// handleInputMode handles updates when in input mode
func (m Model) handleInputMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.Type {
		case tea.KeyTab:
			m.focusIndex = NavigateForm(m.inputs, m.focusIndex, 1)
			return m, nil
		case tea.KeyShiftTab:
			m.focusIndex = NavigateForm(m.inputs, m.focusIndex, -1)
			return m, nil
		case tea.KeyEnter:
			// Process form submission
			handler := NewCommandHandler()
			return handler.ProcessFormSubmission(&m)
		case tea.KeyEscape:
			// Cancel form input
			m.SetInputMode(false)
			m.ClearError()
			return m, nil
		case tea.KeyUp, tea.KeyDown, tea.KeyLeft, tea.KeyRight:
			// Pass arrow keys to the focused input field
			m.inputs, cmd = UpdateFormInputs(m.inputs, m.focusIndex, msg)
			return m, cmd
		}
	}

	// Update the focused input
	m.inputs, cmd = UpdateFormInputs(m.inputs, m.focusIndex, msg)
	return m, cmd
}

// View renders the TUI
func (m Model) View() string {
	renderer := NewViewRenderer()
	renderer.SetSize(m.width, m.height)
	return renderer.RenderView(&m)
}

// GetExploitStreamCmd returns a command that periodically checks for exploit output
func (m Model) GetExploitStreamCmd() tea.Cmd {
	if m.IsRunningCommand() && strings.HasPrefix(m.activeCommand, "exploit run") {
		return m.cmdRunner.GetExploitOutputCmd()
	}
	return nil
}

// StartTUI launches the TUI application
func StartTUI(banner string) error {
	// Set TUI mode in config
	config.NoTUI = true

	p := tea.NewProgram(
		New(banner),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	_, err := p.Run()
	return err
}
