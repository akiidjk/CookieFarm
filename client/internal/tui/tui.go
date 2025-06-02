package tui

import (
	"strconv"
	"strings"
	"time"

	"github.com/ByteTheCookies/cookieclient/internal/config"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

// New creates a new TUI model
func New(banner string) Model {
	// Initialize menus
	mainMenu, configMenu, exploitMenu := InitializeMenus()

	// Initialize spinner
	s := spinner.New()
	s.Spinner = spinner.Points
	s.Style = SpinnerStyle

	// Initialize exploit table with columns
	exploitTable := table.New(
		table.WithColumns([]table.Column{
			{Title: "ID", Width: 6},
			{Title: "NAME", Width: 40},
			{Title: "PID", Width: 10},
		}),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	// Style the table
	exploitTable.SetStyles(table.Styles{
		Header:   TableHeaderStyle,
		Selected: TableSelectedRowStyle,
		Cell:     TableRowStyle,
	})

	return Model{
		activeView:   "main",
		mainMenu:     mainMenu,
		configMenu:   configMenu,
		exploitMenu:  exploitMenu,
		help:         help.New(),
		banner:       banner,
		cmdRunner:    NewCommandRunner(),
		spinner:      s,
		exploitTable: exploitTable,
	}
}

type tickMsg struct{}

// Init initializes the TUI
func (m Model) Init() tea.Cmd {
	return m.spinner.Tick
}

// TableUpdateMsg is a message to update the table data
type TableUpdateMsg struct {
	Rows []table.Row
	Show bool
}

// Update handles TUI state updates
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	// Handle table update message
	if updateMsg, ok := msg.(TableUpdateMsg); ok {
		m.exploitTable.SetRows(updateMsg.Rows)
		m.showTable = updateMsg.Show
		return m, nil
	}

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

		// Handle table navigation when table is visible
		if m.showTable && (m.activeCommand == "exploit list" || m.activeCommand == "exploit stop") {
			switch msg.String() {
			case "enter":
				if m.activeCommand == "exploit stop" {
					// Process form submission with selected table row
					handler := NewCommandHandler()
					m.SetLoading(true)
					return handler.ProcessFormSubmission(&m)
				}
				return m, nil
			case "esc", "q":
				// Exit table view via handleBackAction for consistent state management
				return m.handleBackAction()
			}
			// Let the table handle other keys like up/down for table navigation
			m.exploitTable, cmd = m.exploitTable.Update(msg)
			return m, cmd
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
		m.SetLoading(false) // Command completed, stop loading spinner
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

	// Handle spinner updates
	if m.loading {
		var spinnerCmd tea.Cmd
		m.spinner, spinnerCmd = m.spinner.Update(msg)
		cmds = append(cmds, spinnerCmd)
	}

	// Handle table updates
	if m.showTable && (m.activeCommand == "exploit list" || m.activeCommand == "exploit stop") {
		var tableCmd tea.Cmd
		m.exploitTable, tableCmd = m.exploitTable.Update(msg)
		cmds = append(cmds, tableCmd)
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
	// First, check if we're in table view mode
	if m.showTable {
		m.showTable = false
		m.ClearError()

		// If it was for stop command, also exit input mode
		if m.activeCommand == "exploit stop" {
			m.SetInputMode(false)
		}

		// Clear running command state if it's not a streaming command
		if !strings.HasPrefix(m.activeCommand, "exploit run") {
			m.SetRunningCommand(false)
		}

		return m, nil
	}

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
		m.SetLoading(true) // Show loading spinner when submitting the form
		newModel, cmd := handler.ProcessFormSubmission(&m)

		// If we're running an exploit command, also start the streaming output
		if strings.HasPrefix(newModel.activeCommand, "exploit run") {
			return newModel, tea.Batch(cmd, newModel.GetExploitStreamCmd(), m.spinner.Tick)
		}

		return newModel, tea.Batch(cmd, m.spinner.Tick)
	}

	selectedItem, ok := m.getSelectedMenuItem()
	if !ok {
		return m, nil
	}

	newModel, cmd := m.processMenuSelection(selectedItem, handler)

	// If we're running an exploit command, also start the streaming output
	if mTyped, ok := newModel.(Model); ok && mTyped.IsRunningCommand() {
		if strings.HasPrefix(mTyped.activeCommand, "exploit run") {
			return newModel, tea.Batch(cmd, mTyped.GetExploitStreamCmd())
		}
		// For exploit list, setup the table
		if mTyped.activeCommand == "exploit list" {
			return newModel, tea.Batch(cmd, mTyped.SetupExploitTableCmd())
		}
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
		m.SetLoading(true)

		// For exploit list command, set up the table view
		if command == "exploit list" {
			return m, tea.Batch(
				handler.HandleCommand(command, nil),
				m.SetupExploitTableCmd(),
			)
		}

		return m, handler.HandleCommand(command, nil)
	}

	// Handle commands that require input
	if RequiresInput(command) {
		handler.SetupFormForCommand(&m, command)
		return m, nil
	}

	// We handle exploit list commands in the direct command section now

	return m, nil
}

// handleInputMode handles updates when in input mode
func (m Model) handleInputMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		// Special handling for process selection list
		if m.IsProcessListVisible() && m.activeCommand == "exploit stop" {
			switch keyMsg.Type {
			case tea.KeyUp:
				m.SelectPreviousProcess()
				return m, nil
			case tea.KeyDown:
				m.SelectNextProcess()
				return m, nil
			case tea.KeyEnter:
				// Process form submission with selected exploit
				handler := NewCommandHandler()
				m.SetLoading(true)
				return handler.ProcessFormSubmission(&m)
			case tea.KeyEscape:
				// Cancel process selection; clear command and reset view to exploit menu.
				m.SetInputMode(false)
				m.SetProcessListVisible(false)
				m.ClearError()
				m.activeCommand = ""
				m.SetActiveView("exploit")
				return m, nil
			}
			return m, nil
		}

		// Normal input handling
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
			m.SetLoading(true) // Show loading spinner during form submission
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

// SetupExploitTableCmd returns a command to update the exploit table with data
func (m Model) SetupExploitTableCmd() tea.Cmd {
	return func() tea.Msg {
		// Short delay to allow command to complete
		time.Sleep(100 * time.Millisecond)

		// Get exploit data from command runner
		processes, err := m.cmdRunner.GetRunningExploits()
		if err != nil {
			return CommandOutput{
				Output: "Error fetching exploit data: " + err.Error(),
				Error:  err,
			}
		}

		if len(processes) == 0 {
			return CommandOutput{
				Output: "No running exploits found",
				Error:  nil,
			}
		}

		rows := []table.Row{}

		// Convert processes to table rows
		for _, p := range processes {
			rows = append(rows, table.Row{
				strconv.Itoa(p.ID),
				p.Name,
				strconv.Itoa(p.PID),
			})
		}

		// Reset the table's selected row to the first row
		modelCopy := m
		modelCopy.exploitTable.SetRows(rows)
		modelCopy.exploitTable.SetCursor(0) // Set cursor to first row
		modelCopy.showTable = true

		// Create a messenger that will update the model
		return func() tea.Msg {
			return TableUpdateMsg{
				Rows: rows,
				Show: true,
			}
		}()
	}
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
