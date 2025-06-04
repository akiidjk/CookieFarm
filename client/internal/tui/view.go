package tui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ViewRenderer handles all rendering operations for the TUI
type ViewRenderer struct {
	width  int
	height int
}

// NewViewRenderer creates a new view renderer
func NewViewRenderer() *ViewRenderer {
	return &ViewRenderer{}
}

// SetSize updates the renderer's dimensions
func (r *ViewRenderer) SetSize(width, height int) {
	r.width = width
	r.height = height
}

// RenderView renders the main view based on the model state
func (r *ViewRenderer) RenderView(m *Model) string {
	if m.quitting {
		return r.renderQuitMessage()
	}

	switch {
	case m.IsRunningCommand():
		return r.renderCommandOutput(m)
	case m.IsInputMode():
		return r.renderInputForm(m)
	default:
		return r.renderMenu(m)
	}
}

// renderExploitTable renders the exploit list in a nice table format
func (r *ViewRenderer) renderExploitTable(m *Model, banner, title string) string {
	// Adjust table size to fit the window
	m.exploitTable.SetWidth(r.width - 10)
	m.exploitTable.SetHeight(r.height - 15)

	// Render the table
	tableView := m.exploitTable.View()

	// Instructions
	instructions := FooterStyle.Render("↑/↓: Navigate table • Enter: Select exploit to stop • ESC: Back")

	// Combine all elements
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		banner,
		title,
		"",
		tableView,
		"",
		instructions,
	)

	// Add error message if present
	if m.err != nil {
		errorMsg := ErrorStyle.Render("Error: " + m.err.Error())
		content = lipgloss.JoinVertical(lipgloss.Left, content, "", errorMsg)
	}

	return content
}

// renderBanner renders the application banner
func (*ViewRenderer) renderBanner(banner string) string {
	if banner == "" {
		return ""
	}

	bannerLines := strings.Split(banner, "\n")
	var formattedLines []string

	for _, line := range bannerLines {
		if strings.TrimSpace(line) != "" {
			formattedLines = append(formattedLines, BannerStyle.Render(line))
		}
	}

	if len(formattedLines) == 0 {
		return ""
	}

	return strings.Join(formattedLines, "\n") + "\n"
}

// renderMenu renders the appropriate menu based on active view
func (r *ViewRenderer) renderMenu(m *Model) string {
	banner := r.renderBanner(m.banner)

	var menuView string
	switch m.GetActiveView() {
	case "main":
		menuView = m.mainMenu.View()
	case "config":
		menuView = m.configMenu.View()
	case "exploit":
		menuView = m.exploitMenu.View()
	default:
		menuView = m.mainMenu.View()
	}

	content := lipgloss.JoinVertical(lipgloss.Left, banner, menuView)
	helpView := r.renderHelpFooter()

	return lipgloss.JoinVertical(lipgloss.Left, content, helpView)
}

// renderInputForm renders the input form interface
func (r *ViewRenderer) renderInputForm(m *Model) string {
	banner := r.renderBanner(m.banner)

	// Command title
	commandTitle := SubtitleStyle.Render("Command: " + m.activeCommand)

	// For exploit stop command, show process selection list
	if m.IsProcessListVisible() && m.activeCommand == "exploit stop" {
		return r.renderProcessSelectionList(m, banner, commandTitle)
	}

	// Render form inputs
	var inputViews []string
	for i, input := range m.inputs {
		if i < len(m.inputLabels) {
			var label string
			if i == m.focusIndex {
				// Highlight focused input label
				label = lipgloss.NewStyle().
					Foreground(primaryColor).
					Bold(true).
					Render(fmt.Sprintf("▶ %s: ", m.inputLabels[i]))
			} else {
				label = InputLabelStyle.Render(fmt.Sprintf("  %s: ", m.inputLabels[i]))
			}

			inputView := lipgloss.JoinHorizontal(lipgloss.Left, label, input.View())
			inputViews = append(inputViews, inputView)
		}
	}

	formContent := lipgloss.JoinVertical(lipgloss.Left, inputViews...)

	// Instructions
	instructions := FooterStyle.Render("Tab: Next field • Shift+Tab: Previous field • Enter: Submit • ESC: Cancel")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		banner,
		commandTitle,
		"",
		formContent,
		"",
		instructions,
	)

	// Add error message if present
	if m.err != nil {
		errorMsg := ErrorStyle.Render("Error: " + m.err.Error())
		content = lipgloss.JoinVertical(lipgloss.Left, content, "", errorMsg)
	}

	return content
}

