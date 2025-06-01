# CookieFarm TUI Architecture Documentation

## Overview

The CookieFarm TUI (Text User Interface) is built using a modular architecture that separates concerns and provides a clean, maintainable codebase. The architecture leverages the Charm libraries (Bubbletea, Bubbles, and Lipgloss) to create an interactive command-line interface.

## Architecture Principles

- **Separation of Concerns**: Each module has a specific responsibility
- **Modularity**: Components can be developed and tested independently
- **Maintainability**: Clear interfaces between modules
- **Extensibility**: Easy to add new commands and features
- **Performance**: Efficient command execution and output capture

## Module Structure

### Core Modules

#### 1. `model.go` - Data Models and State Management
**Purpose**: Defines the core data structures and state management for the TUI.

**Key Components**:
- `Model`: Main TUI state container
- `CommandOutput`: Command execution result structure
- `menuItem`: Menu item representation
- `keyMap`: Keyboard bindings definition
- State management methods (SetError, ClearError, etc.)

**Responsibilities**:
- Maintain application state
- Provide state accessors and mutators
- Define keyboard bindings
- Manage UI component states

#### 2. `menu.go` - Menu Management
**Purpose**: Handles creation and management of all menu interfaces.

**Key Components**:
- `InitializeMenus()`: Creates all menus
- `createMainMenu()`, `createConfigMenu()`, `createExploitMenu()`: Specific menu creators
- `GetSelectedItem()`: Retrieves selected menu items
- `UpdateMenuSize()`: Handles menu resizing
- Menu classification functions (`IsNavigationCommand`, `IsDirectCommand`, etc.)

**Responsibilities**:
- Menu creation and configuration
- Menu item management
- Menu sizing and layout
- Command classification

#### 3. `forms.go` - Form Management
**Purpose**: Manages input forms for commands requiring user input.

**Key Components**:
- `CreateForm()`: Main form creation dispatcher
- `FormData`: Form data container
- Validation functions (`ValidateForm`, `validateLoginForm`, etc.)
- Navigation helpers (`NavigateForm`, `UpdateFormInputs`)

**Responsibilities**:
- Dynamic form creation based on commands
- Input validation
- Form navigation and focus management
- Data extraction from forms

#### 4. `commandrunner.go` - Command Execution
**Purpose**: Handles execution of Cobra commands within the TUI environment.

**Key Components**:
- `CommandRunner`: Main command executor
- `ExecuteCommand()`: Generic command execution
- Specialized execution methods (`ExecuteLogin`, `ExecuteConfigUpdate`, etc.)
- Output capture and error handling

**Responsibilities**:
- Isolated command execution
- Output and error capture
- Command argument preparation
- State preservation during execution

#### 5. `handlers.go` - Event and Command Handling
**Purpose**: Orchestrates command processing and event handling.

**Key Components**:
- `CommandHandler`: Main event handler
- `HandleCommand()`: Command dispatch logic
- `ProcessFormSubmission()`: Form processing
- `HandleNavigation()`: Navigation command processing

**Responsibilities**:
- Command routing and dispatch
- Form submission processing
- Navigation state management
- Event coordination

#### 6. `view.go` - Rendering and Display
**Purpose**: Handles all rendering and visual presentation logic.

**Key Components**:
- `ViewRenderer`: Main rendering engine
- `RenderView()`: Main view dispatcher
- Specialized renderers (`renderMenu`, `renderInputForm`, `renderCommandOutput`)
- Output formatting (`formatCommandOutput`, `styleOutputLine`)

**Responsibilities**:
- View rendering coordination
- Output styling and formatting
- Layout management
- Visual state representation

### Supporting Modules

#### 7. `styles.go` - Visual Styling
**Purpose**: Centralizes all visual styling and color definitions.

**Key Components**:
- Color constants (`primaryColor`, `secondaryColor`, etc.)
- Style definitions (`BannerStyle`, `ErrorStyle`, etc.)
- Text styling functions (`ErrorText`, `SuccessText`, etc.)

**Responsibilities**:
- Consistent visual theming
- Style definitions and management
- Color scheme coordination

#### 8. `utils.go` - Utility Functions
**Purpose**: Provides common utility functions used across modules.

