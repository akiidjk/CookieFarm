package tui

import (
	"fmt"
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

// renderBanner renders the application banner
func (r *ViewRenderer) renderBanner(banner string) string {
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
	commandTitle := SubtitleStyle.Render(fmt.Sprintf("Command: %s", m.activeCommand))

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
		errorMsg := ErrorStyle.Render(fmt.Sprintf("Error: %s", m.err.Error()))
		content = lipgloss.JoinVertical(lipgloss.Left, content, "", errorMsg)
	}

	return content
}

// renderCommandOutput renders the command execution output
func (r *ViewRenderer) renderCommandOutput(m *Model) string {
	banner := r.renderBanner(m.banner)

	// Command title
	commandTitle := SubtitleStyle.Render("Command Output:")

	// Format output with styling
	formattedOutput := r.formatCommandOutput(m.commandOutput)

	// Create output box
	outputBox := CommandOutputStyle.
		Width(r.width - 4).
		Render(formattedOutput)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		banner,
		commandTitle,
		"",
		outputBox,
	)

	// Add error message if present
	if m.err != nil {
		errorMsg := ErrorStyle.Render(fmt.Sprintf("Error: %s", m.err.Error()))
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
		return DimmedStyle("No output")
	}

	lines := strings.Split(output, "\n")
	var formattedLines []string

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
			Render(fmt.Sprintf("%d", i+1))

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
func (r *ViewRenderer) styleOutputLine(line string) string {
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
func (r *ViewRenderer) renderHelpFooter() string {
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

// renderLoadingIndicator renders a loading indicator
func (r *ViewRenderer) renderLoadingIndicator(message string) string {
	spinner := "⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏"

	loadingText := lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Render(fmt.Sprintf("%c %s", spinner[0], message))

	return loadingText
}

// renderStatusBar renders a status bar with current state information
func (r *ViewRenderer) renderStatusBar(m *Model) string {
	var statusItems []string

	// Current view
	viewStatus := fmt.Sprintf("View: %s", strings.Title(m.GetActiveView()))
	statusItems = append(statusItems, viewStatus)

	// Command status
	if m.IsRunningCommand() {
		statusItems = append(statusItems, "Status: Running")
	} else if m.IsInputMode() {
		statusItems = append(statusItems, "Status: Input")
	} else {
		statusItems = append(statusItems, "Status: Ready")
	}

	statusText := strings.Join(statusItems, " • ")

	return lipgloss.NewStyle().
		Background(mutedColor).
		Foreground(textColor).
		Padding(0, 1).
		Width(r.width).
		Render(statusText)
}

// renderBox renders content within a styled box
func (r *ViewRenderer) renderBox(title, content string, style lipgloss.Style) string {
	if title != "" {
		titleStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			Padding(0, 1)

		styledTitle := titleStyle.Render(title)
		content = lipgloss.JoinVertical(lipgloss.Left, styledTitle, "", content)
	}

	return style.Render(content)
}

// centerContent centers content within the available space
func (r *ViewRenderer) centerContent(content string) string {
	return lipgloss.Place(
		r.width,
		r.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}
