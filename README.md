# Password Manager

A secure, terminal-based password manager built in Go using GPG encryption and a beautiful TUI interface.

## 🔐 Features

- **GPG-based encryption** using ProtonMail's gopenpgp library with RFC9580 OpenPGP profile
- **Beautiful TUI interface** built with Bubble Tea and Lipgloss with consistent purple styling
- **Complete password management** - add, list, view, and delete passwords
- **Interactive password forms** with multi-field input and validation
- **Secure password generation** with configurable options and strength evaluation
- **Password visibility toggle** - view passwords securely with explicit reveal
- **Safe deletion workflow** with confirmation dialogs and warnings
- **Real-time list updates** - newly added/deleted passwords immediately reflected
- **Master password protection** with validation phrase verification during setup
- **Password requirements enforcement** - 8+ chars for master, 12+ for phrase
- **Secure local storage** in `~/.password-manager-store/` with proper permissions

## 🛡️ Security Model

This password manager uses a dual-layer security approach:

1. **Master Password**: Encrypts all stored data (minimum 8 characters)
2. **Validation Phrase**: Used to verify master password correctness (minimum 12 characters)
3. **GPG Encryption**: All data encrypted using RFC9580 OpenPGP standard
4. **Local Storage**: No cloud dependencies, all data stored locally

## 📋 Requirements

- **Go 1.24.6** or later
- Unix-like system (Linux, macOS)
- Terminal with color support

## 🚀 Installation

1. **Clone the repository**:
   ```bash
   git clone <repository-url>
   cd password-manager
   ```

2. **Install dependencies**:
   ```bash
   go mod tidy
   ```

3. **Build the application**:
   ```bash
   go build -o password-manager
   ```

4. **Run the application**:
   ```bash
   ./password-manager
   ```

## 📖 Usage

### First-Time Setup

1. **Create Master Password**: Enter a secure password (minimum 8 characters) - visible during setup
2. **Create Validation Phrase**: Enter a random phrase (minimum 12 characters) 
3. The application creates the password store and initialization files

### Daily Usage

1. **Launch the application**
2. **Enter Master Password**: The same password created during setup (masked for security)
3. **Main Menu**: Choose from available password management options

### Main Menu Options

- **📋 List All Passwords**: Browse all stored passwords in a scrollable list
- **➕ Add New Password**: Create new password entries with comprehensive forms
- **🔍 Search Passwords**: Find specific passwords (coming soon)
- **📝 Edit Password**: Modify existing password entries (coming soon)
- **🗑️ Change Master Password**: Update your master password (coming soon)
- **📤 Export Passwords**: Export password data (coming soon)
- **❌ Quit**: Exit the application safely

### Password Management Workflow

#### Adding New Passwords
1. Select "Add New Password" from main menu
2. Fill out the multi-field form:
   - **Site/Service Name** (required): e.g., "Gmail", "GitHub"
   - **Username** (optional): Your username for the service
   - **Email** (optional): Associated email address
   - **URL** (optional): Website URL
   - **Password** (required): The actual password to store
3. Navigate between fields with Tab/Enter/Arrow keys
4. Submit when complete or cancel with Esc

#### Viewing Passwords
1. Select "List All Passwords" from main menu
2. Navigate the list with arrow keys or j/k
3. Press Enter/Space to view details of selected password
4. In detail view:
   - **v/Space**: Toggle password visibility (hidden by default)
   - **d**: Delete the password (requires confirmation)
   - **Esc/q**: Return to password list

#### Deleting Passwords
1. From password detail view, press 'd' to request deletion
2. Confirmation dialog appears with clear warnings
3. Choose Yes/No with arrow keys or y/n keys
4. Deleted passwords are immediately removed from the list

### Password Requirements