**Key Components**:
- Terminal utilities (`GetTerminalSize`)
- Output formatting (`FormatOutput`, `FormatCommand`)
- File operations (`OpenEditor`)
- Configuration helpers (`LoadConfigToForms`)

**Responsibilities**:
- Common functionality provision
- System interaction helpers
- Utility function centralization

#### 9. `logger.go` - TUI-Specific Logging
**Purpose**: Manages logging specific to TUI operations.

**Key Components**:
- `TUIWriter`: Custom log writer for TUI
- `SetupTUILogger()`: TUI logger initialization
- `FormatLogMessages()`: Log message formatting

**Responsibilities**:
- TUI-specific log capture
- Log message formatting
- Logger coordination with main application

#### 10. `tui.go` - Main Orchestrator
**Purpose**: Main entry point and orchestration logic for the TUI.

**Key Components**:
- `New()`: TUI model initialization
- `Update()`: Main update loop handler
- `View()`: Main view renderer
- `StartTUI()`: TUI application launcher

**Responsibilities**:
- Application initialization
- Event loop coordination
- Module integration
- Application lifecycle management

## Data Flow

### 1. Initialization Flow
```
main.go → StartTUI() → New() → InitializeMenus() → Model creation
```

### 2. User Input Flow
```
Key Press → Update() → handleKeyPress() → CommandHandler → Command Execution
```

### 3. Command Execution Flow
```
CommandHandler → CommandRunner → Cobra Command → Output Capture → View Update
```

### 4. Form Processing Flow
```
Form Creation → User Input → Validation → Data Extraction → Command Execution
```

### 5. Rendering Flow
```
Model State → ViewRenderer → Specialized Renderers → Styled Output
```

## Key Design Decisions

### 1. Modular Architecture
- **Rationale**: Improves maintainability and testability
- **Implementation**: Each module has a specific responsibility
- **Benefits**: Easy to modify individual components without affecting others

### 2. Command Isolation
- **Rationale**: Prevents interference between TUI and CLI operations
- **Implementation**: CommandRunner creates isolated execution environments
- **Benefits**: Reliable command execution with proper output capture

### 3. State Management
- **Rationale**: Centralized state reduces complexity
- **Implementation**: Model struct with accessor methods
- **Benefits**: Consistent state access and modification

### 4. View Separation
- **Rationale**: Separates presentation logic from business logic
- **Implementation**: Dedicated ViewRenderer module
- **Benefits**: Easy to modify visual appearance without affecting functionality

## Integration Points

### 1. Cobra Commands
- **Interface**: CommandRunner module
- **Method**: Command isolation and execution
- **Considerations**: Preserves original command functionality

### 2. Configuration System
- **Interface**: Direct access to config package
- **Method**: Form pre-population and updates
- **Considerations**: Maintains configuration consistency

### 3. Logging System
- **Interface**: TUIWriter and logger integration
- **Method**: Custom log capture and formatting
- **Considerations**: Coordinated logging between TUI and application

## Extension Guidelines

### Adding New Commands
1. Update menu definitions in `menu.go`
2. Add form creation logic in `forms.go` if input required
3. Add command handler in `handlers.go`
4. Add execution method in `commandrunner.go` if needed
5. Update command classification functions

### Adding New Views
1. Add view state to Model in `model.go`
2. Add view rendering logic in `view.go`
3. Add navigation logic in `handlers.go`
4. Update menu definitions if needed

### Modifying Styles
1. Update style definitions in `styles.go`
2. Modify rendering logic in `view.go` if needed
3. Ensure consistency across all components

## Performance Considerations

- **Command Execution**: Isolated execution prevents blocking UI
- **Output Capture**: Efficient buffer management for large outputs
- **Memory Management**: Proper cleanup of temporary resources
- **Rendering**: Optimized view rendering with minimal redraws

## Testing Strategy

- **Unit Tests**: Each module can be tested independently
- **Integration Tests**: Test module interactions
- **UI Tests**: Test user interaction flows
- **Command Tests**: Verify command execution isolation

This architecture provides a solid foundation for the CookieFarm TUI while maintaining flexibility for future enhancements and modifications.