// renderProcessSelectionList renders the exploit process selection list
func (*ViewRenderer) renderProcessSelectionList(m *Model, banner, commandTitle string) string {
	// Table header
	header := lipgloss.JoinVertical(
		lipgloss.Left,
		SubtitleStyle.Render("Select a running exploit to stop:"),
		"",
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			lipgloss.NewStyle().Width(6).Render("ID"),
			lipgloss.NewStyle().Width(40).Render("Name"),
			lipgloss.NewStyle().Width(10).Render("PID"),
		),
		lipgloss.NewStyle().Foreground(mutedColor).Render(strings.Repeat("─", 56)),
	)

	// Render process list
	var processViews []string
	for i, process := range m.exploitProcesses {
		var style lipgloss.Style
		if i == m.selectedProcess {
			// Highlight selected process
			style = lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true)

			// Add selection indicator
			processViews = append(processViews, lipgloss.JoinHorizontal(
				lipgloss.Left,
				style.Width(6).Render(strconv.Itoa(process.ID)),
				style.Width(40).Render(process.Name),
				style.Width(10).Render(strconv.Itoa(process.PID)),
			))
		} else {
			style = lipgloss.NewStyle()
			processViews = append(processViews, lipgloss.JoinHorizontal(
				lipgloss.Left,
				style.Width(6).Render(strconv.Itoa(process.ID)),
				style.Width(40).Render(process.Name),
				style.Width(10).Render(strconv.Itoa(process.PID)),
			))
		}
	}

	processList := lipgloss.JoinVertical(lipgloss.Left, processViews...)

	// Instructions
	instructions := FooterStyle.Render("↑/↓: Navigate • Enter: Select process to stop • ESC: Cancel")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		banner,
		commandTitle,
		"",
		header,
		processList,
		"",
		instructions,
	)

	// Add error message if present
	if m.err != nil {
		errorMsg := ErrorStyle.Render("Error: " + m.err.Error())
		content = lipgloss.JoinVertical(lipgloss.Left, content, "", errorMsg)
	}

	return content
}

// renderCommandOutput renders the command execution output
func (r *ViewRenderer) renderCommandOutput(m *Model) string {
	banner := r.renderBanner(m.banner)

	// Command title
	commandTitle := SubtitleStyle.Render("Command Output:")

	// For exploit list, show table view instead of text output
	if m.activeCommand == "exploit list" && m.showTable {
		return r.renderExploitTable(m, banner, commandTitle)
	}

	// Format output with styling
	formattedOutput := r.formatCommandOutput(m.commandOutput)

	// Create output box with or without spinner
	var outputContent string
	if m.loading {
		loadingText := LoadingStyle.Render(" Loading...")
		spinner := m.spinner.View() + loadingText
		if m.commandOutput != "" {
			outputContent = spinner + "\n\n" + formattedOutput
		} else {
			outputContent = spinner
		}
	} else {
		outputContent = formattedOutput
	}

	// Create output box
	outputBox := CommandOutputStyle.
		Width(r.width - 4).
		Render(outputContent)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		banner,
		commandTitle,
		"",
		outputBox,
	)

	// Add error message if present
	if m.err != nil {
		errorMsg := ErrorStyle.Render("Error: " + m.err.Error())
		content = lipgloss.JoinVertical(lipgloss.Left, content, "", errorMsg)
	}

	// Instructions
	instructions := FooterStyle.Render("Press ESC to go back")
	content = lipgloss.JoinVertical(lipgloss.Left, content, "", instructions)

	return content
}

// formatCommandOutput applies styling to command output
func (r *ViewRenderer) formatCommandOutput(output string) string {
	if output == "" {
		return DimmedStyle("No output yet")
	}

	lines := strings.Split(output, "\n")
	formattedLines := make([]string, 0, len(lines))

	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			formattedLines = append(formattedLines, "")
			continue
		}

		// Add line numbers
		lineNum := lipgloss.NewStyle().
			Foreground(mutedColor).
			Width(4).
			Align(lipgloss.Right).
			Render(strconv.Itoa(i + 1))

		// Style line based on content
		styledLine := r.styleOutputLine(line)

		formattedLine := lipgloss.JoinHorizontal(
			lipgloss.Left,
			lineNum,
			lipgloss.NewStyle().Foreground(mutedColor).Render(" │ "),
			styledLine,
		)

		formattedLines = append(formattedLines, formattedLine)
	}

	return strings.Join(formattedLines, "\n")
}

// styleOutputLine applies appropriate styling to individual output lines
func (*ViewRenderer) styleOutputLine(line string) string {
	lowerLine := strings.ToLower(line)

	switch {
	case strings.Contains(lowerLine, "error") || strings.Contains(lowerLine, "failed") || strings.Contains(lowerLine, "fatal"):
		return ErrorText(line)
	case strings.Contains(lowerLine, "warning") || strings.Contains(lowerLine, "warn"):
		return WarningText(line)
	case strings.Contains(lowerLine, "success") || strings.Contains(lowerLine, "completed") || strings.Contains(lowerLine, "done"):
		return SuccessText(line)
	case strings.Contains(lowerLine, "info") || strings.Contains(lowerLine, "starting") || strings.Contains(lowerLine, "loading"):
		return InfoText(line)
	case strings.Contains(lowerLine, "debug"):
		return lipgloss.NewStyle().Faint(true).Render(line)
	default:
		return line
	}
}

// renderHelpFooter renders the help/instruction footer
func (*ViewRenderer) renderHelpFooter() string {
	helpText := "↑/↓: Navigate • Enter: Select • q: Quit • ESC: Back"
	return HelpStyle.Render(helpText)
}

// renderQuitMessage renders the quit message
func (r *ViewRenderer) renderQuitMessage() string {
	message := lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Align(lipgloss.Center).
		Width(50).
		Render("Thanks for using CookieFarm!")

	return lipgloss.Place(
		r.width,
		r.height,
		lipgloss.Center,
		lipgloss.Center,
		message,
	)
}