- **Master Password**: Minimum 8 characters for security
- **Validation Phrase**: Minimum 12 characters for better entropy  
- **Stored Passwords**: No minimum length (user's choice)
- All requirements validated during creation and cannot be bypassed

## 🏗️ Architecture

### Directory Structure

```
password-manager/
├── main.go              # Application entry point and menu handling
├── types/               # Shared type definitions  
│   └── types.go
├── fileio/              # File system operations and storage management
│   └── fileio.go
├── encryption/          # GPG encryption/decryption operations
│   └── encryption.go
├── menus/               # Authentication, password management, and flow control
│   └── menus.go
├── utils/               # Utility functions (generation, validation, formatting)
│   └── generator.go
└── ui/                  # User interface components
    ├── textinput/       # Styled text input with masking support
    │   └── textinput.go
    ├── menu/            # Main menu selection interface
    │   └── menu.go
    ├── form/            # Multi-field password entry forms
    │   └── form.go
    ├── list/            # Password list browsing interface
    │   └── list.go
    ├── detail/          # Password detail view with visibility controls
    │   └── detail.go
    └── confirm/         # Confirmation dialogs for destructive actions
        └── confirm.go
```

### Package Overview

- **`main`**: Application entry point, login flow, and menu action handling
- **`types`**: Shared data structures and application state management
- **`fileio`**: Password store management, file operations, and directory handling  
- **`encryption`**: GPG encryption/decryption using ProtonMail library with JSON serialization
- **`menus`**: Authentication logic, password CRUD operations, and UI flow orchestration
- **`utils`**: Password generation, strength evaluation, filename handling, and input sanitization
- **`ui/textinput`**: Styled text input components with password masking support
- **`ui/menu`**: Main menu interface with keyboard navigation and consistent styling
- **`ui/form`**: Multi-field password entry forms with validation and navigation
- **`ui/list`**: Password browsing interface with scrollable lists and entry selection
- **`ui/detail`**: Password detail view with visibility controls and action handling
- **`ui/confirm`**: Confirmation dialogs for destructive operations with safety defaults

### Data Storage

- **Location**: `$HOME/.password-manager-store/`
- **Structure**:
  ```
  ~/.password-manager-store/
  ├── .checker/
  │   └── init.gpg                    # Encrypted validation data
  ├── gmail_20240830_143022.gpg       # Example password entry
  ├── github_20240830_143045.gpg      # Example password entry
  └── banking_20240830_143103.gpg     # Example password entry
  ```
- **File Naming**: `sitename_YYYYMMDD_HHMMSS.gpg` for uniqueness and sorting
- **Permissions**: Directory created with 0750 permissions for security

### Encryption Details

- **Algorithm**: OpenPGP (RFC9580 profile)
- **Library**: ProtonMail gopenpgp v3
- **Key Derivation**: Password-based encryption
- **Data Format**: Armored GPG messages (.gpg files)

## 🎨 User Interface

The application features a modern TUI with:

- **Styled headers** with purple borders
- **Password masking** with bullet characters (•)
- **Error messages** in red with clear descriptions
- **Help text** with keyboard shortcuts
- **Screen clearing** for clean interactions
- **Responsive layout** with proper spacing

### Keyboard Shortcuts

#### Global Controls
- **Ctrl+C / Esc**: Quit application or cancel current action
- **Enter**: Submit input, continue, or confirm selection

#### Main Menu
- **↑↓ / j/k**: Navigate menu options  
- **Enter / Space**: Select menu item

#### Password List
- **↑↓ / j/k**: Navigate password entries
- **Enter / Space**: View selected password details
- **Home/End**: Jump to first/last entry

#### Password Detail View  
- **v / Space**: Toggle password visibility (show/hide)
- **d**: Delete password (with confirmation)
- **Esc / q / Backspace**: Return to password list

#### Password Forms (Add New)
- **Tab / Enter**: Move to next field
- **↑↓ / Shift+Tab**: Move to previous field  
- **Enter on last field**: Submit form
- **Esc**: Cancel and return to menu

#### Confirmation Dialogs
- **←→ / h/l**: Navigate between Yes/No buttons
- **y**: Confirm action (Yes)
- **n / Esc**: Cancel action (No)
- **Enter / Space**: Confirm selection

## 🔧 Development

### Dependencies

- **[gopenpgp](https://github.com/ProtonMail/gopenpgp)**: GPG encryption
- **[bubbletea](https://github.com/charmbracelet/bubbletea)**: TUI framework  
- **[bubbles](https://github.com/charmbracelet/bubbles)**: TUI components
- **[lipgloss](https://github.com/charmbracelet/lipgloss)**: Terminal styling

### Building from Source

```bash
# Install Go 1.24.6+
# Clone repository
go mod download
go build -o password-manager
```

### Code Style

- Follow Go conventions and `gofmt` formatting
- Use meaningful variable names
- Add documentation comments for exported functions
- Keep functions focused and testable

## 📝 To Be Added

### Upcoming Features
- **🔄 Change Master Password**: Update your master password with secure re-encryption of all stored passwords
- **📤 Export Passwords**: Export password data to various formats (CSV, JSON, plain text) with encryption options

### Maybe Features  
- **🔗 Automated GitHub Upload**: Optional encrypted backup to private GitHub repositories with sync capabilities

---

## 🔒 Security Considerations

- **Local Storage**: All data stays on your machine
- **No Network Access**: Application doesn't connect to internet
- **Strong Encryption**: Uses battle-tested GPG implementation
- **Password Requirements**: Enforced minimum lengths
- **Memory Safety**: Passwords cleared appropriately
- **File Permissions**: Store directory created with 0750 permissions

## 🐛 Troubleshooting

### Common Issues

1. **"Incorrect Password"**: Ensure you're using the exact master password from setup
2. **"Password too short"**: Master password needs 8+ characters, phrase needs 12+
3. **Permission errors**: Check that `~/.password-manager-store/` is writable

### Error Messages

- Clear, styled error messages guide you through issues
- Validation errors show exact requirements
- Press Enter after reading error messages to continue

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make changes with proper documentation
4. Test thoroughly
5. Submit a pull request

## 🙏 Acknowledgments

Built with excellent Go libraries:
- ProtonMail team for gopenpgp
- Charm team for Bubble Tea ecosystem
- Go team for the fantastic standard library
