# CookieFarm TUI Guide

The CookieFarm Interactive TUI (Text User Interface) provides an easy way to use the CookieFarm client. This document explains how to use the TUI, its features, and keyboard shortcuts.

## Overview

The TUI provides a visually appealing, interactive interface for all CookieFarm client operations without changing the underlying functionality. It offers:

- Intuitive navigation with keyboard controls
- Colorful, organized command menus
- Form-based input for command parameters
- Visual command output display

## Usage

When launching the CookieFarm client, the TUI will start by default:

```bash
./cookieclient
```

If you prefer the traditional CLI interface (e.g., for scripting), you can disable the TUI with:

```bash
./cookieclient --no-tui
# OR
COOKIECLIENT_NO_TUI=1 ./cookieclient
```

## Navigation

The TUI consists of several views:

1. **Main Menu**: Shows top-level commands (Configuration, Exploits, Quit)
2. **Configuration Menu**: Commands for managing client configuration
3. **Exploit Menu**: Commands for creating and running exploits
4. **Input Forms**: For commands requiring parameters
5. **Output View**: Displays command results

### Keyboard Controls

- **↑/↓** or **j/k**: Navigate menu items
- **Enter**: Select menu item or submit form
- **ESC**: Go back to previous view
- **q** or **Ctrl+C**: Quit the application
- **Tab**: Navigate between input fields (when in form view)

## Main Features

### Configuration Management

The Configuration menu lets you:
- View current config (`Show Config`)
- Update configuration settings (`Update Config`)
- Reset to default configuration (`Reset Config`)
- Login/logout from the server

### Exploit Management

The Exploit menu lets you:
- Run exploits against other teams (`Run Exploit`)
- Create new exploit templates (`Create Exploit`)
- List all running exploits (`List Exploits`)
- Stop running exploits (`Stop Exploit`)
- Remove exploit templates (`Remove Exploit`)

## Form Input

When commands require parameters, the TUI will display input forms with:
- Labeled fields with placeholders
- Error validation for required fields
- Tab navigation between fields
- Submission with Enter key

## Troubleshooting

If the TUI doesn't render correctly:
- Check your terminal supports ANSI colors and Unicode
- Ensure your terminal window is large enough
- Try falling back to CLI mode with `--no-tui` flag

## Architecture

The TUI is built using a modular architecture with the following components:

### Core Modules
- **`model.go`**: Data models and state management
- **`menu.go`**: Menu creation and management
- **`forms.go`**: Dynamic form creation and validation
- **`commandrunner.go`**: Isolated command execution
- **`handlers.go`**: Event and command handling
- **`view.go`**: Rendering and display logic

### Supporting Modules
- **`styles.go`**: Visual styling and theming
- **`utils.go`**: Common utility functions
- **`logger.go`**: TUI-specific logging
- **`tui.go`**: Main orchestrator and entry point

### Key Benefits
- **Modular Design**: Each component has a specific responsibility
- **Easy Maintenance**: Clean separation of concerns
- **Extensible**: Simple to add new commands and features
- **Robust Command Execution**: Isolated execution with proper output capture
- **Consistent Styling**: Centralized theming system

For detailed architectural information, see [`ARCHITECTURE.md`](ARCHITECTURE.md).

## Implementation Notes

The TUI is built using the Charm libraries:
- [Bubbletea](https://github.com/charmbracelet/bubbletea): TUI framework
- [Bubbles](https://github.com/charmbracelet/bubbles): TUI components
- [Lipgloss](https://github.com/charmbracelet/lipgloss): Styling

### Command Execution
Commands are executed in isolated environments to prevent interference with the TUI state. The `CommandRunner` module handles:
- Proper argument passing
- Output and error capture
- State preservation
- Clean execution environments

### Form Management
Dynamic forms are created based on the selected command with:
- Automatic field generation
- Input validation
- Focus management
- Error handling